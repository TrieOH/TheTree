import type { Cycle, Issue } from "./model";
import { isOpen } from "./model";

const DAY = 86_400_000;

export function inCycle(issue: Issue, cycle: Cycle | null): boolean {
  return cycle !== null && issue.cycleId === cycle.id;
}

/** Completed inside the cycle's window, or in the last 14 days without one. */
export function doneThisCycle(issues: readonly Issue[], cycle: Cycle | null, now = new Date()): Issue[] {
  const from = cycle ? Date.parse(cycle.startsAt) : now.getTime() - 14 * DAY;
  return issues.filter((issue) => issue.completedAt && Date.parse(issue.completedAt) >= from);
}

/** Median days from created to completed, over the completed issues given. */
export function medianDaysToClose(issues: readonly Issue[]): number | null {
  const days = issues
    .filter((issue) => issue.completedAt)
    .map((issue) => (Date.parse(issue.completedAt!) - Date.parse(issue.createdAt)) / DAY)
    .sort((a, b) => a - b);
  if (days.length === 0) return null;
  const mid = Math.floor(days.length / 2);
  return days.length % 2 ? days[mid] : (days[mid - 1] + days[mid]) / 2;
}

export function cycleProgress(issues: readonly Issue[], cycle: Cycle | null): { done: number; total: number } {
  const mine = issues.filter((issue) => inCycle(issue, cycle));
  return { done: mine.filter((issue) => !isOpen(issue)).length, total: mine.length };
}

const GROUP_ORDER = ["started", "unstarted", "triage", "backlog", "completed", "canceled"];

/** Issues grouped by their workflow state, in-progress first. */
export function groupByState(issues: readonly Issue[]): { state: string; type: string; issues: Issue[] }[] {
  const groups = new Map<string, { state: string; type: string; issues: Issue[] }>();
  for (const issue of issues) {
    const group = groups.get(issue.state) ?? { state: issue.state, type: issue.stateType, issues: [] };
    group.issues.push(issue);
    groups.set(issue.state, group);
  }
  return [...groups.values()].sort((a, b) => GROUP_ORDER.indexOf(a.type) - GROUP_ORDER.indexOf(b.type));
}

export function formatDay(iso: string): string {
  return new Date(iso).toLocaleDateString("en", { month: "short", day: "numeric" });
}
