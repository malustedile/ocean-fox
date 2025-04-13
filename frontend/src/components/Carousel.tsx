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
  Brasil: <GiBrazilFlag size={30} className="text-[#007090]" />,
  "América do Sul": <GiAztecCalendarSun size={30} className="text-[#007090]" />,
  "América do Norte": <GiMayanPyramid size={30} className="text-[#007090]" />,
  Caribe: <GiPalmTree size={30} className="text-[#007090]" />,
  África: <GiAfrica size={30} className="text-[#007090]" />,
  "Oriente Médio": (
    <TbBuildingBurjAlArab size={30} className="text-[#007090]" />
  ),
  Ásia: <GiJapaneseBridge size={30} className="text-[#007090]" />,
  Mediterrâneo: <GiGreekTemple size={30} className="text-[#007090]" />,
  Escandinávia: <IoSnowOutline size={30} className="text-[#007090]" />,
  Oceania: <GiKangaroo size={30} className="text-[#007090]" />,
};

interface Destino {
  quantidade: string;
  categoria: keyof typeof iconsCategorias;
}

export const Carousel = () => {
  const [loading, setLoading] = useState(false);
  const [destinos, setDestinos] = useState([]);
  const [slide, setSlide] = useState(1);

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
  console.log(slide);
  return (
    <div className="flex flex-col w-full h-full gap-4 items-center justify-center ">
      <div className="relative w-full">
        <div className="carousel rounded-box w-full gap-2 ">
          {destinos.map((destino: Destino, index) => (
            <CarouselItem
              id={index}
              title={destino.categoria}
              subtitle={`${destino.quantidade} destinos`}
              icon={iconsCategorias[destino.categoria]}
              onClick={handleClickItem}
            />
          ))}
          <div className="absolute left-[-24px] right-[-24px] top-1/2 flex -translate-y-1/2 transform justify-between">
            <a
              href={`#slide${slide - 2}`}
              className="btn btn-circle bg-[#007090] border-none"
              onClick={() => setSlide((prev) => (prev - 2 > 2 ? prev - 2 : 2))}
            >
              ❮
            </a>
            <a
              href={`#slide${slide + 2}`}
              className="btn btn-circle bg-[#007090] border-none"
              onClick={() => setSlide((prev) => (prev + 2 < 7 ? prev + 2 : 7))}
            >
              ❯
            </a>
          </div>
        </div>
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
