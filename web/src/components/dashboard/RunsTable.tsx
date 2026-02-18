import { Link } from "react-router-dom";
import { useRuns } from "@/hooks/useRuns";
import { Badge } from "@/components/ui/Badge";
import { formatDuration, formatCost } from "@/lib/format";

const outcomeBadge = (outcome: string) => {
  switch (outcome) {
    case "passed":
      return <Badge variant="success">Passed</Badge>;
    case "failed":
      return <Badge variant="danger">Failed</Badge>;
    default:
      return <Badge variant="warning">Error</Badge>;
  }
};

export function RunsTable() {
  const { data: runs, isLoading } = useRuns();

  if (isLoading) {
    return <p className="py-4 text-gray-500">Loading runs…</p>;
  }

  if (!runs?.length) {
    return <p className="py-4 text-gray-500">No evaluation runs yet.</p>;
  }

  return (
    <table className="w-full text-left text-sm">
      <thead>
        <tr className="border-b border-gray-200 text-gray-500 dark:border-gray-800 dark:text-gray-400">
          <th className="py-2 font-medium">Spec</th>
          <th className="py-2 font-medium">Model</th>
          <th className="py-2 font-medium">Outcome</th>
          <th className="py-2 font-medium">Tasks</th>
          <th className="py-2 font-medium">Duration</th>
          <th className="py-2 font-medium">Cost</th>
        </tr>
      </thead>
      <tbody>
        {runs.map((run) => (
          <tr
            key={run.id}
            className="border-b border-gray-100 dark:border-gray-800/50"
          >
            <td className="py-2">
              <Link
                to={`/runs/${run.id}`}
                className="text-waza-600 hover:underline dark:text-waza-400"
              >
                {run.spec}
              </Link>
            </td>
            <td className="py-2 text-gray-600 dark:text-gray-400">{run.model}</td>
            <td className="py-2">{outcomeBadge(run.outcome)}</td>
            <td className="py-2">
              {run.passCount}/{run.taskCount}
            </td>
            <td className="py-2">{formatDuration(run.duration)}</td>
            <td className="py-2">{formatCost(run.cost)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
