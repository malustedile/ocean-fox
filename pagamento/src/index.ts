import { Elysia } from "elysia";
import amqp from "amqplib";
import { readFileSync } from "fs";
import { createSign } from "crypto";

interface reservaPayload {
  destino: string;
  dataEmbarque: string;
  numeroPassageiros: number;
  numeroCabines: number;
  linkPagamento: string;
  status: string;
  criadoEm: string;
}

const reservaExchange = "reserva-criada-exc";

const rabbit = await amqp.connect("amqp://rabbitmq");
const channelReserva = await rabbit.createChannel();
const channelPagamentoAprovado = await rabbit.createChannel();
const channelPagamentoRecusado = await rabbit.createChannel();

await channelReserva.assertExchange(reservaExchange, "fanout", {
  durable: false,
});
const reservaQueue = await channelReserva.assertQueue("reserva-criada", {
  durable: true,
});
await channelPagamentoAprovado.assertQueue("pagamento-aprovado", {
  durable: true,
});
await channelPagamentoRecusado.assertQueue("pagamento-recusado", {
  durable: true,
});

await channelReserva.bindQueue(reservaQueue.queue, reservaExchange, "");

channelReserva.consume(
  "reserva-criada",
  (msg: any) => {
    if (msg) {
      const reserva = JSON.parse(msg.content.toString()) as reservaPayload;
      console.log("Reserva recebida:", reserva);
      channelReserva.ack(msg);

      const pagamentoAprovado = Math.random() > 0.5;
      reserva.status = pagamentoAprovado ? "aprovado" : "recusado";

      const chavePrivada = readFileSync("./pagamento", "utf-8");
      const signer = createSign("sha256");
      signer.update(JSON.stringify(reserva));
      signer.end();

      const assinatura = signer.sign(
        { key: chavePrivada, passphrase: "your-passphrase" },
        "base64"
      );

      if (pagamentoAprovado) {
        channelPagamentoAprovado.sendToQueue(
          "pagamento-aprovado",
          Buffer.from(
            JSON.stringify({
              reserva,
              assinatura,
            })
          )
        );
        console.log("Pagamento aprovado:", reserva);
      } else {
        channelPagamentoRecusado.sendToQueue(
          "pagamento-recusado",
          Buffer.from(
            JSON.stringify({
              reserva,
              assinatura,
            })
          )
        );
        console.log("Pagamento recusado:", reserva);
      }
    }
  },
  { noAck: false }
);

const app = new Elysia().get("/", () => "Hello Elysia").listen(3001);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
