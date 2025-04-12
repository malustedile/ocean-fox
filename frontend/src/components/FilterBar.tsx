import { IoSearch } from "react-icons/io5";

export const FilterBar = () => {
  return (
    <div className="mt-10 w-full max-w-5xl bg-white rounded-2xl shadow-lg p-6 flex flex-wrap gap-4 justify-between items-center">
      <div className="flex flex-col w-[200px]">
        <label className="text-[#007090] font-semibold">Destino</label>
        <input
          type="text"
          placeholder="Enter your destination"
          className="border-b border-gray-300 py-2 focus:outline-none"
        />
      </div>
      <div className="flex flex-col w-[200px]">
        <label className="text-[#007090] font-semibold">
          Porto de Embarque
        </label>
        <input
          type="text"
          placeholder="Enter your destination"
          className="border-b border-gray-300 py-2 focus:outline-none"
        />
      </div>
      <div className="flex flex-col w-[150px]">
        <label className="text-[#007090] font-semibold">Mês de Embarque</label>
        <input
          type="date"
          className="border-b border-gray-300 py-2 focus:outline-none"
        />
      </div>
      <button className="bg-black text-white p-4 rounded-full mt-6">
        <IoSearch />
      </button>
    </div>
  );
};
