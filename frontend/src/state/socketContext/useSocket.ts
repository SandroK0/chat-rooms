import { useContext } from "react";
import { SocketContext } from "./socketContext";

const useSocket = () => {
  const context = useContext(SocketContext);
  if (!context) {
    throw new Error("useSocket must be used within a SocketContextProvider");
  }
  return context;
}

export default useSocket;