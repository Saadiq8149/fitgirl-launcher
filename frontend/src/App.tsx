import { useState } from "react";

import Sidebar from "./components/Sidebar";
import Home from "./pages/Home";
import Downloads from "./pages/Downloads";
import Library from "./pages/Library";
import Search from "./pages/Search";
import GameDetails from "./pages/GameDetails";

export default function App() {
  const [page, setPage] = useState<string>("home");
  const [gameUrl, setGameUrl] = useState<string>("");

  return (
    <div className="flex h-screen w-screen overflow-hidden bg-black">
      <Sidebar currentPage={page} setPage={setPage} />
      <div className="flex-1 overflow-hidden">
        {page == "home" && <Home setGameUrl={setGameUrl} setPage={setPage} />}
        {page == "downloads" && <Downloads />}
        {page == "library" && <Library />}
        {page == "search" && (
          <Search setGameUrl={setGameUrl} setPage={setPage} />
        )}
        {page == "gameDetails" && (
          <GameDetails url={gameUrl} setPage={setPage} />
        )}
      </div>
    </div>
  );
}
