import { Container, getRandom } from "@cloudflare/containers";

// Required by outboundByHost: the container's outbound requests are routed
// through this entrypoint.
export { ContainerProxy } from "@cloudflare/containers";

// How many container instances share the traffic. Payssage keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 3;

// The container enqueues by POSTing a job envelope here (JOBS_URL). The
// request never leaves Cloudflare: outboundByHost hands it to the queue.
const JOBS_HOST = "jobs.internal";

// The container's job events endpoint (lib/jobs/httpjobs.Path). Only this
// Worker may reach it.
const EVENTS_PATH = "/events";

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
	"*/15 * * * *": ["webhook.redrive"],
};

// Jobs whose retries ran out, mapped to the job that records the give-up.
// A dead job with no entry here is logged and dropped.
const ON_DEAD: Record<string, string> = {
	"webhook.deliver": "webhook.fail",
};

type Envelope = { kind: string; args?: unknown };

export class PayssageContainer extends Container<Env> {
	defaultPort = 8080;
	// Idle time before the container sleeps; the next request cold-starts it.
	sleepAfter = "5m";

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
	async fetch(request, env) {
		if (new URL(request.url).pathname.startsWith(EVENTS_PATH)) {
			return new Response("Not Found", { status: 404 });
		}
		const container = await instance(env);
		return container.fetch(request);
	},

	async queue(batch, env) {
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
		const kinds = PERIODIC[controller.cron] ?? [];
		if (kinds.length === 0) {
			console.error(`no periodic jobs for cron ${controller.cron}`);
			return;
		}
		await env.JOBS.sendBatch(kinds.map((kind) => ({ body: { kind }, contentType: "json" })));
	},
} satisfies ExportedHandler<Env, Envelope>;
