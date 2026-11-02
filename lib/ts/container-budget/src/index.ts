import { Container } from "@cloudflare/containers";

// Keeps the four API containers inside the Workers Paid plan's included
// usage. Cloudflare has no spending cap, so the cut lives here: every
// container counts its awake time into one shared pool, and once the pool
// is spent no container starts again until the billing cycle turns over.
//
// The binding limit is container memory: 25 GiB-hours a month, which a
// `basic` (1 GiB) instance spends in 25 awake hours. Disk (200 GB-hours,
// 50 h at 4 GB), vCPU (active only) and Durable Object duration all run
// out later.

// The apps sharing the pool; each counts under its own key.
export const APPS = ["identityx", "payssage", "univents", "informd"] as const;
export type App = (typeof APPS)[number];

// Awake seconds the pool allows per cycle: 90% of the included 25 GiB-hours
// at 1 GiB per instance, leaving room for counting lag and cold starts.
export const LIMIT_SECONDS = 0.9 * 25 * 3600;

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

// Total awake seconds the pool has spent this cycle, across every app.
export async function spentSeconds(kv: KVNamespace, now = new Date()): Promise<number> {
	const cycle = cycleStart(now);
	const values = await Promise.all(APPS.map((app) => kv.get(key(cycle, app))));
	return values.reduce((sum, v) => sum + (Number(v) || 0), 0);
}

// A container that refuses to start once the shared pool is spent.
export abstract class BudgetedContainer<Env extends BudgetEnv> extends Container<Env> {
	protected abstract readonly app: App;

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
			this.total = { seconds: await spentSeconds(this.env.BUDGET), readAt: now };
		}
		return this.total.seconds >= LIMIT_SECONDS;
	}

	// Adds the time since the last write to this app's count. While the
	// container runs it writes at most once a minute; on stop it writes
	// whatever is left and forgets the start.
	private async flush(stopped: boolean): Promise<void> {
		const since = await this.ctx.storage.get<number>(STARTED_AT);
		if (since === undefined) return;
		const now = Date.now();
		if (!stopped && now - since < FLUSH_MS) return;

		const k = key(cycleStart(new Date(now)), this.app);
		const prev = Number(await this.env.BUDGET.get(k)) || 0;
		await this.env.BUDGET.put(k, String(prev + (now - since) / 1000), {
			// Outlives the cycle, then cleans itself up.
			expirationTtl: 40 * 24 * 3600,
		});
		if (stopped) await this.ctx.storage.delete(STARTED_AT);
		else await this.ctx.storage.put(STARTED_AT, now);
		if (this.total) this.total.seconds += (now - since) / 1000;
	}
}
