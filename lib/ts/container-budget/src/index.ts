import { Container } from "@cloudflare/containers";

// Keeps the four API containers inside the Workers Paid plan's included
// usage. Cloudflare has no spending cap, so the cut lives here: every
// container counts its awake time, weighted by its memory, into one shared
// pool, and once the pool is spent no container starts again until the
// billing cycle turns over.
//
// The pool is the plan's container memory: 25 GiB-hours a month, 100 awake
// hours of a `lite` (256 MiB) instance. Disk (200 GB-hours, 100 h at 2 GB)
// runs out at the same pace; vCPU (active only) and Durable Object duration
// run out later.

// The apps sharing the pool; each counts under its own key.
export const APPS = ["identityx", "payssage", "univents", "informd"] as const;
export type App = (typeof APPS)[number];

// GiB-seconds the pool allows per cycle: 90% of the included 25 GiB-hours,
// leaving room for counting lag and cold starts.
export const LIMIT_GIB_SECONDS = 0.9 * 25 * 3600;

// Day of the month (UTC) the billing cycle starts on.
const CYCLE_DAY = 24;

// A running container's time is written at most this often.
const FLUSH_MS = 60_000;

// How long a read of the pool's total is trusted.
const TOTAL_TTL_MS = 60_000;

const STARTED_AT = "budget:startedAt";

export type BudgetEnv = { BUDGET: KVNamespace };

// The first day of the billing cycle `now` falls in, as YYYY-MM-DD.
export function cycleStart(now: Date): string {
	let y = now.getUTCFullYear();
	let m = now.getUTCMonth();
	if (now.getUTCDate() < CYCLE_DAY) {
		m -= 1;
		if (m < 0) {
			m = 11;
			y -= 1;
		}
	}
	return `${y}-${String(m + 1).padStart(2, "0")}-${CYCLE_DAY}`;
}

function key(cycle: string, app: App): string {
	return `awake:${cycle}:${app}`;
}

// GiB-seconds the pool has spent this cycle, across every app.
export async function spentGiBSeconds(kv: KVNamespace, now = new Date()): Promise<number> {
	const cycle = cycleStart(now);
	const values = await Promise.all(APPS.map((app) => kv.get(key(cycle, app))));
	return values.reduce((sum, v) => sum + (Number(v) || 0), 0);
}

// A container that refuses to start once the shared pool is spent.
export abstract class BudgetedContainer<Env extends BudgetEnv> extends Container<Env> {
	protected abstract readonly app: App;
	// The instance type's memory, which is what Cloudflare bills.
	protected abstract readonly memoryGiB: number;

	private total?: { seconds: number; readAt: number };

	override async fetch(request: Request): Promise<Response> {
		await this.flush(false);
		if (await this.overBudget()) {
			return new Response("Paused: this month's compute budget is spent.", {
				status: 503,
				headers: { "Retry-After": "3600" },
			});
		}
		return super.fetch(request);
	}

	override async onStart(): Promise<void> {
		await this.ctx.storage.put(STARTED_AT, Date.now());
	}

	override async onStop(): Promise<void> {
		await this.flush(true);
	}

	private async overBudget(): Promise<boolean> {
		const now = Date.now();
		if (!this.total || now - this.total.readAt > TOTAL_TTL_MS) {
			this.total = { seconds: await spentGiBSeconds(this.env.BUDGET), readAt: now };
		}
		return this.total.seconds >= LIMIT_GIB_SECONDS;
	}

	// Adds the GiB-seconds since the last write to this app's count. While the
	// container runs it writes at most once a minute; on stop it writes
	// whatever is left and forgets the start.
	private async flush(stopped: boolean): Promise<void> {
		const since = await this.ctx.storage.get<number>(STARTED_AT);
		if (since === undefined) return;
		const now = Date.now();
		if (!stopped && now - since < FLUSH_MS) return;

		const spent = ((now - since) / 1000) * this.memoryGiB;
		const k = key(cycleStart(new Date(now)), this.app);
		const prev = Number(await this.env.BUDGET.get(k)) || 0;
		await this.env.BUDGET.put(k, String(prev + spent), {
			// Outlives the cycle, then cleans itself up.
			expirationTtl: 40 * 24 * 3600,
		});
		if (stopped) await this.ctx.storage.delete(STARTED_AT);
		else await this.ctx.storage.put(STARTED_AT, now);
		if (this.total) this.total.seconds += spent;
	}
}

export { isNoise, notFound } from "./noise";
