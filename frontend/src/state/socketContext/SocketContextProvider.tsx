import { useEffect, useState } from "react";
import type { ReactNode } from "react";
import { SocketContext } from "./socketContext";
import { config } from "../../config";

type SocketContextProviderProps = {
  children: ReactNode;
};

export default function SocketContextProvider({ children }: SocketContextProviderProps) {
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    const socket = new WebSocket(config.websocketUrl);
    setWs(socket);
    
    return () => {
      socket.close();
    };
  }, []);
  
  return (
    <SocketContext.Provider value={{ ws, setWs }}>
      {children}
    </SocketContext.Provider>
  );
}


