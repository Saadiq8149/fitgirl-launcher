import { useEffect, useState } from "react";
import { models } from "../../wailsjs/go/models";
import {
  GetTorrents,
  GetTorrent,
  RemoveTorrent,
} from "../../wailsjs/go/handlers/TorrentHandler";
import { GetGameFromDatabaseByHash } from "../../wailsjs/go/handlers/DatabaseHandler";

export default function Downloads() {
  const [downloads, setDownloads] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    async function fetchDownloads() {
      try {
        const torrents = await GetTorrents();
        let detailedDownloads = await Promise.all(
          torrents.map(async (torrent: any) => {
            const game = await GetGameFromDatabaseByHash(torrent.hash);
            return {
              ...torrent,
              Game: game,
            };
          }),
        );
        detailedDownloads = detailedDownloads.filter(
          (download) => download.Game.status != "installed",
        );
        console.log(detailedDownloads);
        setDownloads(detailedDownloads);
      } catch (err) {
        console.error("Error fetching downloads:", err);
      } finally {
        setLoading(false);
      }
    }

    fetchDownloads();
    const interval = setInterval(fetchDownloads, 1000); // Auto-refresh every 5s
    return () => clearInterval(interval);
  }, []);

  async function handleRemove(magnet: string) {
    try {
      await RemoveTorrent(magnet);
      setDownloads((prev) => prev.filter((d) => d.Magnet !== magnet));
    } catch (err) {
      console.error("Failed to remove torrent:", err);
    }
  }

  function getStatusColor(status: string) {
    console.log(status);

    switch (status.toLowerCase()) {
      case "downloading":
        return "text-blue-400";
      case "installing":
        return "text-yellow-400";
      case "completed":
      case "installed":
        return "text-green-400";
      default:
        return "text-gray-400";
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen bg-black text-white">
        <div className="w-12 h-12 border-4 border-white border-t-transparent rounded-full animate-spin"></div>
      </div>
    );
  }

  if (downloads.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-screen bg-black text-white">
        <h1 className="text-3xl font-bold mb-4">Downloads</h1>
        <p className="text-gray-400">No active downloads at the moment.</p>
      </div>
    );
  }

  return (
    <div className="flex-1 p-10 bg-black text-white min-h-screen">
      <h1 className="text-3xl font-bold mb-6 border-b border-gray-700 pb-4">
        Downloads
      </h1>
      <div className="space-y-6">
        {downloads.map((download) => (
          <div
            key={download.Game.magnet}
            className="flex flex-col md:flex-row items-center justify-between border border-gray-700 rounded-md p-4 bg-gray-900 hover:bg-gray-800 transition"
          >
            {/* Left: Image */}
            <div className="flex items-center gap-4">
              <img
                src={download.Game?.thumbnail || "/placeholder.png"}
                alt={download.Game?.title || "Game"}
                className="w-24 h-32 object-cover rounded border border-gray-600"
              />
              <div>
                <h2 className="text-2xl font-semibold">
                  {download.Game?.title || "Unknown Game"}
                </h2>
                <p
                  className={`mt-2 font-medium ${getStatusColor(download.Game.status)}`}
                >
                  {download.Game.status === "downloading" && download.progress
                    ? `Downloading (${(download.progress * 100).toFixed(1)}%)`
                    : download.Game.status.charAt(0).toUpperCase() +
                      download.Game.status.slice(1)}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-4 mt-4 md:mt-0">
              {download.Game.status === "downloading" && (
                <div className="w-40 h-2 bg-gray-700 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-blue-500"
                    style={{ width: `${(download.progress || 0) * 100}%` }}
                  ></div>
                </div>
              )}
              {download.Game.status === "installing" && (
                <div className="flex items-center gap-2 text-yellow-400">
                  <div className="w-4 h-4 border-2 border-yellow-400 border-t-transparent rounded-full animate-spin"></div>
                  Installing...
                </div>
              )}
              <button
                onClick={() => handleRemove(download.Game.magnet)}
                className="px-4 py-2 bg-gray-800 border border-gray-600 text-white rounded-md hover:bg-red-600 hover:border-red-600 transition"
              >
                Remove
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
