import { env } from "cloudflare:workers";
import { beforeEach, describe, expect, it } from "vitest";
import type { DeliveryPayload } from "../../src/features/delivery/model";
import { hmacSha256Hex } from "../../src/shared/lib/hmac";
import worker from "../../src/server";

let deliveryCounter = 0;

async function github(event: string, payload: object, options: { signature?: string; delivery?: string } = {}) {
  const body = JSON.stringify(payload);
  return worker.fetch(
    new Request("http://localhost/webhooks/github", {
      method: "POST",
      headers: {
        "x-github-event": event,
        "x-github-delivery": options.delivery ?? `delivery-${++deliveryCounter}`,
        "x-hub-signature-256": options.signature ?? `sha256=${await hmacSha256Hex("ghsec_test", body)}`,
      },
      body,
    }),
    env,
  );
}

async function deploy(report: object, secret = "deploy_test") {
  return worker.fetch(
    new Request("http://localhost/webhooks/deploy", {
      method: "POST",
      headers: { Authorization: `Bearer ${secret}`, "Content-Type": "application/json" },
      body: JSON.stringify(report),
    }),
    env,
  );
}

async function readApi(): Promise<DeliveryPayload> {
  const response = await worker.fetch(new Request("http://localhost/api/delivery"), env);
  expect(response.status).toBe(200);
  return response.json();
}

const pr = (overrides: Record<string, unknown> = {}) => ({
  pull_request: {
    number: 72,
    title: "fix(root): let the portal's first release create its Worker",
    user: { login: "MintzyG" },
    state: "open",
    draft: false,
    created_at: new Date(Date.now() - 3 * 3_600_000).toISOString(),
    updated_at: new Date(Date.now() - 3 * 3_600_000).toISOString(),
    merged_at: null,
    merge_commit_sha: null,
    ...overrides,
  },
});

const run = (overrides: Record<string, unknown> = {}) => ({
  workflow_run: {
    id: 41860946529,
    name: "Deploy",
    head_branch: "portal-ui/v0.1.1",
    head_sha: "8d487f9a",
    status: "completed",
    conclusion: "success",
    run_attempt: 1,
    created_at: "2026-11-04T15:44:00Z",
    run_started_at: "2026-11-04T15:44:10Z",
    updated_at: "2026-11-04T15:49:10Z",
    ...overrides,
  },
});

const report = (overrides: Record<string, unknown> = {}) => ({
  component: "portal-ui",
  tag: "portal-ui/v0.2.0",
  sha: "cc44557c",
  actor: "MintzyG",
  status: "success",
  run_id: 41863939045,
  ...overrides,
});

beforeEach(async () => {
  await env.DB.batch(
    ["pull_requests", "workflow_runs", "tags", "deploys", "webhook_deliveries"].map((table) =>
      env.DB.prepare(`DELETE FROM ${table}`),
    ),
  );
});

describe("POST /webhooks/github", () => {
  it("answers GitHub's ping", async () => {
    expect((await github("ping", { zen: "Keep it logically awesome." })).status).toBe(204);
  });

  it("rejects a bad signature", async () => {
    const response = await github("pull_request", pr(), { signature: "sha256=00" });

    expect(response.status).toBe(401);
    expect((await readApi()).pulls).toEqual([]);
  });

  it("tracks a pull request through review and merge", async () => {
    await github("pull_request", pr());
    await github("pull_request_review", { action: "submitted", review: { state: "APPROVED" }, pull_request: { number: 72 } });
    await github("pull_request_review", { action: "submitted", review: { state: "COMMENTED" }, pull_request: { number: 72 } });

    expect((await readApi()).pulls).toEqual([
      expect.objectContaining({ number: 72, state: "open", reviewState: "approved", author: "MintzyG" }),
    ]);

    const mergedAt = new Date().toISOString();
    await github("pull_request", pr({ state: "closed", merged_at: mergedAt, merge_commit_sha: "8d487f9a", updated_at: mergedAt }));

    expect((await readApi()).pulls).toEqual([expect.objectContaining({ number: 72, state: "merged", mergedAt })]);
  });

  it("keeps the latest attempt of a workflow run", async () => {
    await github("workflow_run", run({ status: "in_progress", conclusion: null, updated_at: "2026-11-04T15:45:00Z" }));
    await github("workflow_run", run({ conclusion: "failure" }));
    await github("workflow_run", run({ conclusion: "success", run_attempt: 2, updated_at: "2026-11-04T15:59:00Z" }));
    // A late event from the first attempt doesn't win.
    await github("workflow_run", run({ conclusion: "failure" }));

    expect((await readApi()).runs).toEqual([
      expect.objectContaining({ workflow: "Deploy", branch: "portal-ui/v0.1.1", conclusion: "success", attempt: 2 }),
    ]);
  });

  it("records release tags from pushes and forgets deleted ones", async () => {
    const push = (ref: string, deleted = false) => ({
      ref,
      after: "cc44557c",
      deleted,
      head_commit: { message: "Merge pull request #73\n\nbody" },
    });
    await github("push", push("refs/tags/portal-ui/v0.2.0"));
    await github("push", push("refs/tags/not-a-release"));
    await github("push", push("refs/heads/main"));

    expect((await readApi()).tags).toEqual([
      expect.objectContaining({ name: "portal-ui/v0.2.0", component: "portal-ui", version: "v0.2.0", message: "Merge pull request #73" }),
    ]);

    await github("push", push("refs/tags/portal-ui/v0.2.0", true));
    expect((await readApi()).tags).toEqual([]);
  });

  it("applies a redelivery once", async () => {
    await github("pull_request", pr(), { delivery: "same" });
    await github("pull_request_review", { action: "submitted", review: { state: "APPROVED" }, pull_request: { number: 72 } }, { delivery: "x" });
    await github("pull_request_review", { action: "dismissed", review: { state: "DISMISSED" }, pull_request: { number: 72 } }, { delivery: "x" });

    expect((await readApi()).pulls[0].reviewState).toBe("approved");
  });
});

describe("POST /webhooks/deploy", () => {
  it("needs the shared secret", async () => {
    expect((await deploy(report(), "wrong")).status).toBe(401);
  });

  it("rejects a malformed report", async () => {
    expect((await deploy(report({ status: "great" }))).status).toBe(400);
    expect((await deploy(report({ run_id: "41863939045" }))).status).toBe(400);
  });

  it("serves the latest deploy of each component as production", async () => {
    await deploy(report({ tag: "portal-ui/v0.1.1" }));
    await deploy(report());
    await deploy(report({ component: "identityx", tag: "identityx/v0.42.0", status: "failure" }));

    const { production, deploys } = await readApi();
    expect(production).toEqual([
      expect.objectContaining({ component: "identityx", status: "failure" }),
      expect.objectContaining({ component: "portal-ui", tag: "portal-ui/v0.2.0", runId: 41863939045 }),
    ]);
    expect(deploys).toHaveLength(3);
  });
});
