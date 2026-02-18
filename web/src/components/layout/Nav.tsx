import { Link } from "react-router-dom";
import { Badge } from "@/components/ui/Badge";

export function Nav() {
  return (
    <header
      className="flex h-14 items-center justify-between border-b px-6"
      style={{
        backgroundColor: "var(--bg-surface)",
        borderColor: "var(--border-default)",
      }}
    >
      <Link
        to="/"
        className="flex items-center gap-2 text-base font-semibold"
        style={{ color: "var(--text-primary)" }}
      >
        <span>🔥</span>
        <span>waza</span>
      </Link>
      <Badge variant="info">v0.4.0-alpha.1</Badge>
    </header>
  );
}
