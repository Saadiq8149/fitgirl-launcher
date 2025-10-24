import { useEffect, useState } from "react";
import { models } from "../../wailsjs/go/models";
import { LoadDatabase } from "../../wailsjs/go/handlers/DatabaseHandler";
import { LaunchGame } from "../../wailsjs/go/handlers/GameHandler";

export default function Library() {
  const [installedGames, setInstalledGames] = useState<models.Game[]>([]);

  useEffect(() => {
    async function fetchInstalledGames() {
      try {
        const allGames = await LoadDatabase();
        // ✅ Only include installed games
        const filtered = allGames.games.filter(
          (game: models.Game) => game.status === "installed",
        );
        setInstalledGames(filtered);
      } catch (error) {
        console.error("Failed to load games:", error);
      }
    }
    fetchInstalledGames();
  }, []);

  if (installedGames.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen bg-black text-white">
        <p className="text-gray-400 text-lg">No installed games found.</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-black text-white p-10">
      <h1 className="text-4xl font-bold mb-8 border-b border-white pb-4">
        Installed Games
      </h1>

      {/* Grid Layout */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
        {installedGames.map((game) => (
          <div
            key={game.url}
            className="bg-white/5 border border-white/10 rounded-lg overflow-hidden shadow-lg flex flex-col hover:scale-[1.02] hover:border-white/20 transition-transform duration-200"
          >
            {/* Game Image */}
            <div className="relative w-full h-56 overflow-hidden">
              <img
                src={game.thumbnail || "placeholder.jpg"}
                alt={game.title}
                className="object-cover w-full h-full"
              />
            </div>

            {/* Game Info */}
            <div className="flex flex-col flex-1 p-4">
              <h2 className="text-xl font-semibold mb-4 truncate">
                {game.title}
              </h2>
              <div className="mt-auto flex justify-center">
                <button
                  className="w-full bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-md font-semibold transition"
                  onClick={() => launchGame(game)}
                >
                  Play
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// Example launch function
async function launchGame(game: models.Game) {
  try {
    await LaunchGame(game);
  } catch (error) {
    alert("Game already running");
  }
}
