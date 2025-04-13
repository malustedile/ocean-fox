import { MdOutlineEventAvailable, MdOutlinePinDrop } from "react-icons/md";
import { Trip } from "./FilteredTrips";
import { LuShip } from "react-icons/lu";
import { IoIosNotificationsOutline } from "react-icons/io";
import { inscrever } from "../api/marketing";
import { useState } from "react";

export const CruiseDetails = ({
  trip,
  passageiros,
  showSubscription,
  subscripted,
}: {
  trip: Trip;
  passageiros?: number;
  showSubscription?: boolean;
  subscripted?: boolean;
}) => {
  const [inscrito, setInscrito] = useState(subscripted);
  const subscribe = async () => {
    await inscrever(trip.nome);
    setInscrito(true);
  };

  return (
    <div className="flex flex-col w-full gap-2">
      <div className="flex flex-col gap-2">
        <div>
          <span className="flex flex-row items-center text-slate-500 text-sm gap-2">
            <MdOutlinePinDrop size={14} /> {trip.descricao.embarque}
          </span>
          <span className="flex flex-row items-center text-slate-500 text-sm gap-2">
            <MdOutlineEventAvailable size={14} />{" "}
            {trip.descricao.datasDisponiveis
              .map((data) =>
                new Date(data).toLocaleDateString("pt-BR", {
                  day: "2-digit",
                  month: "2-digit",
                  year: "numeric",
                })
              )
              .join(" - ")}
          </span>
          <span className="flex flex-row items-center text-slate-500 text-sm gap-2">
            <LuShip size={14} /> {trip.descricao.navio}
          </span>
        </div>
        <div className="flex flex-col w-full">
          <div className="text-slate-500 text-xs">Itinerário:</div>
          <div className="text-slate-600 text-sm truncate">
            {trip.descricao.lugaresVisitados.join(", ").toLocaleString()}
          </div>
        </div>
        <div className="flex justify-between w-full">
          <div className="flex items-center ">
            {showSubscription && !inscrito && (
              <div
                className="flex items-center gap-2 bg-[#f1f3f5] text-[#495057] text-xs p-1 px-4 rounded-lg  hover:bg-[#dee2e6]"
                onClick={(e) => {
                  e.stopPropagation();
                  subscribe();
                }}
              >
                <IoIosNotificationsOutline size={16} />
                Inscrever-se
              </div>
            )}

            {showSubscription && inscrito && (
              <div
                className="flex items-center gap-2 bg-[#d3f9d8] text-[#37b24d] text-xs p-1 px-4 rounded-lg "
                onClick={(e) => {
                  e.stopPropagation();
                  subscribe();
                }}
              >
                Inscrito
              </div>
            )}
          </div>
          <div className="flex flex-col items-end">
            <div className="text-slate-600 text-xl font-medium">
              {passageiros ? (
                <span className="text-slate-600 text-xl font-medium">
                  R$ {Number(trip.descricao.valorPorPessoa) * passageiros}
                </span>
              ) : (
                <span className="text-slate-600 text-xl font-medium">
                  R$ {trip.descricao.valorPorPessoa}
                </span>
              )}
            </div>
            {!passageiros && (
              <div className="text-slate-500 text-xs">por pessoa</div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
