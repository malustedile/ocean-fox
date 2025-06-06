import React from "react";
import { Screens } from "../App";

interface NavBarProps {
  activeScreen: Screens;
  setScreen: (screen: Screens) => any;
}

export const NavBar: React.FC<NavBarProps> = ({ activeScreen, setScreen }) => {
  const classNameButtons =
    "hover:text-[#007090] hover:bg-[#cce2e9] font-bold hover:cursor-pointer px-4 py-1 rounded-full";

  return (
    <div className="flex items-center justify-center w-full mb-6 gap-2">
      <div
        onClick={() => setScreen(Screens.Itinerarios)}
        className={
          activeScreen === Screens.Itinerarios
            ? "bg-[#cce2e9] text-[#007090] font-bold hover:cursor-pointer px-4 py-1 rounded-full"
            : classNameButtons
        }
      >
        Itinerários
      </div>
      <div>•</div>
      <div
        onClick={() => setScreen(Screens.MinhasReservas)}
        className={
          activeScreen === Screens.MinhasReservas
            ? "bg-[#cce2e9] text-[#007090] font-bold hover:cursor-pointer px-4 py-1 rounded-full"
            : classNameButtons
        }
      >
        Minhas Reservas
      </div>
      <div>•</div>
      <div
        className={classNameButtons}
        onClick={() => setScreen(Screens.InscrevaSe)}
      >
        Promoções
      </div>
    </div>
  );
};
