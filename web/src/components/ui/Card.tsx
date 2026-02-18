import { type ReactNode } from "react";

interface CardProps {
  title?: string;
  children: ReactNode;
  className?: string;
  hover?: boolean;
}

export function Card({ title, children, className = "", hover = false }: CardProps) {
  return (
    <div
      className={`rounded-lg border p-4 ${hover ? "transition-colors hover:border-[var(--accent-primary)]/40" : ""} ${className}`}
      style={{
        backgroundColor: "var(--bg-surface)",
        borderColor: "var(--border-default)",
      }}
    >
      {title && (
        <h3 className="mb-2 text-sm font-medium" style={{ color: "var(--text-secondary)" }}>
          {title}
        </h3>
      )}
      {children}
    </div>
  );
}
