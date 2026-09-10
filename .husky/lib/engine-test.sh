#!/bin/bash
# Regression suite for .husky/lib/engine.sh — synthetic diffs, no repo state
# touched. Run from the repo root:
#
#   .husky/lib/engine-test.sh
#
# Prints one "ok <case> → <result>" line per case and exits non-zero on any
# FAIL.
set -u
cd /home/sophia/GolandProjects/TheTree
ENGINE_TMP=$(mktemp -d); trap 'rm -rf "$ENGINE_TMP"' EXIT
. .husky/lib/engine.sh

fail=0
check() { # check <desc> <expected> <actual>
    if [ "$2" = "$3" ]; then echo "ok   $1 → $3"; else echo "FAIL $1 → got '$3', want '$2'"; fail=1; fi
}

scope_of() { printf '%s\n' "$@" >"$ENGINE_TMP/l"; engine_scope "$ENGINE_TMP/l"; }
targets_of() { printf '%s\n' "$@" >"$ENGINE_TMP/l"; engine_targets "$ENGINE_TMP/l"; }

check "single api"            "payssage"          "$(scope_of api/payssage/main.go)"
check "two apis"              "apis"              "$(scope_of api/payssage/a.go api/identityx/b.go)"
check "single front"          "univents-ui"       "$(scope_of front/univents/src/x.ts)"
check "two fronts"            "fronts"            "$(scope_of front/univents/a.ts front/payssage/b.ts)"
check "lib-go"                "lib-go"            "$(scope_of lib/go/errx/e.go)"
check "one lib/ts member"     "lib-ts"            "$(scope_of lib/ts/front-core/i.ts)"
check "two lib/ts members"    "libs-ts"           "$(scope_of lib/ts/front-core/i.ts lib/ts/ui-base/b.ts)"
check "lib mix go+ts"         "libs"              "$(scope_of lib/go/e.go lib/ts/ui-base/b.ts)"
check "single sdk-ts"         "sdk-ts-identityx"  "$(scope_of sdk/ts/identityx/s.ts)"
check "two sdk-ts"            "sdks-ts"           "$(scope_of sdk/ts/identityx/s.ts sdk/ts/payssage/p.ts)"
check "two sdk-go"            "sdks-go"           "$(scope_of sdk/go/IdentityX/a.go sdk/go/Payssage/b.go)"
check "sdk mix go+ts"         "sdks"              "$(scope_of sdk/go/IdentityX/a.go sdk/ts/payssage/p.ts)"
check "case-kept Payssage"    "sdk-ts-Payssage"   "$(scope_of sdk/ts/Payssage/x.ts)"
check "api+front"             "root"              "$(scope_of api/payssage/a.go front/univents/b.ts)"
check "docs only"             "docs"              "$(scope_of docs/adr/x.md)"
check "README only"           "docs"              "$(scope_of README.md)"
check "infra only"            "infra"             "$(scope_of .forgejo/workflows/deploy.yml)"
check "docs+infra"            "root"              "$(scope_of docs/x.md .forgejo/y.yml)"
check "go.work alone"         "root"              "$(scope_of go.work)"
check "go.work + api"         "root"              "$(scope_of go.work api/payssage/a.go)"
check "misc root file"        "root"              "$(scope_of orval.config.ts)"
check "univents-solid scope"  "univents-solid-ui" "$(scope_of front/univents-solid/src/a.ts)"

t=$(targets_of go.work)
echo "$t" | grep -q '^go:api/payssage$'      && echo "ok   broad-go hits api"    || { echo "FAIL broad-go api"; fail=1; }
echo "$t" | grep -q '^go:lib/go$'            && echo "ok   broad-go hits lib/go" || { echo "FAIL broad-go lib/go"; fail=1; }
echo "$t" | grep -q '^ts:'                   && echo "FAIL broad-go leaked ts"   || echo "ok   broad-go has no ts"
[ "$(echo "$t" | grep -c '^go:')" -eq 7 ]    && echo "ok   broad-go = 7 go members" || { echo "FAIL broad-go count: $(echo "$t"|grep -c '^go:')"; fail=1; }

t=$(targets_of pnpm-lock.yaml)
echo "$t" | grep -q '^ts:front/identityx$'   && echo "ok   broad-ts hits front"  || { echo "FAIL broad-ts front"; fail=1; }
echo "$t" | grep -q 'univents-solid'         && { echo "FAIL broad-ts hit univents-solid"; fail=1; } || echo "ok   broad-ts skips univents-solid"
echo "$t" | grep -q '^go:'                   && { echo "FAIL broad-ts leaked go"; fail=1; } || echo "ok   broad-ts has no go"

t=$(targets_of front/univents-solid/src/a.ts)
[ "$t" = "ts:front/univents-solid" ]         && echo "ok   solid direct target"  || { echo "FAIL solid: '$t'"; fail=1; }

t=$(targets_of api/payssage/x.go lib/ts/orval/y.ts front/univents-solid/z.ts)
echo "$t" | sort | tr '\n' ' '; echo
[ "$(echo "$t" | wc -l)" -eq 3 ]             && echo "ok   union targets"        || { echo "FAIL union: $t"; fail=1; }

check "valid: payssage"        "0" "$(engine_valid_scope payssage >/dev/null 2>&1; echo $?)"
check "valid: univents-ui"     "0" "$(engine_valid_scope univents-ui >/dev/null 2>&1; echo $?)"
check "valid: univents-solid-ui" "0" "$(engine_valid_scope univents-solid-ui >/dev/null 2>&1; echo $?)"
check "valid: sdk-go-IdentityX" "0" "$(engine_valid_scope sdk-go-IdentityX >/dev/null 2>&1; echo $?)"
check "valid: sdk-ts-Payssage" "0" "$(engine_valid_scope sdk-ts-Payssage >/dev/null 2>&1; echo $?)"
check "valid: infra"           "0" "$(engine_valid_scope infra >/dev/null 2>&1; echo $?)"
check "invalid: nope"          "1" "$(engine_valid_scope nope >/dev/null 2>&1; echo $?)"
check "invalid: not-a-scope" "1" "$(engine_valid_scope not-a-scope >/dev/null 2>&1; echo $?)"

exit $fail
