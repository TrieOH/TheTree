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
    // recent: the CI tab only lists workflows that ran in the last 30 days
    created_at: new Date(Date.now() - 3_600_000).toISOString(),
    run_started_at: new Date(Date.now() - 3_590_000).toISOString(),
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
    [
      "pull_requests",
      "workflow_runs",
      "tags",
      "deploys",
      "webhook_deliveries",
      "ci_runs",
      "ci_jobs",
      "ci_steps",
      "ci_metrics",
      "ci_substeps",
    ].map((table) =>
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

  it("groups runs by their workflow, not by GitHub's per-run titles", async () => {
    const dependabot = (id: number) => ({
      ...run({ id, name: `github_actions in /., /.github/actions/** - Update #${id}` }),
      workflow: { name: "Dependabot Updates" },
    });
    await github("workflow_run", dependabot(1684151908));
    await github("workflow_run", dependabot(1684151989));

    expect((await readApi()).runs.map((r) => r.workflow)).toEqual(["Dependabot Updates", "Dependabot Updates"]);
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

  it("keeps a release blocked at checks out of production and the DORA list", async () => {
    await deploy(report({ stage: "publish" }));
    await deploy(report({ tag: "portal-ui/v0.3.0", status: "failure", stage: "checks" }));

    const { production, deploys } = await readApi();
    expect(production).toEqual([expect.objectContaining({ tag: "portal-ui/v0.2.0", status: "success" })]);
    expect(deploys).toHaveLength(1);
  });

  it("rejects an unknown stage", async () => {
    expect((await deploy(report({ stage: "lint" }))).status).toBe(400);
  });
});

async function ci(body: object, secret = "deploy_test") {
  return worker.fetch(
    new Request("http://localhost/webhooks/ci", {
      method: "POST",
      headers: { Authorization: `Bearer ${secret}`, "Content-Type": "application/json" },
      body: JSON.stringify(body),
    }),
    env,
  );
}

const at = (seconds: number) => new Date(Date.UTC(2026, 10, 5, 19, 0, seconds)).toISOString();

const step = (number: number, name: string, from: number, to: number, status = "success") => ({
  number,
  name,
  status,
  started_at: at(from),
  completed_at: at(to),
});

const ciReport = (overrides: Record<string, unknown> = {}) => ({
  run_id: 42091066478,
  attempt: 1,
  component: "identityx",
  tag: "identityx/v0.48.0",
  sha: "859c2fe6",
  type: "service",
  lane: "full",
  status: "success",
  stage: "publish",
  jobs: [
    {
      name: "service / checks",
      status: "success",
      queued_at: at(0),
      started_at: at(4),
      completed_at: at(60),
      runner: "blacksmith-1",
      steps: [step(1, "Set up job", 4, 5), step(2, "Setup Go", 5, 11), step(3, "Lint", 11, 11, "skipped")],
      substeps: [
        { step: 2, action_id: "__self.__actions_setup-go", name: "Toolchain", depth: 1, completed_at: at(9), duration_ms: 3652 },
        { step: 2, action_id: "__self.__useblacksmith_stickydisk", name: "Mount module cache", depth: 1, completed_at: at(6), duration_ms: 449 },
        // post steps can reuse an id
        { step: 3, action_id: "__post.__useblacksmith_stickydisk", name: "Commit", depth: 1, completed_at: at(11), duration_ms: 155 },
        { step: 3, action_id: "__post.__useblacksmith_stickydisk", name: "Commit", depth: 1, completed_at: at(11), duration_ms: 118 },
      ],
    },
  ],
  // 20 metrics: more than fit in one INSERT under D1's 100-parameter limit
  metrics: Array.from({ length: 20 }, (_, i) => ({
    job: "service / checks",
    source: "go/test",
    name: "slow_test_seconds",
    value: i / 10,
    labels: { test: `Test${i}`, mode: "short" },
  })),
  ...overrides,
});

describe("POST /webhooks/ci", () => {
  it("needs the shared secret", async () => {
    expect((await ci(ciReport(), "wrong")).status).toBe(401);
  });

  it("rejects a malformed report", async () => {
    expect((await ci(ciReport({ stage: "deploy" }))).status).toBe(400);
    expect((await ci(ciReport({ metrics: [{ job: "x", source: "y", name: "z", value: "1", labels: {} }] }))).status).toBe(400);
  });

  it("serves step times, queue time and metrics per pipeline", async () => {
    expect((await ci(ciReport())).status).toBe(204);

    const { ci: payload } = await readApi();
    expect(payload.runs).toEqual([
      expect.objectContaining({ runId: 42091066478, stage: "publish", startedAt: at(4), completedAt: at(60) }),
    ]);
    expect(payload.steps.map((s) => [s.job, s.step, s.medianMs])).toEqual([
      ["checks", "Waiting for a runner", 4000],
      ["checks", "Set up job", 1000],
      ["checks", "Setup Go", 6000],
    ]);
    expect(payload.metrics).toHaveLength(20);
    expect(payload.metrics[0]).toEqual(
      expect.objectContaining({ job: "checks", source: "go/test", labels: { mode: "short", test: "Test0" }, last: 0 }),
    );
  });

  it("keeps the same metric from two jobs apart", async () => {
    const pnpm = (job: string, value: number) => ({ job, source: "setup/node", name: "pnpm_downloaded", value, labels: {} });
    await ci(ciReport({ metrics: [pnpm("service / checks", 1129), pnpm("service / publish", 3)] }));

    const { ci: payload } = await readApi();
    expect(payload.metrics.map((m) => [m.job, m.last])).toEqual([
      ["checks", 1129],
      ["publish", 3],
    ]);
  });

  it("replaces a resent attempt instead of doubling it", async () => {
    await ci(ciReport());
    await ci(ciReport({ status: "failure", stage: "checks" }));

    const { ci: payload } = await readApi();
    expect(payload.runs).toEqual([expect.objectContaining({ status: "failure", stage: "checks" })]);
    expect(payload.metrics).toHaveLength(20);
    expect(payload.steps.find((s) => s.step === "Setup Go")?.runs).toBe(1);
  });
});

async function api(path: string): Promise<Response> {
  return worker.fetch(new Request(`http://localhost${path}`), env);
}

describe("GET /api/ci", () => {
  it("describes itself", async () => {
    const response = await api("/api/ci");
    expect(response.headers.get("access-control-allow-origin")).toBe("*");
    expect(Object.keys((await response.json<{ endpoints: object }>()).endpoints)).toContain("/api/ci/substeps");
  });

  it("filters runs and gives their duration", async () => {
    await ci(ciReport());
    await ci(ciReport({ run_id: 2, component: "univents-ui", tag: "univents-ui/v0.16.1", type: "front" }));

    const runs = await (await api("/api/ci/runs?type=front")).json<{ tag: string; duration_ms: number }[]>();
    expect(runs).toEqual([expect.objectContaining({ tag: "univents-ui/v0.16.1", duration_ms: 56_000 })]);
  });

  it("serves steps, substeps and metrics as rows, by job and name", async () => {
    await ci(ciReport());

    const steps = await (await api("/api/ci/steps?job=checks&name=Setup%20Go")).json<{ duration_ms: number }[]>();
    expect(steps).toEqual([expect.objectContaining({ job: "service / checks", duration_ms: 6000 })]);

    const subs = await (await api("/api/ci/substeps?step=Setup%20Go&name=Toolchain")).json<object[]>();
    expect(subs).toEqual([expect.objectContaining({ step: "Setup Go", step_number: 2, duration_ms: 3652 })]);

    const metrics = await (await api("/api/ci/metrics?name=slow_test_seconds&label.test=Test3")).json<{ value: number }[]>();
    expect(metrics).toEqual([expect.objectContaining({ value: 0.3, labels: { mode: "short", test: "Test3" } })]);
  });

  it("answers in CSV", async () => {
    await ci(ciReport());

    const response = await api("/api/ci/steps?format=csv&name=Setup%20Go");
    expect(response.headers.get("content-type")).toContain("text/csv");
    const [header, row] = (await response.text()).trim().split("\n");
    expect(header.split(",")).toContain("duration_ms");
    expect(row).toContain("Setup Go");
  });

  it("serves one run with its steps nesting their substeps", async () => {
    await ci(ciReport());

    const run = await (await api("/api/ci/runs/42091066478")).json<{
      jobs: { queue_ms: number; steps: { name: string; substeps: { name: string }[] }[] }[];
    }>();
    expect(run.jobs[0].queue_ms).toBe(4000);
    expect(run.jobs[0].steps.find((s) => s.name === "Setup Go")?.substeps.map((s) => s.name)).toEqual([
      "Mount module cache",
      "Toolchain",
    ]);
    expect((await api("/api/ci/runs/1")).status).toBe(404);
  });
});
