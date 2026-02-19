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

/** Align tool calls from two transcripts by tool name + position. */
export function alignToolCalls(
  transcriptA: TranscriptEvent[],
  transcriptB: TranscriptEvent[],
): DiffEntry[] {
  const callsA = extractToolCalls(transcriptA);
  const callsB = extractToolCalls(transcriptB);
  const usedB = new Set<number>();
  const entries: DiffEntry[] = [];

  for (let i = 0; i < callsA.length; i++) {
    const a = callsA[i]!;
    const name = a.toolName ?? "unknown";
    const bAtI = i < callsB.length ? callsB[i]! : undefined;

    // Perfect match: same name at same index
    if (bAtI && bAtI.toolName === name && !usedB.has(i)) {
      usedB.add(i);
      const same = jsonEqual(a.arguments, bAtI.arguments) && jsonEqual(a.toolResult, bAtI.toolResult);
      entries.push({ kind: same ? "matched" : "changed", index: i, toolName: name, a, b: bAtI });
      continue;
    }
    // Fallback: find same name elsewhere in B
    const alt = callsB.findIndex((b, j) => !usedB.has(j) && b.toolName === name);
    if (alt !== -1) {
      const bAlt = callsB[alt]!;
      usedB.add(alt);
      const same = jsonEqual(a.arguments, bAlt.arguments) && jsonEqual(a.toolResult, bAlt.toolResult);
      entries.push({ kind: same ? "matched" : "changed", index: i, toolName: name, a, b: bAlt });
    } else {
      entries.push({ kind: "insertion", index: i, toolName: name, a });
    }
  }

  // Remaining unmatched in B → deletions
  for (let j = 0; j < callsB.length; j++) {
    if (!usedB.has(j)) {
      const bJ = callsB[j]!;
      entries.push({ kind: "deletion", index: j, toolName: bJ.toolName ?? "unknown", b: bJ });
    }
  }

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
