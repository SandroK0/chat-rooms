import { useRef, useState } from "react";
import { useNavigate } from "react-router";
import useRoom from "../../state/roomContext/useRoom";

function Home() {
  const navigate = useNavigate();
  const {
    rooms,
    currentRoom,
    username,
    error,
    joinRoom,
    createRoom,
    leaveRoom,
    setUsername,
  } = useRoom();
  
  const [createRoomName, setCreateRoomName] = useState<string>("");

  const usernameInputRef = useRef<HTMLInputElement | null>(null);

  const handleCreateRoom = () => {
    const usernameValue = usernameInputRef.current?.value;
    if (usernameValue && createRoomName.trim()) {
      setUsername(usernameValue);
      createRoom(createRoomName, usernameValue);
      setCreateRoomName("");
      navigate("/room");
    }
  };

  const handleJoinRoom = (roomName: string) => {
    const usernameValue = usernameInputRef.current?.value;
    if (usernameValue) {
      setUsername(usernameValue);
      joinRoom(roomName, usernameValue);
      navigate("/room");
    }
  };

  return (
    <div className="flex gap-10 flex-col items-center justify-center relative w-screen min-h-screen ">
      {error && (
        <div className="text-red-800 absolute bg-amber-50 top-10">
          {error.code}:{error.message}
        </div>
      )}
      <div className="flex items-center gap-5">
        <input
          placeholder="Enter username"
          className="focus:ring-0 outline-0 border rounded p-2"
          ref={usernameInputRef}
          defaultValue={username || ""}
        />
      </div>
      <div className="mt-40 flex gap-10 ">
        <div>
          <div className="flex justify-between">
            <input
              type="text"
              className="focus:ring-0 outline-0 border rounded p-2"
              onChange={(e) => setCreateRoomName(e.target.value)}
              value={createRoomName}
            />
            <button onClick={handleCreateRoom}>Create Room</button>
            <button onClick={leaveRoom}>Leave Room</button>
          </div>
          <div className="flex flex-col gap-2 border rounded w-100 h-100 mt-5">
            {rooms &&
              Object.keys(rooms).map((key) => (
                <div
                  key={key}
                  className="flex justify-between border-b px-4 py-1 items-center"
                >
                  <div>{key}</div>
                  <button onClick={() => handleJoinRoom(key)}>Join</button>
                </div>
              ))}
          </div>
          <div>Current Room: {currentRoom || "None"}</div>
          <div>Username: {username || "Not set"}</div>
        </div>
      </div>
    </div>
  );
}

export default Home;
