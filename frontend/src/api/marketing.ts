import axios from "axios";

const client = axios.create({
  baseURL: "http://localhost:3004",
  withCredentials: true,
});

export const puxarPromocoes = async () => {
  const response = await client.get("/minhas-inscricoes");
  return response.data;
};

export const inscrever = async () => {
  const response = await client.post("/inscrever");
  return response.data;
};

export const cancelarInscricao = async () => {
  const response = await client.post("/cancelar");
  return response.data;
};
