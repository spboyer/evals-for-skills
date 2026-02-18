import { type ReactNode } from "react";
import { Nav } from "./Nav";
import { Footer } from "./Footer";

export function Shell({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col" style={{ backgroundColor: "var(--bg-base)" }}>
      <Nav />
      <main className="mx-auto w-full max-w-7xl flex-1 px-6 py-8">{children}</main>
      <Footer />
    </div>
  );
}
