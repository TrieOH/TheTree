import { getRandom } from "@cloudflare/containers";
import { BudgetedContainer } from "@trieoh/container-budget";

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

// Public signing keys, cached at the edge for this long. Short, so a
// rotated-in key is served within minutes.
const JWKS_PATH = "/.well-known/jwks.json";
const JWKS_TTL_SECONDS = 300;

// Periodic jobs by cron, mirroring jobs.Periodic in internal/jobs/jobs.go.
// Keep these crons and wrangler.jsonc's triggers in step with it.
const PERIODIC: Record<string, string[]> = {
	"0 * * * *": ["cleanup_blacklist", "cleanup_action_tokens"],
	"0 3 * * *": ["keys.rotate"],
};

type Envelope = { kind: string; args?: unknown };

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

export default {
	async fetch(request, env, ctx) {
		const { pathname } = new URL(request.url);
		if (pathname === EVENTS_PATH || pathname.startsWith(`${EVENTS_PATH}/`)) {
			return new Response("Not Found", { status: 404 });
		}
		// Other services fetch the JWKS on every cold start; serving it from
		// the edge cache keeps those fetches from waking the container.
		if (pathname === JWKS_PATH && request.method === "GET") {
			const cache = caches.default;
			const key = new Request(request.url, { method: "GET" });
			const hit = await cache.match(key);
			if (hit) return hit;
			const res = await (await instance(env)).fetch(request);
			if (res.ok) {
				const cached = new Response(res.body, res);
				cached.headers.set("Cache-Control", `public, max-age=${JWKS_TTL_SECONDS}`);
				ctx.waitUntil(cache.put(key, cached.clone()));
				return cached;
			}
			return res;
		}
		const container = await instance(env);
		return container.fetch(request);
	},

	async queue(batch, env) {
		for (const message of batch.messages) {
			try {
				if (await deliver(env, message.body)) message.ack();
				else message.retry();
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
