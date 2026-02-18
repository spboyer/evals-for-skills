import type { RunSummary, Run, RunDetail } from "@/types";

const BASE_URL = window.location.origin;

async function fetchJSON<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`);
  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }
  return res.json() as Promise<T>;
}

export function fetchSummary(): Promise<RunSummary> {
  return fetchJSON<RunSummary>("/api/summary");
}

export function fetchRuns(): Promise<Run[]> {
  return fetchJSON<Run[]>("/api/runs");
}

export function fetchRunDetail(id: string): Promise<RunDetail> {
  return fetchJSON<RunDetail>(`/api/runs/${id}`);
}

export function fetchHealth(): Promise<{ status: string }> {
  return fetchJSON<{ status: string }>("/api/health");
}
