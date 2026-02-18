import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Shell } from "@/components/layout/Shell";
import { Dashboard } from "@/components/dashboard/Dashboard";
import { RunDetail } from "@/components/dashboard/RunDetail";

export default function App() {
  return (
    <BrowserRouter>
      <Shell>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/runs/:id" element={<RunDetail />} />
        </Routes>
      </Shell>
    </BrowserRouter>
  );
}
