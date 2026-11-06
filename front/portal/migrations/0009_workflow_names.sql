-- Dependabot Updates runs were stored under their per-run titles
-- ("github_actions in /., /.github/actions/** - Update #1684151908"), one CI
-- row each; the GitHub webhook now stores the workflow's own name.
UPDATE workflow_runs SET workflow = 'Dependabot Updates'
WHERE workflow LIKE '% - Update #%';
