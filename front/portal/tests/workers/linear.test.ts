import { env } from "cloudflare:workers";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import type { IssuesPayload } from "../../src/features/issues/model";
import { syncLinear } from "../../src/features/issues/server/linear-sync";
import worker from "../../src/server";

const encoder = new TextEncoder();

async function sign(body: string, secret = "whsec_test"): Promise<string> {
  const key = await crypto.subtle.importKey("raw", encoder.encode(secret), { name: "HMAC", hash: "SHA-256" }, false, [
    "sign",
  ]);
  const mac = await crypto.subtle.sign("HMAC", key, encoder.encode(body));
  return [...new Uint8Array(mac)].map((byte) => byte.toString(16).padStart(2, "0")).join("");
}

async function deliver(payload: object, options: { signature?: string } = {}): Promise<Response> {
  const body = JSON.stringify({ webhookTimestamp: Date.now(), ...payload });
  return worker.fetch(
    new Request("http://localhost/webhooks/linear", {
      method: "POST",
      headers: { "linear-signature": options.signature ?? (await sign(body)) },
      body,
    }),
    env,
  );
}

async function readApi(): Promise<IssuesPayload> {
  const response = await worker.fetch(new Request("http://localhost/api/issues"), env);
  expect(response.status).toBe(200);
  return response.json();
}

function issueData(overrides: Record<string, unknown> = {}) {
  return {
    id: "issue-1",
    identifier: "TRI-1",
    title: "Retry MercadoPago webhook on 5xx",
    priority: 2,
    estimate: 3,
    url: "https://linear.app/trieoh/issue/TRI-1",
    cycleId: null,
    projectId: "project-payssage",
    createdAt: "2026-11-01T10:00:00.000Z",
    updatedAt: "2026-11-02T10:00:00.000Z",
    completedAt: null,
    archivedAt: null,
    trashed: null,
    state: { name: "Todo", type: "unstarted" },
    assignee: { name: "Sophia Hoffmann" },
    labels: [{ name: "Bug" }],
    ...overrides,
  };
}

const payssage = {
  id: "project-payssage",
  name: "Payssage",
  targetDate: null,
  updatedAt: "2026-11-01T00:00:00.000Z",
  archivedAt: null,
  trashed: null,
  status: { type: "started" },
};

async function count(table: string): Promise<number> {
  const row = await env.DB.prepare(`SELECT COUNT(*) AS n FROM ${table}`).first<{ n: number }>();
  return row!.n;
}

beforeEach(async () => {
  await env.DB.batch(
    ["issues", "cycles", "projects", "webhook_deliveries", "sync_state"].map((table) =>
      env.DB.prepare(`DELETE FROM ${table}`),
    ),
  );
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe("POST /webhooks/linear", () => {
  it("stores an issue and serves it with its service from the project", async () => {
    expect((await deliver({ type: "Project", action: "create", data: payssage })).status).toBe(204);
    expect((await deliver({ type: "Issue", action: "create", data: issueData() })).status).toBe(204);

    const { issues, projects } = await readApi();
    expect(issues).toEqual([
      expect.objectContaining({
        identifier: "TRI-1",
        state: "Todo",
        stateType: "unstarted",
        assignee: "Sophia Hoffmann",
        labels: ["Bug"],
        project: "Payssage",
        service: "payssage",
      }),
    ]);
    expect(projects).toEqual([expect.objectContaining({ name: "Payssage", total: 1, done: 0 })]);
  });

  it("rejects a bad signature and writes nothing", async () => {
    const response = await deliver(
      { type: "Issue", action: "create", data: issueData() },
      { signature: await sign("something else") },
    );

    expect(response.status).toBe(401);
    expect(await count("issues")).toBe(0);
  });

  it("rejects a delivery whose signed timestamp is over a minute old", async () => {
    const response = await deliver({
      type: "Issue",
      action: "create",
      data: issueData(),
      webhookTimestamp: Date.now() - 5 * 60_000,
    });

    expect(response.status).toBe(401);
    expect(await count("issues")).toBe(0);
  });

  it("applies a redelivery once", async () => {
    const body = JSON.stringify({ type: "Issue", action: "create", data: issueData(), webhookTimestamp: Date.now() });
    const signature = await sign(body);
    const send = () =>
      worker.fetch(
        new Request("http://localhost/webhooks/linear", {
          method: "POST",
          headers: { "linear-signature": signature },
          body,
        }),
        env,
      );

    expect((await send()).status).toBe(204);
    expect((await send()).status).toBe(204);
    expect(await count("webhook_deliveries")).toBe(1);
    expect(await count("issues")).toBe(1);
  });

  it("ignores an update older than what is stored", async () => {
    await deliver({ type: "Issue", action: "update", data: issueData({ title: "New", updatedAt: "2026-11-03T00:00:00.000Z" }) });
    await deliver({ type: "Issue", action: "update", data: issueData({ title: "Old", updatedAt: "2026-11-02T00:00:00.000Z" }) });

    const row = await env.DB.prepare("SELECT title FROM issues").first<{ title: string }>();
    expect(row?.title).toBe("New");
  });

  it("drops removed and archived issues", async () => {
    await deliver({ type: "Issue", action: "create", data: issueData() });
    await deliver({ type: "Issue", action: "create", data: issueData({ id: "issue-2", identifier: "TRI-2" }) });

    await deliver({ type: "Issue", action: "remove", data: issueData() });
    await deliver({
      type: "Issue",
      action: "update",
      data: issueData({ id: "issue-2", identifier: "TRI-2", archivedAt: "2026-11-03T00:00:00.000Z" }),
    });

    expect(await count("issues")).toBe(0);
  });

  it("accepts events it doesn't keep", async () => {
    const response = await deliver({ type: "Comment", action: "create", data: { id: "c1" } });

    expect(response.status).toBe(204);
  });
});

describe("GET /api/issues", () => {
  it("picks the running cycle and leaves out issues completed over 30 days ago", async () => {
    const now = Date.now();
    const iso = (offsetDays: number) => new Date(now + offsetDays * 86_400_000).toISOString();
    await deliver({ type: "Cycle", action: "create", data: { id: "c-old", number: 1, startsAt: iso(-30), endsAt: iso(-16), completedAt: iso(-16), archivedAt: null } });
    await deliver({ type: "Cycle", action: "create", data: { id: "c-now", number: 2, startsAt: iso(-2), endsAt: iso(12), completedAt: null, archivedAt: null } });
    await deliver({
      type: "Issue",
      action: "create",
      data: issueData({ id: "recent", identifier: "TRI-2", state: { name: "Done", type: "completed" }, completedAt: iso(-1), updatedAt: iso(-1) }),
    });
    await deliver({
      type: "Issue",
      action: "create",
      data: issueData({ id: "ancient", identifier: "TRI-3", state: { name: "Done", type: "completed" }, completedAt: iso(-40), updatedAt: iso(-40) }),
    });

    const payload = await readApi();
    expect(payload.cycle).toMatchObject({ id: "c-now", number: 2 });
    expect(payload.issues.map((issue) => issue.identifier)).toEqual(["TRI-2"]);
  });
});

describe("syncLinear", () => {
  function graphql(pages: Record<string, unknown[][]>) {
    const calls: { query: string; variables: { since: string; after: string | null } }[] = [];
    vi.spyOn(globalThis, "fetch").mockImplementation(async (_input, init) => {
      const body = JSON.parse(String(init?.body));
      calls.push(body);
      const kind = /items: (\w+)/.exec(body.query)![1];
      const list = pages[kind] ?? [[]];
      const index = body.variables.after ? Number(body.variables.after) : 0;
      return Response.json({
        data: {
          items: {
            nodes: list[index] ?? [],
            pageInfo: { hasNextPage: index + 1 < list.length, endCursor: String(index + 1) },
          },
        },
      });
    });
    return calls;
  }

  const issueNode = (id: string, extra: Record<string, unknown> = {}) => ({
    ...issueData({ id, identifier: id.toUpperCase() }),
    cycle: null,
    project: { id: "project-payssage" },
    labels: { nodes: [{ name: "Bug" }] },
    ...extra,
  });

  it("backfills everything on its first run, across pages", async () => {
    const calls = graphql({
      issues: [[issueNode("tri-1")], [issueNode("tri-2"), issueNode("tri-3", { trashed: true })]],
      projects: [[payssage]],
      cycles: [[]],
    });

    const result = await syncLinear(env, new Date("2026-11-04T16:00:00.000Z"));

    expect(result).toEqual({ issues: 3, cycles: 0, projects: 1 });
    expect(calls.every((call) => call.variables.since === "1970-01-01T00:00:00.000Z")).toBe(true);
    expect(calls.every((call) => call.query.includes('key: { eq: "TRI" }'))).toBe(true);
    expect(await count("issues")).toBe(2);
    const state = await env.DB.prepare("SELECT synced_at FROM sync_state").first<{ synced_at: string }>();
    expect(state?.synced_at).toBe("2026-11-04T16:00:00.000Z");
  });

  it("resumes from the last run, a few minutes early", async () => {
    await env.DB.prepare("INSERT INTO sync_state (source, synced_at) VALUES ('linear', '2026-11-04T16:00:00.000Z')").run();
    const calls = graphql({});

    await syncLinear(env);

    expect(calls[0].variables.since).toBe("2026-11-04T15:55:00.000Z");
  });

  it("keeps its place when Linear answers with errors", async () => {
    vi.spyOn(globalThis, "fetch").mockImplementation(async () =>
      Response.json({ errors: [{ message: "Authentication required" }] }),
    );

    await expect(syncLinear(env)).rejects.toThrow("Authentication required");
    expect(await count("sync_state")).toBe(0);
  });
});
