import { MdOutlineEventAvailable, MdOutlinePinDrop } from "react-icons/md";
import { Trip } from "./FilteredTrips";
import { LuShip } from "react-icons/lu";

export const CruiseDetails = ({
  trip,
  passageiros,
}: {
  trip: Trip;
  passageiros?: number;
}) => {
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
          <div className="text-slate-600 text-sm">
            {trip.descricao.lugaresVisitados.join(", ").toLocaleString()}
          </div>
        </div>
        <div className="flex flex-col w-full items-end">
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
  );
};
