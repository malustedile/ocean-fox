import "./App.css";
import { FilterBar } from "./components/FilterBar";
import { NavBar } from "./components/NavBar";
import { useEffect, useState } from "react";
import axios from "axios";
import { MinhasReservas } from "./screens/MinhasReservas";
import { Promocoes, Subscriptions } from "./screens/Promocoes";
import { Itinerarios } from "./screens/Itinerarios";
import { puxarPromocoes } from "./api/marketing";
import { minhasReservas } from "./api/reserva";
export enum Screens {
  Itinerarios = "ITINERARIOS",
  MinhasReservas = "MINHAS_RESERVAS",
  InscrevaSe = "INSCREVA_SE",
}

function App() {
  const [filter, setFilter] = useState();
  const [reservas, setReservas] = useState<any[]>([]);

  const [promocoes, setPromocoes] = useState<Subscriptions>({
    hasSubscription: false,
    promotions: [],
  });
  const fetchPromotions = async () => {
    const res = await puxarPromocoes();
    setPromocoes(res);
  };
  const [activeScreen, setActiveScreen] = useState<Screens>(
    Screens.Itinerarios
  );

  const fetchReservas = async () => {
    const res = await minhasReservas();
    setReservas(res);
  };

  useEffect(() => {
    fetchReservas();
  }, []);
  useEffect(() => {
    axios.get("http://localhost:3005/session", { withCredentials: true });
  }, []);

  useEffect(() => {
    fetchPromotions();
    const eventSource = new EventSource("http://localhost:3000/sse", {
      withCredentials: true,
    });
    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log({ data });
      if (data.eventType === "promocao") {
        setPromocoes((prev) => ({
          ...prev,
          promotions: [{ mensagem: data.msg }, ...prev.promotions],
        }));
      } else if (data.eventType === "UPDATE_PAYMENT_STATUS") {
        setReservas((prev) =>
          prev.map((res) => {
            if (res.id !== data.data.id) return res;
            if (data.data.canceled) {
              res.status = "cancelado";
              return res;
            }
            res.statusPagamento = data.data.status;
            console.log({ res });
            return res;
          })
        );
      }
    };
    return () => {
      eventSource.close();
    };
  }, []);

  const renderPanel = () => {
    switch (activeScreen) {
      case Screens.Itinerarios:
        return <Itinerarios filter={filter} />;
      case Screens.MinhasReservas:
        return <MinhasReservas reservas={reservas} />;
      case Screens.InscrevaSe:
        return <Promocoes promocoes={promocoes} />;
      default:
        return <Itinerarios filter={filter} />;
    }
  };

  return (
    <div className="flex w-screen justify-center text-slate-800 min-h-screen bg-gray-100 p-4">
      <div className="flex flex-col w-[1000px] h-full gap-4">
        <div
          className="flex flex-col w-full rounded-2xl shadow-md p-6 py-4 static relative bg-no-repeat bg-center bg-cover"
          style={{ backgroundImage: "url(../public/hero.png)" }}
        >
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
              <p className="mt-4 font-bold">Encontre sua próxima viagem</p>
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
