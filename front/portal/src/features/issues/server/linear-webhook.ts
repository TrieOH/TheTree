import {
  deleteCycle,
  deleteIssue,
  deleteProject,
  upsertCycle,
  upsertIssue,
  upsertProject,
} from "./store";

// POST /webhooks/linear. Verified the way @linear/sdk does it: the
// `linear-signature` header is the hex HMAC-SHA256 of the raw body, and the
// signed body's `webhookTimestamp` must be within a minute. Redeliveries
// carry the same signature, so it doubles as the delivery id for dedupe.

export const LINEAR_WEBHOOK_PATH = "/webhooks/linear";
const MAX_SKEW_MS = 60_000;

interface Named {
  name: string;
}

interface IssueData {
  id: string;
  identifier: string;
  title: string;
  priority: number;
  estimate: number | null;
  url: string;
  cycleId: string | null;
  projectId: string | null;
  createdAt: string;
  updatedAt: string;
  completedAt: string | null;
  archivedAt: string | null;
  trashed: boolean | null;
  state: Named & { type: string };
  assignee: Named | null;
  labels: Named[];
}

interface CycleData {
  id: string;
  number: number;
  startsAt: string;
  endsAt: string;
  completedAt: string | null;
  archivedAt: string | null;
}

interface ProjectData {
  id: string;
  name: string;
  targetDate: string | null;
  updatedAt: string;
  archivedAt: string | null;
  trashed: boolean | null;
  status: { type: string } | null;
}

type Payload =
  | { type: "Issue"; action: string; data: IssueData; webhookTimestamp: number }
  | { type: "Cycle"; action: string; data: CycleData; webhookTimestamp: number }
  | { type: "Project"; action: string; data: ProjectData; webhookTimestamp: number }
  | { type: string; action: string; data: unknown; webhookTimestamp: number };

const encoder = new TextEncoder();

async function hmacHex(secret: string, body: string): Promise<string> {
  const key = await crypto.subtle.importKey(
    "raw",
    encoder.encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign"],
  );
  const mac = await crypto.subtle.sign("HMAC", key, encoder.encode(body));
  return [...new Uint8Array(mac)].map((byte) => byte.toString(16).padStart(2, "0")).join("");
}

function timingSafeEqual(a: string, b: string): boolean {
  if (a.length !== b.length) return false;
  let diff = 0;
  for (let i = 0; i < a.length; i++) diff |= a.charCodeAt(i) ^ b.charCodeAt(i);
  return diff === 0;
}

/** The statements that apply one event, or none for events we don't keep. */
export function statementsFor(db: D1Database, payload: Payload): D1PreparedStatement[] {
  const removed = payload.action === "remove";

  if (payload.type === "Issue") {
    const issue = payload.data as IssueData;
    if (removed || issue.archivedAt || issue.trashed) return [deleteIssue(db, issue.id)];
    return [
      upsertIssue(db, {
        id: issue.id,
        identifier: issue.identifier,
        title: issue.title,
        state: issue.state.name,
        stateType: issue.state.type,
        priority: issue.priority,
        estimate: issue.estimate,
        assignee: issue.assignee?.name ?? null,
        labels: (issue.labels ?? []).map((label) => label.name),
        cycleId: issue.cycleId,
        projectId: issue.projectId,
        url: issue.url,
        createdAt: issue.createdAt,
        completedAt: issue.completedAt,
        updatedAt: issue.updatedAt,
      }),
    ];
  }

  if (payload.type === "Cycle") {
    const cycle = payload.data as CycleData;
    if (removed || cycle.archivedAt) return [deleteCycle(db, cycle.id)];
    return [
      upsertCycle(db, {
        id: cycle.id,
        number: cycle.number,
        startsAt: cycle.startsAt,
        endsAt: cycle.endsAt,
        completedAt: cycle.completedAt,
      }),
    ];
  }

  if (payload.type === "Project") {
    const project = payload.data as ProjectData;
    if (removed || project.archivedAt || project.trashed) return [deleteProject(db, project.id)];
    return [
      upsertProject(db, {
        id: project.id,
        name: project.name,
        state: project.status?.type ?? null,
        targetDate: project.targetDate,
        updatedAt: project.updatedAt,
      }),
    ];
  }

  return [];
}

export async function handleLinearWebhook(request: Request, env: Env): Promise<Response> {
  const signature = request.headers.get("linear-signature");
  if (!signature) return new Response("missing signature", { status: 400 });

  const body = await request.text();
  if (!timingSafeEqual(await hmacHex(env.LINEAR_WEBHOOK_SECRET, body), signature)) {
    return new Response("invalid signature", { status: 401 });
  }

  let payload: Payload;
  try {
    payload = JSON.parse(body) as Payload;
  } catch {
    return new Response("invalid body", { status: 400 });
  }
  if (
    typeof payload.webhookTimestamp !== "number" ||
    Math.abs(Date.now() - payload.webhookTimestamp) > MAX_SKEW_MS
  ) {
    return new Response("stale delivery", { status: 401 });
  }

  const seen = await env.DB.prepare(
    "SELECT 1 FROM webhook_deliveries WHERE source = 'linear' AND delivery_id = ?",
  )
    .bind(signature)
    .first();
  if (seen) return new Response(null, { status: 204 });

  // One transaction: the delivery is only recorded if the change applied,
  // so a failed write is retried by Linear instead of being skipped.
  await env.DB.batch([
    env.DB.prepare(
      "INSERT OR IGNORE INTO webhook_deliveries (source, delivery_id, received_at) VALUES ('linear', ?, ?)",
    ).bind(signature, new Date().toISOString()),
    ...statementsFor(env.DB, payload),
  ]);

  return new Response(null, { status: 204 });
}
