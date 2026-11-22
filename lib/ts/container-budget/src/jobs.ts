import { PURGE_RETRY_SECONDS, cached, isPurgeJob } from "./edge";

// The Worker half of background jobs (ADR-0011): the container enqueues
// through the Worker, the queue hands each job back to the Worker, and the
// Worker delivers it to the container's job endpoint. The Worker is only a
// carrier: retries, backoff and give-ups are decided by the container
// (lib/jobs.Registry.Handle) and come back as an Outcome. Every API Worker
// with jobs uses these; what's left in each Worker is what only that API does.

// The container enqueues by POSTing a job here (JOBS_URL). The request never
// leaves Cloudflare: the outbound handler hands it to the queue.
export const JOBS_HOST = "jobs.internal";

// The container's job endpoint (lib/jobs/httpjobs.Path). Only the Worker
// may reach it.
export const JOBS_PATH = "/_jobs";

// A delayed job carries its delay in this header, in whole seconds
// (lib/jobs/httpjobs.DelayHeader).
export const DELAY_HEADER = "Job-Delay-Seconds";

// A delivery to the container's job endpoint carries the jobs secret in this
// header (lib/jobs/httpjobs.SecretHeader).
export const SECRET_HEADER = "Jobs-Secret";

// The API's JOBS_SECRET, when it's been set (`wrangler secret put
// JOBS_SECRET`). It reaches the container with every other secret
// (containerEnv), which then refuses deliveries without it. Optional, so a
// deploy before it's set still works.
export function jobsSecret(env: object): string | undefined {
	const secret = (env as { JOBS_SECRET?: unknown }).JOBS_SECRET;
	return typeof secret === "string" && secret !== "" ? secret : undefined;
}

// trace is the W3C trace context of whatever enqueued the job
// (lib/jobs.Envelope.Trace). The Worker passes it through untouched.
export type Job = { kind: string; args?: unknown; trace?: Record<string, string> };
export type JobsEnv = { JOBS: Queue };

// True for the job endpoint, which the public must never reach.
export function isJobsPath(pathname: string): boolean {
	return pathname === JOBS_PATH || pathname.startsWith(`${JOBS_PATH}/`);
}

// The container's environment: every string var and secret, plus JOBS_URL
// when the API has jobs. A running container keeps what it started with
// until it restarts.
export function containerEnv(env: object, { jobs }: { jobs: boolean }): Record<string, string> {
	const vars: Record<string, string> = {};
	for (const [key, value] of Object.entries(env)) {
		if (typeof value === "string") vars[key] = value;
	}
	if (jobs) vars.JOBS_URL = `http://${JOBS_HOST}/enqueue`;
	return vars;
}

// The container's outboundByHost entry for enqueueing. intercept runs first
// and returns true for a job it handled itself, which then skips the queue.
export function jobsOutbound<E extends JobsEnv>(intercept?: (job: Job, env: E) => Promise<boolean>) {
	return {
		[JOBS_HOST]: async (req: Request, env: E) => {
			const job = await req.json<Job>();
			if (intercept && (await intercept(job, env))) return new Response(null, { status: 202 });
			const delaySeconds = Number(req.headers.get(DELAY_HEADER) ?? 0) || 0;
			await env.JOBS.send(job, { contentType: "json", delaySeconds });
			return new Response(null, { status: 202 });
		},
	};
}

// A delivery's attempt (1 on the first), and "true" on one from the
// dead-letter queue (lib/jobs/httpjobs.AttemptHeader, DeadHeader).
export const ATTEMPT_HEADER = "Job-Attempt";
export const DEAD_HEADER = "Job-Dead";

// What the container decided for a delivered job (lib/jobs/httpjobs.Outcome).
// The Worker applies it and decides nothing about the job itself.
export type Outcome = { outcome: "done" | "retry" | "drop"; retry_after_seconds?: number };

// Retry delay when the container never answered with an Outcome (it was
// unreachable or crashed): the carrier's own concern, not the job's.
export const UNREACHABLE_RETRY_SECONDS = 30;

// One delivery of a job to the container.
export type Delivery = { attempt: number; dead: boolean };

// Delivers one job to a container, with the jobs secret when there is one
// (jobsSecret), and returns the container's Outcome. A 4xx (a body that
// isn't a job, a wrong secret) can never succeed and is a drop; anything
// else that isn't an Outcome is a retry.
export async function deliver(
	container: Fetcher | { fetch(req: Request): Promise<Response> },
	job: Job,
	delivery: Delivery,
	secret?: string,
): Promise<Outcome> {
	const headers: Record<string, string> = {
		"Content-Type": "application/json",
		[ATTEMPT_HEADER]: String(delivery.attempt),
	};
	if (delivery.dead) headers[DEAD_HEADER] = "true";
	if (secret) headers[SECRET_HEADER] = secret;
	const res = await container.fetch(
		new Request(`http://container${JOBS_PATH}`, {
			method: "POST",
			headers,
			body: JSON.stringify(job),
		}),
	);
	if (res.status >= 400 && res.status < 500) {
		console.error(`job ${job.kind} rejected: ${res.status} ${await res.text()}`);
		return { outcome: "drop" };
	}
	if (res.ok && res.headers.get("Content-Type")?.startsWith("application/json")) {
		return res.json<Outcome>();
	}
	// A container from before Outcomes answers 204 for done.
	if (res.status === 204) return { outcome: "done" };
	return { outcome: "retry", retry_after_seconds: UNREACHABLE_RETRY_SECONDS };
}

export type JobsOptions = {
	// Delivers one job; see deliver.
	deliver: (job: Job, delivery: Delivery) => Promise<Outcome>;
	// The dead-letter queue: jobs the queue gave up redelivering arrive
	// here, and are delivered dead so the container runs their give-up.
	deadLetterQueue?: string;
};

// The queue handler: cache purges (edge work) run here; every other job is
// delivered to the container, and its Outcome applied to the message.
export async function runJobs(batch: MessageBatch<Job>, ctx: ExecutionContext, opts: JobsOptions): Promise<void> {
	const dead = batch.queue === opts.deadLetterQueue;
	for (const message of batch.messages) {
		if (isPurgeJob(message.body)) {
			if (await cached(ctx).purge(message.body.args.tags)) message.ack();
			else message.retry({ delaySeconds: PURGE_RETRY_SECONDS });
			continue;
		}
		let out: Outcome;
		try {
			out = await opts.deliver(message.body, { attempt: message.attempts, dead });
		} catch (err) {
			console.error("job delivery failed", err);
			out = { outcome: "retry", retry_after_seconds: UNREACHABLE_RETRY_SECONDS };
		}
		if (out.outcome === "retry") message.retry({ delaySeconds: out.retry_after_seconds ?? UNREACHABLE_RETRY_SECONDS });
		else message.ack();
	}
}

// The jobs a cron runs. Each Worker's table mirrors jobs.Periodic in its
// internal/jobs/jobs.go; keep it and wrangler.jsonc's triggers in step.
export function periodicJobs(periodic: Record<string, string[]>, cron: string): string[] {
	const kinds = periodic[cron] ?? [];
	if (kinds.length === 0) console.error(`no periodic jobs for cron ${cron}`);
	return kinds;
}

// Enqueues argument-less jobs, as crons start them.
export async function enqueue(queue: Queue, kinds: string[]): Promise<void> {
	if (kinds.length === 0) return;
	await queue.sendBatch(kinds.map((kind) => ({ body: { kind }, contentType: "json" as const })));
}
