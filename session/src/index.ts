import { Elysia } from "elysia";
import { MongoClient, ObjectId } from "mongodb";
import amqp from "amqplib";

const reservaExchange = "reserva-criada-exc";

// Conexão com MongoDB
const client = new MongoClient("mongodb://root:exemplo123@meu_mongodb:27017");
const db = client.db("ocean-fox");
const sessions = db.collection("sessions");
const reservas = db.collection("reservas");

const rabbit = await amqp.connect("amqp://rabbitmq");
const channel = await rabbit.createChannel();
const reservaQueue = await channel.assertQueue("reserva-criada-session", {
  durable: true,
});

await channel.assertExchange(reservaExchange, "fanout", {
  durable: false,
});

await channel.bindQueue(reservaQueue.queue, reservaExchange, "");

channel.consume(
  "reserva-criada-session",
  (msg: any) => {
    const reserva = JSON.parse(msg.content.toString());
    console.log({ reserva });
  },
  { noAck: false }
);

const app = new Elysia()
  .get("/", async () => {})

  .get("/session", async ({ cookie }) => {
    let sessionId = cookie.sessionId.value;
    let sessionData;

    if (sessionId) {
      sessionData = await sessions.findOne({ _id: new ObjectId(sessionId) });
    }

    // Se não encontrou a sessão ou não havia cookie
    if (!sessionData) {
      // Cria nova sessão
      const newSession = {
        createdAt: new Date(),
        data: {},
      };

      const result = await sessions.insertOne(newSession);
      sessionId = result.insertedId.toString();

      // Seta o cookie no navegador
      cookie.sessionId.set({
        value: sessionId,
        httpOnly: true,
        path: "/",
        maxAge: 60 * 60 * 24 * 7, // 7 dias
      });

      sessionData = newSession;
    }

    return {
      mensagem: "Sessão ativa",
      sessionId,
      sessionData,
    };
  })

  .listen(3005);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
