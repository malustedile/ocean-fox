import { useState } from "react";
import { iconsCategorias } from "./Carousel";
import { TripCard } from "./TripCard";
import { BookingModal } from "./BookingModal";

export interface Trip {
  nome: string;
  categoria: keyof typeof iconsCategorias;
  descricao: {
    embarque: string;
    desembarque: string;
    datasDisponiveis: string[];
    lugaresVisitados: string[];
    navio: string;
    noites: string;
    valorPorPessoa: string;
  };
}

export const FilteredTrips = ({ filter }: { filter: Trip[] }) => {
  const [openBookingModal, setOpenBookingModal] = useState(false);
  const [trip, setTrip] = useState<Trip>();

  return (
    <div className="flex flex-wrap h-full w-full gap-4 items-center justify-start overflow-x-hidden">
      {filter.map((trip: Trip, index: number) => (
        <TripCard
          index={index}
          trip={trip}
          setTrip={setTrip}
          setOpenBookingModal={setOpenBookingModal}
        />
      ))}
      {openBookingModal && (
        <BookingModal
          trip={trip as Trip}
          openBookingModal={openBookingModal}
          setOpenBookingModal={setOpenBookingModal}
        />
      )}
    </div>
  );
};
