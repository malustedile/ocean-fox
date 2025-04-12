import { useState } from "react";
import { Trip } from "./FilteredTrips";
import { reservarViagem } from "../api/destinos";
import { CruiseDetails } from "./CruiseDetails";

interface BookingModalProps {
  trip: Trip;
  openBookingModal: boolean;
  setOpenBookingModal: React.Dispatch<React.SetStateAction<boolean>>;
}

export const BookingModal = ({
  trip,
  openBookingModal,
  setOpenBookingModal,
}: BookingModalProps) => {
  const [dataEmbarque, setDataEmbarque] = useState<string>("");
  const [numeroPassageiros, setNumeroPassageiros] = useState<number>(1);
  const [numeroCabines, setNumeroCabines] = useState<number>(1);

  const bookTrip = async () => {
    const [, err] = await reservarViagem({
      destino: trip.nome,
      dataEmbarque,
      numeroPassageiros,
      numeroCabines,
      valorTotal: Number(trip.descricao.valorPorPessoa) * numeroPassageiros,
    });
    if (err) return;
  };

  return (
    <dialog
      open={openBookingModal}
      className="modal"
      onClick={(e) => {
        if (e.target instanceof HTMLDialogElement) {
          setOpenBookingModal(false);
        }
      }}
    >
      <div
        className="bg-white rounded-lg p-6 flex flex-col gap-4"
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-slate-700 text-xl font-medium">{trip.nome}</h2>
        <CruiseDetails trip={trip} passageiros={numeroPassageiros} />
        <hr className="border-slate-300" />
        <div className="flex flex-col w-full">
          <h2 className="text-slate-700 text-lg font-medium">
            Informações da Reserva
          </h2>
        </div>
        <div className="flex flex-row items-start gap-4">
          <div className="flex flex-col w-[200px]">
            <label className="text-slate-500 text-xs">Data de Embarque</label>
            <input
              type="date"
              value={dataEmbarque}
              onChange={(e) => setDataEmbarque(e.target.value)}
              min={trip.descricao.datasDisponiveis[0]}
              max={trip.descricao.datasDisponiveis[1]}
              className="input-sm rounded border border-gray-300 p-2 focus:outline-none"
            />
          </div>
          <div className="flex flex-col w-[200px]">
            <label className="text-slate-500 text-xs">
              Número de Passageiros
            </label>
            <input
              type="number"
              value={numeroPassageiros}
              onChange={(e) => setNumeroPassageiros(Number(e.target.value))}
              min={1}
              max={10}
              className="input-sm rounded border border-gray-300 p-2 focus:outline-none"
            />
          </div>
          <div className="flex flex-col w-[200px]">
            <label className="text-slate-500 text-xs">Número de Cabines</label>
            <input
              type="number"
              value={numeroCabines}
              onChange={(e) => setNumeroCabines(Number(e.target.value))}
              min={1}
              max={4}
              className="input-sm rounded border border-gray-300 p-2 focus:outline-none"
            />
          </div>
        </div>
        <div className="flex justify-end gap-2">
          <button
            className="btn px-4 py-2 rounded-md"
            onClick={() => setOpenBookingModal(false)}
          >
            Cancelar
          </button>
          <button
            className="btn btn px-4 py-2 rounded-md"
            onClick={() => {
              bookTrip();
              setOpenBookingModal(false);
            }}
          >
            Confirmar
          </button>
        </div>
      </div>
    </dialog>
  );
};
