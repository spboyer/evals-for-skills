import { Card } from "@/components/ui/Card";
import { useSummary } from "@/hooks/useSummary";
import { formatPercent, formatNumber, formatCost, formatDuration } from "@/lib/format";
import { Activity, CheckCircle, Coins, Timer } from "lucide-react";

export function KPICards() {
  const { data, isLoading, isError, error } = useSummary();

  if (isError) {
    return (
      <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-red-400">
        <p>Failed to load data: {error?.message || 'Unknown error'}</p>
      </div>
    );
  }

  if (isLoading || !data) {
    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <Card key={i}>
            <div className="h-12 animate-pulse rounded bg-gray-200 dark:bg-gray-800" />
          </Card>
        ))}
      </div>
    );
  }

  const kpis = [
    { label: "Total Runs", value: formatNumber(data.totalRuns), icon: Activity },
    { label: "Pass Rate", value: formatPercent(data.passRate), icon: CheckCircle },
    { label: "Avg Cost", value: formatCost(data.avgCost), icon: Coins },
    { label: "Avg Duration", value: formatDuration(data.avgDuration), icon: Timer },
  ];

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {kpis.map((kpi) => (
        <Card key={kpi.label}>
          <div className="flex items-center gap-3">
            <kpi.icon className="h-5 w-5 text-waza-500" />
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">{kpi.label}</p>
              <p className="text-xl font-semibold">{kpi.value}</p>
            </div>
          </div>
        </Card>
      ))}
    </div>
  );
}
