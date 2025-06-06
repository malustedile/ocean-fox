import { useEffect, useState } from "react";
import { cancelarInscricao, inscrever, puxarPromocoes } from "../api/marketing";

export const Promocoes = () => {
  const [promocoes, setPromocoes] = useState<{
    subscriptions: any[];
    promotions: any[];
  }>({
    subscriptions: [],
    promotions: [],
  });
  const fetchPromotions = async () => {
    const res = await puxarPromocoes();
    setPromocoes(res);
  };

  useEffect(() => {
    fetchPromotions();
    const interval = setInterval(() => {
      fetchPromotions();
    }, 500);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="">
      <PromotionsScreen promotions={promocoes} />
    </div>
  );
};

export const PromotionsScreen = ({ promotions }: { promotions: any }) => {
  const subscribe = async () => {
    await inscrever();
  };
  const unsubscribe = async () => {
    await cancelarInscricao();
  };

  return (
    <div>
      <div className="flex flex-col gap-2">
        <div className="font-bold">Minhas Promoções</div>
        <div className="mx-2 mb-4">
          <button
            className="btn rounded-xl"
            onClick={
              promotions.subscriptions?.length > 0 ? unsubscribe : subscribe
            }
          >
            {promotions.subscriptions?.length > 0
              ? "Cancelar inscrição"
              : "Inscrever-se"}
          </button>
        </div>
      </div>
      <div>
        <div className="font-bold">Promoções</div>
        <div className="flex flex-col gap-8 my-2 mt-8">
          {promotions.promotions?.map((s: any) => (
            <div className="flex justify-between bg-white p-4 px-12 pr-8 rounded-xl relative w-full">
              <div className="grid grid-cols-[1fr_2fr] gap-14 w-full">
                <div>
                  <div className="font-bold">Descrição:</div>
                  <div>{s.mensagem}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
