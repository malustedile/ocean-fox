import { MdOutlineEventAvailable, MdOutlinePinDrop } from "react-icons/md";
import { Tag } from "./Tag";
import { LuShip } from "react-icons/lu";
import { Trip } from "./FilteredTrips";

const coresCategorias = {
  Brasil: "bg-red-100 text-red-800",
  "América do Sul": "bg-orange-100 text-orange-800",
  "América do Norte": "bg-yellow-100 text-yellow-800",
  Caribe: "bg-green-100 text-green-800",
  África: "bg-blue-100 text-blue-800",
  "Oriente Médio": "bg-purple-100 text-purple-800",
  Ásia: "bg-pink-100 text-pink-800",
  Mediterrâneo: "bg-indigo-100 text-indigo-800",
  Escandinávia: "bg-slate-100 text-slate-800",
  Oceania: "bg-teal-100 text-teal-800",
};

export const TripCard = ({ index, trip }: { index: number; trip: Trip }) => {
  return (
    <div
      key={index}
      className="flex flex-col bg-white rounded-xl p-3 w-[300px] h-[240px] hover:cursor-pointer hover:bg-slate-50 border border-slate-300"
    >
      <div className="flex flex-row justify-between">
        <div className="text-slate-500 text-xs">
          {trip.descricao.noites} noites
        </div>
        <Tag tag={trip.categoria} tagColor={coresCategorias[trip.categoria]} />
      </div>
      <div className="flex flex-col w-full gap-2">
        <div className="text-slate-700 text-md font-medium">{trip.nome}</div>
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
              R$ {trip.descricao.valorPorPessoa}
            </div>
            <div className="text-slate-500 text-xs">por pessoa</div>
          </div>
        </div>
      </div>
    </div>
  );
};
