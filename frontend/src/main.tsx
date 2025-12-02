import { createRoot } from "react-dom/client";
import "@/index.css";
import Router from "@/router/Router";
import { SocketContextProvider, RoomContextProvider } from "@/state";

createRoot(document.getElementById("root")!).render(<SocketContextProvider>
    <RoomContextProvider>
        <Router />
    </RoomContextProvider>
</SocketContextProvider>);