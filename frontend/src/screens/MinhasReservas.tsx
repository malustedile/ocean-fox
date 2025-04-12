import { minhasReservas } from "../api/destinos";

export const MinhasReservas = () => {
  const fetchReservas = async () => {
    const reservas = await minhasReservas();
    console.log(reservas);
  };

  return (
    <div className="flex justify-center ">
      <div className="flex flex-col gap-8 text-slate-700 py-4 w-[1000px]">
        <ItemReserva />
        <ItemReserva />

        <ItemReserva />
        <ItemReserva />
      </div>
    </div>
  );
};

const ItemReserva = () => {
  return (
    <div className="flex bg-white gap-16 p-4 px-12 rounded-xl relative">
      <div className="absolute top-[-20px] left-12  ">
        <div className=" bg-[#b3d4de] px-2 py-1 rounded-lg text-[#007090] font-bold">
          #123545
        </div>
      </div>
      <div>
        <div className="font-bold">Data do pedido:</div>
        <div className="text-slate-500">30/03/2025</div>
      </div>
      <div>
        <div className="font-bold">Destino:</div>
        <div className="text-slate-500">Cruzeiro Ilhas Gregas</div>
      </div>
      <div>
        <div className="font-bold">Status:</div>
        <div className="bg-[#d3f9d8] text-sm px-2 rounded-lg text-[#2b8a3e] font-bold ">
          Pagamento Aprovado
        </div>
      </div>
      <div>
        <div className="font-bold">Passageiros:</div>
        <div className="text-slate-500">3 passageiros</div>
      </div>

      <div>
        <div className="font-bold">Valor:</div>
        <div className="text-black">R$ 9.539,00</div>
      </div>
    </div>
  );
};
