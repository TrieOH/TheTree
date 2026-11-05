import { getRandom } from "@cloudflare/containers";
import { BudgetedContainer, isNoise, notFound } from "@trieoh/container-budget";
import {
	CachedEntrypoint,
	PURGE_RETRY_SECONDS,
	cached,
	gateway,
	isPurgeJob,
	type PurgeJob,
} from "@trieoh/container-budget/edge";

// Required by outboundByHost: the container's outbound requests are routed
// through this entrypoint.
export { ContainerProxy } from "@cloudflare/containers";

// How many container instances share the traffic. IdentityX keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 1;

// The container enqueues by POSTing a job envelope here (JOBS_URL). The
// request never leaves Cloudflare: outboundByHost hands it to the queue.
const JOBS_HOST = "jobs.internal";

// The container's job endpoint (lib/jobs/httpjobs.Path). Only this Worker
// may reach it.
const EVENTS_PATH = "/_jobs";

// Public signing keys. Every other API fetches them on its cold start, so
// the Worker answers from the edge cache, then from KV (shared by every
// location), and only wakes the container when both miss. Keys only change
// while the container runs: on boot (EnsureAll), in the daily keys.rotate
// job, and when a project is created. The KV copy is dropped at each of
// those, so it can otherwise live as long as a key does.
const JWKS_PATH = "/.well-known/jwks.json";
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

type Envelope = { kind: string; args?: unknown };

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

	static outboundByHost = {
		[JOBS_HOST]: async (req: Request, env: Env) => {
			const envelope = await req.json<Envelope>();
			await env.JOBS.send(envelope, { contentType: "json" });
			return new Response(null, { status: 202 });
		},
	};

	constructor(ctx: DurableObjectState<{}>, env: Env) {
		super(ctx, env);
		// Vars and secrets become the container's environment; a running
		// container keeps what it started with until it restarts.
		this.envVars = { ...stringBindings(env), JOBS_URL: `http://${JOBS_HOST}/enqueue` };
	}

	// Boot runs EnsureAll, which may rotate keys.
	override async onStart(): Promise<void> {
		await super.onStart();
		await forgetJWKS(this.env);
	}
}

function stringBindings(env: Env): Record<string, string> {
	const vars: Record<string, string> = {};
	for (const [key, value] of Object.entries(env)) {
		if (typeof value === "string") vars[key] = value;
	}
	return vars;
}

async function instance(env: Env) {
	return getRandom(env.IDENTITYX, INSTANCES);
}

// Delivers one job to a container. True when the job is done or can never
// succeed (4xx); false when it should be retried.
async function deliver(env: Env, envelope: Envelope): Promise<boolean> {
	const container = await instance(env);
	const res = await container.fetch(
		new Request(`http://identityx${EVENTS_PATH}`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(envelope),
		}),
	);
	if (res.status >= 400 && res.status < 500) {
		console.error(`job ${envelope.kind} rejected: ${res.status} ${await res.text()}`);
		return true;
	}
	return res.ok;
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
		if (pathname === EVENTS_PATH || pathname.startsWith(`${EVENTS_PATH}/`)) {
			return new Response("Not Found", { status: 404 });
		}
		if (isNoise(pathname)) return notFound();
		if (pathname === JWKS_PATH && request.method === "GET") {
			return jwks(request, env, ctx);
		}
		if (request.method === "POST" && CREATES_PROJECT.test(pathname)) {
			ctx.waitUntil(forgetJWKS(env));
		}
		return gateway(request, env, ctx, async () => (await instance(env)).fetch(request), env.JOBS as Queue<PurgeJob>);
	},

	async queue(batch, env, ctx) {
		for (const message of batch.messages) {
			if (isPurgeJob(message.body)) {
				if (await cached(ctx).purge(message.body.args.tags)) message.ack();
				else message.retry({ delaySeconds: PURGE_RETRY_SECONDS });
				continue;
			}
			try {
				if (await deliver(env, message.body)) {
					if (message.body.kind === "keys.rotate") await forgetJWKS(env);
					message.ack();
				} else message.retry();
			} catch (err) {
				console.error("job delivery failed", err);
				message.retry();
			}
		}
	},

	async scheduled(controller, env) {
		const kinds = PERIODIC[controller.cron] ?? [];
		if (kinds.length === 0) {
			console.error(`no periodic jobs for cron ${controller.cron}`);
			return;
		}
		await env.JOBS.sendBatch(kinds.map((kind) => ({ body: { kind }, contentType: "json" })));
	},
} satisfies ExportedHandler<Env, Envelope>;
