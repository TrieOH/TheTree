// Writes Linear entities into D1. The webhook and the catch-up sync both go
// through here, so both agree on the shape. An upsert only lands when it is
// at least as new as the stored row, so a late webhook can't undo a newer
// sync and vice versa.

export interface IssueRecord {
  id: string;
  identifier: string;
  title: string;
  state: string;
  stateType: string;
  priority: number;
  estimate: number | null;
  assignee: string | null;
  labels: string[];
  cycleId: string | null;
  projectId: string | null;
  url: string | null;
  createdAt: string;
  completedAt: string | null;
  updatedAt: string;
}

export interface CycleRecord {
  id: string;
  number: number;
  startsAt: string;
  endsAt: string;
  completedAt: string | null;
}

export interface ProjectRecord {
  id: string;
  name: string;
  state: string | null;
  targetDate: string | null;
  updatedAt: string;
}

export function upsertIssue(db: D1Database, issue: IssueRecord): D1PreparedStatement {
  return db
    .prepare(
      `INSERT INTO issues (id, identifier, title, state, state_type, priority, estimate, assignee,
         labels, cycle_id, project_id, url, created_at, completed_at, updated_at)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
       ON CONFLICT (id) DO UPDATE SET
         identifier = excluded.identifier, title = excluded.title, state = excluded.state,
         state_type = excluded.state_type, priority = excluded.priority, estimate = excluded.estimate,
         assignee = excluded.assignee, labels = excluded.labels, cycle_id = excluded.cycle_id,
         project_id = excluded.project_id, url = excluded.url, completed_at = excluded.completed_at,
         updated_at = excluded.updated_at
       WHERE excluded.updated_at >= issues.updated_at`,
    )
    .bind(
      issue.id,
      issue.identifier,
      issue.title,
      issue.state,
      issue.stateType,
      issue.priority,
      issue.estimate,
      issue.assignee,
      JSON.stringify(issue.labels),
      issue.cycleId,
      issue.projectId,
      issue.url,
      issue.createdAt,
      issue.completedAt,
      issue.updatedAt,
    );
}

export function deleteIssue(db: D1Database, id: string): D1PreparedStatement {
  return db.prepare("DELETE FROM issues WHERE id = ?").bind(id);
}

export function upsertCycle(db: D1Database, cycle: CycleRecord): D1PreparedStatement {
  return db
    .prepare(
      `INSERT INTO cycles (id, number, starts_at, ends_at, completed_at) VALUES (?, ?, ?, ?, ?)
       ON CONFLICT (id) DO UPDATE SET number = excluded.number, starts_at = excluded.starts_at,
         ends_at = excluded.ends_at, completed_at = excluded.completed_at`,
    )
    .bind(cycle.id, cycle.number, cycle.startsAt, cycle.endsAt, cycle.completedAt);
}

export function deleteCycle(db: D1Database, id: string): D1PreparedStatement {
  return db.prepare("DELETE FROM cycles WHERE id = ?").bind(id);
}

export function upsertProject(db: D1Database, project: ProjectRecord): D1PreparedStatement {
  return db
    .prepare(
      `INSERT INTO projects (id, name, state, target_date, updated_at) VALUES (?, ?, ?, ?, ?)
       ON CONFLICT (id) DO UPDATE SET name = excluded.name, state = excluded.state,
         target_date = excluded.target_date, updated_at = excluded.updated_at
       WHERE excluded.updated_at >= projects.updated_at`,
    )
    .bind(project.id, project.name, project.state, project.targetDate, project.updatedAt);
}

export function deleteProject(db: D1Database, id: string): D1PreparedStatement {
  return db.prepare("DELETE FROM projects WHERE id = ?").bind(id);
}
