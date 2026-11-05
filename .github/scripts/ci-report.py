"""ci-report.py: the last job of every release run.

Reads this run's jobs and steps from the Actions API (start, end and result of
each, plus how long each job waited for a runner), adds the metrics the jobs
recorded with .github/scripts/metric, and tells the dev portal:

  POST /webhooks/deploy  what the release did (component, tag, status, and the
                         stage it stopped at: parse, checks or publish)
  POST /webhooks/ci      the run's jobs, steps and metrics, into D1

Never fails the job: a release is not worse because the portal didn't hear
about it.
"""

import glob
import json
import os
import sys
import urllib.request
from datetime import datetime

env = os.environ
PORTAL = env.get("PORTAL_URL", "https://portal.trieoh.com")
STAGES = ["parse", "checks", "publish"]


def request(url, *, token=None, body=None, bearer=None):
    headers = {"Accept": "application/vnd.github+json", "User-Agent": "thetree-ci-report"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    if bearer:
        headers["Authorization"] = f"Bearer {bearer}"
    data = None
    if body is not None:
        data = json.dumps(body).encode()
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(url, data=data, headers=headers, method="POST" if data else "GET")
    with urllib.request.urlopen(req, timeout=20) as response:
        raw = response.read()
        return json.loads(raw) if raw else None


def fetch_jobs():
    base = (
        f"https://api.github.com/repos/{env['GITHUB_REPOSITORY']}/actions/runs/"
        f"{env['GITHUB_RUN_ID']}/attempts/{env['GITHUB_RUN_ATTEMPT']}/jobs?per_page=100"
    )
    jobs = []
    for page in range(1, 4):
        batch = request(f"{base}&page={page}", token=env["GH_TOKEN"])["jobs"]
        jobs += batch
        if len(batch) < 100:
            break
    return jobs


def stage_of(job_name):
    # reusable workflows name their jobs "<caller> / <job>"
    return job_name.rsplit(" / ", 1)[-1]


def read_metrics(jobs):
    by_stage = {stage_of(job["name"]): job["name"] for job in jobs}
    metrics = []
    for path in sorted(glob.glob(os.path.join(env.get("METRICS_DIR", ""), "*.jsonl"))):
        with open(path) as lines:
            for line in lines:
                try:
                    m = json.loads(line)
                    metrics.append({
                        "job": by_stage.get(m["job"], m["job"]),
                        "source": str(m["source"]),
                        "name": str(m["name"]),
                        "value": float(m["value"]),
                        "labels": {str(k): str(v) for k, v in m.get("labels", {}).items()},
                    })
                except (ValueError, KeyError, TypeError):
                    print(f"::warning::skipping bad metric line in {os.path.basename(path)}: {line.strip()}")
    return metrics


def outcome(jobs):
    """(status, stage): the first stage that didn't succeed, else publish."""
    for stage in STAGES:
        for job in jobs:
            if stage_of(job["name"]) != stage:
                continue
            if job["conclusion"] == "cancelled":
                return "cancelled", stage
            if job["conclusion"] not in ("success", "skipped"):
                return "failure", stage
    ran = {stage_of(job["name"]) for job in jobs if job["conclusion"] == "success"}
    return "success", "publish" if "publish" in ran else max(ran, key=STAGES.index, default="parse")


def summarize(run):
    steps = [
        (job["name"], step["name"], seconds(step))
        for job in run["jobs"]
        for step in job["steps"]
        if step["status"] != "skipped"
    ]
    lines = ["### Slowest steps", "", "| Job | Step | Seconds |", "|---|---|---:|"]
    lines += [f"| {job} | {name} | {secs} |" for job, name, secs in sorted(steps, key=lambda s: -s[2])[:10]]
    lines += ["", f"{len(run['metrics'])} metrics recorded."]
    with open(env.get("GITHUB_STEP_SUMMARY", os.devnull), "a") as summary:
        summary.write("\n".join(lines) + "\n")


def seconds(step):
    if not step["started_at"] or not step["completed_at"]:
        return 0
    start, end = (datetime.fromisoformat(step[k].replace("Z", "+00:00")) for k in ("started_at", "completed_at"))
    return int((end - start).total_seconds())


def main():
    tag = env["TAG"]
    component = env.get("ARTIFACT") or tag.rpartition("/v")[0] or tag
    report_job = env.get("REPORT_JOB", "report")

    jobs = [
        job for job in fetch_jobs()
        if job["name"] != report_job and job["status"] == "completed" and job["conclusion"] != "skipped"
    ]
    status, stage = outcome(jobs)
    run_id = int(env["GITHUB_RUN_ID"])

    run = {
        "run_id": run_id,
        "attempt": int(env["GITHUB_RUN_ATTEMPT"]),
        "component": component,
        "tag": tag,
        "sha": env["SHA"],
        "type": env.get("TYPE") or "unknown",
        "lane": env.get("LANE") or "unknown",
        "status": status,
        "stage": stage,
        "jobs": [
            {
                "name": job["name"],
                "status": job["conclusion"],
                "queued_at": job.get("created_at"),
                "started_at": job.get("started_at"),
                "completed_at": job.get("completed_at"),
                "runner": job.get("runner_name"),
                "steps": [
                    {
                        "number": step["number"],
                        "name": step["name"],
                        "status": step.get("conclusion") or step["status"],
                        "started_at": step.get("started_at"),
                        "completed_at": step.get("completed_at"),
                    }
                    for step in job.get("steps", [])
                ],
            }
            for job in jobs
        ],
        "metrics": read_metrics(jobs),
    }
    print(f"[report] {component} {tag}: {status} at {stage}; {len(jobs)} jobs, {len(run['metrics'])} metrics")
    summarize(run)

    secret = env.get("SECRET", "")
    deploy = {
        "component": component, "tag": tag, "sha": env["SHA"], "actor": env["ACTOR"],
        "status": status, "stage": stage, "run_id": run_id,
    }
    for path, body in (("/webhooks/deploy", deploy), ("/webhooks/ci", run)):
        try:
            request(f"{PORTAL}{path}", body=body, bearer=secret)
            print(f"[report] {path} ok")
        except Exception as error:  # noqa: BLE001 - reporting must never fail the run
            print(f"::warning::could not send {path} to the portal: {error}")

    if env.get("DUMP"):
        json.dump(run, sys.stdout, indent=1)


if __name__ == "__main__":
    try:
        main()
    except Exception as error:  # noqa: BLE001
        print(f"::warning::release report failed: {error}")
