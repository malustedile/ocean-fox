import { Elysia } from "elysia";
import { MongoClient, ObjectId } from "mongodb";
import amqp from "amqplib";
import cors from "@elysiajs/cors";
import { createVerify } from "crypto";
import { readFileSync } from "fs";

const client = new MongoClient("mongodb://root:exemplo123@meu_mongodb:27017");
const db = client.db("ocean-fox");
const destinos = db.collection("destinos");
const reservas = db.collection("reservas");

const reservaExchange = "reserva-criada-exc";

const rabbit = await amqp.connect("amqp://rabbitmq");
const channel = await rabbit.createChannel();
await channel.assertQueue("reserva-criada", { durable: true });

await channel.assertExchange(reservaExchange, "fanout", {
  durable: false,
});

const pagamentoAprovadoExchange = "pagamento-aprovado-exc";
const channelPagamentoAprovado = await rabbit.createChannel();
await channelPagamentoAprovado.assertExchange(
  pagamentoAprovadoExchange,
  "direct",
  {
    durable: true,
  }
);
const q = await channelPagamentoAprovado.assertQueue("", {
  durable: true,
});

const channelPagamentoRecusado = await rabbit.createChannel();
await channelPagamentoRecusado.assertQueue("pagamento-recusado", {
  durable: true,
});

const channelBilheteGerado = await rabbit.createChannel();
await channelBilheteGerado.assertQueue("bilhete-gerado", {
  durable: true,
});

await client.connect();

channelPagamentoAprovado.bindQueue(
  q.queue,
  pagamentoAprovadoExchange,
  "reserva"
);

channelPagamentoAprovado.consume(
  q.queue,
  async (msg: any) => {
    const pedido = JSON.parse(msg.content.toString());
    const verifier = createVerify("sha256");
    verifier.update(JSON.stringify(pedido.reserva));
    verifier.end();

    const chavePublica = readFileSync("./pagamento.pub", "utf-8");
    const isValid = verifier.verify(chavePublica, pedido.assinatura, "base64");
    console.log("Assinatura válida:", isValid);
    console.log("Pedido recebido:", pedido);
    return await reservas.updateOne(
      {
        _id: new ObjectId(pedido.reserva.id),
      },
      {
        $set: {
          pagamentoValido: isValid,
        },
      }
    );
  },
  { noAck: true }
);

enum Categorias {
  BRAZIL = "Brasil",
  SOUTH_AMERICA = "América do Sul",
  CARIBBEAN = "Caribe",
  NORTH_AMERICA = "América do Norte",
  AFRICA = "África",
  MIDDLE_EAST = "Oriente Médio",
  ASIA = "Ásia",
  MEDITERRANEAN = "Mediterrâneo",
  SCANDINAVIA = "Escandinávia",
  OCEANIA = "Oceania",
}

interface destinosDto {
  nome: string;
  categoria: Categorias;
  descricao: {
    datasDisponiveis: string[];
    navio: string;
    embarque: string;
    desembarque: string;
    lugaresVisitados: string[];
    noites: number;
    valorPorPessoa: number;
  };
}

interface reservaDto {
  destino: string;
  dataEmbarque: string;
  numeroPassageiros: number;
  numeroCabines: number;
  valorTotal: number;
}

interface filtrosDto {
  destino?: string;
  mes?: string;
  embarque?: string;
  categoria?: Categorias;
}

const app = new Elysia()
  .use(cors())
  .get("/", () => {
    return "Hello Elysia";
  })

  .get("/minhas-reservas", async ({ cookie }) => {
    return await reservas
      .find({
        sessionId: cookie.sessionId.value,
      })
      .toArray();
  })

  // Endpoint de cadastro

  .post("/destinos", async ({ body }) => {
    const { nome, categoria, descricao } = (body as any) ?? {};
    if (!nome || !descricao || !categoria) {
      return {
        erro: "Campos 'nome', 'descricao' e 'categoria' são obrigatórios.",
      };
    }

    const {
      datasDisponiveis,
      navio,
      embarque,
      desembarque,
      lugaresVisitados,
      noites,
      valorPorPessoa,
    } = descricao;

    if (
      !Array.isArray(datasDisponiveis) ||
      !navio ||
      !embarque ||
      !desembarque ||
      !Array.isArray(lugaresVisitados) ||
      typeof noites !== "number" ||
      typeof valorPorPessoa !== "number"
    ) {
      return { erro: "Dados incompletos ou inválidos na descrição." };
    }

    const resultado = await destinos.insertOne({
      nome,
      categoria,
      descricao: {
        datasDisponiveis,
        navio,
        embarque,
        desembarque,
        lugaresVisitados,
        noites,
        valorPorPessoa,
      },
    });

    return {
      mensagem: "Destino adicionado com sucesso",
      id: resultado.insertedId,
    };
  })

  // Endpoint de consulta por nome, mês e porto de embarque
  .post("/destinos/buscar", async ({ body }: { body: filtrosDto }) => {
    const { destino, mes, embarque, categoria } = body;

    const filtro: any = {};

    if (destino) {
      filtro["descricao.lugaresVisitados"] = {
        $elemMatch: { $regex: destino, $options: "i" },
      };
    }

    if (embarque) {
      filtro["descricao.embarque"] = { $regex: embarque, $options: "i" };
    }

    if (mes) {
      const mesNum = parseInt(mes);
      if (!isNaN(mesNum) && mesNum >= 1 && mesNum <= 12) {
        filtro["descricao.datasDisponiveis"] = {
          $elemMatch: {
            $regex: `-${String(mesNum).padStart(2, "0")}-`, // ex: "-06-"
          },
        };
      }
    }

    if (categoria) {
      filtro["categoria"] = { $regex: categoria, $options: "i" };
    }

    const resultados = await destinos.find(filtro).toArray();

    return resultados;
  })

  .get("/destinos-por-categoria", async () => {
    const categorias = Object.values(Categorias);

    const resultados = await Promise.all(
      categorias.map(async (categoria) => {
        const count = await destinos.countDocuments({ categoria });
        return { categoria, quantidade: count };
      })
    );

    return resultados;
  })

  // Endpoint de efetuar reserva
  .post("/destinos/reservar", async ({ body, cookie }) => {
    const {
      destino,
      dataEmbarque,
      numeroPassageiros,
      numeroCabines,
      valorTotal,
    } = (body as reservaDto) ?? {};
    const sessionId = cookie.sessionId.value;
    if (
      !destino ||
      !dataEmbarque ||
      !numeroPassageiros ||
      !numeroCabines ||
      !valorTotal
    ) {
      return {
        erro: "Campos 'destino', 'dataEmbarque', 'numeroPassageiros', 'numeroCabines' e 'valorTotal' são obrigatórios.",
      };
    }

    // 🔗 Simulação do link de pagamento (mock)
    const linkPagamento = `https://pagamento.fake/checkout?token=${crypto.randomUUID()}`;

    // 📤 Publica na fila
    const reservaPayload = {
      destino,
      sessionId,
      dataEmbarque,
      numeroPassageiros,
      numeroCabines,
      valorTotal,
      linkPagamento,
      status: "AGUARDANDO_PAGAMENTO",
      bilhete: null,
      criadoEm: new Date().toISOString(),
    };

    const reserva = await reservas.insertOne(reservaPayload);

    channel.publish(
      reservaExchange,
      "",
      Buffer.from(
        JSON.stringify({ id: reserva.insertedId, ...reservaPayload })
      ),
      { persistent: true }
    );

    return {
      mensagem: "Reserva registrada. Link de pagamento gerado.",
      linkPagamento,
      reserva,
    };
  })

  .listen(3000);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
