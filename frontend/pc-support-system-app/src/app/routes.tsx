import { BrowserRouter, Routes, Route } from "react-router-dom";

import Register from "../features/auth/pages/Register";

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<h1>Home</h1>} />
        <Route path="/register" element={<Register />} />
      </Routes>
    </BrowserRouter>
  );
}
