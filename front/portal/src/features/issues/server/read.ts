import type { Cycle, Issue, IssuesPayload, Priority, Project, StateType } from "../model";
import { serviceForProject } from "../model";

// GET /api/issues: one small payload that the Issues page, the service pages
// and Home all read. Open issues plus the last 30 days of completed ones is a
// few hundred rows at most.

interface IssueRow {
  id: string;
  identifier: string;
  title: string;
  state: string;
  state_type: StateType;
  priority: Priority;
  estimate: number | null;
  assignee: string | null;
  labels: string;
  project: string | null;
  url: string | null;
  cycle_id: string | null;
  created_at: string;
  completed_at: string | null;
}

interface ProjectRow {
  id: string;
  name: string;
  state: string | null;
  target_date: string | null;
  total: number;
  done: number;
}

export async function readIssues(db: D1Database, now = new Date()): Promise<IssuesPayload> {
  const nowIso = now.toISOString();
  const since = new Date(now.getTime() - 30 * 86_400_000).toISOString();

  const [issues, cycle, projects, sync] = await db.batch<unknown>([
    db.prepare(
      `SELECT i.id, i.identifier, i.title, i.state, i.state_type, i.priority, i.estimate, i.assignee,
              i.labels, p.name AS project, i.url, i.cycle_id, i.created_at, i.completed_at
       FROM issues i LEFT JOIN projects p ON p.id = i.project_id
       WHERE i.state_type IN ('triage', 'backlog', 'unstarted', 'started')
          OR (i.state_type = 'completed' AND i.completed_at >= ?)
       ORDER BY i.priority = 0, i.priority, i.updated_at DESC`,
    ).bind(since),
    db.prepare(
      `SELECT id, number, starts_at, ends_at FROM cycles
       WHERE starts_at <= ? AND ends_at > ? ORDER BY starts_at DESC LIMIT 1`,
    ).bind(nowIso, nowIso),
    db.prepare(
      `SELECT p.id, p.name, p.state, p.target_date,
              COUNT(i.id) AS total,
              COALESCE(SUM(i.state_type = 'completed'), 0) AS done
       FROM projects p LEFT JOIN issues i ON i.project_id = p.id AND i.state_type != 'canceled'
       WHERE p.state IS NULL OR p.state NOT IN ('completed', 'canceled')
       GROUP BY p.id ORDER BY p.target_date IS NULL, p.target_date, p.name`,
    ),
    db.prepare("SELECT synced_at FROM sync_state WHERE source = 'linear'"),
  ]);

  const cycleRow = (cycle.results as { id: string; number: number; starts_at: string; ends_at: string }[])[0];
  const currentCycle: Cycle | null = cycleRow
    ? { id: cycleRow.id, number: cycleRow.number, startsAt: cycleRow.starts_at, endsAt: cycleRow.ends_at }
    : null;

  return {
    cycle: currentCycle,
    issues: (issues.results as IssueRow[]).map(
      (row): Issue => ({
        id: row.id,
        identifier: row.identifier,
        title: row.title,
        state: row.state,
        stateType: row.state_type,
        priority: row.priority,
        estimate: row.estimate,
        assignee: row.assignee,
        labels: JSON.parse(row.labels) as string[],
        project: row.project,
        service: serviceForProject(row.project),
        url: row.url,
        cycleId: row.cycle_id,
        createdAt: row.created_at,
        completedAt: row.completed_at,
      }),
    ),
    projects: (projects.results as ProjectRow[]).map(
      (row): Project => ({
        id: row.id,
        name: row.name,
        state: row.state,
        targetDate: row.target_date,
        total: row.total,
        done: row.done,
      }),
    ),
    syncedAt: (sync.results as { synced_at: string }[])[0]?.synced_at ?? null,
  };
}
