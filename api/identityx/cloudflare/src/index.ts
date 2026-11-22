import { getRandom } from "@cloudflare/containers";
import { BudgetedContainer, isRoute, notFound } from "@trieoh/container-budget";
import { CachedEntrypoint, gateway, type PurgeJob } from "@trieoh/container-budget/edge";
import {
	type Job,
	containerEnv,
	deliver,
	enqueue,
	jobsOutbound,
	jobsSecret,
	periodicJobs,
	runJobs,
} from "@trieoh/container-budget/jobs";
import { ROUTES } from "./routes";

// Required by outboundByHost: the container's outbound requests are routed
// through this entrypoint.
export { ContainerProxy } from "@cloudflare/containers";

// How many container instances share the traffic. IdentityX keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 1;

// Public signing keys. Every other API fetches them to verify tokens, so
// the Worker answers from the edge cache, then from KV (shared by every
// location), and only wakes the container when both miss. Keys only change
// while the container runs, and the container says so: whenever it creates,
// rotates or retires a key it sends a JWKS_CHANGED envelope, and the KV copy
// is dropped then. Project creation also drops it, as a backstop. Otherwise
// the copy lives as long as a key does; a boot alone changes nothing.
const JWKS_PATH = "/.well-known/jwks.json";
// Matches JWKSChangedArgs.Kind() in internal/jobs.
const JWKS_CHANGED = "jwks.changed";
const JWKS_EDGE_TTL_SECONDS = 300;
const JWKS_KV_TTL_SECONDS = 7 * 24 * 3600;
const JWKS_KV_PREFIX = "jwks:";

// Routes that create a project, and with it the project's keys.
const CREATES_PROJECT = /^(\/organizations\/[^/]+)?\/projects$/;

// Periodic jobs by cron, mirroring jobs.Periodic in internal/jobs/jobs.go.
// Keep these crons and wrangler.jsonc's triggers in step with it.
const PERIODIC: Record<string, string[]> = {
	// One daily run, so the cron wakes the container once a day.
	"0 3 * * *": ["cleanup_blacklist", "cleanup_action_tokens", "keys.rotate"],
};

// Reads, with Workers Caching in front (ADR-0020).
export class Cached extends CachedEntrypoint<Env> {
	protected async origin(request: Request): Promise<Response> {
		return (await instance(this.env)).fetch(request);
	}
}

export class IdentityXContainer extends BudgetedContainer<Env> {
	protected readonly app = "identityx";
	protected readonly memoryGiB = 0.25; // lite
	defaultPort = 8080;
	// Idle time before the container sleeps; the next request cold-starts it.
	sleepAfter = "2m";

	// JWKS_CHANGED is not work for the queue: the keys changed, so the
	// cached JWKS is dropped right away.
	// Assigned in a static block, not declared as a static field: the
	// containers package registers outbound handlers through the inherited
	// `outboundByHost` setter, and a class field defines its own property
	// instead, so the handlers never registered and every enqueue went to
	// the internet (530).
	static {
		this.outboundByHost = jobsOutbound<Env>(async (job, env) => {
			if (job.kind !== JWKS_CHANGED) return false;
			await forgetJWKS(env);
			return true;
		});
	}

	constructor(ctx: DurableObjectState<{}>, env: Env) {
		super(ctx, env);
		this.envVars = containerEnv(env, { jobs: true });
	}
}

async function instance(env: Env) {
	return getRandom(env.IDENTITYX, INSTANCES);
}

async function jwks(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
	const cache = caches.default;
	const cacheKey = new Request(request.url, { method: "GET" });
	const hit = await cache.match(cacheKey);
	if (hit) return hit;

	const kvKey = JWKS_KV_PREFIX + new URL(request.url).search;
	let body = await env.BUDGET.get(kvKey);
	if (body === null) {
		// Fill KV even when the caller gives up during a slow cold start, so
		// the next cold start elsewhere does not wake the container again.
		const fill = (async () => {
			const res = await (await instance(env)).fetch(request);
			if (!res.ok) return res;
			const text = await res.text();
			await env.BUDGET.put(kvKey, text, { expirationTtl: JWKS_KV_TTL_SECONDS });
			return text;
		})();
		ctx.waitUntil(fill);
		const filled = await fill;
		if (filled instanceof Response) return filled;
		body = filled;
	}
	const res = new Response(body, {
		headers: {
			"Content-Type": "application/json",
			"Cache-Control": `public, max-age=${JWKS_EDGE_TTL_SECONDS}`,
		},
	});
	ctx.waitUntil(cache.put(cacheKey, res.clone()));
	return res;
}

// Drops every stored JWKS copy, so the next fetch reads the rotated keys.
async function forgetJWKS(env: Env): Promise<void> {
	const { keys } = await env.BUDGET.list({ prefix: JWKS_KV_PREFIX });
	await Promise.all(keys.map((k) => env.BUDGET.delete(k.name)));
}

export default {
	async fetch(request, env, ctx) {
		const { pathname } = new URL(request.url);
		if (!isRoute(pathname, ROUTES)) return notFound();
		if (pathname === JWKS_PATH && request.method === "GET") {
			return jwks(request, env, ctx);
		}
		if (request.method === "POST" && CREATES_PROJECT.test(pathname)) {
			ctx.waitUntil(forgetJWKS(env));
		}
		return gateway(request, env, ctx, async () => (await instance(env)).fetch(request), env.JOBS as Queue<PurgeJob>);
	},

	async queue(batch, env, ctx) {
		await runJobs(batch, ctx, { deliver: async (job, delivery) => deliver(await instance(env), job, delivery, jobsSecret(env)) });
	},

	async scheduled(controller, env) {
		await enqueue(env.JOBS, periodicJobs(PERIODIC, controller.cron));
	},
} satisfies ExportedHandler<Env, Job>;
