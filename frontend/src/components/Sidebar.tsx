import { Download, Home, Library, Search } from "lucide-react";

export default function Sidebar(props: {
  currentPage: string;
  setPage: (page: string) => void;
}) {
  function getClass(selected: boolean) {
    if (selected) {
      return "bg-white text-black flex items-center cursor-pointer space-x-2 px-2 py-4 rounded-md";
    } else {
      return " hover:bg-white hover:text-black flex items-center cursor-pointer space-x-2 px-2 py-4 rounded-md";
    }
  }

  return (
    <div className="basis-1/6 text-white py-10 px-4 bg-black">
      <div>
        <h1 className="text-2xl">Fitgirl Launcher</h1>
      </div>
      <div className="mt-10">
        <nav className="flex-col space-y-2">
          <div
            className={
              props.currentPage == "home" ? getClass(true) : getClass(false)
            }
            onClick={() => props.setPage("home")}
          >
            <Home></Home>
            <span>Home</span>
          </div>
          <div
            className={
              props.currentPage == "library" ? getClass(true) : getClass(false)
            }
            onClick={() => props.setPage("library")}
          >
            <Library></Library>
            <span>Library</span>
          </div>
          <div
            className={
              props.currentPage == "downloads"
                ? getClass(true)
                : getClass(false)
            }
            onClick={() => props.setPage("downloads")}
          >
            <Download></Download>
            <span>Downloads</span>
          </div>
          <div
            className={
              props.currentPage == "search" ? getClass(true) : getClass(false)
            }
            onClick={() => props.setPage("search")}
          >
            <Search></Search>
            <span>Search</span>
          </div>
        </nav>
      </div>
    </div>
  );
}
