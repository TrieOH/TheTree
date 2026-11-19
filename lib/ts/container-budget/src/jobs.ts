import { PURGE_RETRY_SECONDS, cached, isPurgeJob } from "./edge";

// The Worker half of background jobs (ADR-0011): the container enqueues
// through the Worker, the queue hands each job back to the Worker, and the
// Worker delivers it to the container's job endpoint. Every API Worker with
// jobs uses these; what's left in each Worker is what only that API does.

// The container enqueues by POSTing a job here (JOBS_URL). The request never
// leaves Cloudflare: the outbound handler hands it to the queue.
export const JOBS_HOST = "jobs.internal";

// The container's job endpoint (lib/jobs/httpjobs.Path). Only the Worker
// may reach it.
export const JOBS_PATH = "/_jobs";

// A delayed job carries its delay in this header, in whole seconds
// (lib/jobs/httpjobs.DelayHeader).
export const DELAY_HEADER = "Job-Delay-Seconds";

export type Job = { kind: string; args?: unknown };
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

// Retry backoff: 30s, doubling per attempt, capped at an hour.
export function backoff(attempts: number): number {
	return Math.min(30 * 2 ** (attempts - 1), 60 * 60);
}

// Delivers one job to a container. True when the job is done or can never
// succeed (4xx); false when it should be retried.
export async function deliver(container: Fetcher | { fetch(req: Request): Promise<Response> }, job: Job): Promise<boolean> {
	const res = await container.fetch(
		new Request(`http://container${JOBS_PATH}`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(job),
		}),
	);
	if (res.status >= 400 && res.status < 500) {
		console.error(`job ${job.kind} rejected: ${res.status} ${await res.text()}`);
		return true;
	}
	return res.ok;
}

export type JobsOptions = {
	// Delivers one job; see deliver.
	deliver: (job: Job) => Promise<boolean>;
	// Seconds before a failed job is retried. Unset, the queue's own retry
	// delay applies.
	retryDelay?: (attempts: number) => number;
	// Jobs whose retries ran out arrive on this dead-letter queue, and
	// onDead maps each to the job that records the give-up. A dead job with
	// no entry is logged and dropped.
	deadLetter?: { queue: string; onDead: Record<string, string> };
};

// The queue handler: cache purges retried, dead jobs turned into their
// give-up job, everything else delivered to the container.
export async function runJobs(batch: MessageBatch<Job>, ctx: ExecutionContext, opts: JobsOptions): Promise<void> {
	for (const message of batch.messages) {
		if (isPurgeJob(message.body)) {
			if (await cached(ctx).purge(message.body.args.tags)) message.ack();
			else message.retry({ delaySeconds: PURGE_RETRY_SECONDS });
			continue;
		}
		const job = jobFor(batch.queue, message.body, opts.deadLetter);
		if (!job) {
			message.ack();
			continue;
		}
		const retry = opts.retryDelay ? { delaySeconds: opts.retryDelay(message.attempts) } : undefined;
		try {
			if (await opts.deliver(job)) message.ack();
			else message.retry(retry);
		} catch (err) {
			console.error("job delivery failed", err);
			message.retry(retry);
		}
	}
}

function jobFor(queue: string, job: Job, deadLetter: JobsOptions["deadLetter"]): Job | null {
	if (!deadLetter || queue !== deadLetter.queue) return job;
	const kind = deadLetter.onDead[job.kind];
	if (!kind) {
		console.error(`dead job ${job.kind} dropped: no give-up job`);
		return null;
	}
	return { kind, args: job.args };
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
