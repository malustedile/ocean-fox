import { Elysia } from "elysia";
import amqp from "amqplib";
import { createVerify } from "crypto";
import { readFileSync } from "fs";

interface reservaPayload {
  destino: string;
  dataEmbarque: string;
  numeroPassageiros: number;
  numeroCabines: number;
  linkPagamento: string;
  status: string;
  criadoEm: string;
}
interface pedidoPayload {
  reserva: reservaPayload;
  assinatura: string;
}

interface bilhete {
  idReserva: string;
  destino: string;
  dataEmbarque: string;
  numeroPassageiros: number;
  numeroCabines: number;
  criadoEm: string;
}

const rabbit = await amqp.connect("amqp://rabbitmq");
const channelPagamentoAprovado = await rabbit.createChannel();
const channelBilheteGerado = await rabbit.createChannel();

await channelPagamentoAprovado.assertQueue("pagamento-aprovado", {
  durable: true,
});
await channelBilheteGerado.assertQueue("bilhete-gerado", {
  durable: true,
});

channelPagamentoAprovado.consume(
  "pagamento-aprovado",
  (msg: any) => {
    if (msg) {
      const pedido = JSON.parse(msg.content.toString()) as pedidoPayload;
      console.log(pedido);

      const verifier = createVerify("sha256");
      verifier.update(JSON.stringify(pedido.reserva));
      verifier.end();

      const chavePublica = readFileSync("./pagamento.pub", "utf-8");
      const isValid = verifier.verify(
        chavePublica,
        pedido.assinatura,
        "base64"
      );

      if (isValid) {
        const bilhete: bilhete = {
          idReserva: pedido.reserva.linkPagamento,
          destino: pedido.reserva.destino,
          dataEmbarque: pedido.reserva.dataEmbarque,
          numeroPassageiros: pedido.reserva.numeroPassageiros,
          numeroCabines: pedido.reserva.numeroCabines,
          criadoEm: new Date().toISOString(),
        };

        channelBilheteGerado.sendToQueue(
          "bilhete-gerado",
          Buffer.from(JSON.stringify(bilhete))
        );
        console.log("Bilhete gerado:", bilhete);
      } else {
        console.error("Assinatura inválida");
      }
    }
  },
  {
    noAck: false,
  }
);

const app = new Elysia().get("/", () => "Hello Elysia").listen(3002);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
