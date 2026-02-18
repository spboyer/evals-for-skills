import { useParams, useNavigate, Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { fetchRunDetail } from "@/api/client";
import { Card } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Table } from "@/components/ui/Table";
import { TaskRow } from "./TaskRow";
import { formatDuration, formatCost, formatNumber } from "@/lib/format";
import { ArrowLeft } from "lucide-react";

function LoadingSkeleton() {
  return (
    <div className="mx-auto max-w-5xl animate-pulse space-y-6">
      <div className="h-4 w-40 rounded bg-gray-200 dark:bg-gray-700" />
      <div className="h-8 w-64 rounded bg-gray-200 dark:bg-gray-700" />
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3 md:grid-cols-6">
        {Array.from({ length: 6 }).map((_, i) => (
          <div
            key={i}
            className="h-20 rounded-lg border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900"
          />
        ))}
      </div>
      <div className="h-48 rounded-lg border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900" />
    </div>
  );
}

export function RunDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: run, isLoading, error } = useQuery({
    queryKey: ["run", id],
    queryFn: () => fetchRunDetail(id!),
    enabled: !!id,
    retry: false,
  });

  if (!id) {
    return (
      <div className="mx-auto max-w-5xl p-6">
        <p className="text-gray-500 dark:text-gray-400">Invalid run ID.</p>
      </div>
    );
  }

  if (isLoading) {
    return <LoadingSkeleton />;
  }

  if (error) {
    const is404 = error instanceof Error && error.message.includes("404");
    return (
      <div className="mx-auto max-w-5xl space-y-4 p-6">
        <button
          onClick={() => navigate("/")}
          className="inline-flex items-center gap-1 text-sm text-waza-600 hover:underline dark:text-waza-400"
        >
          <ArrowLeft className="h-4 w-4" /> Back to Dashboard
        </button>
        <Card>
          <p className="text-gray-500 dark:text-gray-400">
            {is404
              ? `Run "${id}" not found.`
              : "Failed to load run details. Please try again."}
          </p>
        </Card>
      </div>
    );
  }

  if (!run) {
    return (
      <div className="mx-auto max-w-5xl p-6">
        <p className="text-gray-500 dark:text-gray-400">Run not found.</p>
      </div>
    );
  }

  const shortId = id.length > 8 ? id.slice(0, 8) : id;

  return (
    <div className="mx-auto max-w-5xl space-y-6">
      {/* Breadcrumb + Back */}
      <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
        <button
          onClick={() => navigate("/")}
          className="inline-flex items-center gap-1 text-waza-600 hover:underline dark:text-waza-400"
        >
          <ArrowLeft className="h-4 w-4" /> Back to Dashboard
        </button>
        <span>/</span>
        <Link to="/" className="hover:underline">
          Dashboard
        </Link>
        <span>&gt;</span>
        <span className="text-gray-900 dark:text-gray-100">
          Run #{shortId}
        </span>
      </div>

      {/* Run Summary Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          {run.spec}
        </h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          Run #{shortId}
        </p>
      </div>

      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-6">
        <Card title="Eval Spec">
          <p className="text-sm font-semibold text-waza-600 dark:text-waza-400">
            {run.spec}
          </p>
        </Card>
        <Card title="Model">
          <p className="text-sm font-semibold text-waza-600 dark:text-waza-400">
            {run.model}
          </p>
        </Card>
        <Card title="Result">
          <Badge variant={run.outcome === "passed" ? "success" : "error"}>
            {run.outcome}
          </Badge>
        </Card>
        <Card title="Total Tokens">
          <p className="text-sm font-semibold text-waza-600 dark:text-waza-400">
            {formatNumber(run.tokens)}
          </p>
        </Card>
        <Card title="Cost">
          <p className="text-sm font-semibold text-waza-600 dark:text-waza-400">
            {formatCost(run.cost)}
          </p>
        </Card>
        <Card title="Duration">
          <p className="text-sm font-semibold text-waza-600 dark:text-waza-400">
            {formatDuration(run.duration)}
          </p>
        </Card>
      </div>

      {/* Tasks Table */}
      <Card title={`Tasks (${run.passCount}/${run.taskCount} passed)`}>
        {run.tasks.length === 0 ? (
          <p className="text-gray-500 dark:text-gray-400">
            No task details available.
          </p>
        ) : (
          <Table>
            <thead>
              <tr className="border-b border-gray-200 text-xs uppercase tracking-wide text-gray-500 dark:border-gray-700 dark:text-gray-400">
                <th className="w-8 py-2 pl-4 pr-2" />
                <th className="py-2 pr-4 text-left">Task Name</th>
                <th className="py-2 pr-4 text-left">Outcome</th>
                <th className="py-2 pr-4 text-left">Score</th>
                <th className="py-2 pr-4 text-left">Duration</th>
              </tr>
            </thead>
            <tbody>
              {run.tasks.map((task) => (
                <TaskRow key={task.name} task={task} />
              ))}
            </tbody>
          </Table>
        )}
      </Card>
    </div>
  );
}
