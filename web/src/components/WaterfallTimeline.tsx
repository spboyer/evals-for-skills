import { useState, useMemo } from "react";
import { CheckCircle2, XCircle, Clock } from "lucide-react";
import type { TranscriptEvent } from "../api/client";
import type { ToolSpan } from "../types/trajectory";
import { extractToolSpans } from "../types/trajectory";
import DetailPanel from "./DetailPanel";

// ---------------------------------------------------------------------------
// Status helpers
// ---------------------------------------------------------------------------

const barColor: Record<ToolSpan["status"], string> = {
  pass: "bg-green-500",
  fail: "bg-red-500",
  pending: "bg-blue-500/70",
};

const borderColor: Record<ToolSpan["status"], string> = {
  pass: "border-l-green-500",
  fail: "border-l-red-500",
  pending: "border-l-blue-400",
};

function StatusIcon({ status }: { status: ToolSpan["status"] }) {
  if (status === "pass")
    return <CheckCircle2 className="h-3.5 w-3.5 shrink-0 text-green-400" />;
  if (status === "fail")
    return <XCircle className="h-3.5 w-3.5 shrink-0 text-red-400" />;
  return <Clock className="h-3.5 w-3.5 shrink-0 text-blue-400" />;
}

// ---------------------------------------------------------------------------
// Column header with tick marks
// ---------------------------------------------------------------------------

function BarHeader({ totalEvents }: { totalEvents: number }) {
  const step =
    totalEvents <= 10
      ? 2
      : totalEvents <= 20
        ? 5
        : totalEvents <= 50
          ? 10
          : 25;
  const ticks: number[] = [];
  for (let i = 0; i <= totalEvents; i += step) ticks.push(i);
  if (ticks[ticks.length - 1] !== totalEvents) ticks.push(totalEvents);

  return (
    <div className="relative flex-1 min-w-0 h-7 px-2">
      {ticks.map((t) => {
        const pct = totalEvents > 0 ? (t / totalEvents) * 100 : 0;
        return (
          <div
            key={t}
            className="absolute top-0 flex flex-col items-center"
            style={{ left: `${pct}%` }}
          >
            <div className="h-2 w-px bg-zinc-600" />
            <span className="text-[9px] text-zinc-500 -translate-x-1/2 mt-px">
              {t}
            </span>
          </div>
        );
      })}
    </div>
  );
}

// ---------------------------------------------------------------------------
// Individual span row
// ---------------------------------------------------------------------------

interface SpanRowProps {
  span: ToolSpan;
  totalEvents: number;
  isSelected: boolean;
  onSelect: (span: ToolSpan) => void;
}

function SpanRow({ span, totalEvents, isSelected, onSelect }: SpanRowProps) {
  const [hovered, setHovered] = useState(false);

  const left = totalEvents > 1 ? (span.startIndex / (totalEvents - 1)) * 100 : 0;
  const width =
    totalEvents > 1
      ? Math.max(((span.duration + 1) / (totalEvents - 1)) * 100, 1.5)
      : 100;

  return (
    <div
      className={`flex items-stretch border-b border-zinc-700/50 last:border-b-0 cursor-pointer border-l-2 ${borderColor[span.status]} transition-colors ${
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
      {/* Tool name column */}
      <div className="flex items-center gap-2 px-3 py-2.5 w-44 shrink-0 border-r border-zinc-700/30">
        <StatusIcon status={span.status} />
        <span className="text-xs font-mono text-zinc-200 truncate">
          {span.toolName}
        </span>
      </div>

      {/* Duration column */}
      <div className="flex items-center justify-end px-3 py-2.5 w-20 shrink-0 border-r border-zinc-700/30">
        <span className="text-xs tabular-nums text-zinc-500">
          {span.duration + 1} events
        </span>
      </div>

      {/* Timeline bar area */}
      <div className="relative flex-1 min-w-0 flex items-center px-2">
        <div className="relative w-full h-4">
          <div
            className={`absolute top-0 h-full rounded-sm ${barColor[span.status]} opacity-80 hover:opacity-100 transition-opacity`}
            style={{
              left: `${left}%`,
              width: `${width}%`,
              minWidth: "6px",
            }}
          />
        </div>
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main component
// ---------------------------------------------------------------------------

interface WaterfallTimelineProps {
  events: TranscriptEvent[];
}

export default function WaterfallTimeline({ events }: WaterfallTimelineProps) {
  const [selectedSpan, setSelectedSpan] = useState<ToolSpan | null>(null);

  const spans = useMemo(() => extractToolSpans(events), [events]);

  return (
    <div className="flex rounded-lg border border-zinc-700 overflow-hidden">
      {/* Main timeline area */}
      <div className="flex-1 min-w-0 bg-zinc-900 overflow-x-auto">
        {/* Column headers */}
        <div className="flex items-stretch border-b border-zinc-700 bg-zinc-800/60 sticky top-0 z-10">
          <div className="flex items-center px-3 py-2 w-44 shrink-0 border-r border-zinc-700/30 border-l-2 border-l-transparent">
            <span className="text-[10px] font-medium uppercase tracking-wider text-zinc-500">
              Tool Call
            </span>
          </div>
          <div className="flex items-center justify-end px-3 py-2 w-20 shrink-0 border-r border-zinc-700/30">
            <span className="text-[10px] font-medium uppercase tracking-wider text-zinc-500">
              Events
            </span>
          </div>
          <BarHeader totalEvents={events.length} />
        </div>

        {/* Span rows */}
        {spans.length > 0 ? (
          spans.map((span) => (
            <SpanRow
              key={span.id}
              span={span}
              totalEvents={events.length}
              isSelected={selectedSpan?.id === span.id}
              onSelect={setSelectedSpan}
            />
          ))
        ) : (
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
  );
}
