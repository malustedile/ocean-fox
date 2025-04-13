import { Carousel } from "../components/Carousel";
import { FilteredTrips } from "../components/FilteredTrips";

interface ItinerariosProps {
  filter: any;
}

export const Itinerarios = ({ filter }: ItinerariosProps) => {
  return (
    <div className="flex flex-col">
      {filter ? (
        <FilteredTrips filter={filter} subscriptions={[]} />
      ) : (
        <Carousel />
      )}
    </div>
  );
};
