#!/usr/bin/env node
// One-off backfill of the Delivery tables from before the GitHub webhook and
// the deploy hook existed. Prints SQL; nothing is written here:
//
//   git fetch --tags
//   node scripts/backfill-github.mjs > /tmp/backfill.sql
//   pnpm exec wrangler d1 execute portal --remote --file /tmp/backfill.sql
//
// Needs `gh` logged in. Every statement only inserts what's missing, so a
// rerun or a race with the live webhook never overwrites newer data. Past
// deploys come from the Deploy workflows' publish jobs; their health is
// unknown (null).

import { execFileSync } from "node:child_process";

const REPO = "TrieOH/TheTree";
const DAYS = 60;
const since = new Date(Date.now() - DAYS * 86_400_000);

const sh = (cmd, args) => execFileSync(cmd, args, { encoding: "utf8", maxBuffer: 64 << 20 });
const gh = (path) => JSON.parse(sh("gh", ["api", "-H", "Accept: application/vnd.github+json", path]));
const ghPages = (path, key, stop) => {
  const all = [];
  for (let page = 1; page <= 10; page++) {
    const body = gh(`${path}${path.includes("?") ? "&" : "?"}per_page=100&page=${page}`);
    const items = key ? body[key] : body;
    all.push(...items);
    if (items.length < 100 || (stop && stop(items.at(-1)))) break;
  }
  return all;
};

// GitHub sends `…:33Z`; the portal stores toISOString()'s `…:33.000Z`.
const iso = (value) => (value ? new Date(value).toISOString() : null);
const q = (value) =>
  value === null || value === undefined ? "NULL" : typeof value === "number" ? String(value) : `'${String(value).replaceAll("'", "''")}'`;
const out = [];
const insert = (table, row) =>
  out.push(`INSERT OR IGNORE INTO ${table} (${Object.keys(row).join(", ")}) VALUES (${Object.values(row).map(q).join(", ")});`);

// Tags, from git: the tags API has no dates.
const RELEASE_TAG = /^(.+)\/(v\d+\.\d+\.\d+.*)$/;
for (const line of sh("git", ["for-each-ref", "refs/tags", "--format=%(refname:short)\t%(*objectname)%(objectname)\t%(creatordate:iso-strict)\t%(subject)"]).split("\n")) {
  if (!line) continue;
  const [name, shas, date, subject] = line.split("\t");
  const match = RELEASE_TAG.exec(name);
  if (!match) continue;
  insert("tags", {
    name,
    component: match[1],
    version: match[2],
    // annotated tags print the commit then the tag object; keep the commit
    sha: shas.slice(0, 40),
    message: subject || null,
    created_at: new Date(date).toISOString(),
  });
}

// Pull requests: open ones, and everything updated in the window.
for (const pr of ghPages(`repos/${REPO}/pulls?state=all&sort=updated&direction=desc`, null, (last) => new Date(last.updated_at) < since)) {
  if (pr.state !== "open" && new Date(pr.updated_at) < since) continue;
  let review = null;
  if (pr.state === "open") {
    const states = gh(`repos/${REPO}/pulls/${pr.number}/reviews?per_page=100`).map((r) => r.state.toLowerCase());
    review = states.findLast((s) => s !== "commented") ?? states.at(-1) ?? null;
  }
  insert("pull_requests", {
    number: pr.number,
    title: pr.title,
    author: pr.user.login,
    state: pr.merged_at ? "merged" : pr.state,
    draft: pr.draft ? 1 : 0,
    review_state: review,
    opened_at: iso(pr.created_at),
    merged_at: iso(pr.merged_at),
    merge_sha: pr.merged_at ? pr.merge_commit_sha : null,
    updated_at: iso(pr.updated_at),
  });
}

// Workflow runs in the window, and the deploys they made.
const runs = ghPages(`repos/${REPO}/actions/runs?created=>=${since.toISOString().slice(0, 10)}`, "workflow_runs");
// a run's name can be per run (Dependabot Updates); the workflow's is stable
const workflowNames = Object.fromEntries(
  ghPages(`repos/${REPO}/actions/workflows`, "workflows").map((w) => [w.id, w.name]),
);
for (const run of runs) {
  const completed = run.status === "completed";
  insert("workflow_runs", {
    id: run.id,
    workflow: workflowNames[run.workflow_id] ?? run.name,
    branch: run.head_branch ?? "",
    sha: run.head_sha,
    conclusion: completed ? run.conclusion : null,
    attempt: run.run_attempt,
    queued_at: iso(run.created_at),
    started_at: iso(run.run_started_at),
    completed_at: completed ? iso(run.updated_at) : null,
  });

  const job = { Deploy: "publish", "Deploy Pages": "deploy" }[run.name];
  const component = RELEASE_TAG.exec(run.head_branch ?? "")?.[1];
  if (!job || !component || !completed || component.includes("-sdk-ts")) continue;
  const publish = gh(`repos/${REPO}/actions/runs/${run.id}/jobs?per_page=100`).jobs.find((j) => j.name === job);
  if (!publish || !["success", "failure", "cancelled"].includes(publish.conclusion)) continue;
  out.push(
    `INSERT INTO deploys (component, tag, sha, actor, status, health_ok, at, run_id) SELECT ${[
      component,
      run.head_branch,
      run.head_sha,
      run.triggering_actor?.login ?? run.actor.login,
      publish.conclusion,
      null,
      iso(publish.completed_at),
      run.id,
    ]
      .map(q)
      .join(", ")} WHERE NOT EXISTS (SELECT 1 FROM deploys WHERE run_id = ${run.id});`,
  );
}

process.stdout.write(`${out.join("\n")}\n`);
console.error(`backfill: ${out.length} statements`);
