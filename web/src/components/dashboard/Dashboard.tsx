import { Card } from "@/components/ui/Card";
import { KPICards } from "./KPICards";
import { RunsTable } from "./RunsTable";

export function Dashboard() {
  return (
    <div className="mx-auto max-w-7xl space-y-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      <KPICards />
      <Card title="Evaluation Runs">
        <RunsTable />
      </Card>
    </div>
  );
}
