import { IoSearch } from "react-icons/io5";
import "./App.css";

function App() {
  return (
    <div className="flex w-screen min-h-screen bg-gray-100 p-4">
      <div className="flex flex-col w-full h-full gap-4">
        <div
          className="flex flex-col w-full rounded-2xl shadow-md p-6 static"
          style={{ backgroundImage: "url('/cruise.png')" }}
        >
          <div className="flex items-center gap-2 text-2xl font-bold mr-auto absolute top-8 left-8">
            <div className="w-8 h-8">
              <img src="src/assets/logo.png" />
            </div>
            ocean fox
          </div>
          <div className="flex items-center justify-center w-full mb-6 gap-2">
            <div className="hover:text-[#007090] hover:bg-[#cce2e9] hover:cursor-pointer px-4 py-1 rounded-full">
              Itinerários
            </div>
            <div>•</div>
            <div className="hover:text-[#007090] hover:bg-[#cce2e9] hover:cursor-pointer px-4 py-1 rounded-full">
              Reservas
            </div>
            <div>•</div>
            <div className="hover:text-[#007090] hover:bg-[#cce2e9] hover:cursor-pointer px-4 py-1 rounded-full">
              Inscreva-se
            </div>
          </div>

          <div className="flex flex-col items-center text-center">
            <h1 className="text-5xl md:text-6xl font-black text-[#002d3a]">
              Explore o mundo
            </h1>
            <p className="mt-4 text-gray-500">Encontre sua próxima viagem</p>

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
                <label className="text-[#007090] font-semibold">
                  Mês de Embarque
                </label>
                <input
                  type="date"
                  className="border-b border-gray-300 py-2 focus:outline-none"
                />
              </div>
              <button className="bg-black text-white p-4 rounded-full mt-6">
                <IoSearch />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
