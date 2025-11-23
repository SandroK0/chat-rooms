import { createRoot } from "react-dom/client";
import "./index.css";
import Router from "./router/Router.tsx";
import SocketContextProvider from "./state/socketContext/SocketContextProvider.tsx";
import { RoomContextProvider } from "./state/roomContext/index.ts";

createRoot(document.getElementById("root")!).render(<SocketContextProvider>
    <RoomContextProvider>
        <Router />
    </RoomContextProvider>
</SocketContextProvider>);