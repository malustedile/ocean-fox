import React from "react";
import { minhasReservas } from "../api/destinos";
import { Screens } from "../App";

interface NavBarProps {
  activeScreen: Screens;
  setScreen: (screen: Screens) => any;
}

export const NavBar: React.FC<NavBarProps> = ({ activeScreen, setScreen }) => {
  const handleMyReservationsClick = async () => {
    const res = await minhasReservas();
    console.log(res);
  };

  const classNameButtons =
    "hover:text-[#007090] hover:bg-[#cce2e9] hover:cursor-pointer px-4 py-1 rounded-full";

  return (
    <div className="flex items-center justify-center w-full mb-6 gap-2">
      <div
        onClick={() => setScreen(Screens.Itinerarios)}
        className={
          activeScreen === Screens.Itinerarios
            ? "bg-[#cce2e9] text-[#007090] hover:cursor-pointer px-4 py-1 rounded-full"
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
            ? "bg-[#cce2e9] text-[#007090] hover:cursor-pointer px-4 py-1 rounded-full"
            : classNameButtons
        }
      >
        Minhas Reservas
      </div>
      <div>•</div>
      <div className={classNameButtons}>Inscreva-se</div>
    </div>
  );
};
