import { useEffect, useState } from "react";

import {
  GetLatestRepacks,
  GetPopularRepacks,
} from "../../wailsjs/go/handlers/FitgirlScraperHandler";

import { models } from "../../wailsjs/go/models";
import GameCard from "../components/GameCard";
export default function Home(props: {
  setGameUrl: (url: string) => void;
  setPage: (page: string) => void;
}) {
  const [latestRepacks, setLatestRepacks] = useState<models.PopularRepacks>();
  const [popularRepacks, setPopularRepacks] = useState<models.PopularRepacks>();

  useEffect(() => {
    async function fetchData() {
      const latest = await GetLatestRepacks();
      setLatestRepacks(latest);

      const popular = await GetPopularRepacks();
      setPopularRepacks(popular);
    }
    fetchData();
  }, []);

  return (
    <div className="flex flex-col h-full overflow-y-auto p-10 scrollbar-custom">
      <div>
        <h1 className="text-white text-2xl font-bold">Latest Repacks</h1>
        <div className="cursor-pointer flex space-x-4 w-full overflow-x-auto py-4 scrollbar-custom mb-10">
          {(!latestRepacks || latestRepacks.Repacks.length === 0) && (
            <div className="flex justify-center items-center py-10">
              <div className="w-8 h-8 border-4 border-gray-300 border-t-indigo-500 rounded-full animate-spin"></div>
            </div>
          )}

          {latestRepacks?.Repacks.map((repack, index) => (
            <GameCard
              coverImage={repack.CoverImage}
              link={repack.Link}
              setGameUrl={props.setGameUrl}
              setPage={props.setPage}
            />
          ))}
        </div>
      </div>

      <div>
        <h1 className="text-white text-2xl font-bold">Popular Repacks</h1>
        <div className="cursor-pointer flex space-x-4 w-full overflow-x-auto py-4 scrollbar-custom mb-10">
          {(!popularRepacks || popularRepacks.Repacks.length === 0) && (
            <div className="flex justify-center items-center py-10">
              <div className="w-8 h-8 border-4 border-gray-300 border-t-indigo-500 rounded-full animate-spin"></div>
            </div>
          )}

          {popularRepacks?.Repacks.map((repack, index) => (
            <GameCard
              coverImage={repack.CoverImage}
              link={repack.Link}
              setGameUrl={props.setGameUrl}
              setPage={props.setPage}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
