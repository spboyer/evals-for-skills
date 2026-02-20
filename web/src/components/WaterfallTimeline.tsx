import { useState, useMemo } from "react";
import { CheckCircle2, XCircle, Clock, Wrench } from "lucide-react";
import type { TranscriptEvent } from "../api/client";
import type { ToolSpan } from "../types/trajectory";
import { extractToolSpans } from "../types/trajectory";
import DetailPanel from "./DetailPanel";

const LABEL_WIDTH = 200;

interface WaterfallTimelineProps {
  events: TranscriptEvent[];
}

// Track how many times each tool has been called (for the #N badge)
function buildCallIndex(spans: ToolSpan[]): Map<string, number> {
  const counter = new Map<string, number>();
  const callIndex = new Map<string, number>();
  for (const span of spans) {
    const n = (counter.get(span.toolName) ?? 0) + 1;
    counter.set(span.toolName, n);
    callIndex.set(span.id, n);
  }
  return callIndex;
}

const statusIcon: Record<
  ToolSpan["status"],
  React.ComponentType<{ className?: string }>
> = {
  pass: CheckCircle2,
  fail: XCircle,
  pending: Clock,
};

const statusIconColor: Record<ToolSpan["status"], string> = {
  pass: "text-green-400",
  fail: "text-red-400",
  pending: "text-yellow-400",
};

const barColor: Record<ToolSpan["status"], string> = {
  pass: "bg-teal-500 hover:bg-teal-400",
  fail: "bg-red-500 hover:bg-red-400",
  pending: "bg-yellow-500 hover:bg-yellow-400",
};

// Axis tick marks
function TraceAxis({ totalEvents }: { totalEvents: number }) {
  const ticks = [0, 25, 50, 75, 100];
  return (
    <div
      className="relative flex-1 h-6 border-b border-zinc-700 select-none"
    >
      {ticks.map((pct) => {
        const evNum = Math.round((pct / 100) * totalEvents);
        return (
          <div
            key={pct}
            className="absolute top-0 flex flex-col items-center"
            style={{ left: `${pct}%` }}
          >
            <div className="h-2 w-px bg-zinc-600" />
            <span className="text-[10px] text-zinc-500 -translate-x-1/2 whitespace-nowrap">
              {evNum}
            </span>
          </div>
        );
      })}
    </div>
  );
}

// Single span row
function SpanRow({
  span,
  callNum,
  totalEvents,
  isSelected,
  onSelect,
}: {
  span: ToolSpan;
  callNum: number;
  totalEvents: number;
  isSelected: boolean;
  onSelect: (span: ToolSpan) => void;
}) {
  const [hovered, setHovered] = useState(false);
  const Icon = statusIcon[span.status];

  const left = totalEvents > 0 ? (span.startIndex / totalEvents) * 100 : 0;
  const width =
    totalEvents > 0 ? Math.max((span.duration / totalEvents) * 100, 1) : 1;

  return (
    <button
      className={`flex w-full items-stretch border-b border-zinc-700/50 text-left transition-colors last:border-b-0 ${
        isSelected
          ? "bg-zinc-700/60"
          : hovered
            ? "bg-zinc-800/80"
            : "bg-transparent"
      }`}
      onClick={() => onSelect(span)}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      {/* Label column */}
      <div
        className="shrink-0 flex items-center gap-2 px-3 py-2.5 border-r border-zinc-700"
        style={{ width: LABEL_WIDTH }}
      >
        <Icon className={`h-3.5 w-3.5 shrink-0 ${statusIconColor[span.status]}`} />
        <span className="text-xs font-mono text-zinc-300 truncate flex-1">
          {span.toolName}
        </span>
        {callNum > 1 && (
          <span className="shrink-0 rounded bg-zinc-700 px-1 py-0.5 text-[10px] text-zinc-400">
            #{callNum}
          </span>
        )}
      </div>

      {/* Bar column */}
      <div className="relative flex-1 py-2">
        {/* duration bar */}
        <div
          className={`absolute top-2 bottom-2 rounded-sm transition-opacity ${barColor[span.status]} ${
            isSelected ? "ring-1 ring-white/50" : ""
          }`}
          style={{
            left: `${left}%`,
            width: `${width}%`,
            minWidth: "6px",
          }}
        />

        {/* hover tooltip */}
        {hovered && (
          <div
            className="pointer-events-none absolute z-20 -top-7 rounded bg-zinc-700 px-2 py-1 text-xs text-zinc-200 whitespace-nowrap shadow-lg"
            style={{ left: `${Math.min(left, 70)}%` }}
          >
            <Wrench className="inline h-3 w-3 mr-1 text-zinc-400" />
            {span.toolName} — {span.duration} events
          </div>
        )}
      </div>
    </button>
  );
}

export default function WaterfallTimeline({ events }: WaterfallTimelineProps) {
  const [selectedSpan, setSelectedSpan] = useState<ToolSpan | null>(null);

  const spans = useMemo(() => extractToolSpans(events), [events]);
  const callIndex = useMemo(() => buildCallIndex(spans), [spans]);

  // Unique tool names for summary
  const toolSummary = useMemo(() => {
    const counts = new Map<string, number>();
    for (const s of spans) counts.set(s.toolName, (counts.get(s.toolName) ?? 0) + 1);
    return Array.from(counts.entries());
  }, [spans]);

  return (
    <div className="flex flex-col rounded-lg border border-zinc-700 overflow-hidden">
      {/* Trace header */}
      <div className="flex items-center gap-3 px-4 py-2 bg-zinc-800 border-b border-zinc-700">
        <span className="text-xs font-semibold text-zinc-300 uppercase tracking-wide">
          Trace
        </span>
        <span className="text-xs text-zinc-500">·</span>
        <span className="text-xs text-zinc-400">{spans.length} spans</span>
        {toolSummary.length > 0 && (
          <>
            <span className="text-xs text-zinc-600">·</span>
            <span className="text-xs text-zinc-500 truncate">
              {toolSummary.map(([name, count]) => `${name} × ${count}`).join("  ")}
            </span>
          </>
        )}
      </div>

      {/* Main area */}
      <div className="flex flex-1 overflow-hidden">
        {/* Timeline panel */}
        <div className="flex-1 min-w-0 bg-zinc-900 overflow-x-auto">
          {/* Column headers */}
          <div className="sticky top-0 z-10 flex bg-zinc-900 border-b border-zinc-700">
            <div
              className="shrink-0 flex items-end px-3 pb-1 border-r border-zinc-700"
              style={{ width: LABEL_WIDTH }}
            >
              <span className="text-[10px] text-zinc-500 uppercase tracking-wide">
                Tool
              </span>
            </div>
            <TraceAxis totalEvents={events.length} />
          </div>

          {/* Span rows */}
          <div>
            {spans.map((span) => (
              <SpanRow
                key={span.id}
                span={span}
                callNum={callIndex.get(span.id) ?? 1}
                totalEvents={events.length}
                isSelected={selectedSpan?.id === span.id}
                onSelect={setSelectedSpan}
              />
            ))}
          </div>

          {spans.length === 0 && (
            <div className="p-6 text-center text-sm text-zinc-500">
              No tool calls found in transcript
            </div>
          )}
        </div>

        {/* Detail sidebar */}
        {selectedSpan && (
          <DetailPanel
            span={selectedSpan}
            onClose={() => setSelectedSpan(null)}
          />
        )}
      </div>
    </div>
  );
}
