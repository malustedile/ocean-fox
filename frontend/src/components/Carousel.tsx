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

export const Carousel = () => {
  return (
    <div className="carousel rounded-box w-full gap-2">
      <CarouselItem
        title="Brasil"
        subtitle="3 destinos"
        icon={<GiBrazilFlag size={30} />}
        tag="Nacional"
        tagColor="text-pink-500 bg-pink-100"
      />
      <CarouselItem
        title="América do Sul"
        subtitle="3 destinos"
        icon={<GiAztecCalendarSun size={30} />}
      />
      <CarouselItem
        title="América do Norte"
        subtitle="3 destinos"
        icon={<GiMayanPyramid size={30} />}
      />
      <CarouselItem
        title="Caribe"
        subtitle="3 destinos"
        icon={<GiPalmTree size={30} />}
        tag="Popular"
        tagColor="text-blue-500 bg-blue-100"
      />
      <CarouselItem
        title="África"
        subtitle="3 destinos"
        icon={<GiAfrica size={30} />}
      />
      <CarouselItem
        title="Oriente Médio"
        subtitle="3 destinos"
        icon={<TbBuildingBurjAlArab size={30} />}
        tag="Novo"
        tagColor="text-green-500 bg-green-100"
      />
      <CarouselItem
        title="Ásia"
        subtitle="3 destinos"
        icon={<GiJapaneseBridge size={30} />}
      />
      <CarouselItem
        title="Mediterrâneo"
        subtitle="3 destinos"
        icon={<GiGreekTemple size={30} />}
        tag="Popular"
        tagColor="text-blue-500 bg-blue-100"
      />
      <CarouselItem
        title="Escandinávia"
        subtitle="3 destinos"
        icon={<IoSnowOutline size={30} />}
        tag="Verão"
        tagColor="text-yellow-500 bg-yellow-100"
      />
      <CarouselItem
        title="Oceania"
        subtitle="3 destinos"
        icon={<GiKangaroo size={30} />}
      />
    </div>
  );
};
