# e2e

End-to-end tests for TheTree in [Robot Framework](https://robotframework.org),
written BDD-style: scenarios read as English, the code lives underneath.
API tests use RequestsLibrary; browser tests use the
[Browser library](https://robotframework-browser.org), which is Playwright.

Only IdentityX is covered so far. Scenarios assert the contract in
`api/identityx/api-spec.yml` (and what a user should see), so a red test is
a bug report: they are not adjusted to match what the code happens to do.

## Run

```bash
just e2e                          # up → every suite → down; browsers visible
just e2e api                      # API suites only (the front is not started)
just e2e ui --headless            # browser suites, no window
just e2e identityx/api/03__sessions.robot
just e2e all -t "*password*"      # robot args pass through
```

The first argument picks what runs: `all` (default), `api` or `ui` (by tag,
across every service), or any path under `suites/`, down to one file.
Then flags:

| Flag         | Effect |
|--------------|--------|
| `--headless` | no browser window; also forced with `E2E_HEADLESS=1` or when there is no display |
| `--slow`     | 250ms between browser actions, to follow a visible run by eye |

Browser suites open a visible Chromium by default so you can watch them.

While iterating, keep the stack up between runs:

```bash
just e2e-up                       # --no-front: API and mailpit only
just e2e-run identityx/ui --slow
just e2e-down
```

Needs Docker, Go (for goose), pnpm (root `pnpm install` done) and
[uv](https://docs.astral.sh/uv/). `just init` (from `e2e/`, or
`just e2e/init` from the root) installs the Robot deps, Playwright and
Chromium into `e2e/.venv`; `e2e-run` runs it too, so it's only needed to set
up the venv ahead of time (e.g. for the IDE). Every root `e2e*` recipe also
works from `e2e/`. Reports land in `e2e/results/` (`log.html` has every step).

## The stack

`just e2e-up` starts compose project `thetree-e2e` from `e2e/compose.yml`,
apart from the dev stack:

| What           | Where                    |
|----------------|--------------------------|
| IdentityX API  | http://localhost:18080   |
| IdentityX app  | http://localhost:13001   |
| mailpit (UI + API) | http://localhost:18025 |

Postgres is tmpfs and the secrets are minted per run, so every `e2e-up`
starts from an empty platform. The suite setup creates the super admin via
`/auth/setup`; every scenario then registers its own users with unique
emails, so suites can rerun against a stack that is already up.

## Layout

```
suites/identityx/          # the English: Given/When/Then scenarios
  api/  ui/                # tagged `api` / `ui`
resources/
  identityx/steps.resource # imports the sentences below; suites import only this
  identityx/steps/         # each sentence → code (state in test variables):
                           #   common, accounts, tenants, project_settings, app
  identityx/api.resource   # one keyword per HTTP operation, no assertions
  identityx/ui.resource    # Playwright selectors and page actions
  mailpit.resource         # read the inbox, pull tokens out of links
  common.resource          # URLs, unique data, fun.Response envelope helpers
```

Writing a scenario: use the sentences in `steps.resource`. When one is
missing, add it to the matching file in `steps/`, built from the code-level keywords — a `.robot`
file under `suites/` should never call `api.resource` or Browser directly.
