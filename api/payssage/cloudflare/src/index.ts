import { getRandom } from "@cloudflare/containers";
import { BudgetedContainer, isNoise, notFound } from "@trieoh/container-budget";
import { CachedEntrypoint, gateway } from "@trieoh/container-budget/edge";

// Required by outboundByHost: the container's outbound requests are routed
// through this entrypoint.
export { ContainerProxy } from "@cloudflare/containers";

// How many container instances share the traffic. Payssage keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 1;

// The container enqueues by POSTing a job envelope here (JOBS_URL). The
// request never leaves Cloudflare: outboundByHost hands it to the queue.
const JOBS_HOST = "jobs.internal";

// The container's job endpoint (lib/jobs/httpjobs.Path). Only this Worker
// may reach it.
const EVENTS_PATH = "/_jobs";

// Queue names, as in wrangler.jsonc.
const DEAD_LETTER_QUEUE = "payssage-jobs-dlq";

// Retry backoff: 30s, doubling per attempt, capped at an hour. With the
// consumer's max_retries a tenant endpoint gets a few hours to come back
// before the delivery is marked failed.
const RETRY_BASE_SECONDS = 30;
const RETRY_MAX_SECONDS = 60 * 60;

function retryDelay(attempts: number): number {
	return Math.min(RETRY_BASE_SECONDS * 2 ** (attempts - 1), RETRY_MAX_SECONDS);
}

// Periodic jobs by cron, mirroring jobs.Periodic in internal/jobs/jobs.go.
// Keep these crons and wrangler.jsonc's triggers in step with it.
const PERIODIC: Record<string, string[]> = {
	"0 */6 * * *": ["webhook.redrive"],
};

// Jobs whose retries ran out, mapped to the job that records the give-up.
// A dead job with no entry here is logged and dropped.
const ON_DEAD: Record<string, string> = {
	"webhook.deliver": "webhook.fail",
};

// webhook.redrive only finds deliveries whose enqueue was lost, and a
// delivery is only written while the container handles a write request or
// a job. The Worker records when it last passed one on, and the cron skips
// the redrive (and the wake) when nothing happened since the previous one.
// The margin covers the sweep's own 2-minute redriveAfter and KV lag.
const LAST_WRITE = "payssage:lastWrite";
const LAST_REDRIVE = "payssage:lastRedrive";
const REDRIVE_MARGIN_MS = 15 * 60_000;
const REDRIVE_KIND = "webhook.redrive";

// A write is recorded at most this often per isolate.
const MARK_EVERY_MS = 60_000;
let lastMark = 0;

function markWrite(env: Env, ctx: ExecutionContext): void {
	const now = Date.now();
	if (now - lastMark < MARK_EVERY_MS) return;
	lastMark = now;
	ctx.waitUntil(env.BUDGET.put(LAST_WRITE, String(now)).catch((err) => console.error("mark write", err)));
}

// True when a write since the previous redrive could have lost an enqueue.
async function redriveNeeded(env: Env): Promise<boolean> {
	const [write, redrive] = await Promise.all([env.BUDGET.get(LAST_WRITE), env.BUDGET.get(LAST_REDRIVE)]);
	if (redrive === null) return true;
	return Number(write) >= Number(redrive) - REDRIVE_MARGIN_MS;
}

const READS = new Set(["GET", "HEAD", "OPTIONS"]);

type Envelope = { kind: string; args?: unknown };

// Reads, with Workers Caching in front (ADR-0020).
export class Cached extends CachedEntrypoint<Env> {
	protected async origin(request: Request): Promise<Response> {
		return (await instance(this.env)).fetch(request);
	}
}

export class PayssageContainer extends BudgetedContainer<Env> {
	protected readonly app = "payssage";
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
	return getRandom(env.PAYSSAGE, INSTANCES);
}

// Delivers one job to a container. True when the job is done or can never
// succeed (4xx); false when it should be retried.
async function deliver(env: Env, envelope: Envelope): Promise<boolean> {
	const container = await instance(env);
	const res = await container.fetch(
		new Request(`http://payssage${EVENTS_PATH}`, {
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

// The job to run for a message: itself on the jobs queue; on the dead-letter
// queue, the job that records the give-up, or nothing.
function jobFor(queue: string, envelope: Envelope): Envelope | null {
	if (queue !== DEAD_LETTER_QUEUE) return envelope;
	const kind = ON_DEAD[envelope.kind];
	if (!kind) {
		console.error(`dead job ${envelope.kind} dropped: no give-up job`);
		return null;
	}
	return { kind, args: envelope.args };
}

export default {
	async fetch(request, env, ctx) {
		const { pathname } = new URL(request.url);
		if (pathname === EVENTS_PATH || pathname.startsWith(`${EVENTS_PATH}/`)) {
			return new Response("Not Found", { status: 404 });
		}
		if (isNoise(pathname)) return notFound();
		if (!READS.has(request.method)) markWrite(env, ctx);
		// Payssage caches only its spec and problem docs, which change per
		// deploy, so no write ever purges.
		return gateway(request, env, ctx, async () => (await instance(env)).fetch(request));
	},

	async queue(batch, env, ctx) {
		if (batch.messages.some((m) => m.body.kind !== REDRIVE_KIND)) markWrite(env, ctx);
		for (const message of batch.messages) {
			const job = jobFor(batch.queue, message.body);
			if (!job) {
				message.ack();
				continue;
			}
			const retry = { delaySeconds: retryDelay(message.attempts) };
			try {
				if (await deliver(env, job)) message.ack();
				else message.retry(retry);
			} catch (err) {
				console.error("job delivery failed", err);
				message.retry(retry);
			}
		}
	},

	async scheduled(controller, env) {
		let kinds = PERIODIC[controller.cron] ?? [];
		if (kinds.length === 0) {
			console.error(`no periodic jobs for cron ${controller.cron}`);
			return;
		}
		if (kinds.includes(REDRIVE_KIND)) {
			if (await redriveNeeded(env)) {
				await env.BUDGET.put(LAST_REDRIVE, String(Date.now()));
			} else {
				kinds = kinds.filter((k) => k !== REDRIVE_KIND);
			}
		}
		if (kinds.length === 0) return;
		await env.JOBS.sendBatch(kinds.map((kind) => ({ body: { kind }, contentType: "json" })));
	},
} satisfies ExportedHandler<Env, Envelope>;
