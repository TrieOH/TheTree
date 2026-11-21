import type { CiMetricStat, CiStepStat } from "./model";

// Turns the raw per-step timings and recorded numbers of one pipeline into
// what a person reads: the steps that take time, and the numbers that say
// something. GitHub reports ~30 steps per pipeline, most of them runner
// bookkeeping that takes no time; the actions record ~100 numbers, most of
// them zero.

/** A step shows when it typically takes this long, or sometimes much longer. */
const SHOWN_MEDIAN_MS = 2_000;
const SHOWN_P90_MS = 5_000;

const JOB_ORDER = ["parse", "checks", "publish"];

export interface StepRow extends CiStepStat {
  /** The step's name without commit pins: `Run useblacksmith/checkout@2b9f…` → `Checkout`. */
  label: string;
}

export interface JobTimes {
  job: string;
  /** Sum of the steps' medians: roughly how long the job takes. */
  medianMs: number;
  steps: StepRow[];
  /** Steps too quick to list, folded into one line. */
  hidden: { count: number; medianMs: number };
}

export function stepLabel(step: string): string {
  const name = step.replace(/@[0-9a-f]{7,40}\b/, "");
  if (/^(Post )?Run \S*checkout$/.test(name)) return name.startsWith("Post ") ? "Post checkout" : "Checkout";
  return name.replace(/^Run (\S+\/\S+)$/, "$1");
}

/** One pipeline's steps, by job in run order, with the quick ones folded away. */
export function jobTimes(steps: readonly CiStepStat[]): JobTimes[] {
  const byJob = new Map<string, CiStepStat[]>();
  for (const step of steps) byJob.set(step.job, [...(byJob.get(step.job) ?? []), step]);
  const rank = (job: string) => (JOB_ORDER.includes(job) ? JOB_ORDER.indexOf(job) : JOB_ORDER.length);
  return [...byJob.entries()]
    .sort(([a], [b]) => rank(a) - rank(b))
    .map(([job, all]) => {
      const shown = all.filter((s) => s.medianMs >= SHOWN_MEDIAN_MS || s.p90Ms >= SHOWN_P90_MS);
      const quick = all.filter((s) => !shown.includes(s));
      return {
        job,
        medianMs: all.reduce((sum, s) => sum + s.medianMs, 0),
        steps: shown.map((s) => ({ ...s, label: stepLabel(s.step) })),
        hidden: { count: quick.length, medianMs: quick.reduce((sum, s) => sum + s.medianMs, 0) },
      };
    });
}

/** A step whose last run took at least twice its median and 20 s more. */
export function isSlowLastRun(step: CiStepStat): boolean {
  return step.runs >= 5 && step.lastMs >= 2 * step.medianMs && step.lastMs - step.medianMs >= 20_000;
}

const SEVERITIES = ["CRITICAL", "HIGH", "MEDIUM", "LOW"] as const;

export interface PipelineNumbers {
  /** Trivy findings in the last run, summed over every scanned path; null when nothing was scanned. */
  findings: { severity: (typeof SEVERITIES)[number]; count: number }[] | null;
  scannedPaths: number;
  /** Test counts of the last run, per mode (short, integration). */
  tests: { mode: string; pass: number; fail: number; skip: number }[];
  /** The five slowest tests that take a second or more. */
  slowTests: { test: string; pkg: string; mode: string; medianS: number; lastS: number }[];
  /** Everything else that isn't zero in every run. */
  other: CiMetricStat[];
}

export function pipelineNumbers(metrics: readonly CiMetricStat[]): PipelineNumbers {
  const findings = metrics.filter((m) => m.source === "scan/trivy" && m.name === "findings");
  const tests = metrics.filter((m) => m.source === "go/test" && /^tests_(pass|fail|skip)$/.test(m.name));
  const slow = metrics.filter((m) => m.name === "slow_test_seconds");
  const summarized = new Set([...findings, ...tests, ...slow]);

  const modes = [...new Set(tests.map((m) => m.labels.mode ?? ""))].sort();
  const count = (mode: string, name: string) =>
    tests.find((m) => (m.labels.mode ?? "") === mode && m.name === name)?.last ?? 0;

  return {
    findings: findings.length
      ? SEVERITIES.map((severity) => ({
          severity,
          count: findings.filter((m) => m.labels.severity === severity).reduce((sum, m) => sum + m.last, 0),
        }))
      : null,
    scannedPaths: new Set(findings.map((m) => m.labels.path)).size,
    tests: modes.map((mode) => ({
      mode,
      pass: count(mode, "tests_pass"),
      fail: count(mode, "tests_fail"),
      skip: count(mode, "tests_skip"),
    })),
    slowTests: slow
      .filter((m) => m.median >= 1)
      .sort((a, b) => b.median - a.median)
      .slice(0, 5)
      .map((m) => ({
        test: m.labels.test ?? "",
        pkg: m.labels.package ?? "",
        mode: m.labels.mode ?? "",
        medianS: m.median,
        lastS: m.last,
      })),
    other: metrics.filter((m) => !summarized.has(m) && !(m.last === 0 && m.median === 0)),
  };
}
