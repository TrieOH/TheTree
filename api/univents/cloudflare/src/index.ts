import { Container, getRandom } from "@cloudflare/containers";

// Required by outboundByHost: the container's outbound requests are routed
// through this entrypoint.
export { ContainerProxy } from "@cloudflare/containers";

// How many container instances share the traffic. Univents keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 3;

// The container enqueues by POSTing a job envelope here (JOBS_URL). The
// request never leaves Cloudflare: outboundByHost hands it to the queue.
const JOBS_HOST = "jobs.internal";

// A delayed job carries its delay in this header, in whole seconds
// (lib/jobs/httpjobs.DelayHeader).
const DELAY_HEADER = "Job-Delay-Seconds";

// The container's job events endpoint (lib/jobs/httpjobs.Path). Only this
// Worker may reach it.
const EVENTS_PATH = "/events";

// Retry backoff: 30s, doubling per attempt, capped at an hour.
const RETRY_BASE_SECONDS = 30;
const RETRY_MAX_SECONDS = 60 * 60;

function retryDelay(attempts: number): number {
	return Math.min(RETRY_BASE_SECONDS * 2 ** (attempts - 1), RETRY_MAX_SECONDS);
}

type Envelope = { kind: string; args?: unknown };

export class UniventsContainer extends Container<Env> {
	defaultPort = 8080;
	// Idle time before the container sleeps; the next request cold-starts it.
	sleepAfter = "5m";

	static outboundByHost = {
		[JOBS_HOST]: async (req: Request, env: Env) => {
			const envelope = await req.json<Envelope>();
			const delaySeconds = Number(req.headers.get(DELAY_HEADER) ?? 0) || 0;
			await env.JOBS.send(envelope, { contentType: "json", delaySeconds });
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
	return getRandom(env.UNIVENTS, INSTANCES);
}

// Delivers one job to a container. True when the job is done or can never
// succeed (4xx); false when it should be retried.
async function deliver(env: Env, envelope: Envelope): Promise<boolean> {
	const container = await instance(env);
	const res = await container.fetch(
		new Request(`http://univents${EVENTS_PATH}`, {
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
	async fetch(request, env) {
		if (new URL(request.url).pathname.startsWith(EVENTS_PATH)) {
			return new Response("Not Found", { status: 404 });
		}
		const container = await instance(env);
		return container.fetch(request);
	},

	async queue(batch, env) {
		for (const message of batch.messages) {
			const retry = { delaySeconds: retryDelay(message.attempts) };
			try {
				if (await deliver(env, message.body)) message.ack();
				else message.retry(retry);
			} catch (err) {
				console.error("job delivery failed", err);
				message.retry(retry);
			}
		}
	},
} satisfies ExportedHandler<Env, Envelope>;
