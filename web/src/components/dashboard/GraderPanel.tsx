import { Badge } from "@/components/ui/Badge";
import type { GraderResult } from "@/types";

interface GraderPanelProps {
  graders: GraderResult[];
}

export function GraderPanel({ graders }: GraderPanelProps) {
  if (graders.length === 0) {
    return (
      <p className="py-2 pl-8 text-sm text-gray-500 dark:text-gray-400">
        No grader results available.
      </p>
    );
  }

  return (
    <div className="ml-8 mr-4 mb-2 space-y-2 rounded-md border border-gray-200 bg-gray-50 p-3 dark:border-gray-700 dark:bg-gray-800/60">
      <h4 className="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
        Grader Results
      </h4>
      {graders.map((g) => (
        <div
          key={g.name}
          className="flex flex-col gap-1 rounded border border-gray-200 bg-white p-2 text-sm dark:border-gray-700 dark:bg-gray-900"
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="font-medium text-gray-900 dark:text-gray-100">
                {g.name}
              </span>
              <span className="text-xs text-gray-500 dark:text-gray-400">
                ({g.type})
              </span>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-xs text-gray-500 dark:text-gray-400">
                Score: {g.score.toFixed(2)}
              </span>
              <Badge variant={g.passed ? "success" : "error"}>
                {g.passed ? "pass" : "fail"}
              </Badge>
            </div>
          </div>
          {g.message && (
            <p className="text-xs text-gray-600 dark:text-gray-400">
              {g.message}
            </p>
          )}
        </div>
      ))}
    </div>
  );
}
