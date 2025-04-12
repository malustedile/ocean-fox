import "./App.css";
import { FilterBar } from "./components/FilterBar";
import { NavBar } from "./components/NavBar";
import { Carousel } from "./components/Carousel";
import { useEffect, useState } from "react";
import { FilteredTrips } from "./components/FilteredTrips";
import axios from "axios";
import { MinhasReservas } from "./screens/MinhasReservas";
import { Promocoes } from "./screens/Promocoes";

export enum Screens {
  Itinerarios = "ITINERARIOS",
  MinhasReservas = "MINHAS_RESERVAS",
  InscrevaSe = "INSCREVA_SE",
}

function App() {
  const [filter, setFilter] = useState();
  const [activeScreen, setActiveScreen] = useState<Screens>(
    Screens.Itinerarios
  );
  useEffect(() => {
    axios.get("http://localhost:3005/session", { withCredentials: true });
  }, []);

  const renderPanel = () => {
    switch (activeScreen) {
      case Screens.Itinerarios:
        return (
          <div className="flex flex-col overflow-y-scroll">
            {filter ? <FilteredTrips filter={filter} /> : <Carousel />}
          </div>
        );
      case Screens.MinhasReservas:
        return <MinhasReservas />;
      case Screens.InscrevaSe:
        return <Promocoes />;
      default:
        return (
          <div className="flex flex-col overflow-y-scroll">
            {filter ? <FilteredTrips filter={filter} /> : <Carousel />}
          </div>
        );
    }
  };

  return (
    <div className="flex w-screen justify-center text-slate-800 min-h-screen bg-gray-100 p-4">
      <div className="flex flex-col w-[1000px] h-full gap-4">
        <div className="flex flex-col w-full rounded-2xl shadow-md p-6 py-4 static relative">
          <div className="flex items-center justify-center gap-2 text-2xl font-bold mr-auto relative top-8 left-8">
            <div className="w-8 h-8">
              <img src="src/assets/logo.png" />
            </div>
            <div>Ocean Fox</div>
          </div>
          <NavBar activeScreen={activeScreen} setScreen={setActiveScreen} />
          {activeScreen === Screens.Itinerarios && (
            <div className="flex flex-col items-center text-center">
              <h1 className="text-5xl md:text-6xl font-black text-[#002d3a]">
                Explore o mundo
              </h1>
              <p className="mt-4 text-gray-500">Encontre sua próxima viagem</p>
              <FilterBar setFilter={setFilter} />
            </div>
          )}
        </div>
        {renderPanel()}
      </div>
    </div>
  );
}

export default App;
