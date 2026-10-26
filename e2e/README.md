# e2e

End-to-end tests for TheTree in [Robot Framework](https://robotframework.org),
written BDD-style: scenarios read as English, the code lives underneath.
API tests use RequestsLibrary; browser tests use the
[Browser library](https://robotframework-browser.org), which is Playwright.

Only IdentityX is covered so far.

## Run

```bash
just e2e                          # up → every IdentityX suite → down
just e2e identityx/api            # API suites only
just e2e identityx --include ui   # robot args pass through
```

While iterating, keep the stack up between runs:

```bash
just e2e-up
just e2e-run identityx/api/03__sessions.robot
just e2e-down
```

Needs Docker, Go (for goose), pnpm (root `pnpm install` done) and
[uv](https://docs.astral.sh/uv/). The first `e2e-run` also runs
`rfbrowser init chromium`, which installs Playwright and Chromium into
`e2e/.venv`. Reports land in `e2e/results/` (`log.html` has every step).

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
  identityx/steps.resource # each sentence → code (state in test variables)
  identityx/api.resource   # one keyword per HTTP operation, no assertions
  identityx/ui.resource    # Playwright selectors and page actions
  mailpit.resource         # read the inbox, pull tokens out of links
  common.resource          # URLs, unique data, fun.Response envelope helpers
```

Writing a scenario: use the sentences in `steps.resource`. When one is
missing, add it there, built from the code-level keywords — a `.robot`
file under `suites/` should never call `api.resource` or Browser directly.
