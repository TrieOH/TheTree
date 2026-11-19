#!/usr/bin/env bash
# Turns a service's .example.env into an env file the image can boot with in
# CI: the database points at the throwaway Postgres, empty secrets get random
# values, and nothing points at a real service (the example's URLs are compose
# hostnames, which don't resolve on the runner, and no Backend calls out at
# boot). Telemetry stays off so the run sends nothing to Grafana.
#
#   env.sh <example.env> <db-host> <stub-url>
#
# <stub-url> answers for every other service the Backend calls (IdentityX,
# Payssage): Univents checks its Payssage wallet at boot, and pointing those
# URLs at a stub keeps the smoke run from ever reaching a real service.
set -euo pipefail
src=$1 db_host=$2 stub=$3
seen_project=

while IFS= read -r line || [ -n "$line" ]; do
	line="${line%%[[:space:]]#*}"                 # inline comment
	[[ "$line" =~ ^[[:space:]]*(#|$) ]] && continue
	[[ "$line" == *=* ]] || continue
	key="${line%%=*}" val="${line#*=}"
	key="${key//[[:space:]]/}"
	val="${val%"${val##*[![:space:]]}"}"          # trailing spaces
	val="${val#\'}" val="${val%\'}" val="${val#\"}" val="${val%\"}"
	case "$key" in
	POSTGRES_*) continue ;;                       # the Postgres image's own vars
	*_POSTGRES_HOST | *_POSTGRES_DIRECT_HOST) val=$db_host ;;
	*_POSTGRES_PORT | *_POSTGRES_DIRECT_PORT) val=5432 ;;
	*_POSTGRES_DB | *_POSTGRES_USER | *_POSTGRES_PASSWORD) val=smoke ;;
	*_POSTGRES_SSLMODE) val=disable ;;
	OTEL_*) continue ;;
	IDENTITY_X_URL | PAYSSAGE_URL) val=$stub ;;
	IDENTITY_X_PROJECT_ID) seen_project=1; [ -n "$val" ] || val=$(cat /proc/sys/kernel/random/uuid) ;;
	PORT) val=8080 ;;
	PROFILE_PORT) val= ;;
	*PASSWORD | *SECRET | *_KEY | *_TOKEN) [ -n "$val" ] || val=$(openssl rand -hex 32) ;;
	esac
	printf '%s=%s\n' "$key" "$val"
done <"$src"
# required by the IdentityX SDK, not always in the example
[ -n "$seen_project" ] || printf 'IDENTITY_X_PROJECT_ID=%s\n' "$(cat /proc/sys/kernel/random/uuid)"
