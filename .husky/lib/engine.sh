# engine.sh — the diff∩scope engine shared by commit-msg, pre-commit and pre-push.
#
# Principles (plan §3.7):
#   - Dumb: recompute everything from the raw diff, every time. Zero trust in
#     the stated scope. No caching, no cleverness. CI never uses this (deploy.yml
#     routes by tag prefix alone).
#   - Scopes are DISCOVERED from the repo layout — a new api/front/lib/sdk is a
#     valid scope with zero edits here.
#   - Naming = closest common ancestor of everything changed.
#   - Checks = actual per-member diffs ∩ scope, never the scope's size.
#
# Sourcing contract (hooks do this before sourcing):
#   ENGINE_TMP=$(mktemp -d "${TMPDIR:-/tmp}/thetree-hooks.XXXXXX")
#   trap 'rm -rf "$ENGINE_TMP"' EXIT
#   . "$(dirname "$0")/lib/engine.sh"

# --- classification ---------------------------------------------------------
#
# engine_classify FILE
#   Prints one line:  <category>\t<member-dir>\t<scope-name>
#   or nothing if the file is at a category root (e.g. api/ itself).
#
# Categories:
#   apis / fronts / libs / sdks   — real checkable members
#   docs                          — zero checks
#   infra                         — zero checks (dot-dirs + compose.yml + justfile, Q16)
#   root-go / root-ts             — workspace-level files; broad escape hatch
#                                   (checks run for all members they can affect)
#   root-misc                     — other root files; zero checks, naming scope root
engine_classify() {
    f=$1
    case "$f" in
        api/*)
            rest=${f#api/}; name=${rest%%/*}
            [ "$name" = "$rest" ] && return 0
            printf 'apis\tapi/%s\t%s\n' "$name" "$name" ;;
        lib/go/*)
            printf 'libs\tlib/go\tlib-go\n' ;;
        lib/ts/*)
            rest=${f#lib/ts/}; name=${rest%%/*}
            [ "$name" = "$rest" ] && return 0
            printf 'libs\tlib/ts/%s\tlib-ts\n' "$name" ;;
        sdk/go/*)
            rest=${f#sdk/go/}; name=${rest%%/*}
            [ "$name" = "$rest" ] && return 0
            printf 'sdks\tsdk/go/%s\tsdk-go-%s\n' "$name" "$name" ;;
        sdk/ts/*)
            rest=${f#sdk/ts/}; name=${rest%%/*}
            [ "$name" = "$rest" ] && return 0
            printf 'sdks\tsdk/ts/%s\tsdk-ts-%s\n' "$name" "$name" ;;
        front/*)
            rest=${f#front/}; name=${rest%%/*}
            [ "$name" = "$rest" ] && return 0
            printf 'fronts\tfront/%s\t%s-ui\n' "$name" "$name" ;;
        docs/*|README.md)
            printf 'docs\tdocs\tdocs\n' ;;
        .husky/*|.forgejo/*|.agents/*|compose.yml|justfile)
            printf 'infra\tinfra\tinfra\n' ;;
        go.work|go.work.sum|.golangci.yml)
            printf 'root-go\t%s\troot\n' "$f" ;;
        pnpm-lock.yaml|pnpm-workspace.yaml|package.json|patches/*)
            printf 'root-ts\t%s\troot\n' "$f" ;;
        *)
            printf 'root-misc\t%s\troot\n' "$f" ;;
    esac
}

# engine_classify_all LISTFILE
#   Classifies every file in LISTFILE into $ENGINE_TMP/classes (tab-separated).
engine_classify_all() {
    list=$1
    : >"$ENGINE_TMP/classes"
    while IFS= read -r f; do
        [ -n "$f" ] || continue
        engine_classify "$f" >>"$ENGINE_TMP/classes"
    done <"$list"
}

# --- naming scope (closest common ancestor) ---------------------------------
#
# engine_scope LISTFILE
#   Prints the closest-common-ancestor naming scope of all files in LISTFILE.
#
#   - any root file (go/go-ts/misc) present            → root
#   - docs only → docs · infra only → infra
#   - one category, one member   → that member's scope (payssage, univents-ui,
#                                    lib-go, lib-ts, sdk-ts-identityx, …)
#   - one category, many members → the category's group scope (apis, fronts,
#                                    libs, libs-ts, sdks-go, …)
#   - more than one category     → root
engine_scope() {
    list=$1
    engine_classify_all "$list"

    [ -s "$ENGINE_TMP/classes" ] || { printf 'root\n'; return 0; }

    cut -f1 "$ENGINE_TMP/classes" | sort -u >"$ENGINE_TMP/cats"

    # any root-level file wins: it can belong to anything
    if grep -q '^root' "$ENGINE_TMP/cats"; then
        printf 'root\n'; return 0
    fi

    # remaining categories are all top-level dirs: more than one present and
    # their closest common ancestor is root
    n=$(wc -l <"$ENGINE_TMP/cats")
    [ "$n" -gt 1 ] && { printf 'root\n'; return 0; }

    if grep -qx 'apis' "$ENGINE_TMP/cats"; then
        m=$(grep '^apis' "$ENGINE_TMP/classes" | cut -f2 | sort -u)
        [ "$(printf '%s\n' "$m" | wc -l)" -eq 1 ] && awk -F'\t' '$1=="apis"{print $3; exit}' "$ENGINE_TMP/classes" \
            || printf 'apis\n'
    elif grep -qx 'fronts' "$ENGINE_TMP/cats"; then
        m=$(grep '^fronts' "$ENGINE_TMP/classes" | cut -f2 | sort -u)
        [ "$(printf '%s\n' "$m" | wc -l)" -eq 1 ] && awk -F'\t' '$1=="fronts"{print $3; exit}' "$ENGINE_TMP/classes" \
            || printf 'fronts\n'
    elif grep -qx 'libs' "$ENGINE_TMP/cats"; then
        members=$(grep '^libs' "$ENGINE_TMP/classes" | cut -f2 | sort -u)
        scopes=$(grep '^libs' "$ENGINE_TMP/classes" | cut -f3 | sort -u)
        if [ "$(printf '%s\n' "$scopes" | wc -l)" -gt 1 ]; then
            printf 'libs\n'                       # lib/go + lib/ts mixed
        elif printf '%s\n' "$scopes" | grep -qx 'lib-go'; then
            printf 'lib-go\n'                     # single module, always singular
        elif [ "$(printf '%s\n' "$members" | wc -l)" -eq 1 ]; then
            printf 'lib-ts\n'
        else
            printf 'libs-ts\n'
        fi
    elif grep -qx 'sdks' "$ENGINE_TMP/cats"; then
        go_n=$(grep '^sdks' "$ENGINE_TMP/classes" | grep 'sdk/go/' | cut -f2 | sort -u | wc -l)
        ts_n=$(grep '^sdks' "$ENGINE_TMP/classes" | grep 'sdk/ts/' | cut -f2 | sort -u | wc -l)
        if [ "$go_n" -gt 0 ] && [ "$ts_n" -gt 0 ]; then
            printf 'sdks\n'                       # mixed languages
        elif [ "$go_n" -gt 1 ] || [ "$ts_n" -gt 1 ]; then
            [ "$go_n" -gt 1 ] && printf 'sdks-go\n' || printf 'sdks-ts\n'
        else
            awk -F'\t' '$1=="sdks"{print $3; exit}' "$ENGINE_TMP/classes"
        fi
    elif grep -qx 'infra' "$ENGINE_TMP/cats"; then
        printf 'infra\n'
    elif grep -qx 'docs' "$ENGINE_TMP/cats"; then
        printf 'docs\n'
    else
        printf 'root\n'
    fi
}

# --- check targets (diff ∩ scope, never scope size) -------------------------
#
# engine_targets LISTFILE
#   Prints the members that have runnable checks for the given diff:
#     go:<dir>   — golangci-lint + gotestsum -short
#     ts:<dir>   — package.json scripts (check/tsc|typecheck/test) run if present
#
#   The broad escape hatch (plan §3.7): a workspace-level file can break any
#   module it affects, so it fans out — but by what it can actually break:
#     root-go → all Go members · root-ts → all root-workspace TS members.
#   (front/univents-solid has its own lockfile/workspace — root lock files
#   don't affect it, so it only checks when directly changed.)
engine_targets() {
    list=$1
    engine_classify_all "$list"
    : >"$ENGINE_TMP/targets"

    add_target() {
        grep -qxF "$1" "$ENGINE_TMP/targets" || printf '%s\n' "$1" >>"$ENGINE_TMP/targets"
    }

    if grep -q '^root-go' "$ENGINE_TMP/classes"; then
        for d in api/*/; do [ -d "$d" ] && add_target "go:${d%/}"; done
        [ -d lib/go ] && add_target "go:lib/go"
        for d in sdk/go/*/; do [ -d "$d" ] && add_target "go:${d%/}"; done
    fi
    if grep -q '^root-ts' "$ENGINE_TMP/classes"; then
        for d in front/*/; do
            d=${d%/}
            [ -d "$d" ] || continue
            case "$d" in front/univents-solid) continue ;; esac   # independent workspace
            add_target "ts:$d"
        done
        for d in lib/ts/*/; do [ -d "$d" ] && add_target "ts:${d%/}"; done
        for d in sdk/ts/*/; do [ -d "$d" ] && add_target "ts:${d%/}"; done
    fi

    while IFS= read -r f; do
        [ -n "$f" ] || continue
        case "$f" in
            api/*)
                rest=${f#api/}; name=${rest%%/*}
                [ -d "api/$name" ] && add_target "go:api/$name" ;;
            lib/go/*) [ -d lib/go ] && add_target "go:lib/go" ;;
            sdk/go/*)
                rest=${f#sdk/go/}; name=${rest%%/*}
                [ -d "sdk/go/$name" ] && add_target "go:sdk/go/$name" ;;
            front/*)
                rest=${f#front/}; name=${rest%%/*}
                [ -d "front/$name" ] && add_target "ts:front/$name" ;;
            lib/ts/*)
                rest=${f#lib/ts/}; name=${rest%%/*}
                [ -d "lib/ts/$name" ] && add_target "ts:lib/ts/$name" ;;
            sdk/ts/*)
                rest=${f#sdk/ts/}; name=${rest%%/*}
                [ -d "sdk/ts/$name" ] && add_target "ts:sdk/ts/$name" ;;
        esac
    done <"$list"

    cat "$ENGINE_TMP/targets"
}

# --- scope validation -------------------------------------------------------
#
# engine_valid_scope SCOPE
#   Exit 0 if SCOPE is a valid discovered scope. Group scopes are always
#   accepted (they name a category, not a dir); singular scopes must exist
#   in the current layout.
engine_valid_scope() {
    scope=$1
    case "$scope" in
        root|docs|infra|apis|fronts|libs|lib-go|lib-ts|libs-ts|sdks|sdks-go|sdks-ts)
            return 0 ;;
    esac
    for d in api/*/; do
        [ "${d%/}" = "api/$scope" ] && return 0
    done
    for d in front/*/; do
        [ -d "$d" ] || continue
        [ "$(basename "$d")-ui" = "$scope" ] && return 0
    done
    for d in sdk/go/*/; do
        [ -d "$d" ] || continue
        [ "sdk-go-$(basename "$d")" = "$scope" ] && return 0
    done
    for d in sdk/ts/*/; do
        [ -d "$d" ] || continue
        [ "sdk-ts-$(basename "$d")" = "$scope" ] && return 0
    done
    return 1
}

# engine_all_scopes
#   Prints every currently valid scope (for hard-fail messages).
engine_all_scopes() {
    printf '  root docs infra apis fronts libs lib-go lib-ts libs-ts sdks sdks-go sdks-ts\n'
    for d in api/*/;  do [ -d "$d" ] && printf '  %s\n' "$(basename "$d")"; done
    for d in front/*/; do [ -d "$d" ] && printf '  %s\n' "$(basename "$d")-ui"; done
    for d in sdk/go/*/; do [ -d "$d" ] && printf '  sdk-go-%s\n' "$(basename "$d")"; done
    for d in sdk/ts/*/; do [ -d "$d" ] && printf '  sdk-ts-%s\n' "$(basename "$d")"; done
}

# engine_offenders LISTFILE
#   Prints "  <file>  →  <scope>" for each file (Q14: the hook says why).
engine_offenders() {
    list=$1
    while IFS= read -r f; do
        [ -n "$f" ] || continue
        line=$(engine_classify "$f")
        if [ -n "$line" ]; then
            scope=$(printf '%s\n' "$line" | cut -f3)
        else
            scope=root   # file at a category root (e.g. api/ itself)
        fi
        printf '  %s  →  %s\n' "$f" "$scope"
    done <"$list"
}
