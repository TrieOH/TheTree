import { SERVICES, type ServiceId } from "@/features/catalog/services";

// Shared by the Worker (which writes and reads D1) and the pages (which only
// read GET /api/issues). Everything here mirrors Linear: one team (TRI), one
// project per service.

/** Linear's workflow state types. */
export type StateType = "triage" | "backlog" | "unstarted" | "started" | "completed" | "canceled";

/** 0 none, 1 urgent, 2 high, 3 medium, 4 low. */
export type Priority = 0 | 1 | 2 | 3 | 4;

export interface Issue {
  id: string;
  identifier: string;
  title: string;
  state: string;
  stateType: StateType;
  priority: Priority;
  estimate: number | null;
  assignee: string | null;
  labels: string[];
  project: string | null;
  service: ServiceId | null;
  url: string | null;
  cycleId: string | null;
  createdAt: string;
  completedAt: string | null;
}

export interface Cycle {
  id: string;
  number: number;
  startsAt: string;
  endsAt: string;
}

export interface Project {
  id: string;
  name: string;
  state: string | null;
  targetDate: string | null;
  /** Issues in the project, canceled ones left out. */
  total: number;
  done: number;
}

export interface IssuesPayload {
  /** The cycle running now, if any. */
  cycle: Cycle | null;
  /** Every open issue, plus those completed in the last 30 days. */
  issues: Issue[];
  /** Projects not yet completed or canceled. */
  projects: Project[];
  /** When the last catch-up sync finished; null before the first one. */
  syncedAt: string | null;
}

export const OPEN_STATES: readonly StateType[] = ["triage", "backlog", "unstarted", "started"];

export function isOpen(issue: Pick<Issue, "stateType">): boolean {
  return OPEN_STATES.includes(issue.stateType);
}

/** A project is a service when its name is one, case aside ("IdentityX"). */
export function serviceForProject(name: string | null | undefined): ServiceId | null {
  const id = name?.trim().toLowerCase();
  return SERVICES.find((service) => service.id === id)?.id ?? null;
}

export const PRIORITY_LABEL: Record<Priority, string> = {
  0: "No priority",
  1: "Urgent",
  2: "High",
  3: "Medium",
  4: "Low",
};
