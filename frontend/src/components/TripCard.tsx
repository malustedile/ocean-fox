import { Tag } from "./Tag";
import { Trip } from "./FilteredTrips";
import React from "react";
import { CruiseDetails } from "./CruiseDetails";

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

interface TripCardProps {
  index: number;
  trip: Trip;
  setTrip: React.Dispatch<React.SetStateAction<Trip | undefined>>;
  setOpenBookingModal: React.Dispatch<React.SetStateAction<boolean>>;
  subscripted: boolean;
}

export const TripCard = ({
  index,
  trip,
  subscripted,
  setTrip,
  setOpenBookingModal,
}: TripCardProps) => {
  const openModal = () => {
    setTrip(trip);
    setOpenBookingModal(true);
  };

  return (
    <div
      key={index}
      className="flex flex-col bg-white rounded-xl p-3 w-[300px] h-[240px] hover:cursor-pointer hover:bg-slate-50 border border-slate-300"
      onClick={openModal}
    >
      <div className="flex flex-row justify-between">
        <div className="text-slate-500 text-xs">
          {trip.descricao.noites} noites
        </div>
        <Tag tag={trip.categoria} tagColor={coresCategorias[trip.categoria]} />
      </div>
      <div className="flex flex-col gap-2">
        <div className="text-slate-700 text-md font-medium">{trip.nome}</div>
        <CruiseDetails trip={trip} showSubscription subscripted={subscripted} />
      </div>
    </div>
  );
};
