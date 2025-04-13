import {
  GiAfrica,
  GiAztecCalendarSun,
  GiBrazilFlag,
  GiGreekTemple,
  GiJapaneseBridge,
  GiKangaroo,
  GiMayanPyramid,
  GiPalmTree,
} from "react-icons/gi";
import { CarouselItem } from "./CarouselItem";
import { TbBuildingBurjAlArab } from "react-icons/tb";
import { IoSnowOutline } from "react-icons/io5";
import { destinosPorCategoria, puxarDestinos } from "../api/destinos";
import { useEffect, useState } from "react";
import { FilteredTrips, Trip } from "./FilteredTrips";
import { puxarPromocoes } from "../api/marketing";

export const iconsCategorias = {
  Brasil: <GiBrazilFlag size={30} />,
  "América do Sul": <GiAztecCalendarSun size={30} />,
  "América do Norte": <GiMayanPyramid size={30} />,
  Caribe: <GiPalmTree size={30} />,
  África: <GiAfrica size={30} />,
  "Oriente Médio": <TbBuildingBurjAlArab size={30} />,
  Ásia: <GiJapaneseBridge size={30} />,
  Mediterrâneo: <GiGreekTemple size={30} />,
  Escandinávia: <IoSnowOutline size={30} />,
  Oceania: <GiKangaroo size={30} />,
};

interface Destino {
  quantidade: string;
  categoria: keyof typeof iconsCategorias;
}

export const Carousel = () => {
  const [loading, setLoading] = useState(false);
  const [destinos, setDestinos] = useState([]);
  const [promocoes, setPromocoes] = useState({
    subscriptions: [],
    promotions: [],
  });
  const [destinosDeCategoria, setDestinosDeCategoria] = useState<Trip[]>([]);

  const categoriasDeDestinos = async () => {
    const response = await destinosPorCategoria();
    setDestinos(response);
  };

  const handleClickItem = async (categoria: keyof typeof iconsCategorias) => {
    setLoading(true);
    const response = await puxarDestinos({ categoria });
    setDestinosDeCategoria(response);
    fetchPromotions();

    setLoading(false);
  };

  const fetchPromotions = async () => {
    const p = await puxarPromocoes();
    setPromocoes(p);
  };

  useEffect(() => {
    const fetchCategorias = async () => {
      try {
        await categoriasDeDestinos();
      } catch (error) {
        console.error("Erro ao buscar destinos:", error);
      }
    };

    fetchCategorias();
    fetchPromotions();
  }, []);

  return (
    <div className="flex flex-col w-full h-full gap-4 items-center justify-center">
      <div className="carousel rounded-box w-full gap-2">
        {destinos.map((destino: Destino) => (
          <CarouselItem
            title={destino.categoria}
            subtitle={`${destino.quantidade} destinos`}
            icon={iconsCategorias[destino.categoria]}
            onClick={handleClickItem}
          />
        ))}
      </div>
      {loading && <div className="loading loading-spinner"></div>}
      {destinosDeCategoria.length != 0 && !loading && (
        <div className="flex flex-col w-full h-full p-4 gap-4">
          <div>
            <h2 className="text-slate-700 text-lg font-medium">
              {destinosDeCategoria[0].categoria}
            </h2>
            <div className="text-slate-500 text-sm">
              {destinosDeCategoria.length} destinos encontrados
            </div>
          </div>
          <FilteredTrips
            filter={destinosDeCategoria}
            subscriptions={promocoes.subscriptions}
          />
        </div>
      )}
    </div>
  );
};
