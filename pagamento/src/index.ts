import { Elysia } from "elysia";
import amqp from "amqplib";
import { readFileSync } from "fs";
import { createSign } from "crypto";
import { MongoClient, ObjectId } from "mongodb";

interface reservaPayload {
  id: string;
  destino: string;
  dataEmbarque: string;
  numeroPassageiros: number;
  numeroCabines: number;
  linkPagamento: string;
  status: string;
  criadoEm: string;
}

const client = new MongoClient("mongodb://root:exemplo123@meu_mongodb:27017");
const db = client.db("ocean-fox");
const reservas = db.collection("reservas");

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
  async (msg: any) => {
    if (msg) {
      const reserva = JSON.parse(msg.content.toString()) as reservaPayload;
      console.log("Reserva recebida:", reserva);
      channelReserva.ack(msg);

      const pagamentoAprovado = Math.random() > 0.5;
      reserva.status = pagamentoAprovado
        ? "PAGAMENTO_APROVADO"
        : "PAGAMENTO_REPROVADO";

      const chavePrivada = readFileSync("./pagamento", "utf-8");
      const signer = createSign("sha256");
      signer.update(JSON.stringify(reserva));
      signer.end();

      const assinatura = signer.sign(
        { key: chavePrivada, passphrase: "your-passphrase" },
        "base64"
      );
      const payload = Buffer.from(
        JSON.stringify({
          reserva,
          assinatura,
        })
      );
      Buffer.from(
        JSON.stringify({
          reserva,
          assinatura,
        })
      );
      if (pagamentoAprovado) {
        channelPagamentoAprovado.sendToQueue("pagamento-aprovado", payload);
        console.log("Pagamento aprovado:", reserva);
      } else {
        channelPagamentoRecusado.sendToQueue("pagamento-recusado", payload);

        console.log("Pagamento recusado:", reserva);
      }

      const response = await reservas.updateOne(
        { _id: new ObjectId(reserva.id) },
        {
          $set: {
            status: reserva.status,
            assinatura,
          },
        }
      );

      console.log(response);
    }
  },
  { noAck: false }
);

const app = new Elysia().get("/", () => "Hello Elysia").listen(3001);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
