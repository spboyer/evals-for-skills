import { type ReactNode } from "react";

interface TableProps {
  children: ReactNode;
  className?: string;
}

export function Table({ children, className = "" }: TableProps) {
  return (
    <div className="overflow-x-auto">
      <table
        className={`w-full text-left text-sm ${className}`}
      >
        {children}
      </table>
    </div>
  );
}
