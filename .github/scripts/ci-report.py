"""ci-report.py: the last job of every release run.

Reads this run's jobs and steps from the Actions API (start, end and result of
each, plus how long each job waited for a runner), what ran inside each
composite step (from the job logs: the runner marks every action with
`##[end-action id=...;duration_ms=N]`), and the metrics the jobs recorded
with .github/scripts/metric (handed over as job outputs, NEEDS_JSON). Then it
tells the dev portal:

  POST /webhooks/deploy  what the release did (component, tag, status, and the
                         stage it stopped at: parse, checks or publish)
  POST /webhooks/ci      the run's jobs, steps, substeps and metrics, into D1

Never fails the job: a release is not worse because the portal didn't hear
about it.
"""

import json
import os
import re
import sys
import urllib.error
import urllib.request
from datetime import datetime, timedelta

env = os.environ
PORTAL = env.get("PORTAL_URL", "https://portal.trieoh.com")
STAGES = ["parse", "checks", "publish"]
MAX_SUBSTEPS = 300  # per job; a real one has ~60

START = re.compile(r"^(\S+) ##\[start-action display=(.*);id=([^;\]]+)\]")
END = re.compile(r"^(\S+) ##\[end-action id=([^;]+);.*duration_ms=(\d+)")


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, *args, **kwargs):
        return None


def request(url, *, token=None, body=None, bearer=None, raw=False):
    headers = {"Accept": "application/vnd.github+json", "User-Agent": "thetree-ci-report"}
    if token or bearer:
        headers["Authorization"] = f"Bearer {token or bearer}"
    data = None
    if body is not None:
        data = json.dumps(body).encode()
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(url, data=data, headers=headers, method="POST" if data else "GET")
    with urllib.request.urlopen(req, timeout=20) as response:
        content = response.read()
        return content if raw else (json.loads(content) if content else None)


def api(path):
    return f"https://api.github.com/repos/{env['GITHUB_REPOSITORY']}{path}"


def fetch_jobs():
    base = api(f"/actions/runs/{env['GITHUB_RUN_ID']}/attempts/{env['GITHUB_RUN_ATTEMPT']}/jobs?per_page=100")
    jobs = []
    for page in range(1, 4):
        batch = request(f"{base}&page={page}", token=env["GH_TOKEN"])["jobs"]
        jobs += batch
        if len(batch) < 100:
            break
    return jobs


def fetch_log(job_id):
    """The job's plain-text log. The API answers with a redirect to a signed
    URL, which must be fetched without the token."""
    req = urllib.request.Request(
        api(f"/actions/jobs/{job_id}/logs"),
        headers={"Authorization": f"Bearer {env['GH_TOKEN']}", "User-Agent": "thetree-ci-report"},
    )
    try:
        urllib.request.build_opener(NoRedirect).open(req, timeout=20)
    except urllib.error.HTTPError as redirect:
        if redirect.code not in (301, 302, 307):
            raise
        return request(redirect.headers["Location"], raw=True).decode("utf-8", "replace")
    return ""


def parse_time(stamp):
    # the log has 7 decimals; fromisoformat takes up to 6
    stamp = re.sub(r"(\.\d{6})\d*", r"\1", stamp.replace("Z", "+00:00"))
    return datetime.fromisoformat(stamp)


def substeps(job):
    """Every action the runner marked in the job's log, placed in the
    top-level step whose time window holds its end."""
    names, found = {}, []
    for line in fetch_log(job["id"]).splitlines():
        if match := START.match(line):
            names[match[3]] = match[2]
        elif match := END.match(line):
            found.append((match[2], parse_time(match[1]), int(match[3])))
    windows = [
        (step["number"], parse_time(step["started_at"]), parse_time(step["completed_at"]))
        for step in job.get("steps", [])
        if step.get("started_at") and step.get("completed_at")
    ]
    out = []
    for action_id, end, ms in found[:MAX_SUBSTEPS]:
        # step times are whole seconds; give the window a second of slack
        step = next((n for n, a, b in windows if a <= end <= b + timedelta(seconds=1)), None)
        out.append({
            "step": step,
            "action_id": action_id,
            "name": names.get(action_id, action_id)[:200],
            "depth": action_id.count("."),
            "completed_at": end.isoformat().replace("+00:00", "Z"),
            "duration_ms": ms,
        })
    return out


def stage_of(job_name):
    # reusable workflows name their jobs "<caller> / <job>"
    return job_name.rsplit(" / ", 1)[-1]


def needs():
    """{caller job: {result, outputs}} from deploy.yml's toJSON(needs)."""
    try:
        return json.loads(env.get("NEEDS_JSON") or "{}")
    except json.JSONDecodeError:
        return {}


def read_metrics(jobs, called):
    """Metrics come as job outputs (<job>-metrics: JSON lines) of the reusable workflow."""
    by_stage = {stage_of(job["name"]): job["name"] for job in jobs}
    metrics = []
    for key, text in called.get("outputs", {}).items():
        if not key.endswith("-metrics"):
            continue
        for line in (text or "").splitlines():
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
                print(f"::warning::skipping bad metric line from {key}: {line.strip()}")
    return metrics


def outcome(jobs, superseded):
    """(status, stage): the first stage that didn't succeed, else publish.
    A failed "Parse tag" step counts as the parse stage."""
    for job in sorted(jobs, key=lambda j: STAGES.index(stage_of(j["name"])) if stage_of(j["name"]) in STAGES else 9):
        stage = stage_of(job["name"])
        if job["conclusion"] == "cancelled":
            return "cancelled", stage
        if job["conclusion"] not in ("success", "skipped"):
            parse_failed = any(s["name"] == "Parse tag" and s.get("conclusion") == "failure" for s in job.get("steps", []))
            return "failure", "parse" if parse_failed else stage
    if superseded:
        return "cancelled", "publish"
    ran = {stage_of(job["name"]) for job in jobs if job["conclusion"] == "success"}
    return "success", "publish" if "publish" in ran else "checks" if "checks" in ran else "parse"


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
    return int((parse_time(step["completed_at"]) - parse_time(step["started_at"])).total_seconds())


def main():
    tag = env["TAG"]
    report_job = env.get("REPORT_JOB", "report")

    # the reusable workflow that ran: service, front, sdk or pages
    called_name, called = next(
        ((name, n) for name, n in needs().items() if n.get("result") not in (None, "skipped")), (None, {})
    )
    outputs = called.get("outputs", {})
    component = outputs.get("artifact") or tag.rpartition("/v")[0] or tag
    lane = outputs.get("lane") or ("express" if "-hotfix." in tag else "full")

    jobs = [
        job for job in fetch_jobs()
        if job["name"] != report_job and job["status"] == "completed" and job["conclusion"] != "skipped"
    ]
    status, stage = outcome(jobs, outputs.get("superseded") == "true")
    run_id = int(env["GITHUB_RUN_ID"])

    report_jobs = []
    for job in jobs:
        try:
            subs = substeps(job)
        except Exception as error:  # noqa: BLE001 - substeps are a bonus
            print(f"::warning::no substeps for {job['name']}: {error}")
            subs = []
        report_jobs.append({
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
            "substeps": subs,
        })

    run = {
        "run_id": run_id,
        "attempt": int(env["GITHUB_RUN_ATTEMPT"]),
        "component": component,
        "tag": tag,
        "sha": env["SHA"],
        "type": called_name or "unknown",
        "lane": lane,
        "status": status,
        "stage": stage,
        "jobs": report_jobs,
        "metrics": read_metrics(jobs, called),
    }
    substep_count = sum(len(job["substeps"]) for job in report_jobs)
    print(
        f"[report] {component} {tag}: {status} at {stage}; {len(jobs)} jobs, "
        f"{substep_count} substeps, {len(run['metrics'])} metrics"
    )
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
