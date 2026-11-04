import { hmacSha256Hex, timingSafeEqual } from "@/shared/lib/hmac";
import { parseTag } from "../model";

// POST /webhooks/github, from the repo webhook on TrieOH/TheTree. GitHub signs
// the raw body: `X-Hub-Signature-256: sha256=<hex HMAC-SHA256>`. Every
// delivery has its own `X-GitHub-Delivery` id, kept for dedupe (redeliveries
// reuse it).
//
// Events: pull_request, pull_request_review, workflow_run, and push (only
// tag pushes are kept; GitHub's `create` event carries no sha).

export const GITHUB_WEBHOOK_PATH = "/webhooks/github";

/** GitHub sends `…:33Z`; D1 compares dates as text, so store one format. */
const iso = <T extends string | null>(value: T): T => (value === null ? value : (new Date(value).toISOString() as T));

interface PullRequestEvent {
  pull_request: {
    number: number;
    title: string;
    user: { login: string };
    state: "open" | "closed";
    draft: boolean;
    created_at: string;
    updated_at: string;
    merged_at: string | null;
    merge_commit_sha: string | null;
  };
}

interface ReviewEvent {
  action: string;
  review: { state: string };
  pull_request: { number: number };
}

interface WorkflowRunEvent {
  workflow_run: {
    id: number;
    name: string;
    head_branch: string | null;
    head_sha: string;
    status: string;
    conclusion: string | null;
    run_attempt: number;
    created_at: string;
    run_started_at: string | null;
    updated_at: string;
  };
}

interface PushEvent {
  ref: string;
  after: string;
  deleted: boolean;
  head_commit: { message: string } | null;
}

/** The statements that apply one event, or none for events we don't keep. */
export function statementsFor(db: D1Database, event: string, body: unknown, now = new Date()): D1PreparedStatement[] {
  if (event === "pull_request") {
    const pr = (body as PullRequestEvent).pull_request;
    return [
      db
        .prepare(
          `INSERT INTO pull_requests (number, title, author, state, draft, opened_at, merged_at, merge_sha, updated_at)
           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
           ON CONFLICT (number) DO UPDATE SET title = excluded.title, author = excluded.author,
             state = excluded.state, draft = excluded.draft, merged_at = excluded.merged_at,
             merge_sha = excluded.merge_sha, updated_at = excluded.updated_at
           WHERE excluded.updated_at >= pull_requests.updated_at`,
        )
        .bind(
          pr.number,
          pr.title,
          pr.user.login,
          pr.merged_at ? "merged" : pr.state,
          pr.draft ? 1 : 0,
          iso(pr.created_at),
          iso(pr.merged_at),
          pr.merged_at ? pr.merge_commit_sha : null,
          iso(pr.updated_at),
        ),
    ];
  }

  if (event === "pull_request_review") {
    const { action, review, pull_request } = body as ReviewEvent;
    const state = action === "dismissed" ? "dismissed" : review.state.toLowerCase();
    // A plain comment doesn't replace an approval or a change request.
    return [
      db
        .prepare(
          `UPDATE pull_requests SET review_state = ?
           WHERE number = ? AND (? != 'commented' OR review_state IS NULL)`,
        )
        .bind(state, pull_request.number, state),
    ];
  }

  if (event === "workflow_run") {
    const run = (body as WorkflowRunEvent).workflow_run;
    const completed = run.status === "completed";
    return [
      db
        .prepare(
          `INSERT INTO workflow_runs (id, workflow, branch, sha, conclusion, attempt, queued_at, started_at, completed_at)
           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
           ON CONFLICT (id) DO UPDATE SET conclusion = excluded.conclusion, attempt = excluded.attempt,
             queued_at = excluded.queued_at, started_at = excluded.started_at, completed_at = excluded.completed_at
           WHERE excluded.attempt > workflow_runs.attempt
              OR (excluded.attempt = workflow_runs.attempt AND workflow_runs.completed_at IS NULL)`,
        )
        .bind(
          run.id,
          run.name,
          run.head_branch ?? "",
          run.head_sha,
          completed ? run.conclusion : null,
          run.run_attempt,
          iso(run.created_at),
          iso(run.run_started_at),
          completed ? iso(run.updated_at) : null,
        ),
    ];
  }

  if (event === "push") {
    const push = body as PushEvent;
    if (!push.ref.startsWith("refs/tags/")) return [];
    const name = push.ref.slice("refs/tags/".length);
    if (push.deleted) return [db.prepare("DELETE FROM tags WHERE name = ?").bind(name)];
    const parsed = parseTag(name);
    if (!parsed) return [];
    return [
      db
        .prepare(
          `INSERT INTO tags (name, component, version, sha, message, created_at) VALUES (?, ?, ?, ?, ?, ?)
           ON CONFLICT (name) DO UPDATE SET sha = excluded.sha, message = excluded.message`,
        )
        .bind(
          name,
          parsed.component,
          parsed.version,
          push.after,
          push.head_commit?.message.split("\n")[0] ?? null,
          now.toISOString(),
        ),
    ];
  }

  return [];
}

export async function handleGitHubWebhook(request: Request, env: Env): Promise<Response> {
  const signature = request.headers.get("x-hub-signature-256");
  const delivery = request.headers.get("x-github-delivery");
  const event = request.headers.get("x-github-event");
  if (!signature || !delivery || !event) return new Response("missing headers", { status: 400 });

  const body = await request.text();
  const expected = `sha256=${await hmacSha256Hex(env.GITHUB_WEBHOOK_SECRET, body)}`;
  if (!timingSafeEqual(expected, signature)) return new Response("invalid signature", { status: 401 });

  if (event === "ping") return new Response(null, { status: 204 });

  let payload: unknown;
  try {
    payload = JSON.parse(body);
  } catch {
    return new Response("invalid body", { status: 400 });
  }

  const seen = await env.DB.prepare(
    "SELECT 1 FROM webhook_deliveries WHERE source = 'github' AND delivery_id = ?",
  )
    .bind(delivery)
    .first();
  if (seen) return new Response(null, { status: 204 });

  await env.DB.batch([
    env.DB.prepare(
      "INSERT OR IGNORE INTO webhook_deliveries (source, delivery_id, received_at) VALUES ('github', ?, ?)",
    ).bind(delivery, new Date().toISOString()),
    ...statementsFor(env.DB, event, payload),
  ]);

  return new Response(null, { status: 204 });
}
