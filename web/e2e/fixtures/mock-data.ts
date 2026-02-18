import type { RunSummary, Run, RunDetail } from '../../src/types';

export const mockSummary: RunSummary = {
  totalRuns: 42,
  totalTasks: 168,
  passRate: 0.857,
  avgTokens: 12450,
  avgCost: 0.0342,
  avgDuration: 4523,
};

export const mockRuns: Run[] = [
  {
    id: 'run-001',
    spec: 'code-explainer',
    model: 'gpt-4o',
    outcome: 'failed',
    passCount: 3,
    taskCount: 4,
    tokens: 15200,
    cost: 0.0456,
    duration: 5230,
    timestamp: '2025-01-15T10:30:00Z',
  },
  {
    id: 'run-002',
    spec: 'bug-fixer',
    model: 'claude-sonnet-4-20250514',
    outcome: 'failed',
    passCount: 2,
    taskCount: 4,
    tokens: 11800,
    cost: 0.0298,
    duration: 3890,
    timestamp: '2025-01-15T09:15:00Z',
  },
  {
    id: 'run-003',
    spec: 'test-writer',
    model: 'gpt-4o-mini',
    outcome: 'passed',
    passCount: 3,
    taskCount: 3,
    tokens: 8700,
    cost: 0.0187,
    duration: 2450,
    timestamp: '2025-01-14T16:45:00Z',
  },
  {
    id: 'run-004',
    spec: 'code-reviewer',
    model: 'gpt-4o',
    outcome: 'error',
    passCount: 0,
    taskCount: 3,
    tokens: 4200,
    cost: 0.0126,
    duration: 1200,
    timestamp: '2025-01-14T14:20:00Z',
  },
];

export const mockRunDetail: RunDetail = {
  id: 'run-001',
  spec: 'code-explainer',
  model: 'gpt-4o',
  outcome: 'failed',
  passCount: 3,
  taskCount: 4,
  tokens: 15200,
  cost: 0.0456,
  duration: 5230,
  timestamp: '2025-01-15T10:30:00Z',
  tasks: [
    {
      name: 'explain-function',
      outcome: 'passed',
      score: 0.95,
      duration: 1200,
      graderResults: [
        { name: 'accuracy', type: 'llm', passed: true, score: 0.92, message: 'Explanation accurately describes the function behavior.' },
        { name: 'completeness', type: 'llm', passed: true, score: 0.98, message: 'All key aspects covered.' },
      ],
    },
    {
      name: 'explain-class',
      outcome: 'passed',
      score: 0.88,
      duration: 1500,
      graderResults: [
        { name: 'accuracy', type: 'llm', passed: true, score: 0.85, message: 'Class explanation is accurate.' },
      ],
    },
    {
      name: 'explain-module',
      outcome: 'failed',
      score: 0.42,
      duration: 1100,
      graderResults: [
        { name: 'accuracy', type: 'llm', passed: false, score: 0.4, message: 'Module explanation missed key dependencies.' },
      ],
    },
    {
      name: 'explain-api',
      outcome: 'passed',
      score: 0.91,
      duration: 1430,
      graderResults: [
        { name: 'accuracy', type: 'llm', passed: true, score: 0.91, message: 'API explanation is thorough and correct.' },
      ],
    },
  ],
};

export const mockHealth = { status: 'ok', version: '0.5.0' };

export const mockEmptyRuns: Run[] = [];

export const mockEmptySummary: RunSummary = {
  totalRuns: 0,
  totalTasks: 0,
  passRate: 0,
  avgTokens: 0,
  avgCost: 0,
  avgDuration: 0,
};
