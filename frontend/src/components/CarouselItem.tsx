interface CarouselItemProps {
  title: string;
  subtitle: string;
  icon: React.ReactNode;
  tag?: string;
  tagColor?: string;
}

export const CarouselItem = ({
  title,
  subtitle,
  icon,
  tag,
  tagColor,
}: CarouselItemProps) => {
  return (
    <div className="carousel-item">
      <div className="flex flex-col w-[200px] h-[130px] bg-white rounded-xl p-6 border border-slate-300 gap-3">
        <div className="flex w-full justify-between items-start">
          <div>{icon}</div>
          {tag && (
            <div className={`px-2 py-1 rounded-full text-xs ${tagColor}`}>
              {tag}
            </div>
          )}
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
