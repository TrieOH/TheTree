import { describe, expect, it } from "vitest";
import { needsAttention } from "../../src/features/delivery/attention";
import type { CiRun, CiStepStat, DeliveryPayload } from "../../src/features/delivery/model";
import { jobTimes, pipelineNumbers, stepLabel } from "../../src/features/delivery/pipeline";
import { groupReleases, releaseOutcome, runSpan } from "../../src/features/delivery/releases";

const now = new Date("2026-11-21T23:00:00Z");

const run = (tag: string, startedAt: string, minutes: number, extra: Partial<CiRun> = {}): CiRun => ({
  runId: Date.parse(startedAt),
  attempt: 1,
  component: tag.split("/v")[0],
  tag,
  type: tag.includes("-ui/") ? "front" : "service",
  lane: "full",
  status: "success",
  stage: "publish",
  startedAt,
  completedAt: new Date(Date.parse(startedAt) + minutes * 60_000).toISOString(),
  ...extra,
});

const step = (extra: Partial<CiStepStat>): CiStepStat => ({
  type: "service",
  job: "checks",
  step: "Build",
  runs: 20,
  medianMs: 4000,
  p90Ms: 6000,
  lastMs: 4000,
  recentMs: [4000],
  ...extra,
});

const base: DeliveryPayload = {
  production: [],
  deploys: [],
  pulls: [],
  runs: [],
  tags: [],
  ci: { runs: [], steps: [], metrics: [] },
};

describe("needsAttention", () => {
  it("stays quiet when nothing is off", () => {
    expect(needsAttention(base, now)).toEqual([]);
  });

  it("forgets a failed release once a newer one of the same component ships", () => {
    const ci = {
      ...base.ci,
      runs: [
        run("identityx/v0.49.1", "2026-11-21T20:22:02Z", 2),
        run("identityx/v0.49.0", "2026-11-21T20:18:01Z", 1, { status: "failure", stage: "checks" }),
      ],
    };
    expect(needsAttention({ ...base, ci }, now)).toEqual([]);
  });

  it("flags a component stuck on a failed release, a silent tag, a slow step and a PR left waiting, worst first", () => {
    const items = needsAttention(
      {
        ...base,
        tags: [
          { name: "payssage/v0.21.0", component: "payssage", version: "v0.21.0", sha: "b", message: null, createdAt: "2026-11-21T20:30:00Z" },
          { name: "informd/v0.15.0", component: "informd", version: "v0.15.0", sha: "c", message: null, createdAt: "2026-11-21T22:50:00Z" },
        ],
        pulls: [
          { number: 132, title: "bump checkout", author: "dependabot[bot]", state: "open", draft: false, reviewState: null, openedAt: "2026-11-19T14:55:19Z", mergedAt: null },
          { number: 141, title: "wip", author: "MintzyG", state: "open", draft: true, reviewState: null, openedAt: "2026-11-01T00:00:00Z", mergedAt: null },
        ],
        ci: {
          runs: [run("identityx/v0.49.0", "2026-11-21T20:18:01Z", 1, { status: "failure", stage: "checks" })],
          steps: [step({ job: "publish", step: "Build & push image", medianMs: 22_000, lastMs: 61_000 })],
          metrics: [],
        },
      },
      now,
    );

    expect(items.map((i) => [i.tone, i.tab, i.text])).toEqual([
      ["destructive", "releases", "identityx v0.49.0 was blocked at checks 3 h ago, and nothing newer has shipped."],
      ["warning", "releases", "payssage/v0.21.0 was tagged 3 h ago but no release run reported back."],
      ["warning", "pipelines", "Build & push image took 61 s in the last service release, against a usual 22 s."],
      ["warning", "pulls", "#132 “bump checkout” has waited 2 days for a review."],
    ]);
  });

  it("flags critical and high findings in the newest scan", () => {
    const finding = (severity: string, last: number) => ({
      type: "front",
      job: "checks",
      source: "scan/trivy",
      name: "findings",
      labels: { path: "front/portal", severity },
      runs: 3,
      last,
      median: 0,
    });
    const items = needsAttention({ ...base, ci: { ...base.ci, metrics: [finding("CRITICAL", 1), finding("HIGH", 2)] } }, now);
    expect(items).toMatchObject([
      { tone: "destructive", text: "The last front release scan found 1 critical and 2 high severity vulnerabilities." },
    ]);
  });
});

describe("pipeline", () => {
  it("names steps without their commit pins", () => {
    expect(stepLabel("Run useblacksmith/checkout@2b9fe87d796feea153b455ec2f3125b471a7b0e7")).toBe("Checkout");
    expect(stepLabel("Post Run useblacksmith/checkout@2b9fe87d796feea153b455ec2f3125b471a7b0e7")).toBe("Post checkout");
    expect(stepLabel("Build & push image")).toBe("Build & push image");
  });

  it("lists the steps that take time per job and folds the rest", () => {
    const jobs = jobTimes([
      step({ job: "publish", step: "Deploy to Cloudflare", medianMs: 22_000 }),
      step({ step: "Waiting for a runner", medianMs: 7000 }),
      step({ step: "Complete job", medianMs: 0, p90Ms: 0 }),
      step({ step: "Set up runner", medianMs: 0, p90Ms: 1000 }),
      step({ step: "Integration tests", medianMs: 1000, p90Ms: 14_000 }),
    ]);
    expect(jobs.map((j) => [j.job, j.steps.map((s) => s.label), j.hidden])).toEqual([
      ["checks", ["Waiting for a runner", "Integration tests"], { count: 2, medianMs: 0 }],
      ["publish", ["Deploy to Cloudflare"], { count: 0, medianMs: 0 }],
    ]);
  });

  it("sums scans and tests and drops numbers that never left zero", () => {
    const metric = (source: string, name: string, labels: Record<string, string>, last: number, median = last) => ({
      type: "service",
      job: "checks",
      source,
      name,
      labels,
      runs: 6,
      last,
      median,
    });
    const numbers = pipelineNumbers([
      metric("scan/trivy", "findings", { path: "api/identityx", severity: "HIGH" }, 0),
      metric("scan/trivy", "findings", { path: "lib/go", severity: "LOW" }, 2),
      metric("go/test", "tests_pass", { mode: "short" }, 132),
      metric("go/test", "tests_skip", { mode: "short" }, 85, 6),
      metric("go/test", "slow_test_seconds", { mode: "integration", package: "organizations", test: "TestCreate" }, 5.19, 6.74),
      metric("go/test", "slow_test_seconds", { mode: "short", package: "app", test: "TestFast" }, 0.01),
      metric("go/lint", "issues", {}, 0),
      metric("setup/go", "build_cache_bytes", {}, 759_795_712, 546_611_200),
    ]);
    expect(numbers.findings).toEqual([
      { severity: "CRITICAL", count: 0 },
      { severity: "HIGH", count: 0 },
      { severity: "MEDIUM", count: 0 },
      { severity: "LOW", count: 2 },
    ]);
    expect(numbers.scannedPaths).toBe(2);
    expect(numbers.tests).toEqual([{ mode: "short", pass: 132, fail: 0, skip: 85 }]);
    expect(numbers.slowTests.map((t) => t.test)).toEqual(["TestCreate"]);
    expect(numbers.other.map((m) => m.name)).toEqual(["build_cache_bytes"]);
  });
});

describe("groupReleases", () => {
  it("puts the runs of one merge together under its pull request's title, on one clock", () => {
    const sha = "f03a958";
    const runs = [
      run("identityx-ui/v0.10.0", "2026-11-21T20:34:22Z", 1),
      run("informd/v0.14.0", "2026-11-21T20:25:07Z", 2),
      run("identityx/v0.49.1", "2026-11-21T20:22:02Z", 2),
      run("identityx/v0.49.0", "2026-11-21T20:18:01Z", 1, { status: "failure", stage: "checks" }),
    ];
    const tag = (name: string, commit: string, message: string) => ({
      name,
      component: name.split("/v")[0],
      version: name.split("/")[1],
      sha: commit,
      message,
      createdAt: "2026-11-21T20:00:00Z",
    });
    const releases = groupReleases({
      ci: { runs, steps: [], metrics: [] },
      tags: [
        tag("identityx/v0.49.1", sha, "Merge pull request #138 from TrieOH/lib-go/fix/validator-first-use"),
        tag("informd/v0.14.0", sha, "Merge pull request #138 from TrieOH/lib-go/fix/validator-first-use"),
        tag("identityx/v0.49.0", "3e46a03", "fix(lib-go): rate limit by /64"),
      ],
      deploys: [{ component: "identityx-ui", tag: "identityx-ui/v0.10.0", sha, actor: "MintzyG", status: "success", at: "", runId: null }],
      pulls: [{ number: 138, title: "Build the shared validator on first use", author: "MintzyG", state: "merged", draft: false, reviewState: null, openedAt: "", mergedAt: "" }],
    });

    expect(releases.map((r) => [r.title, r.pr, r.runs.map((x) => x.component), releaseOutcome(r)])).toEqual([
      ["Build the shared validator on first use", 138, ["identityx", "informd", "identityx-ui"], "shipped"],
      ["fix(lib-go): rate limit by /64", null, ["identityx"], "failed"],
    ]);
    const [first] = releases;
    expect(runSpan(first, first.runs[0])).toEqual({ start: 0, width: 2 / 13.333333333333334 });
  });
});
