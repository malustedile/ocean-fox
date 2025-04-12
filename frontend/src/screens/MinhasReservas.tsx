import { useEffect, useState } from "react";
import { minhasReservas } from "../api/destinos";
import { MdOutlineFileDownload } from "react-icons/md";

export const MinhasReservas = () => {
  const [reservas, setReservas] = useState<any[]>([]);
  const fetchReservas = async () => {
    const res = await minhasReservas();
    setReservas(res);
  };

  useEffect(() => {
    fetchReservas();
  }, []);
  return (
    <div className="flex justify-center ">
      <div className="flex flex-col gap-8 text-slate-700 py-4 w-[1000px]">
        {reservas.map((r) => (
          <ItemReserva
            id={r._id}
            data={r.criadoEm}
            destino={r.destino}
            status={r.status}
            passageiros={r.numeroPassageiros}
            valor={r.valorTotal || "0"}
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
  status: string;
  passageiros: number;
  valor: string;
}

const ItemReserva: React.FC<ItemReservaProps> = ({
  id,
  data,
  destino,
  status,
  passageiros,
  valor,
}) => {
  return (
    <div className="flex justify-between bg-white p-4 px-12 pr-8 rounded-xl relative">
      <div className=" flex gap-14">
        <div className="absolute top-[-20px] left-12  ">
          <div className=" bg-[#b3d4de] px-2 py-1 rounded-lg text-[#007090] font-bold">
            #{id.slice(-6)}
          </div>
        </div>
        <div>
          <div className="font-bold">Data do pedido:</div>
          <div className="text-slate-500">{formatDate(new Date(data))}</div>
        </div>
        <div>
          <div className="font-bold">Destino:</div>
          <div className="text-slate-500">{destino}</div>
        </div>
        <div>
          <div className="font-bold">Status:</div>
          {status === "PAGAMENTO_APROVADO" && (
            <div className="bg-[#d3f9d8] text-sm px-2 rounded-lg text-[#2b8a3e] font-bold ">
              Pagamento Aprovado
            </div>
          )}
          {status === "PAGAMENTO_REPROVADO" && (
            <div className="bg-[#ffe3e3] text-sm px-2 rounded-lg text-[#f03e3e] font-bold ">
              Pagamento Reprovado
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
      </div>
      {status === "PAGAMENTO_APROVADO" && (
        <div className="flex items-center justify-center">
          <div className="cursor-pointer">
            <MdOutlineFileDownload className="text-xl" />
          </div>
        </div>
      )}
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
