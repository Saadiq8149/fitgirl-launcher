export default function GameCard(props: {
  coverImage: string;
  link: string;
  setPage: (page: string) => void;
  setGameUrl: (url: string) => void;
}) {
  let length = props.link.split("/").length;
  let name = props.link.split("/")[length - 2].replace(/-/g, " ");

  for (let i = 0; i < name.length; i++) {
    if (i == 0 || name[i - 1] == " ") {
      name =
        name.substring(0, i) +
        name.charAt(i).toUpperCase() +
        name.substring(i + 1);
    }
  }

  return (
    <div
      className="w-48 flex-shrink-0 cursor-pointer group"
      onClick={() => {
        props.setGameUrl(props.link);
        props.setPage("gameDetails");
      }}
    >
      {/* Image container */}
      <div className="relative h-64 overflow-hidden rounded-md">
        <img
          src={props.coverImage}
          alt={name}
          className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-110"
        />

        {/* Hover overlay */}
        <div
          className="
               absolute inset-0
               bg-black/60
               opacity-0
               group-hover:opacity-100
               flex items-center justify-center
               rounded-md
             "
        ></div>
      </div>

      {/* Static title below image */}
      <p className="text-white text-center mt-2 text-sm font-medium truncate">
        {name}
      </p>
    </div>
  );
}
