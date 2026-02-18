import { Link } from "react-router-dom";
import { Zap } from "lucide-react";

export function Nav() {
  return (
    <header className="border-b border-gray-200 bg-white px-6 py-3 dark:border-gray-800 dark:bg-gray-950">
      <div className="flex items-center gap-2">
        <Link to="/" className="flex items-center gap-2 text-lg font-semibold">
          <Zap className="h-5 w-5 text-waza-500" />
          waza
        </Link>
      </div>
    </header>
  );
}
