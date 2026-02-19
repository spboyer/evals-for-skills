import type { Page } from "@playwright/test";

/**
 * Mock API responses that match the Go server's JSON shape.
 * Keeping this in one place so every spec file shares the same data.
 */

export const SUMMARY = {
  totalRuns: 12,
  totalTasks: 48,
  passRate: 85,
  avgTokens: 15230,
  avgCost: 1.47,
  avgDuration: 42,
};

export const RUNS = [
  {
    id: "run-001",
    spec: "code-explainer",
    model: "gpt-4o",
    outcome: "pass",
    passCount: 4,
    taskCount: 4,
    tokens: 12400,
    cost: 1.24,
    duration: 38,
    timestamp: new Date(Date.now() - 3600_000).toISOString(), // 1 hour ago
  },
  {
    id: "run-002",
    spec: "skill-checker",
    model: "claude-sonnet-4",
    outcome: "fail",
    passCount: 2,
    taskCount: 5,
    tokens: 18100,
    cost: 1.81,
    duration: 55,
    timestamp: new Date(Date.now() - 7200_000).toISOString(), // 2 hours ago
  },
  {
    id: "run-003",
    spec: "doc-writer",
    model: "gpt-4o",
    outcome: "pass",
    passCount: 3,
    taskCount: 3,
    tokens: 9800,
    cost: 0.98,
    duration: 27,
    timestamp: new Date(Date.now() - 86400_000).toISOString(), // 1 day ago
  },
];

export const RUN_DETAIL = {
  ...RUNS[0],
  tasks: [
    {
      name: "explain-fibonacci",
      outcome: "pass",
      score: 1.0,
      duration: 12,
      graderResults: [
        {
          name: "output-exists",
          type: "code",
          passed: true,
          score: 1.0,
          message: "Output file found",
        },
        {
          name: "mentions-recursion",
          type: "regex",
          passed: true,
          score: 1.0,
          message: 'Matched pattern "recursion|recursive"',
        },
      ],
    },
    {
      name: "explain-quicksort",
      outcome: "pass",
      score: 1.0,
      duration: 10,
      graderResults: [
        {
          name: "output-exists",
          type: "code",
          passed: true,
          score: 1.0,
          message: "Output file found",
        },
      ],
    },
    {
      name: "explain-binary-search",
      outcome: "fail",
      score: 0.0,
      duration: 8,
      graderResults: [
        {
          name: "output-exists",
          type: "code",
          passed: true,
          score: 1.0,
          message: "Output file found",
        },
        {
          name: "mentions-divide-conquer",
          type: "regex",
          passed: false,
          score: 0.0,
          message: "Pattern not matched",
        },
      ],
    },
    {
      name: "explain-merge-sort",
      outcome: "pass",
      score: 1.0,
      duration: 8,
      graderResults: [
        {
          name: "output-exists",
          type: "code",
          passed: true,
          score: 1.0,
          message: "Output file found",
        },
      ],
    },
  ],
};

export const HEALTH = { status: "ok", version: "0.1.0" };
