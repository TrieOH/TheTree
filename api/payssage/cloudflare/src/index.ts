import { getRandom } from "@cloudflare/containers";
import { BudgetedContainer, isRoute, notFound } from "@trieoh/container-budget";
import { CachedEntrypoint, gateway } from "@trieoh/container-budget/edge";
import {
	type Job,
	backoff,
	containerEnv,
	deliver,
	enqueue,
	jobsOutbound,
	periodicJobs,
	runJobs,
} from "@trieoh/container-budget/jobs";
import { ROUTES } from "./routes";

// Required by outboundByHost: the container's outbound requests are routed
// through this entrypoint.
export { ContainerProxy } from "@cloudflare/containers";

// How many container instances share the traffic. Payssage keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 1;

// Queue names, as in wrangler.jsonc.
const DEAD_LETTER_QUEUE = "payssage-jobs-dlq";

// Failed jobs retry with backoff (30s, doubling, capped at an hour). With
// the consumer's max_retries a tenant endpoint gets a few hours to come back
// before the delivery is marked failed.

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

	static outboundByHost = jobsOutbound<Env>();

	constructor(ctx: DurableObjectState<{}>, env: Env) {
		super(ctx, env);
		this.envVars = containerEnv(env, { jobs: true });
	}
}

async function instance(env: Env) {
	return getRandom(env.PAYSSAGE, INSTANCES);
}

export default {
	async fetch(request, env, ctx) {
		const { pathname } = new URL(request.url);
		if (!isRoute(pathname, ROUTES)) return notFound();
		if (!READS.has(request.method)) markWrite(env, ctx);
		// Payssage caches only its spec and problem docs, which change per
		// deploy, so no write ever purges.
		return gateway(request, env, ctx, async () => (await instance(env)).fetch(request));
	},

	async queue(batch, env, ctx) {
		if (batch.messages.some((m) => m.body.kind !== REDRIVE_KIND)) markWrite(env, ctx);
		await runJobs(batch, ctx, {
			deliver: async (job) => deliver(await instance(env), job),
			retryDelay: backoff,
			deadLetter: { queue: DEAD_LETTER_QUEUE, onDead: ON_DEAD },
		});
	},

	async scheduled(controller, env) {
		let kinds = periodicJobs(PERIODIC, controller.cron);
		if (kinds.includes(REDRIVE_KIND)) {
			if (await redriveNeeded(env)) {
				await env.BUDGET.put(LAST_REDRIVE, String(Date.now()));
			} else {
				kinds = kinds.filter((k) => k !== REDRIVE_KIND);
			}
		}
		await enqueue(env.JOBS, kinds);
	},
} satisfies ExportedHandler<Env, Job>;
