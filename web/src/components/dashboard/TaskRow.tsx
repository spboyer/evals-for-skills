import { useState } from "react";
import { Badge } from "@/components/ui/Badge";
import { GraderPanel } from "./GraderPanel";
import { formatDuration } from "@/lib/format";
import type { Task } from "@/types";
import { ChevronRight, ChevronDown } from "lucide-react";

interface TaskRowProps {
  task: Task;
}

export function TaskRow({ task }: TaskRowProps) {
  const [expanded, setExpanded] = useState(false);

  return (
    <>
      <tr
        className="cursor-pointer border-b border-gray-100 hover:bg-gray-50 dark:border-gray-800 dark:hover:bg-gray-800/50"
        onClick={() => setExpanded(!expanded)}
      >
        <td className="py-3 pl-4 pr-2">
          {expanded ? (
            <ChevronDown className="h-4 w-4 text-gray-400" />
          ) : (
            <ChevronRight className="h-4 w-4 text-gray-400" />
          )}
        </td>
        <td className="py-3 pr-4 font-medium text-gray-900 dark:text-gray-100">
          {task.name}
        </td>
        <td className="py-3 pr-4">
          <Badge
            variant={
              task.outcome === "passed"
                ? "success"
                : task.outcome === "error"
                  ? "warning"
                  : "error"
            }
          >
            {task.outcome}
          </Badge>
        </td>
        <td className="py-3 pr-4 text-gray-600 dark:text-gray-400">
          {task.score.toFixed(2)}
        </td>
        <td className="py-3 pr-4 text-gray-600 dark:text-gray-400">
          {formatDuration(task.duration)}
        </td>
      </tr>
      {expanded && (
        <tr>
          <td colSpan={5} className="pb-2">
            <GraderPanel graders={task.graderResults} />
          </td>
        </tr>
      )}
    </>
  );
}
