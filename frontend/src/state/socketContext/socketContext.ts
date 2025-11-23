import { createContext } from "react";

export type SocketContextType = {
  ws: WebSocket | null;
  setWs: React.Dispatch<React.SetStateAction<WebSocket | null>>;  
};

export const SocketContext = createContext<SocketContextType | null>(null);