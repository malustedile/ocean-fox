import React, { useState } from "react";
import { IoSearch } from "react-icons/io5";
import { puxarDestinos } from "../api/destinos";

interface FilterBarProps {
  setFilter: React.Dispatch<React.SetStateAction<any>>;
}

export const FilterBar = ({ setFilter }: FilterBarProps) => {
  const [destino, setDestino] = useState("");
  const [mes, setMes] = useState("");
  const [embarque, setEmbarque] = useState("");

  const handleFilter = async () => {
    const response = await puxarDestinos(destino, mes, embarque);
    setFilter(response);
  };

  return (
    <div className="mt-10 w-full max-w-5xl bg-white rounded-2xl shadow-lg p-6 flex flex-wrap gap-4 justify-between items-center">
      <div className="flex flex-col w-[200px]">
        <label className="text-[#007090] font-semibold">Destino</label>
        <input
          type="text"
          placeholder="Para onde você quer ir?"
          value={destino}
          onChange={(e) => setDestino(e.target.value)}
          className="border-b border-gray-300 py-2 focus:outline-none"
        />
      </div>
      <div className="flex flex-col w-[200px]">
        <label className="text-[#007090] font-semibold">
          Porto de Embarque
        </label>
        <input
          type="text"
          placeholder="De onde você quer sair?"
          value={embarque}
          onChange={(e) => setEmbarque(e.target.value)}
          className="border-b border-gray-300 py-2 focus:outline-none"
        />
      </div>
      <div className="flex flex-col w-[150px]">
        <label className="text-[#007090] font-semibold">Mês de Embarque</label>
        <input
          type="number"
          className="border-b border-gray-300 py-2 focus:outline-none"
          value={mes}
          min={1}
          max={12}
          placeholder="Mês (1-12)"
          onChange={(e) => setMes(e.target.value)}
        />
      </div>
      <button onClick={handleFilter}>
        <IoSearch />
      </button>
    </div>
  );
};
