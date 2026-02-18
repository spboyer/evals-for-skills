import { type LucideIcon } from "lucide-react";
import {
  Activity,
  ListChecks,
  CheckCircle2,
  Zap,
  DollarSign,
  Clock,
} from "lucide-react";
import { Skeleton } from "@/components/ui/Skeleton";
import { useSummary } from "@/hooks/useSummary";
import {
  formatPercent,
  formatNumber,
  formatCost,
  formatDuration,
} from "@/lib/format";

interface KPIItem {
  label: string;
  value: string;
  icon: LucideIcon;
  color?: string;
}

function passRateColor(rate: number): string {
  const pct = rate * 100;
  if (pct >= 70) return "text-emerald-400";
  if (pct >= 50) return "text-yellow-400";
  return "text-red-400";
}

function KPICard({ label, value, icon: Icon, color }: KPIItem) {
  return (
    <div className="rounded-lg border border-white/10 bg-[#111118] p-4">
      <Icon className="mb-2 h-5 w-5 text-cyan-400" />
      <p className={`text-2xl font-bold tracking-tight ${color ?? "text-white"}`}>
        {value}
      </p>
      <p className="mt-1 text-xs text-gray-400">{label}</p>
    </div>
  );
}

function KPICardSkeleton() {
  return (
    <div className="rounded-lg border border-white/10 bg-[#111118] p-4">
      <Skeleton className="mb-2 h-5 w-5 rounded-full" />
      <Skeleton className="h-7 w-20" />
      <Skeleton className="mt-1 h-3 w-16" />
    </div>
  );
}

export function KPICards() {
  const { data, isLoading, error } = useSummary();

  if (error) {
    return (
      <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-400">
        Failed to load summary data.
      </div>
    );
  }

  if (isLoading || !data) {
    return (
      <div className="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-6">
        {Array.from({ length: 6 }).map((_, i) => (
          <KPICardSkeleton key={i} />
        ))}
      </div>
    );
  }

  const kpis: KPIItem[] = [
    { label: "Total Runs", value: formatNumber(data.totalRuns), icon: Activity },
    { label: "Tasks", value: formatNumber(data.totalTasks), icon: ListChecks },
    {
      label: "Pass Rate",
      value: formatPercent(data.passRate),
      icon: CheckCircle2,
      color: passRateColor(data.passRate),
    },
    { label: "Avg Tokens", value: formatNumber(data.avgTokens), icon: Zap },
    { label: "Avg Cost", value: formatCost(data.avgCost), icon: DollarSign },
    { label: "Avg Duration", value: formatDuration(data.avgDuration), icon: Clock },
  ];

  return (
    <div className="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-6">
      {kpis.map((kpi) => (
        <KPICard key={kpi.label} {...kpi} />
      ))}
    </div>
  );
}
