import { type ReactNode } from "react";
import { Nav } from "./Nav";
import { Footer } from "./Footer";

export function Shell({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col bg-white dark:bg-gray-950">
      <Nav />
      <main className="flex-1 p-6">{children}</main>
      <Footer />
    </div>
  );
}
