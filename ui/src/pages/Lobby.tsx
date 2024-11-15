import { Room } from "../components/room/Room.tsx";
import { Controls } from "../components/Controls.tsx";
import { Console } from "../components/console";
import { useSelector } from "react-redux";
import { getConsoleState } from "../state/room.reducer.ts";

const Lobby = () => {
  const displayConsole = useSelector(getConsoleState);

  return (
    <div className="relative w-full h-full">
      {displayConsole && (
        <div className="absolute w-full h-[90%] flex items-center justify-center">
          <Console />
        </div>
      )}
      <Room />
      <Controls />
    </div>
  );
};

export default Lobby;
