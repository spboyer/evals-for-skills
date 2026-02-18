import { useParams, Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { fetchRunDetail } from "@/api/client";
import { Card } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { ArrowLeft } from "lucide-react";

export function RunDetail() {
  const { id } = useParams<{ id: string }>();

  const { data: run, isLoading } = useQuery({
    queryKey: ["run", id],
    queryFn: () => fetchRunDetail(id!),
    enabled: !!id,
  });

  if (isLoading) {
    return <p className="p-6 text-gray-500">Loading run…</p>;
  }

  if (!run) {
    return <p className="p-6 text-gray-500">Run not found.</p>;
  }

  return (
    <div className="mx-auto max-w-5xl space-y-6">
      <Link
        to="/"
        className="inline-flex items-center gap-1 text-sm text-waza-600 hover:underline dark:text-waza-400"
      >
        <ArrowLeft className="h-4 w-4" /> Back to dashboard
      </Link>

      <h1 className="text-2xl font-bold">{run.spec}</h1>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <Card title="Model">
          <p className="text-lg font-semibold">{run.model}</p>
        </Card>
        <Card title="Outcome">
          <Badge variant={run.outcome === "passed" ? "success" : "error"}>
            {run.outcome}
          </Badge>
        </Card>
        <Card title="Tasks">
          <p className="text-lg font-semibold">
            {run.passCount}/{run.taskCount}
          </p>
        </Card>
      </div>

      <Card title="Tasks">
        {run.tasks.length === 0 ? (
          <p className="text-gray-500">No task details available.</p>
        ) : (
          <ul className="space-y-2">
            {run.tasks.map((task) => (
              <li
                key={task.name}
                className="flex items-center justify-between border-b border-gray-100 py-2 dark:border-gray-800/50"
              >
                <span>{task.name}</span>
                <Badge variant={task.outcome === "passed" ? "success" : "error"}>
                  {task.outcome}
                </Badge>
              </li>
            ))}
          </ul>
        )}
      </Card>
    </div>
  );
}
