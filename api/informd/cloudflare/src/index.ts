import { getRandom } from "@cloudflare/containers";
import { BudgetedContainer, isRoute, notFound } from "@trieoh/container-budget";
import { CachedEntrypoint, gateway } from "@trieoh/container-budget/edge";
import { containerEnv } from "@trieoh/container-budget/jobs";
import { ROUTES } from "./routes";

// How many container instances share the traffic. Informd keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 1;

// Informd has no background jobs, so unlike the other APIs there is no
// queue, no cron and no job endpoint: the Worker only fronts the container.
export class InformdContainer extends BudgetedContainer<Env> {
	protected readonly app = "informd";
	protected readonly memoryGiB = 0.25; // lite
	defaultPort = 8080;
	// Idle time before the container sleeps; the next request cold-starts it.
	sleepAfter = "2m";

	constructor(ctx: DurableObjectState<{}>, env: Env) {
		super(ctx, env);
		this.envVars = containerEnv(env, { jobs: false });
	}
}

// Reads, with Workers Caching in front (ADR-0020).
export class Cached extends CachedEntrypoint<Env> {
	protected async origin(request: Request): Promise<Response> {
		return (await getRandom(this.env.INFORMD, INSTANCES)).fetch(request);
	}
}

export default {
	async fetch(request, env, ctx) {
		if (!isRoute(new URL(request.url).pathname, ROUTES)) return notFound();
		// Informd caches only its spec and problem docs, which change per
		// deploy, so no write ever purges.
		return gateway(request, env, ctx, async () => (await getRandom(env.INFORMD, INSTANCES)).fetch(request));
	},
} satisfies ExportedHandler<Env>;
