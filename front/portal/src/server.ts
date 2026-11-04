// One Worker for the whole portal: the SPA (ASSETS), its read API over D1,
// the webhooks that write to D1 and, later, the crons that pull rollups.
// Nothing here may call the four APIs or Postgres: they scale to zero and the
// portal must not wake them.

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

  return json({ error: "not found" }, { status: 404 });
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const path = new URL(request.url).pathname;

    if (path === "/api" || path.startsWith("/api/")) {
      return handleApi(request, env, path);
    }

    if (path.startsWith("/webhooks/")) {
      // Linear, GitHub and the deploy hook arrive in phases 2 and 3.
      return json({ error: "not found" }, { status: 404 });
    }

    return env.ASSETS.fetch(request);
  },
} satisfies ExportedHandler<Env>;
