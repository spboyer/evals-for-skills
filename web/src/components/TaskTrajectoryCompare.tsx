import type { TaskResult, SessionDigest } from "../api/client";
import TrajectoryDiff from "./TrajectoryDiff";
import { formatNumber } from "../lib/format";
import {
  MessageSquare,
  Wrench,
  Sigma,
  ArrowDownToLine,
  ArrowUpFromLine,
} from "lucide-react";

function DigestDelta({ label, a, b, icon: Icon }: {
  label: string;
  a: number;
  b: number;
  icon: React.ComponentType<{ className?: string }>;
}) {
  const diff = b - a;
  const color = diff === 0 ? "text-zinc-400" : diff > 0 ? "text-red-400" : "text-green-400";
  const arrow = diff > 0 ? "↑" : diff < 0 ? "↓" : "";

  return (
    <div className="flex items-center gap-2">
      <Icon className="h-4 w-4 text-zinc-500" />
      <div className="min-w-0">
        <p className="text-xs text-zinc-500">{label}</p>
        <div className="flex items-baseline gap-2">
          <span className="text-sm font-medium text-zinc-200">{formatNumber(a)}</span>
          <span className="text-zinc-600">→</span>
          <span className="text-sm font-medium text-zinc-200">{formatNumber(b)}</span>
          {diff !== 0 && (
            <span className={`text-xs font-medium ${color}`}>
              {arrow}{formatNumber(Math.abs(diff))}
            </span>
          )}
        </div>
      </div>
    </div>
  );
}

function DigestComparison({ a, b }: { a: SessionDigest; b: SessionDigest }) {
  return (
    <div className="rounded-lg border border-zinc-700 bg-zinc-800 p-4 space-y-3">
      <h4 className="text-xs font-medium uppercase tracking-wider text-zinc-400">
        Session Digest Comparison
      </h4>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
        <DigestDelta label="Turns" a={a.totalTurns} b={b.totalTurns} icon={MessageSquare} />
        <DigestDelta label="Tool Calls" a={a.toolCallCount} b={b.toolCallCount} icon={Wrench} />
        <DigestDelta label="Tokens In" a={a.tokensIn} b={b.tokensIn} icon={ArrowDownToLine} />
        <DigestDelta label="Tokens Out" a={a.tokensOut} b={b.tokensOut} icon={ArrowUpFromLine} />
        <DigestDelta label="Tokens Total" a={a.tokensTotal} b={b.tokensTotal} icon={Sigma} />
      </div>
    </div>
  );
}

interface TaskTrajectoryCompareProps {
  taskName: string;
  taskA: TaskResult;
  taskB: TaskResult;
}

export default function TaskTrajectoryCompare({ taskName, taskA, taskB }: TaskTrajectoryCompareProps) {
  const hasTranscripts = taskA.transcript?.length && taskB.transcript?.length;
  const hasDigests = taskA.sessionDigest && taskB.sessionDigest;

  return (
    <div className="space-y-4">
      <h3 className="text-sm font-medium text-zinc-100">
        Trajectory: <span className="text-blue-400">{taskName}</span>
      </h3>

      {hasDigests && (
        <DigestComparison a={taskA.sessionDigest!} b={taskB.sessionDigest!} />
      )}

      {hasTranscripts ? (
        <TrajectoryDiff
          transcriptA={taskA.transcript!}
          transcriptB={taskB.transcript!}
        />
      ) : (
        <p className="text-sm text-zinc-500">
          Transcript data not available for{" "}
          {!taskA.transcript?.length && !taskB.transcript?.length
            ? "either run"
            : !taskA.transcript?.length
              ? "Run A"
              : "Run B"}
          .
        </p>
      )}
    </div>
  );
}
