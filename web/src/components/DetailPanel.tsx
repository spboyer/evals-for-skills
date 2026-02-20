import { useEffect, useState } from "react";
import { X, CheckCircle2, XCircle, Clock, ChevronDown, ChevronUp } from "lucide-react";
import type { ToolSpan } from "../types/trajectory";

// ---------------------------------------------------------------------------
// Expandable JSON section (Input / Output)
// ---------------------------------------------------------------------------

function ExpandableSection({
  label,
  data,
  defaultOpen = true,
}: {
  label: string;
  data: unknown;
  defaultOpen?: boolean;
}) {
  const [open, setOpen] = useState(defaultOpen);
  if (data === undefined || data === null) return null;

  const text =
    typeof data === "string" ? data : JSON.stringify(data, null, 2);

  return (
    <div className="rounded border border-zinc-700/60 overflow-hidden">
      <button
        onClick={() => setOpen(!open)}
        className="flex w-full items-center justify-between px-3 py-2 bg-zinc-800/80 text-xs font-medium text-zinc-300 hover:bg-zinc-700/60 transition-colors"
      >
        <span>{label}</span>
        {open ? (
          <ChevronUp className="h-3.5 w-3.5 text-zinc-500" />
        ) : (
          <ChevronDown className="h-3.5 w-3.5 text-zinc-500" />
        )}
      </button>
      {open && (
        <pre className="max-h-60 overflow-auto bg-zinc-900/60 p-3 text-xs text-zinc-300 leading-relaxed">
          <code>{text}</code>
        </pre>
      )}
    </div>
  );
}

// ---------------------------------------------------------------------------
// Status badge
// ---------------------------------------------------------------------------

const statusConfig: Record<
  ToolSpan["status"],
  { label: string; icon: React.ComponentType<{ className?: string }>; cls: string }
> = {
  pass: {
    label: "Passed",
    icon: CheckCircle2,
    cls: "bg-green-500/10 text-green-400 border-green-500/20",
  },
  fail: {
    label: "Failed",
    icon: XCircle,
    cls: "bg-red-500/10 text-red-400 border-red-500/20",
  },
  pending: {
    label: "Pending",
    icon: Clock,
    cls: "bg-blue-500/10 text-blue-400 border-blue-500/20",
  },
};

// ---------------------------------------------------------------------------
// Metadata row
// ---------------------------------------------------------------------------

function MetaRow({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex items-start justify-between gap-4 py-1.5 border-b border-zinc-700/30 last:border-b-0">
      <span className="text-xs text-zinc-500 shrink-0">{label}</span>
      <span className="text-xs text-zinc-300 text-right font-mono break-all">
        {value}
      </span>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Main component
// ---------------------------------------------------------------------------

interface DetailPanelProps {
  span: ToolSpan;
  onClose: () => void;
}

export default function DetailPanel({ span, onClose }: DetailPanelProps) {
  useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      if (e.key === "Escape") onClose();
    }
    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, [onClose]);

  const cfg = statusConfig[span.status];
  const StatusIcon = cfg.icon;

  return (
    <div className="w-80 shrink-0 border-l border-zinc-700 bg-zinc-800 flex flex-col overflow-hidden">
      {/* Header */}
      <div className="flex items-center justify-between gap-2 px-4 py-3 border-b border-zinc-700 bg-zinc-800/80 shrink-0">
        <div className="flex items-center gap-2 min-w-0">
          <StatusIcon className={`h-4 w-4 shrink-0 ${cfg.cls.split(" ")[1]}`} />
          <h4 className="text-sm font-medium text-zinc-100 truncate">
            {span.toolName}
          </h4>
        </div>
        <button
          onClick={onClose}
          className="shrink-0 rounded p-1 text-zinc-500 hover:bg-zinc-700 hover:text-zinc-200 transition-colors"
          aria-label="Close detail panel"
        >
          <X className="h-4 w-4" />
        </button>
      </div>

      {/* Scrollable body */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {/* Status badge */}
        <div>
          <span
            className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-medium ${cfg.cls}`}
          >
            <StatusIcon className="h-3 w-3" />
            {cfg.label}
          </span>
        </div>

        {/* Metadata */}
        <div className="rounded border border-zinc-700/60 px-3 py-1">
          <MetaRow label="Call ID" value={span.toolCallId} />
          <MetaRow
            label="Event range"
            value={`#${span.startIndex} → #${span.endIndex}`}
          />
          {(() => {
            const eventCount = span.duration + 1;
            return (
              <MetaRow
                label="Duration"
                value={`${eventCount} event${eventCount !== 1 ? "s" : ""}`}
              />
            );
          })()}
        </div>

        {/* Input / Output */}
        <ExpandableSection label="Input" data={span.arguments} defaultOpen />
        <ExpandableSection label="Output" data={span.toolResult} defaultOpen />
      </div>
    </div>
  );
}
