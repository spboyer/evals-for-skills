import { useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";
import type { TranscriptEvent } from "../api/client";

/** A matched pair or solo entry from the alignment algorithm. */
export interface DiffEntry {
  kind: "matched" | "changed" | "insertion" | "deletion";
  index: number;
  toolName: string;
  a?: TranscriptEvent;
  b?: TranscriptEvent;
}

/** Extract ToolExecutionStart events preserving order. */
function extractToolCalls(transcript: TranscriptEvent[]): TranscriptEvent[] {
  return transcript.filter((e) => e.type === "ToolExecutionStart");
}

/** Deep-equal check for JSON-serialisable values. */
function jsonEqual(a: unknown, b: unknown): boolean {
  return JSON.stringify(a) === JSON.stringify(b);
}

/**
 * Align tool calls from two transcripts using LCS on tool names.
 * Produces a correct diff even when the same tool is called multiple times.
 */
export function alignToolCalls(
  transcriptA: TranscriptEvent[],
  transcriptB: TranscriptEvent[],
): DiffEntry[] {
  const callsA = extractToolCalls(transcriptA);
  const callsB = extractToolCalls(transcriptB);
  const m = callsA.length;
  const n = callsB.length;

  // Build LCS table on tool names
  const dp: number[][] = Array.from({ length: m + 1 }, () => Array(n + 1).fill(0));
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      if ((callsA[i - 1]!.toolName ?? "") === (callsB[j - 1]!.toolName ?? "")) {
        dp[i]![j] = dp[i - 1]![j - 1]! + 1;
      } else {
        dp[i]![j] = Math.max(dp[i - 1]![j]!, dp[i]![j - 1]!);
      }
    }
  }

  // Backtrack to produce diff entries in order
  const entries: DiffEntry[] = [];
  let i = m, j = n;
  const stack: DiffEntry[] = [];

  while (i > 0 || j > 0) {
    if (i > 0 && j > 0 && (callsA[i - 1]!.toolName ?? "") === (callsB[j - 1]!.toolName ?? "")) {
      const a = callsA[i - 1]!;
      const b = callsB[j - 1]!;
      const name = a.toolName ?? "unknown";
      const same = jsonEqual(a.arguments, b.arguments) && jsonEqual(a.toolResult, b.toolResult);
      stack.push({ kind: same ? "matched" : "changed", index: i - 1, toolName: name, a, b });
      i--;
      j--;
    } else if (j > 0 && (i === 0 || dp[i]![j - 1]! >= dp[i - 1]![j]!)) {
      const b = callsB[j - 1]!;
      stack.push({ kind: "deletion", index: j - 1, toolName: b.toolName ?? "unknown", b });
      j--;
    } else {
      const a = callsA[i - 1]!;
      stack.push({ kind: "insertion", index: i - 1, toolName: a.toolName ?? "unknown", a });
      i--;
    }
  }

  // Reverse (backtrack produces entries in reverse order)
  while (stack.length > 0) entries.push(stack.pop()!);

  return entries;
}

const kindColors: Record<DiffEntry["kind"], string> = {
  matched: "border-green-500/40 bg-green-500/5",
  changed: "border-yellow-500/40 bg-yellow-500/5",
  insertion: "border-red-500/40 bg-red-500/5",
  deletion: "border-red-500/40 bg-red-500/5",
};

const kindLabels: Record<DiffEntry["kind"], { text: string; cls: string }> = {
  matched: { text: "Match", cls: "text-green-500" },
  changed: { text: "Changed", cls: "text-yellow-500" },
  insertion: { text: "Only in A", cls: "text-red-500" },
  deletion: { text: "Only in B", cls: "text-red-500" },
};

function JsonExpander({ label, data }: { label: string; data: unknown }) {
  const [open, setOpen] = useState(false);
  if (data === undefined || data === null) return null;
  const text = typeof data === "string" ? data : JSON.stringify(data, null, 2);

  return (
    <div>
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-1 text-xs text-zinc-500 hover:text-zinc-300"
      >
        {open ? <ChevronDown className="h-3 w-3" /> : <ChevronRight className="h-3 w-3" />}
        {label}
      </button>
      {open && (
        <pre className="mt-1 max-h-48 overflow-auto rounded bg-zinc-900 p-2 text-xs text-zinc-300">
          <code>{text}</code>
        </pre>
      )}
    </div>
  );
}

function DiffRow({ entry }: { entry: DiffEntry }) {
  const { kind, toolName, a, b } = entry;
  const label = kindLabels[kind];

  return (
    <div className={`rounded-lg border p-3 space-y-2 ${kindColors[kind]}`}>
      <div className="flex items-center gap-3">
        <span className="text-sm font-medium text-zinc-200">{toolName}</span>
        <span className={`text-xs font-medium ${label.cls}`}>{label.text}</span>
      </div>

      {kind === "changed" && (
        <div className="grid gap-3 sm:grid-cols-2">
          <div>
            <p className="mb-1 text-xs font-medium text-zinc-500">Run A</p>
            <JsonExpander label="Arguments" data={a?.arguments} />
            <JsonExpander label="Result" data={a?.toolResult} />
          </div>
          <div>
            <p className="mb-1 text-xs font-medium text-zinc-500">Run B</p>
            <JsonExpander label="Arguments" data={b?.arguments} />
            <JsonExpander label="Result" data={b?.toolResult} />
          </div>
        </div>
      )}

      {kind === "matched" && (
        <div>
          <JsonExpander label="Arguments" data={a?.arguments} />
          <JsonExpander label="Result" data={a?.toolResult} />
        </div>
      )}

      {kind === "insertion" && (
        <div>
          <JsonExpander label="Arguments (A)" data={a?.arguments} />
          <JsonExpander label="Result (A)" data={a?.toolResult} />
        </div>
      )}

      {kind === "deletion" && (
        <div>
          <JsonExpander label="Arguments (B)" data={b?.arguments} />
          <JsonExpander label="Result (B)" data={b?.toolResult} />
        </div>
      )}
    </div>
  );
}

interface TrajectoryDiffProps {
  transcriptA: TranscriptEvent[];
  transcriptB: TranscriptEvent[];
}

export default function TrajectoryDiff({ transcriptA, transcriptB }: TrajectoryDiffProps) {
  const entries = alignToolCalls(transcriptA, transcriptB);

  if (entries.length === 0) {
    return (
      <p className="text-sm text-zinc-500">No tool calls found in either transcript.</p>
    );
  }

  const counts = { matched: 0, changed: 0, insertion: 0, deletion: 0 };
  for (const e of entries) counts[e.kind]++;

  return (
    <div className="space-y-3">
      {/* Legend */}
      <div className="flex flex-wrap gap-4 text-xs">
        <span className="text-green-500">● {counts.matched} matched</span>
        <span className="text-yellow-500">● {counts.changed} changed</span>
        <span className="text-red-500">● {counts.insertion + counts.deletion} missing</span>
      </div>

      {entries.map((entry, i) => (
        <DiffRow key={`${entry.kind}-${entry.toolName}-${i}`} entry={entry} />
      ))}
    </div>
  );
}
