import { useEffect, useState, useRef, useCallback } from "react";
import { GetRepacks } from "../../wailsjs/go/handlers/FitgirlScraperHandler";
import { models } from "../../wailsjs/go/models";

export default function Search(props: {
  setGameUrl: (url: string) => void;
  setPage: (page: string) => void;
}) {
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(1);
  const [repacks, setRepacks] = useState<string[]>([]);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(false);

  const observer = useRef<IntersectionObserver | null>(null);

  const lastItemRef = useCallback(
    (node: HTMLDivElement | null) => {
      if (loading) return;
      if (observer.current) observer.current.disconnect();

      observer.current = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting && hasMore) {
          setPage((prev) => prev + 1);
        }
      });

      if (node) observer.current.observe(node);
    },
    [loading, hasMore],
  );

  useEffect(() => {
    if (query.trim() === "") {
      setRepacks([]);
      setHasMore(false);
      return;
    }

    const fetchRepacks = async () => {
      setLoading(true);
      try {
        const data: models.FitgirlPage = await GetRepacks(query, page);

        setRepacks((prev) =>
          page === 1 ? data.Results : [...prev, ...data.Results],
        );

        const totalLoaded = page * data.Results.length;
        setHasMore(totalLoaded < data.Total);
      } catch (err) {
        console.error("Error fetching repacks:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchRepacks();
  }, [query, page]);

  const handleSearch = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setPage(1);
  };

  function linkToName(link: string): string {
    let length = link.split("/").length;
    let name = link.split("/")[length - 2].replace(/-/g, " ");

    for (let i = 0; i < name.length; i++) {
      if (i == 0 || name[i - 1] == " ") {
        name =
          name.substring(0, i) +
          name.charAt(i).toUpperCase() +
          name.substring(i + 1);
      }
    }
    return name;
  }

  return (
    <div className="flex-1 p-10 flex flex-col h-full overflow-hidden bg-black">
      <h1 className="text-3xl font-bold mb-6 text-white">Search</h1>

      {/* Search Bar */}
      <form onSubmit={handleSearch} className="mb-6 flex space-x-3">
        <input
          type="text"
          placeholder="Search for a game..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="flex-1 bg-black border border-white text-white p-3 rounded-md outline-none focus:ring-1 focus:ring-white"
        />
        <button
          type="submit"
          className="bg-white text-black px-5 py-2 rounded-md font-semibold hover:bg-gray-200"
        >
          Search
        </button>
      </form>

      {/* Results list */}
      <div className="flex-1 overflow-y-auto scrollbar-thin scrollbar-thumb-white scrollbar-track-black">
        {repacks.length === 0 && !loading && (
          <p className="text-white text-center mt-20">
            No results. Try searching for something!
          </p>
        )}

        <div className="flex flex-col space-y-4">
          {repacks.map((title, index) => {
            if (repacks.length === index + 1) {
              return (
                <div
                  ref={lastItemRef}
                  key={index}
                  className="p-4 border border-white rounded-md text-white cursor-pointer hover:bg-white hover:text-black transition"
                  onClick={() => {
                    props.setGameUrl(title);
                    props.setPage("gameDetails");
                  }}
                >
                  {linkToName(title)}
                </div>
              );
            } else {
              return (
                <div
                  key={index}
                  className="p-4 border border-white rounded-md text-white cursor-pointer hover:bg-white hover:text-black transition"
                  onClick={() => {
                    props.setGameUrl(title);
                    props.setPage("gameDetails");
                  }}
                >
                  {linkToName(title)}
                </div>
              );
            }
          })}
        </div>

        {loading && (
          <div className="flex justify-center items-center py-8">
            <div className="w-8 h-8 border-4 border-white border-t-black rounded-full animate-spin"></div>
          </div>
        )}
      </div>
    </div>
  );
}
