import { useState, useEffect, useCallback } from "react";
import type { ReactNode } from "react";
import axios from "axios";
import { RoomContext, type Message, type Rooms, type Error } from "./roomContext";
import { useSocket } from "@/state/socketContext";
import { config } from "@/config";

type RoomContextProviderProps = {
  children: ReactNode;
};

export default function RoomContextProvider({ children }: RoomContextProviderProps) {
  const [currentRoom, setCurrentRoom] = useState<string | null>(null);
  const [username, setUsernameState] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [rooms, setRooms] = useState<Rooms | null>(null);
  const [error, setError] = useState<Error | null>(null);
  
  const { ws } = useSocket(); // Uses the socket from SocketContext
  
  const resetRoomState = useCallback(() => {
    setCurrentRoom(null);
    setMessages([]);
    setUsernameState(null);
    setError(null);
    localStorage.removeItem("token");
  }, []);

  // Listen for incoming messages and handle reconnection
  useEffect(() => {
    if (!ws) return;

    const handleOpen = () => {
      console.log("Connected to server");
      const token = localStorage.getItem("token");
      if (token) {
        ws.send(JSON.stringify({
          eventType: "reconnect_room",
          data: { token },
        }));
      }
    };

    const handleMessage = (event: MessageEvent) => {
      try {
        const eventData = JSON.parse(event.data);
        
        switch (eventData.eventType) {
          case "room_joined":
            localStorage.setItem("token", eventData.data.token);
            setCurrentRoom(eventData.data.roomName);
            break;
          case "room_left":
            resetRoomState();
            break;
          case "room_reconnected":
            console.log("setting current room on reconnect:", eventData.data.roomName);
            setCurrentRoom(eventData.data.roomName);
            setUsernameState(eventData.data.username);
            break;
          case "message_received":
            setMessages(prev => [...prev, {
              id: eventData.data.id || `${Date.now()}-${Math.random().toString(36).slice(2)}`,
              author: eventData.data.username,
              body: eventData.data.body
            }]);
            break;
          case "invalid_token":
            resetRoomState();
            break;
          case "error":
            setError(eventData.data);
            break;
          default:
            break;
        }
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    const handleClose = () => {
      console.log("Disconnected");
    };

    ws.addEventListener("open", handleOpen);
    ws.addEventListener("message", handleMessage);
    ws.addEventListener("close", handleClose);

    // Fetch rooms on connection
    fetchRooms();

    return () => {
      ws.removeEventListener("open", handleOpen);
      ws.removeEventListener("message", handleMessage);
      ws.removeEventListener("close", handleClose);
    };
  }, [ws]);

  const fetchRooms = useCallback(async () => {
    try {
      const response = await axios.get(`${config.apiUrl}/rooms`);
      setRooms(response.data);
    } catch (error) {
      console.error("Failed to fetch rooms:", error);
    }
  }, []);

  const setUsername = useCallback((newUsername: string) => {
    setUsernameState(newUsername);
  }, []);

  const joinRoom = useCallback((roomName: string, userName: string) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      setMessages([]); // Clear messages when joining a new room
      
      ws.send(JSON.stringify({
        eventType: "join_room",
        data: {
          roomName,
          username: userName,
        },
      }));
    }
  }, [ws]);

  const createRoom = useCallback((roomName: string, userName: string) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        eventType: "create_room",
        data: {
          roomName,
          username: userName,
        },
      }));
    }
    fetchRooms();
  }, [ws, fetchRooms]);

  const leaveRoom = useCallback(() => {
    const token = localStorage.getItem("token");
    
    if (!token) {
      console.error("no token");
      return;
    }

    if (ws && ws.readyState === WebSocket.OPEN && currentRoom && username) {
      ws.send(JSON.stringify({
        eventType: "leave_room",
        data: {
          roomName: currentRoom,
          username,
          token,
        },
      }));
    }
    setCurrentRoom(null);
    setMessages([]);
    fetchRooms();
  }, [ws, currentRoom, username, fetchRooms]);

  const sendMessage = useCallback((messageBody: string) => {
    if (ws && ws.readyState === WebSocket.OPEN && currentRoom && username) {
      ws.send(JSON.stringify({
        eventType: "send_message",
        data: {
          roomName: currentRoom,
          username,
          body: messageBody,
        },
      }));
    }
  }, [ws, currentRoom, username]);

  return (
    <RoomContext.Provider value={{
      currentRoom,
      username,
      rooms,
      error,
      messages,
      joinRoom,
      createRoom,
      leaveRoom,
      sendMessage,
      setUsername,
      fetchRooms,
    }}>
      {children}
    </RoomContext.Provider>
  );
}