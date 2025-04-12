export const Tag = ({
  tag,
  tagColor = "bg-blue-100 text-blue-800",
}: {
  tag?: string;
  tagColor?: string;
}) => {
  return (
    <>
      {tag && (
        <div className={`px-2 py-1 rounded-full text-xs ${tagColor}`}>
          {tag}
        </div>
      )}
    </>
  );
};
