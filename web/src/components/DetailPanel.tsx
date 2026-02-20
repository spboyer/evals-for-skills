import { useEffect, useState } from "react";
import {
  X,
  CheckCircle2,
  XCircle,
  Clock,
  ChevronDown,
  ChevronRight,
  Hash,
  AlignLeft,
  Activity,
} from "lucide-react";
import type { ToolSpan } from "../types/trajectory";

function CollapsibleSection({
  label,
  icon: Icon,
  children,
  defaultOpen = false,
}: {
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  children: React.ReactNode;
  defaultOpen?: boolean;
}) {
  const [open, setOpen] = useState(defaultOpen);
  return (
    <div className="border border-zinc-700 rounded-md overflow-hidden">
      <button
        onClick={() => setOpen(!open)}
        className="flex w-full items-center gap-2 px-3 py-2 bg-zinc-800/60 hover:bg-zinc-700/60 transition-colors text-left"
      >
        {open ? (
          <ChevronDown className="h-3.5 w-3.5 text-zinc-400 shrink-0" />
        ) : (
          <ChevronRight className="h-3.5 w-3.5 text-zinc-400 shrink-0" />
        )}
        <Icon className="h-3.5 w-3.5 text-zinc-400 shrink-0" />
        <span className="text-xs font-medium text-zinc-300">{label}</span>
      </button>
      {open && (
        <div className="bg-zinc-900 p-2.5">
          {children}
        </div>
      )}
    </div>
  );
}

const statusBadge: Record<ToolSpan["status"], string> = {
  pass: "bg-green-500/15 text-green-400 border border-green-500/30",
  fail: "bg-red-500/15 text-red-400 border border-red-500/30",
  pending: "bg-yellow-500/15 text-yellow-400 border border-yellow-500/30",
};

const statusIcon: Record<
  ToolSpan["status"],
  React.ComponentType<{ className?: string }>
> = {
  pass: CheckCircle2,
  fail: XCircle,
  pending: Clock,
};

const statusLabel: Record<ToolSpan["status"], string> = {
  pass: "Passed",
  fail: "Failed",
  pending: "In progress",
};

interface AttributeRowProps {
  label: string;
  value: React.ReactNode;
}
function AttributeRow({ label, value }: AttributeRowProps) {
  return (
    <div className="flex items-start justify-between gap-3 py-1.5 border-b border-zinc-700/50 last:border-b-0">
      <span className="text-xs text-zinc-500 shrink-0">{label}</span>
      <span className="text-xs text-zinc-300 text-right font-mono break-all">
        {value}
      </span>
    </div>
  );
}

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

  const Icon = statusIcon[span.status];

  const argsText =
    span.arguments === undefined
      ? null
      : typeof span.arguments === "string"
        ? span.arguments
        : JSON.stringify(span.arguments, null, 2);

  const resultText =
    span.toolResult === undefined
      ? null
      : typeof span.toolResult === "string"
        ? span.toolResult
        : JSON.stringify(span.toolResult, null, 2);

  return (
    <div className="w-80 shrink-0 border-l border-zinc-700 bg-zinc-900 flex flex-col overflow-hidden">
      {/* Header */}
      <div className="flex items-center justify-between gap-2 px-3 py-2.5 border-b border-zinc-700 bg-zinc-800">
        <div className="flex items-center gap-2 min-w-0">
          <Icon className={`h-4 w-4 shrink-0 ${
            span.status === "pass"
              ? "text-green-400"
              : span.status === "fail"
                ? "text-red-400"
                : "text-yellow-400"
          }`} />
          <h4 className="text-sm font-semibold text-zinc-100 truncate">
            {span.toolName}
          </h4>
        </div>
        <button
          onClick={onClose}
          className="shrink-0 p-1 rounded hover:bg-zinc-700 text-zinc-400 hover:text-zinc-200 transition-colors"
          aria-label="Close detail panel"
        >
          <X className="h-4 w-4" />
        </button>
      </div>

      {/* Body */}
      <div className="flex-1 overflow-y-auto p-3 space-y-4">
        {/* Status badge */}
        <div>
          <span
            className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ${statusBadge[span.status]}`}
          >
            <Icon className="h-3 w-3" />
            {statusLabel[span.status]}
          </span>
        </div>

        {/* Attributes */}
        <CollapsibleSection label="Attributes" icon={Activity} defaultOpen>
          <AttributeRow label="Duration" value={`${span.duration} events`} />
          <AttributeRow
            label="Event range"
            value={`${span.startIndex} → ${span.endIndex}`}
          />
          <AttributeRow
            label="Call ID"
            value={
              <span className="truncate max-w-[160px] block">
                {span.toolCallId}
              </span>
            }
          />
        </CollapsibleSection>

        {/* Arguments */}
        {argsText !== null && (
          <CollapsibleSection label="Arguments" icon={Hash} defaultOpen>
            <pre className="text-xs text-zinc-300 overflow-auto max-h-48 whitespace-pre-wrap break-words">
              <code>{argsText}</code>
            </pre>
          </CollapsibleSection>
        )}

        {/* Result */}
        {resultText !== null && (
          <CollapsibleSection label="Result" icon={AlignLeft} defaultOpen={span.status === "fail"}>
            <pre className="text-xs text-zinc-300 overflow-auto max-h-48 whitespace-pre-wrap break-words">
              <code>{resultText}</code>
            </pre>
          </CollapsibleSection>
        )}
      </div>
    </div>
  );
}
