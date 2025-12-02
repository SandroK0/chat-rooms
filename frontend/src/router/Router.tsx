import { BrowserRouter, Route, Routes } from "react-router";
import { Home, Room } from "@/pages";

export default function Router() {
  return (
    <BrowserRouter>
      <Routes>
        <Route index element={<Home />} />
        <Route path="room" element={<Room />} />
      </Routes>
    </BrowserRouter>
  );
}
