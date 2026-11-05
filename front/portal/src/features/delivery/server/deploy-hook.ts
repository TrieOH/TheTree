import { timingSafeEqual } from "@/shared/lib/hmac";

// POST /webhooks/deploy, sent by the last step of deploy.yml and
// deploy-pages.yml with `Authorization: Bearer $DEPLOY_HOOK_SECRET`:
// { component, tag, sha, actor, status, run_id }. A deploy is judged by its
// run status alone; neither the workflow nor the portal calls the service.
// Older reports also sent health_ok, which is ignored.

export const DEPLOY_HOOK_PATH = "/webhooks/deploy";

const STATUSES = new Set(["success", "failure", "cancelled"]);

interface DeployReport {
  component: string;
  tag: string;
  sha: string;
  actor: string;
  status: string;
  run_id: number | null;
}

function isReport(value: unknown): value is DeployReport {
  const report = value as Partial<DeployReport> | null;
  return (
    typeof report?.component === "string" &&
    typeof report.tag === "string" &&
    typeof report.sha === "string" &&
    typeof report.actor === "string" &&
    typeof report.status === "string" &&
    STATUSES.has(report.status) &&
    (report.run_id === null || Number.isInteger(report.run_id))
  );
}

export async function handleDeployHook(request: Request, env: Env): Promise<Response> {
  const auth = request.headers.get("authorization") ?? "";
  if (!timingSafeEqual(auth, `Bearer ${env.DEPLOY_HOOK_SECRET}`)) {
    return new Response("unauthorized", { status: 401 });
  }

  let report: unknown;
  try {
    report = await request.json();
  } catch {
    return new Response("invalid body", { status: 400 });
  }
  if (!isReport(report)) return new Response("invalid report", { status: 400 });

  await env.DB.prepare(
    `INSERT INTO deploys (component, tag, sha, actor, status, at, run_id)
     VALUES (?, ?, ?, ?, ?, ?, ?)`,
  )
    .bind(
      report.component,
      report.tag,
      report.sha,
      report.actor,
      report.status,
      new Date().toISOString(),
      report.run_id,
    )
    .run();

  return new Response(null, { status: 204 });
}
