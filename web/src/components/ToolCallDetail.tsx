import { useState } from "react";
import {
  CheckCircle2,
  XCircle,
  ChevronDown,
  ChevronRight,
} from "lucide-react";

interface ToolCallDetailProps {
  toolName: string;
  toolCallId?: string;
  args?: unknown;
  result?: unknown;
  success?: boolean;
}

function JsonViewer({ label, data }: { label: string; data: unknown }) {
  const [open, setOpen] = useState(false);
  if (data === undefined || data === null) return null;

  const text =
    typeof data === "string" ? data : JSON.stringify(data, null, 2);

  return (
    <div>
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-1 text-xs text-zinc-500 hover:text-zinc-300"
      >
        {open ? (
          <ChevronDown className="h-3 w-3" />
        ) : (
          <ChevronRight className="h-3 w-3" />
        )}
        {label}
      </button>
      {open && (
        <pre className="mt-1 max-h-64 overflow-auto rounded bg-zinc-900 p-2.5 text-xs text-zinc-300">
          <code>{text}</code>
        </pre>
      )}
    </div>
  );
}

export default function ToolCallDetail({
  toolName,
  toolCallId,
  args,
  result,
  success,
}: ToolCallDetailProps) {
  return (
    <div className="space-y-2 rounded-lg border border-zinc-700 bg-zinc-800/60 p-3">
      <div className="flex items-center gap-2">
        <span className="text-sm font-medium text-zinc-200">{toolName}</span>
        {success !== undefined &&
          (success ? (
            <span className="inline-flex items-center gap-0.5 rounded-full bg-green-500/10 px-1.5 py-0.5 text-xs font-medium text-green-400">
              <CheckCircle2 className="h-3 w-3" /> pass
            </span>
          ) : (
            <span className="inline-flex items-center gap-0.5 rounded-full bg-red-500/10 px-1.5 py-0.5 text-xs font-medium text-red-400">
              <XCircle className="h-3 w-3" /> fail
            </span>
          ))}
        {toolCallId && (
          <span className="text-xs text-zinc-600 font-mono">{toolCallId}</span>
        )}
      </div>
      <JsonViewer label="Arguments" data={args} />
      <JsonViewer label="Result" data={result} />
    </div>
  );
}
