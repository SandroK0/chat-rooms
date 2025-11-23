import { createContext } from "react";

export type Room = {
  name: string;
};

export type Rooms = {
  [roomName: string]: Room;
};

export type Error = {
  code: string;
  message: string;
};

export type RoomContextType = {
  currentRoom: string | null;
  username: string | null;
  rooms: Rooms | null;
  error: Error | null;
  joinRoom: (roomName: string, username: string) => void;
  createRoom: (roomName: string, username: string) => void;
  leaveRoom: () => void;
  sendMessage: (messageBody: string) => void;
  setUsername: (username: string) => void;
  fetchRooms: () => void;
  messages: Message[];
};

export type Message = {
  author: string;
  body: string;
  id: string;
};

export const RoomContext = createContext<RoomContextType | null>(null);