import { Elysia } from "elysia";
import amqp from "amqplib";
import { MongoClient } from "mongodb";
import cors from "@elysiajs/cors";

const client = new MongoClient("mongodb://root:exemplo123@meu_mongodb:27017");
const db = client.db("ocean-fox");
const inscricoes = db.collection("inscricoes");

const promocoes = db.collection("promocoes");

const rabbit = await amqp.connect("amqp://rabbitmq");
const channel = await rabbit.createChannel();

await channel.assertQueue("pagamento-aprovado", {
  durable: true,
});

const incricoesAtivas = await inscricoes.find().toArray();

const consumeSubscription = (msg: any, sessionId: string) => {
  const promocao = JSON.parse(msg.content.toString());
  console.log({
    sessionId,
    mensagem: promocao.mensagem,
    criadoEm: promocao.criadoEm,
    destino: promocao.destino,
  });
  promocoes.insertOne({
    sessionId,
    mensagem: promocao.mensagem,
    criadoEm: promocao.criadoEm,
    destino: promocao.destino,
  });
};

for (const inscricao of incricoesAtivas) {
  const exchange = `promocoes-${inscricao.destino}`;
  await channel.assertExchange(exchange, "fanout", { durable: false });
  const { queue } = await channel.assertQueue("", { exclusive: true });
  await channel.bindQueue(queue, exchange, "");
  channel.consume(queue, (msg) => {
    if (msg !== null) consumeSubscription(msg, inscricao.sessionId);
  });
}

const app = new Elysia()
  .use(cors())
  .get("/", () => "Hello Elysia")
  .get("/minhas-inscricoes", async ({ cookie }) => {
    const mySubscriptions = await inscricoes
      .find({ sessionId: cookie.sessionId.value })
      .toArray();
    const myPromotions = await promocoes
      .find({ sessionId: cookie.sessionId.value })
      .toArray();

    return { subscriptions: mySubscriptions, promotions: myPromotions };
  })
  .post("/inscrever", async ({ body, cookie }) => {
    const { destino } = body as any;
    const exchange = `promocoes-${destino}`;
    await channel.assertExchange(exchange, "fanout", { durable: false });
    await inscricoes.insertOne({
      sessionId: cookie.sessionId.value,
      destino,
      criadoEm: new Date(),
    });
    // Cria uma fila exclusiva e temporária para o cliente
    const { queue } = await channel.assertQueue("", { exclusive: true });
    await channel.bindQueue(queue, exchange, "");
    channel.consume(queue, (msg) => {
      if (msg !== null)
        consumeSubscription(msg, cookie.sessionId.value as string);
    });

    return { success: true };
  })
  .post("/promocao", async ({ body }) => {
    const { destino, mensagem } = body as any;
    const exchange = `promocoes-${destino}`;
    const conn = await amqp.connect("amqp://rabbitmq");
    const channel = await conn.createChannel();
    const promocao = { mensagem, destino, criadoEm: new Date() };
    await channel.assertExchange(exchange, "fanout", { durable: false });
    channel.publish(exchange, "", Buffer.from(JSON.stringify(promocao)));
    console.log(
      `[Publisher] Promoção enviada para exchange ${exchange}: ${mensagem}`
    );

    setTimeout(() => conn.close(), 500);
  })
  .listen(3004);
console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
