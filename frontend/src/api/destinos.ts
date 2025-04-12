import axios from "axios";

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
  const response = await axios.post(`http://localhost:3000/destinos/buscar`, {
    destino,
    mes,
    embarque,
    categoria,
  });
  return response.data;
};

export const destinosPorCategoria = async () => {
  const response = await axios.get(
    `http://localhost:3000/destinos-por-categoria`
  );
  return response.data;
};
