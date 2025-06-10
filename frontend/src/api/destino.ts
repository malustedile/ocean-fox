import axios from "axios";

const client = axios.create({
  baseURL: "http://localhost:8080",
  withCredentials: true,
});

export const puxarDestinos = async ({
  destino,
  mes,
  embarque,
  categoria,
}: {
  destino?: string;
  mes?: string;
  embarque?: string;
  categoria?: string;
}) => {
  const response = await client.post(`/destinos/buscar`, {
    destino,
    mes,
    embarque,
    categoria,
  });
  return response.data;
};

export const destinosPorCategoria = async () => {
  const response = await client.get(`/destinos-por-categoria`);
  return response.data;
};
