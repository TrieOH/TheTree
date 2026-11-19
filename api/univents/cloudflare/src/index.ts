import { getRandom } from "@cloudflare/containers";
import { BudgetedContainer, isRoute, notFound } from "@trieoh/container-budget";
import { CachedEntrypoint, gateway, type PurgeJob } from "@trieoh/container-budget/edge";
import {
	type Job,
	backoff,
	containerEnv,
	deliver,
	jobsOutbound,
	runJobs,
} from "@trieoh/container-budget/jobs";
import { ROUTES } from "./routes";

// Required by outboundByHost: the container's outbound requests are routed
// through this entrypoint.
export { ContainerProxy } from "@cloudflare/containers";

// How many container instances share the traffic. Univents keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 1;

// Reads, with Workers Caching in front (ADR-0020).
export class Cached extends CachedEntrypoint<Env> {
	protected async origin(request: Request): Promise<Response> {
		return (await instance(this.env)).fetch(request);
	}
}

export class UniventsContainer extends BudgetedContainer<Env> {
	protected readonly app = "univents";
	protected readonly memoryGiB = 0.25; // lite
	defaultPort = 8080;
	// Idle time before the container sleeps; the next request cold-starts it.
	sleepAfter = "2m";

	// Delayed jobs (purchase expiry) keep their delay on the queue.
	static outboundByHost = jobsOutbound<Env>();

	constructor(ctx: DurableObjectState<{}>, env: Env) {
		super(ctx, env);
		this.envVars = containerEnv(env, { jobs: true });
	}
}

async function instance(env: Env) {
	return getRandom(env.UNIVENTS, INSTANCES);
}

export default {
	async fetch(request, env, ctx) {
		const { pathname } = new URL(request.url);
		if (!isRoute(pathname, ROUTES)) return notFound();
		return gateway(request, env, ctx, async () => (await instance(env)).fetch(request), env.JOBS as Queue<PurgeJob>);
	},

	async queue(batch, env, ctx) {
		await runJobs(batch, ctx, { deliver: async (job) => deliver(await instance(env), job), retryDelay: backoff });
	},
} satisfies ExportedHandler<Env, Job>;
