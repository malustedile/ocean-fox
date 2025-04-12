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
import { destinosPorCategoria } from "../api/destinos";
import { useEffect, useState } from "react";

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
  const [destinos, setDestinos] = useState([]);

  const puxandoDestinos = async () => {
    const response = await destinosPorCategoria();
    setDestinos(response);
  };

  useEffect(() => {
    const fetchDestinos = async () => {
      try {
        await puxandoDestinos();
      } catch (error) {
        console.error("Erro ao buscar destinos:", error);
      }
    };

    fetchDestinos();
  }, []);

  return (
    <div className="carousel rounded-box w-full gap-2">
      {destinos.map((destino: Destino) => (
        <CarouselItem
          title={destino.categoria}
          subtitle={`${destino.quantidade} destinos`}
          icon={iconsCategorias[destino.categoria]}
        />
      ))}
    </div>
  );
};
