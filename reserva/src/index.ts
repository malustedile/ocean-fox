import { Elysia } from "elysia";
import { MongoClient } from "mongodb";
import amqp from "amqplib";

const client = new MongoClient("mongodb://root:exemplo123@localhost:27017");
const db = client.db("ocean-fox");
const destinos = db.collection("destinos");

const rabbit = await amqp.connect("amqp://localhost");
const channel = await rabbit.createChannel();
await channel.assertQueue("reserva-criada", { durable: true });

const channelPagamentoAprovado = await rabbit.createChannel();
await channelPagamentoAprovado.assertQueue("pagamento-aprovado", {
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
}

const app = new Elysia()
  .get("/", () => "Hello Elysia")

  // Endpoint de cadastro
  .post("/destinos", async ({ body }: { body: destinosDto }) => {
    const { nome, categoria, descricao } = body ?? {};

    if (!nome || !descricao || !categoria) {
      return {
        erro: "Campos 'nome', 'categoria' e 'descricao' são obrigatórios.",
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
  .get("/destinos/buscar", async ({ query }) => {
    const { destino, mes, embarque } = query;

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

    const resultados = await destinos.find(filtro).toArray();

    return resultados;
  })

  // Endpoint de efetuar reserva
  .post("/destinos/reservar", async ({ body }: { body: reservaDto }) => {
    const { destino, dataEmbarque, numeroPassageiros, numeroCabines } =
      body ?? {};

    if (!destino || !dataEmbarque || !numeroPassageiros || !numeroCabines) {
      return {
        erro: "Campos 'destino', 'dataEmbarque', 'numeroPassageiros' e 'numeroCabines' são obrigatórios.",
      };
    }

    // 🔗 Simulação do link de pagamento (mock)
    const linkPagamento = `https://pagamento.fake/checkout?token=${crypto.randomUUID()}`;

    // 📤 Publica na fila
    const reservaPayload = {
      destino,
      dataEmbarque,
      numeroPassageiros,
      numeroCabines,
      linkPagamento,
      status: "AGUARDANDO_PAGAMENTO",
      criadoEm: new Date().toISOString(),
    };

    channel.sendToQueue(
      "reserva-criada",
      Buffer.from(JSON.stringify(reservaPayload)),
      { persistent: true }
    );

    return {
      mensagem: "Reserva registrada. Link de pagamento gerado.",
      linkPagamento,
    };
  })

  .listen(3000);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
