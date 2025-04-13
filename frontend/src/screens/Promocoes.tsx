import { useEffect, useState } from "react";
import { puxarPromocoes } from "../api/marketing";

export const Promocoes = () => {
  const [promocoes, setPromocoes] = useState({
    subscriptions: [],
    promotions: [],
  });
  const fetchPromotions = async () => {
    const res = await puxarPromocoes();
    setPromocoes(res);
  };

  useEffect(() => {
    fetchPromotions();
  }, []);
  return (
    <div className="">
      {promocoes?.subscriptions?.length === 0 && (
        <div className="h-[40vh] flex items-center justify-center">
          <div className="flex items-center flex-col">
            <div className="font-bold text-2xl">Nenhuma Promoção</div>
            <div>Se inscreva e comece a receber promoções</div>
          </div>
        </div>
      )}
      {promocoes?.promotions?.length > 0 && (
        <PromotionsScreen promotions={promocoes} />
      )}
    </div>
  );
};

export const PromotionsScreen = ({ promotions }: { promotions: any }) => {
  return (
    <div>
      <div>
        <div className="font-bold ">Minhas Inscrições</div>
        <div className="mx-2 mb-4">
          {promotions.subscriptions?.map((s: any) => (
            <div className="bg-[#99e9f2] text-[#0c8599] text-xs p-1 px-4 rounded-lg font-bold inline-block mr-2 mt-1">
              {s.destino}
            </div>
          ))}
        </div>
      </div>
      <div>
        <div className="font-bold ">Promoçoes</div>
        <div className="flex flex-col gap-8 my-2 mt-8">
          {promotions.promotions?.map((s: any) => (
            <div className="flex justify-between bg-white p-4 px-12 pr-8 rounded-xl relative w-full">
              <div className=" grid grid-cols-[1fr_2fr] gap-14 w-full">
                <div className="absolute top-[-20px] left-12  ">
                  <div className=" bg-[#b3d4de] px-2 py-1 rounded-lg text-[#007090] font-bold">
                    #{s._id.slice(-6)}
                  </div>
                </div>
                <div className="">
                  <div className="font-bold">Destino</div>
                  <div className="">{s.destino}</div>
                </div>
                <div className="">
                  <div className="font-bold">Descrição:</div>
                  <div className="">{s.mensagem}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
