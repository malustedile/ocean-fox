import { minhasReservas } from "../api/destinos";

export const NavBar = () => {
  const handleMyReservationsClick = async () => {
    const res = await minhasReservas();
    console.log(res);
  };

  return (
    <div className="flex items-center justify-center w-full mb-6 gap-2">
      <div className="hover:text-[#007090] hover:bg-[#cce2e9] hover:cursor-pointer px-4 py-1 rounded-full">
        Itinerários
      </div>
      <div>•</div>
      <div
        onClick={handleMyReservationsClick}
        className="hover:text-[#007090] hover:bg-[#cce2e9] hover:cursor-pointer px-4 py-1 rounded-full"
      >
        Minhas Reservas
      </div>
      <div>•</div>
      <div className="hover:text-[#007090] hover:bg-[#cce2e9] hover:cursor-pointer px-4 py-1 rounded-full">
        Inscreva-se
      </div>
    </div>
  );
};
