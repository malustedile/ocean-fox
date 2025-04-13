import axios from "axios";

const client = axios.create({
  baseURL: "http://localhost:3004",
  withCredentials: true,
});

export const puxarPromocoes = async () => {
  const response = await client.get("/minhas-inscricoes");
  return response.data;
};

export const inscrever = async (destino: string) => {
  const response = await client.post("/inscrever", { destino });
  return response.data;
};
