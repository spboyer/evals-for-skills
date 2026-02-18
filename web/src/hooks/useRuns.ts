import { useQuery } from "@tanstack/react-query";
import { fetchRuns } from "@/api/client";

export function useRuns() {
  return useQuery({
    queryKey: ["runs"],
    queryFn: fetchRuns,
  });
}
