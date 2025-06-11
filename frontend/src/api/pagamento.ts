import axios from "axios";

const client = axios.create({
  baseURL: "http://localhost:8000",
  withCredentials: true,
});

interface PagarReservaDto {
  idReserva: string;
  valorTotal: number;
}

export const pagarReserva = async (dto: PagarReservaDto) => {
  const response = await client.post("/pagar", dto);
  return response.data;
};
