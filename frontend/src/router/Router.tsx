import { BrowserRouter, Route, Routes } from "react-router";
import Home from "../pages/home/Home";
import Room from "../pages/room/Room";

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
