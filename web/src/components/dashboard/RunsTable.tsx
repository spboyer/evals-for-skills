import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  flexRender,
  type ColumnDef,
  type SortingState,
} from "@tanstack/react-table";
import { useRuns } from "@/hooks/useRuns";
import { Badge } from "@/components/ui/Badge";
import {
  formatDuration,
  formatCost,
  formatNumber,
  formatRelativeTime,
} from "@/lib/format";
import type { Run } from "@/types";

function resultBadge(run: Run) {
  const allPassed = run.passCount === run.taskCount;
  return (
    <Badge variant={allPassed ? "success" : "error"}>
      {`${run.passCount}/${run.taskCount} ${allPassed ? "✅" : "❌"}`}
    </Badge>
  );
}

function SortIcon({ direction }: { direction: false | "asc" | "desc" }) {
  if (!direction) return <span className="ml-1 text-gray-600">⇅</span>;
  return (
    <span className="ml-1">
      {direction === "asc" ? "↑" : "↓"}
    </span>
  );
}

function SkeletonRows() {
  return (
    <>
      {Array.from({ length: 4 }).map((_, i) => (
        <tr
          key={i}
          className="border-b border-gray-800/50 animate-pulse"
        >
          {Array.from({ length: 8 }).map((_, j) => (
            <td key={j} className="py-3 px-3">
              <div className="h-4 rounded bg-gray-700/50" />
            </td>
          ))}
        </tr>
      ))}
    </>
  );
}

export function RunsTable() {
  const { data: runs, isLoading, isError, error } = useRuns();
  const navigate = useNavigate();
  const [sorting, setSorting] = useState<SortingState>([
    { id: "timestamp", desc: true },
  ]);

  const columns = useMemo<ColumnDef<Run>[]>(
    () => [
      {
        accessorKey: "id",
        header: "Run ID",
        cell: ({ getValue }) => (
          <span className="font-mono text-xs">
            {(getValue<string>()).slice(0, 8)}
          </span>
        ),
      },
      {
        accessorKey: "spec",
        header: "Eval Spec",
      },
      {
        accessorKey: "model",
        header: "Model",
      },
      {
        id: "result",
        header: "Result",
        accessorFn: (row) => row.passCount / (row.taskCount || 1),
        cell: ({ row }) => resultBadge(row.original),
      },
      {
        accessorKey: "tokens",
        header: "Tokens",
        cell: ({ getValue }) => formatNumber(getValue<number>()),
      },
      {
        accessorKey: "cost",
        header: "Cost",
        cell: ({ getValue }) => formatCost(getValue<number>()),
      },
      {
        accessorKey: "duration",
        header: "Duration",
        cell: ({ getValue }) => formatDuration(getValue<number>()),
      },
      {
        accessorKey: "timestamp",
        header: "Timestamp",
        cell: ({ getValue }) => (
          <span title={new Date(getValue<string>()).toLocaleString()}>
            {formatRelativeTime(getValue<string>())}
          </span>
        ),
      },
    ],
    [],
  );

  const table = useReactTable({
    data: runs ?? [],
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    enableSortingRemoval: true,
  });

  if (isError) {
    return (
      <div className="rounded-lg border border-red-800 bg-red-900/20 p-4 text-red-400">
        Error loading runs: {(error as Error).message}
      </div>
    );
  }

  return (
    <div className="overflow-x-auto rounded-lg border border-gray-800 bg-gray-900/50">
      <table className="w-full text-left text-sm">
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr
              key={headerGroup.id}
              className="border-b border-gray-700 bg-gray-800/60 text-gray-400"
            >
              {headerGroup.headers.map((header) => (
                <th
                  key={header.id}
                  className="select-none px-3 py-2.5 font-medium cursor-pointer hover:text-gray-200 transition-colors"
                  onClick={header.column.getToggleSortingHandler()}
                >
                  <span className="inline-flex items-center">
                    {flexRender(
                      header.column.columnDef.header,
                      header.getContext(),
                    )}
                    <SortIcon direction={header.column.getIsSorted()} />
                  </span>
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {isLoading ? (
            <SkeletonRows />
          ) : table.getRowModel().rows.length === 0 ? (
            <tr>
              <td
                colSpan={columns.length}
                className="py-12 text-center text-gray-500"
              >
                No evaluation runs yet. Run{" "}
                <code className="rounded bg-gray-800 px-1.5 py-0.5 font-mono text-xs text-gray-300">
                  waza run
                </code>{" "}
                to get started.
              </td>
            </tr>
          ) : (
            table.getRowModel().rows.map((row) => (
              <tr
                key={row.id}
                onClick={() => navigate(`/runs/${row.original.id}`)}
                className="cursor-pointer border-b border-gray-800/50 text-gray-300 transition-colors hover:bg-gray-800/40"
              >
                {row.getVisibleCells().map((cell) => (
                  <td key={cell.id} className="px-3 py-2.5">
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
