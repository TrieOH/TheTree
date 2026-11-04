import {
  deleteCycle,
  deleteIssue,
  deleteProject,
  upsertCycle,
  upsertIssue,
  upsertProject,
} from "./store";

// Hourly catch-up from Linear's GraphQL API with the read-only key: whatever
// changed since the last run, archived and trashed items included so their
// rows go too. The first run starts from the epoch, so it is the backfill.
// The webhook keeps things current between runs; this covers deliveries it
// missed or dropped.

const LINEAR_API = "https://api.linear.app/graphql";
const TEAM = "TRI";
/** Re-read a few minutes before the last run, for writes that raced it. */
const OVERLAP_MS = 5 * 60_000;
const PAGE = 100;
const BATCH = 50;

interface PageInfo {
  hasNextPage: boolean;
  endCursor: string | null;
}

interface Connection<T> {
  nodes: T[];
  pageInfo: PageInfo;
}

interface IssueNode {
  id: string;
  identifier: string;
  title: string;
  priority: number;
  estimate: number | null;
  url: string;
  createdAt: string;
  updatedAt: string;
  completedAt: string | null;
  archivedAt: string | null;
  trashed: boolean | null;
  state: { name: string; type: string };
  assignee: { name: string } | null;
  cycle: { id: string } | null;
  project: { id: string } | null;
  labels: { nodes: { name: string }[] };
}

interface CycleNode {
  id: string;
  number: number;
  startsAt: string;
  endsAt: string;
  completedAt: string | null;
  archivedAt: string | null;
}

interface ProjectNode {
  id: string;
  name: string;
  targetDate: string | null;
  updatedAt: string;
  archivedAt: string | null;
  trashed: boolean | null;
  status: { type: string } | null;
}

const ISSUES = `query Issues($since: DateTimeOrDuration!, $after: String) {
  items: issues(first: ${PAGE}, after: $after, includeArchived: true,
    filter: { updatedAt: { gt: $since }, team: { key: { eq: "${TEAM}" } } }) {
    nodes {
      id identifier title priority estimate url createdAt updatedAt completedAt archivedAt trashed
      state { name type } assignee { name } cycle { id } project { id } labels { nodes { name } }
    }
    pageInfo { hasNextPage endCursor }
  }
}`;

const CYCLES = `query Cycles($since: DateTimeOrDuration!, $after: String) {
  items: cycles(first: ${PAGE}, after: $after, includeArchived: true,
    filter: { updatedAt: { gt: $since }, team: { key: { eq: "${TEAM}" } } }) {
    nodes { id number startsAt endsAt completedAt archivedAt }
    pageInfo { hasNextPage endCursor }
  }
}`;

const PROJECTS = `query Projects($since: DateTimeOrDuration!, $after: String) {
  items: projects(first: ${PAGE}, after: $after, includeArchived: true,
    filter: { updatedAt: { gt: $since }, accessibleTeams: { some: { key: { eq: "${TEAM}" } } } }) {
    nodes { id name targetDate updatedAt archivedAt trashed status { type } }
    pageInfo { hasNextPage endCursor }
  }
}`;

async function fetchAll<T>(apiKey: string, query: string, since: string): Promise<T[]> {
  const all: T[] = [];
  let after: string | null = null;
  do {
    const response: Response = await fetch(LINEAR_API, {
      method: "POST",
      headers: { "Content-Type": "application/json", Authorization: apiKey },
      body: JSON.stringify({ query, variables: { since, after } }),
    });
    const body = (await response.json()) as {
      data?: { items: Connection<T> };
      errors?: { message: string }[];
    };
    if (!response.ok || body.errors || !body.data) {
      throw new Error(
        `Linear sync failed (${response.status}): ${body.errors?.map((e) => e.message).join("; ") ?? "no data"}`,
      );
    }
    all.push(...body.data.items.nodes);
    after = body.data.items.pageInfo.hasNextPage ? body.data.items.pageInfo.endCursor : null;
  } while (after);
  return all;
}

function issueStatement(db: D1Database, issue: IssueNode): D1PreparedStatement {
  if (issue.archivedAt || issue.trashed) return deleteIssue(db, issue.id);
  return upsertIssue(db, {
    id: issue.id,
    identifier: issue.identifier,
    title: issue.title,
    state: issue.state.name,
    stateType: issue.state.type,
    priority: issue.priority,
    estimate: issue.estimate,
    assignee: issue.assignee?.name ?? null,
    labels: issue.labels.nodes.map((label) => label.name),
    cycleId: issue.cycle?.id ?? null,
    projectId: issue.project?.id ?? null,
    url: issue.url,
    createdAt: issue.createdAt,
    completedAt: issue.completedAt,
    updatedAt: issue.updatedAt,
  });
}

function cycleStatement(db: D1Database, cycle: CycleNode): D1PreparedStatement {
  if (cycle.archivedAt) return deleteCycle(db, cycle.id);
  return upsertCycle(db, {
    id: cycle.id,
    number: cycle.number,
    startsAt: cycle.startsAt,
    endsAt: cycle.endsAt,
    completedAt: cycle.completedAt,
  });
}

function projectStatement(db: D1Database, project: ProjectNode): D1PreparedStatement {
  if (project.archivedAt || project.trashed) return deleteProject(db, project.id);
  return upsertProject(db, {
    id: project.id,
    name: project.name,
    state: project.status?.type ?? null,
    targetDate: project.targetDate,
    updatedAt: project.updatedAt,
  });
}

export interface SyncResult {
  issues: number;
  cycles: number;
  projects: number;
}

export async function syncLinear(env: Env, now = new Date()): Promise<SyncResult> {
  const last = await env.DB.prepare("SELECT synced_at FROM sync_state WHERE source = 'linear'")
    .first<{ synced_at: string }>();
  const since = last
    ? new Date(Date.parse(last.synced_at) - OVERLAP_MS).toISOString()
    : new Date(0).toISOString();

  const [issues, cycles, projects] = await Promise.all([
    fetchAll<IssueNode>(env.LINEAR_API_KEY, ISSUES, since),
    fetchAll<CycleNode>(env.LINEAR_API_KEY, CYCLES, since),
    fetchAll<ProjectNode>(env.LINEAR_API_KEY, PROJECTS, since),
  ]);

  const statements = [
    ...projects.map((project) => projectStatement(env.DB, project)),
    ...cycles.map((cycle) => cycleStatement(env.DB, cycle)),
    ...issues.map((issue) => issueStatement(env.DB, issue)),
  ];
  for (let i = 0; i < statements.length; i += BATCH) {
    await env.DB.batch(statements.slice(i, i + BATCH));
  }

  // Only advance once everything above landed; a failed run retries the window.
  await env.DB.batch([
    env.DB.prepare(
      `INSERT INTO sync_state (source, synced_at) VALUES ('linear', ?)
       ON CONFLICT (source) DO UPDATE SET synced_at = excluded.synced_at`,
    ).bind(now.toISOString()),
    // Dedupe only needs to outlive Linear's retries.
    env.DB.prepare("DELETE FROM webhook_deliveries WHERE received_at < ?").bind(
      new Date(now.getTime() - 7 * 86_400_000).toISOString(),
    ),
  ]);

  return { issues: issues.length, cycles: cycles.length, projects: projects.length };
}
