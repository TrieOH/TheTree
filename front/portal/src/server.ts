// One Worker for the whole portal: the SPA (ASSETS), its read API over D1,
// the webhooks that write to D1 and, later, the crons that pull rollups.
// Nothing here may call the four APIs or Postgres: they scale to zero and the
// portal must not wake them.

import { DEPLOY_HOOK_PATH, handleDeployHook } from "./features/delivery/server/deploy-hook";
import { GITHUB_WEBHOOK_PATH, handleGitHubWebhook } from "./features/delivery/server/github-webhook";
import { readDelivery } from "./features/delivery/server/read";
import { LINEAR_WEBHOOK_PATH, handleLinearWebhook } from "./features/issues/server/linear-webhook";
import { syncLinear } from "./features/issues/server/linear-sync";
import { readIssues } from "./features/issues/server/read";
import { readMetrics } from "./features/metrics/server/read";
import { syncMetrics } from "./features/metrics/server/sync";

function json(body: unknown, init?: ResponseInit): Response {
  return Response.json(body, {
    ...init,
    headers: { "Cache-Control": "no-store", ...init?.headers },
  });
}

async function handleApi(request: Request, env: Env, path: string): Promise<Response> {
  if (request.method !== "GET") {
    return json({ error: "method not allowed" }, { status: 405, headers: { Allow: "GET" } });
  }

  if (path === "/api/health") {
    // Proves the D1 binding works without touching anything outside Cloudflare.
    const row = await env.DB.prepare("SELECT 1 AS ok").first<{ ok: number }>();
    return json({ ok: row?.ok === 1 });
  }

  if (path === "/api/issues") {
    return json(await readIssues(env.DB));
  }

  if (path === "/api/delivery") {
    return json(await readDelivery(env.DB));
  }

  if (path === "/api/metrics") {
    return json(await readMetrics(env.DB));
  }

  return json({ error: "not found" }, { status: 404 });
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const path = new URL(request.url).pathname;

    if (path === "/api" || path.startsWith("/api/")) {
      return handleApi(request, env, path);
    }

    if (request.method === "POST") {
      if (path === LINEAR_WEBHOOK_PATH) return handleLinearWebhook(request, env);
      if (path === GITHUB_WEBHOOK_PATH) return handleGitHubWebhook(request, env);
      if (path === DEPLOY_HOOK_PATH) return handleDeployHook(request, env);
    }

    if (path.startsWith("/webhooks/")) {
      return json({ error: "not found" }, { status: 404 });
    }

    return env.ASSETS.fetch(request);
  },

  async scheduled(controller: ScheduledController, env: Env, ctx: ExecutionContext): Promise<void> {
    if (controller.cron === "7 * * * *") {
      ctx.waitUntil(syncLinear(env).then((result) => console.log("linear sync", result)));
    } else {
      ctx.waitUntil(syncMetrics(env).then((result) => console.log("metrics sync", result)));
    }
  },
} satisfies ExportedHandler<Env>;
