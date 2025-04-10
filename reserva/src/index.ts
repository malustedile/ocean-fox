import { Elysia } from "elysia";
import { MongoClient } from "mongodb";

const client = new MongoClient("mongodb://root:exemplo123@localhost:27017");
const db = client.db("ocean-fox");
const destinos = db.collection("destinos");

await client.connect();

interface destinosDto {
  nome: string;
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
  dataEmbarque: string;
  numeroPassageiros: number;
  numeroCabines: number;
}

const app = new Elysia()
  .get("/", () => "Hello Elysia")

  // Endpoint de cadastro
  .post("/destinos", async (body: destinosDto) => {
    const { nome, descricao } = body ?? {};

    if (!nome || !descricao) {
      return { erro: "Campos 'nome' e 'descricao' são obrigatórios." };
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
  .post("/destinos/reservar", async (body: reservaDto) => {
    const { dataEmbarque, numeroPassageiros, numeroCabines } = body ?? {};
    if (!dataEmbarque || !numeroPassageiros || !numeroCabines) {
      return {
        erro: "Campos 'dataEmbarque', 'numeroPassageiros' e 'numeroCabines' são obrigatórios.",
      };
    }
  })

  .listen(3000);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
