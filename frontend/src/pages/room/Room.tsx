import { useRef } from "react";
import YouTube from "react-youtube";
import { useRoom } from "@/state";
import { useNavigate } from "react-router";

export default function Room() {
  const { messages, sendMessage, leaveRoom, currentRoom } = useRoom();

  const opts = {
    height: "390",
    width: "640",
    playerVars: {
      autoplay: 1,
    },
  };

  const navigate = useNavigate();

  if (!currentRoom) {
    navigate("/");
  }

  const playerRef = useRef<YouTube>(null);

  return (
    <div className="w-screen h-screen flex items-center justify-center">
      <div className=" flex-col gap-10  flex items-center justify-center">
        <YouTube videoId="2g811Eo7K8U" ref={playerRef} opts={opts} />
        <div className="flex gap-10">
          <button
            onClick={() => {
              playerRef.current?.internalPlayer.playVideo();
            }}
          >
            Play
          </button>
          <button
            onClick={() => {
              playerRef.current?.internalPlayer.pauseVideo();
            }}
          >
            Pause
          </button>
          <button
            onClick={() => {
              leaveRoom();
            }}
          >
            Leave
          </button>
        </div>
      </div>
      <div className="flex flex-col">
        <h1 className="p-2 font-medium">Messages</h1>
        <div className="flex flex-col gap-5 border rounded w-100 h-100 p-5 mt-5">
          {messages.map((message) => (
            <p key={message.id}>
              {message.author}:{message.body}
            </p>
          ))}
        </div>
        <div className="w-full flex justify-between gap-3 mt-2">
          <input
            type="text"
            className="focus:ring-0 outline-0 border rounded p-2 w-full"
            onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => {
              if (e.key === "Enter") {
                const value = (e.target as HTMLInputElement).value;
                if (value.trim() === "") {
                  return;
                }
                sendMessage(value);
                (e.target as HTMLInputElement).value = "";
              }
            }}
          />
        </div>
      </div>
    </div>
  );
}
