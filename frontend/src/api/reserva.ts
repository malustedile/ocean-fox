import axios from "axios";

const client = axios.create({
  baseURL: "http://localhost:3000",
  withCredentials: true,
});

export interface ReservaDto {
  destino: string;
  dataEmbarque: string;
  numeroPassageiros: number;
  numeroCabines: number;
  valorTotal: number;
}

export const minhasReservas = async () => {
  const response = await client.get(`/minhas-reservas`);
  return response.data;
};

export const reservarViagem = async (reserva: ReservaDto) => {
  const response = await client.post(`destinos/reservar`, reserva);
  return response.data;
};
