import { WorkerEntrypoint } from "cloudflare:workers";

// The edge half of ADR-0020. Each API Worker has two entrypoints:
//
// - the default one, the gateway, runs on every request with no cache in
//   front of it: it answers CORS preflights, sends writes straight to the
//   container and purges what they changed;
// - `Cached`, which the gateway calls for reads, has Workers Caching in
//   front of it (per-entrypoint `cache` in wrangler.jsonc). The backend
//   decides what may be stored: every response is `private, no-store`
//   unless its operation is marked x-cache: shared.
//
// Purges are scoped to the entrypoint that stored the entries, so they run
// inside `Cached` too, over RPC.

// Headers the backend writes (lib/go/httpcache). Cloudflare consumes and
// strips Cache-Tag; the gateway strips the purge list.
export const PURGE_HEADER = "Cache-Purge-Tag";

// The queue message a failed purge is retried as.
export const PURGE_JOB = "cache.purge";
export type PurgeJob = { kind: typeof PURGE_JOB; args: { tags: string[] } };

// Delay before a refused purge is retried. Cloudflare refills the purge
// bucket at five requests a minute.
export const PURGE_RETRY_SECONDS = 30;

// How long a browser may reuse a preflight answer. Keep in step with
// lib/go/httpserver.CORSMaxAge, which answers preflights off Cloudflare.
const PREFLIGHT_MAX_AGE_SECONDS = 2 * 3600;
const ALLOWED_METHODS = ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"];
// lib/go/httpserver's set when CORS_ALLOWED_HEADERS is empty.
const DEFAULT_HEADERS = ["Content-Type", "Authorization", "X-Request-ID", "X-API-Key"];

export type CorsEnv = { CORS_ALLOWED_ORIGINS?: string; CORS_ALLOWED_HEADERS?: string };

function csv(value: string | undefined): string[] {
	return (value ?? "")
		.split(",")
		.map((s) => s.trim())
		.filter(Boolean);
}

// rs/cors's origin rules: "*" allows any origin, one "*" inside an entry is
// a wildcard (https://*.example.com), anything else matches exactly.
function originAllowed(origin: string, allowed: string[]): boolean {
	const o = origin.toLowerCase();
	return allowed.some((entry) => {
		const a = entry.toLowerCase();
		if (a === "*") return true;
		const star = a.indexOf("*");
		if (star === -1) return a === o;
		const prefix = a.slice(0, star);
		const suffix = a.slice(star + 1);
		return o.length >= prefix.length + suffix.length && o.startsWith(prefix) && o.endsWith(suffix);
	});
}

// Answers a CORS preflight the way the backend's CORS middleware would, so
// a preflight never starts a container. Null when the request is not a
// preflight; the backend still answers those itself.
export function preflight(request: Request, env: CorsEnv): Response | null {
	const method = request.headers.get("Access-Control-Request-Method");
	if (request.method !== "OPTIONS" || !method) return null;

	const headers = new Headers({
		Vary: "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		"Cache-Control": "no-store",
	});
	const origin = request.headers.get("Origin");
	const allowedHeaders = csv(env.CORS_ALLOWED_HEADERS);
	const allowed = (allowedHeaders.length ? allowedHeaders : DEFAULT_HEADERS).map((h) => h.toLowerCase());
	const requested = csv(request.headers.get("Access-Control-Request-Headers") ?? "").map((h) => h.toLowerCase());

	const ok =
		origin !== null &&
		originAllowed(origin, csv(env.CORS_ALLOWED_ORIGINS)) &&
		ALLOWED_METHODS.includes(method.toUpperCase()) &&
		requested.every((h) => allowed.includes("*") || allowed.includes(h));
	if (ok) {
		headers.set("Access-Control-Allow-Origin", origin);
		headers.set("Access-Control-Allow-Methods", method.toUpperCase());
		if (requested.length) headers.set("Access-Control-Allow-Headers", requested.join(", "));
		headers.set("Access-Control-Allow-Credentials", "true");
		headers.set("Access-Control-Max-Age", String(PREFLIGHT_MAX_AGE_SECONDS));
	}
	return new Response(null, { status: 204, headers });
}

// Reads that may go through the cache: GETs only (the cache keys GET and
// HEAD alike). Streams and upgrades go straight to the container: they are
// never stored, and a loopback call is one more hop for a long-lived
// connection.
export function cacheable(request: Request): boolean {
	if (request.method !== "GET") return false;
	if (request.headers.get("Upgrade")) return false;
	return !(request.headers.get("Accept") ?? "").includes("text/event-stream");
}

// The shape of `Cached` as the gateway reaches it through ctx.exports.
export interface CachedStub {
	fetch(request: Request): Promise<Response>;
	purge(tags: string[]): Promise<boolean>;
}

// ctx.exports.Cached, typed. ctx.exports is on by default from
// compatibility date 2025-11-17.
export function cached(ctx: ExecutionContext): CachedStub {
	return (ctx as unknown as { exports: { Cached: CachedStub } }).exports.Cached;
}

// The read path. Each API Worker exports a subclass named `Cached` that
// implements origin() (the container); wrangler.jsonc turns caching on for
// it and off for the gateway.
export abstract class CachedEntrypoint<E = unknown> extends WorkerEntrypoint<E> {
	protected abstract origin(request: Request): Promise<Response>;

	override async fetch(request: Request): Promise<Response> {
		const res = await this.origin(request);
		// Without Cache-Control the cache falls back to heuristic TTLs (two
		// hours for a 200). The backend always sets one; a response that
		// did not come from it (a container that failed to start) must
		// not be stored either, and neither may a stream (no-cache would
		// still store it).
		if (res.status === 101) return res;
		const stream = (res.headers.get("Content-Type") ?? "").startsWith("text/event-stream");
		if (res.headers.has("Cache-Control") && !stream) return res;
		const guarded = new Response(res.body, res);
		guarded.headers.set("Cache-Control", "private, no-store");
		return guarded;
	}

	// Purges this entrypoint's cached responses carrying any of the tags.
	// False when Cloudflare refused it (the purge rate limit), so the
	// caller can retry later.
	async purge(tags: string[]): Promise<boolean> {
		const cache = this.ctx.cache;
		if (!cache) {
			console.error("cache purge skipped: caching is off for this entrypoint", tags);
			return false;
		}
		const result = await cache.purge({ tags });
		if (!result.success) console.error("cache purge refused", tags, result.errors);
		return result.success;
	}
}

// Takes the purge list off a write's response. Returns the response the
// client gets (without the header) and the tags to purge, empty unless the
// write succeeded.
export function takePurge(res: Response): { response: Response; tags: string[] } {
	const raw = res.headers.get(PURGE_HEADER);
	if (raw === null || res.status === 101) return { response: res, tags: [] };
	const response = new Response(res.body, res);
	response.headers.delete(PURGE_HEADER);
	const tags = res.ok ? csv(raw) : [];
	return { response, tags };
}

// Purges now; when Cloudflare refuses (rate limit) or the call fails, hands
// the tags to the queue so the purge is retried with backoff. Workers
// without a queue only log.
export async function purgeOrQueue(
	ctx: ExecutionContext,
	tags: string[],
	queue?: Queue<PurgeJob>,
): Promise<void> {
	if (tags.length === 0) return;
	try {
		if (await cached(ctx).purge(tags)) return;
	} catch (err) {
		console.error("cache purge failed", tags, err);
	}
	if (queue) await queue.send({ kind: PURGE_JOB, args: { tags } }, { contentType: "json" });
	else console.error("cache purge lost: no queue to retry it", tags);
}

// The gateway's whole request path: preflights answered here, reads
// through the cache, writes to the container with their purges run after.
export async function gateway(
	request: Request,
	env: CorsEnv,
	ctx: ExecutionContext,
	origin: () => Promise<Response>,
	queue?: Queue<PurgeJob>,
): Promise<Response> {
	const answer = preflight(request, env);
	if (answer) return answer;
	if (cacheable(request)) return cached(ctx).fetch(request);
	const { response, tags } = takePurge(await origin());
	if (tags.length) ctx.waitUntil(purgeOrQueue(ctx, tags, queue));
	return response;
}

// True for a queued retry of a purge.
export function isPurgeJob(body: { kind: string }): body is PurgeJob {
	return body.kind === PURGE_JOB;
}
