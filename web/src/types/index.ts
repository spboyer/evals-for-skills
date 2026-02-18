export interface RunSummary {
  totalRuns: number;
  totalTasks: number;
  passRate: number;
  avgTokens: number;
  avgCost: number;
  avgDuration: number;
}

export interface Run {
  id: string;
  spec: string;
  model: string;
  outcome: "passed" | "failed" | "error";
  passCount: number;
  taskCount: number;
  tokens: number;
  cost: number;
  duration: number;
  timestamp: string;
}

export interface Task {
  name: string;
  outcome: "passed" | "failed" | "error";
  score: number;
  duration: number;
  graderResults: GraderResult[];
}

export interface GraderResult {
  name: string;
  type: string;
  passed: boolean;
  score: number;
  message: string;
}

export interface RunDetail extends Run {
  tasks: Task[];
}
