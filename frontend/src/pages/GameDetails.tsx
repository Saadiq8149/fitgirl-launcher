import { models } from "../../wailsjs/go/models";
import { GetRepackDetails } from "../../wailsjs/go/handlers/FitgirlScraperHandler";
import { GetGameFromDatabase } from "../../wailsjs/go/handlers/DatabaseHandler";
import { AddTorrent } from "../../wailsjs/go/handlers/TorrentHandler";
import { useEffect, useState } from "react";

export default function GameDetails(props: {
  url: string;
  setPage: (page: string) => void;
}) {
  const [repackDetails, setRepackDetails] =
    useState<models.FitgirlRepack | null>(null);
  const [status, setStatus] = useState<string>("");

  useEffect(() => {
    async function fetchRepackDetails() {
      try {
        const details = await GetRepackDetails(props.url);
        const gameInDb = await GetGameFromDatabase(props.url);
        console.log(gameInDb);
        if (gameInDb.url) {
          setStatus(gameInDb.status);
        } else {
          setStatus("not installed");
        }
        setRepackDetails(details);
      } catch (error) {
        console.error("Error fetching repack details:", error);
      }
    }

    fetchRepackDetails();
  }, [props.url]);

  function setButtonStatus(status: string) {
    switch (status) {
      case "installed":
        return "Installed";
      case "installing":
        return "Installing...";
      case "downloading":
        return "Downloading...";
      case "not installed":
        return "Download Repack";
      default:
        return "Download Repack";
    }
  }

  function handleAction(status: string) {
    switch (status) {
      case "installed":
        // Open Game
        break;
      case "installing":
        props.setPage("downloads");
        break;
      case "downloading":
        props.setPage("downloads");
        break;
      case "not installed":
        if (repackDetails) {
          AddTorrent(repackDetails.Sources[0], repackDetails, props.url);
          setStatus("downloading");
          props.setPage("downloads");
        }
      default:
        if (repackDetails) {
          AddTorrent(repackDetails.Sources[0], repackDetails, props.url);
          setStatus("downloading");
          props.setPage("downloads");
        }
    }
  }

  function getButtonContent(status: string) {
    switch (status) {
      case "installed":
        return "Installed";
      case "installing":
        return (
          <div className="flex items-center justify-center gap-2">
            <div className="w-4 h-4 border-2 border-white border-t-black rounded-full animate-spin"></div>
            Installing...
          </div>
        );
      case "downloading":
        return (
          <div className="flex items-center justify-center gap-2">
            <div className="w-4 h-4 border-2 border-white border-t-black rounded-full animate-spin"></div>
            Downloading...
          </div>
        );
      case "not installed":
        return "Download Repack";
      default:
        return "Download Repack";
    }
  }

  function getButtonClasses(status: string) {
    switch (status) {
      case "installed":
        return "bg-green-600 text-white px-6 py-3 rounded-md font-semibold cursor-default";
      case "installing":
        return "bg-yellow-500 text-black px-6 py-3 rounded-md font-semibold cursor-not-allowed";
      case "downloading":
        return "bg-blue-500 text-white px-6 py-3 rounded-md font-semibold cursor-not-allowed";
      case "not installed":
        return "bg-white text-black px-6 py-3 rounded-md font-semibold hover:bg-gray-200 transition";
      default:
        return "bg-white text-black px-6 py-3 rounded-md font-semibold hover:bg-gray-200 transition";
    }
  }

  if (!repackDetails) {
    return (
      <div className="flex justify-center items-center h-full bg-black text-white">
        <div className="w-12 h-12 border-4 border-white border-t-black rounded-full animate-spin"></div>
      </div>
    );
  }
  return (
    <div className="min-h-screen flex items-center justify-center bg-black text-white p-4">
      <div className="p-8 bg-black text-white max-w-4xl w-full rounded-md flex flex-col md:flex-row items-center gap-8">
        {/* Image Section */}
        <div className="flex flex-col items-center justify-center md:w-1/3">
          <img
            src={repackDetails.CoverImage}
            alt={repackDetails.Name}
            className="w-full md:w-64 rounded-md object-cover border border-white"
          />
          <div className="mt-4 mb-4 space-y-1 text-center md:text-left">
            <div>
              <span className="font-semibold">Original Size: </span>
              <span>{repackDetails.OriginalSize}</span>
            </div>
            <div>
              <span className="font-semibold">Repack Size: </span>
              <span>{repackDetails.RepackSize}</span>
            </div>
          </div>

          <button
            onClick={() => handleAction(status)}
            className={getButtonClasses(status)}
            disabled={status === "installing" || status === "downloading"}
          >
            {getButtonContent(status)}
          </button>
        </div>

        {/* Details Section */}
        <div className="flex-1 flex flex-col items-center md:items-start">
          <h1 className="text-4xl font-bold mb-6 border-b border-white pb-4 text-center md:text-left">
            {repackDetails.Name}
          </h1>
          {/* You can add more details here if needed */}
        </div>
      </div>
    </div>
  );
}
