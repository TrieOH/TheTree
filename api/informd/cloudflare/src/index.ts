import { Container, getRandom } from "@cloudflare/containers";

// How many container instances share the traffic. Informd keeps no state
// between requests that matters across instances (the in-memory rate limit
// is a per-instance backstop), so any instance can serve any request.
const INSTANCES = 3;

// Informd has no background jobs, so unlike IdentityX and Payssage there is
// no queue, no cron and no /events endpoint: the Worker only fronts the
// container.
export class InformdContainer extends Container<Env> {
	defaultPort = 8080;
	// Idle time before the container sleeps; the next request cold-starts it.
	sleepAfter = "5m";

	constructor(ctx: DurableObjectState<{}>, env: Env) {
		super(ctx, env);
		// Vars and secrets become the container's environment; a running
		// container keeps what it started with until it restarts.
		this.envVars = stringBindings(env);
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
	async fetch(request, env) {
		const container = await getRandom(env.INFORMD, INSTANCES);
		return container.fetch(request);
	},
} satisfies ExportedHandler<Env>;
