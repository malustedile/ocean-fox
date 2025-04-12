import "./App.css";
import { FilterBar } from "./components/FilterBar";
import { NavBar } from "./components/NavBar";
import { Carousel } from "./components/Carousel";
import { useState } from "react";
import { FilteredTrips } from "./components/FilteredTrips";

function App() {
  const [filter, setFilter] = useState();

  console.log("Filter:", filter);

  return (
    <div className="flex w-screen min-h-screen bg-gray-100 p-4">
      <div className="flex flex-col w-full h-full gap-4">
        <div className="flex flex-col w-full rounded-2xl shadow-md p-6 static">
          <div className="flex items-center gap-2 text-2xl font-bold mr-auto absolute top-8 left-8">
            <div className="w-8 h-8">
              <img src="src/assets/logo.png" />
            </div>
            ocean fox
          </div>
          <NavBar />
          <div className="flex flex-col items-center text-center">
            <h1 className="text-5xl md:text-6xl font-black text-[#002d3a]">
              Explore o mundo
            </h1>
            <p className="mt-4 text-gray-500">Encontre sua próxima viagem</p>
            <FilterBar setFilter={setFilter} />
          </div>
        </div>
        <div className="flex flex-col overflow-y-scroll">
          {filter ? <FilteredTrips filter={filter} /> : <Carousel />}
        </div>
      </div>
    </div>
  );
}

export default App;
