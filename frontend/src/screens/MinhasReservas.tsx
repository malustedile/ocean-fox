import { cancelarViagem } from "../api/reserva";
import { FaRegCheckCircle } from "react-icons/fa";
import { FaRegCircleXmark } from "react-icons/fa6";
import { pagarReserva } from "../api/pagamento";
import React from "react";

interface MinhasReservasDto {
  reservas: any;
}

export const MinhasReservas: React.FC<MinhasReservasDto> = ({ reservas }) => {
  return (
    <div className="flex justify-center ">
      <div className="flex flex-col gap-8 text-slate-700 py-4 w-[1000px]">
        {reservas.map((r: any) => (
          <ItemReserva
            id={r.id}
            data={r.criadoEm}
            destino={r.destino}
            statusPagamento={r.statusPagamento}
            status={r.status}
            passageiros={r.numeroPassageiros}
            valor={r.valorTotal || "0"}
            handleCancelar={async () => {
              await cancelarViagem(r.id);
            }}
            handlePagar={async () => {
              await pagarReserva({ idReserva: r.id, valorTotal: r.valorTotal });
              await new Promise((resolve) => setTimeout(resolve, 1000));
            }}
          />
        ))}
      </div>
    </div>
  );
};

interface ItemReservaProps {
  id: string;
  destino: string;
  data: string;
  statusPagamento: string;
  status: string;
  passageiros: number;
  valor: string;
  handleCancelar: () => void;
  handlePagar: () => void;
}

const ItemReserva: React.FC<ItemReservaProps> = ({
  id,
  data,
  destino,
  statusPagamento,
  passageiros,
  valor,
  status,
  handleCancelar,
  handlePagar,
}) => {
  return (
    <div className="flex justify-between bg-white p-4 px-12 pr-8 rounded-xl relative">
      <div className=" grid grid-cols-[8fr_10fr_8fr_6fr_3fr_1fr] gap-6 w-full">
        <div className="absolute top-[-20px] left-12  ">
          <div className=" bg-[#b3d4de] px-2 py-1 rounded-lg text-[#007090] font-bold">
            #{id?.slice(-6)}
          </div>
        </div>
        <div>
          <div className="font-bold">Data do pedido:</div>
          <div className="text-slate-500">{formatDate(new Date(data))}</div>
        </div>
        <div className="truncate">
          <div className="font-bold">Destino:</div>
          <div className="text-slate-500 text-ellipsis">{destino}</div>
        </div>
        <div className="flex flex-col w-full">
          <div className="font-bold">Status:</div>
          {status === "cancelado" && (
            <div className="flex  ">
              <div className="bg-[#ffe3e3] text-sm px-2 rounded-lg text-[#f03e3e] font-bold">
                Cancelado
              </div>
            </div>
          )}
          {status !== "cancelado" && (
            <div>
              {statusPagamento === "PAGAMENTO_APROVADO" && (
                <div className="flex  ">
                  <div className="flex bg-[#d3f9d8] text-sm px-2 rounded-lg text-[#2b8a3e] font-bold ">
                    Pagamento Aprovado
                  </div>
                </div>
              )}
              {statusPagamento === "PAGAMENTO_RECUSADO" && (
                <div className="flex  ">
                  <div className="bg-[#ffe3e3] text-sm px-2 rounded-lg text-[#f03e3e] font-bold">
                    Pagamento Reprovado
                  </div>
                </div>
              )}
              {statusPagamento === "AGUARDANDO_PAGAMENTO" && (
                <div className="flex  ">
                  <div className="bg-[#fff3bf] text-sm px-2 rounded-lg text-[#e67700] font-bold truncate">
                    Aguardando Pagamento
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
        <div>
          <div className="font-bold">Passageiros:</div>
          <div className="text-slate-500">{passageiros} passageiros</div>
        </div>

        <div>
          <div className="font-bold">Valor:</div>
          <div className="text-black">
            {new Intl.NumberFormat("pt-BR", {
              style: "currency",
              currency: "BRL",
            }).format(Number(valor))}
          </div>
        </div>
        <div className="flex items-center justify-center">
          {/* {status !== "cancelado" &&
            statusPagamento === "PAGAMENTO_APROVADO" && (
              <div className="flex items-center justify-center">
                <div className="cursor-pointer">
                  <MdOutlineFileDownload className="text-xl" />
                </div>
              </div>
            )} */}
          {status !== "cancelado" &&
            statusPagamento === "AGUARDANDO_PAGAMENTO" && (
              <div className="flex items-center justify-center mr-2">
                <button
                  className="cursor-pointer"
                  onClick={() => handlePagar()}
                >
                  <FaRegCheckCircle className="text-xl text-green-500" />
                </button>
              </div>
            )}
          {status !== "cancelado" && (
            <div className="flex items-center justify-center">
              <button
                className="cursor-pointer"
                onClick={() => handleCancelar()}
              >
                <FaRegCircleXmark className="text-xl text-red-500" />
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

function formatDate(date: Date) {
  const pad = (n: any) => n.toString().padStart(2, "0");

  const dia = pad(date.getDate());
  const mes = pad(date.getMonth() + 1); // Janeiro = 0
  const ano = date.getFullYear();

  const horas = pad(date.getHours());
  const minutos = pad(date.getMinutes());
  const segundos = pad(date.getSeconds());

  return `${dia}/${mes}/${ano} ${horas}:${minutos}:${segundos}`;
}

// Exemplo de uso
const agora = new Date();
console.log(formatDate(agora)); // ex: 12/04/2025 15:43:21
