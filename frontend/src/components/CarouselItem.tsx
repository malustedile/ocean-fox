import { iconsCategorias } from "./Carousel";
import { Tag } from "./Tag";

interface CarouselItemProps {
  id: number;
  title: string;
  subtitle: string;
  icon: React.ReactNode;
  onClick: (categoria: keyof typeof iconsCategorias) => void;
  tag?: string;
  tagColor?: string;
}

export const CarouselItem = ({
  id,
  title,
  subtitle,
  icon,
  tag,
  tagColor,
  onClick,
}: CarouselItemProps) => {
  const handleClick = () => {
    onClick(title as keyof typeof iconsCategorias);
  };

  return (
    <div id={`slide${id}`} className="carousel-item" onClick={handleClick}>
      <div className="flex flex-col w-[200px] h-[130px] bg-white rounded-xl p-6 border border-slate-300 gap-3 hover:cursor-pointer hover:bg-slate-50">
        <div className="flex w-full justify-between items-start">
          <div>{icon}</div>
          <Tag tag={tag} tagColor={tagColor} />
        </div>
        <div className="flex flex-col w-full">
          <div className="flex w-full items-start text-lg font-medium">
            {title}
          </div>
          <div className="flex w-full text-sm text-gray-500">{subtitle}</div>
        </div>
      </div>
    </div>
  );
};
