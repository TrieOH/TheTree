import { getRandom } from "@cloudflare/containers";
import { BudgetedContainer, isRoute, notFound } from "@trieoh/container-budget";
import { CachedEntrypoint, gateway } from "@trieoh/container-budget/edge";
import { ROUTES } from "./routes";

// How many container instances share the traffic. Informd keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 1;

// Informd has no background jobs, so unlike IdentityX and Payssage there is
// no queue, no cron and no /events endpoint: the Worker only fronts the
// container.
export class InformdContainer extends BudgetedContainer<Env> {
	protected readonly app = "informd";
	protected readonly memoryGiB = 0.25; // lite
	defaultPort = 8080;
	// Idle time before the container sleeps; the next request cold-starts it.
	sleepAfter = "2m";

	constructor(ctx: DurableObjectState<{}>, env: Env) {
		super(ctx, env);
		// Vars and secrets become the container's environment; a running
		// container keeps what it started with until it restarts.
		this.envVars = stringBindings(env);
	}
}

// Reads, with Workers Caching in front (ADR-0020).
export class Cached extends CachedEntrypoint<Env> {
	protected async origin(request: Request): Promise<Response> {
		return (await getRandom(this.env.INFORMD, INSTANCES)).fetch(request);
	}
}

function stringBindings(env: Env): Record<string, string> {
	const vars: Record<string, string> = {};
	for (const [key, value] of Object.entries(env)) {
		if (typeof value === "string") vars[key] = value;
	}
	return vars;
}

export default {
	async fetch(request, env, ctx) {
		if (!isRoute(new URL(request.url).pathname, ROUTES)) return notFound();
		// Informd caches only its spec and problem docs, which change per
		// deploy, so no write ever purges.
		return gateway(request, env, ctx, async () => (await getRandom(env.INFORMD, INSTANCES)).fetch(request));
	},
} satisfies ExportedHandler<Env>;
