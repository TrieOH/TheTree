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

// A container with no new request for this long is stopped, even when the
// library still counts a request in flight: a request whose caller gave up
// (an edge timeout during a slow cold start) never has its response read,
// and the library keeps the container awake for it forever.
const IDLE_LIMIT_MS = 5 * 60_000;

// How often a running container checks that it is still in use.
const WATCHDOG_MS = 60_000;

const STARTED_AT = "budget:startedAt";
const LAST_REQUEST = "budget:lastRequest";

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
		// Not running with a start on record: it stopped and the library has
		// not reported it yet. Close that run before starting a new one.
		await this.flush(!this.ctx.container?.running);
		if (await this.overBudget()) {
			return new Response("Paused: this month's compute budget is spent.", {
				status: 503,
				headers: { "Retry-After": "3600" },
			});
		}
		await this.ctx.storage.put(LAST_REQUEST, Date.now());
		return super.fetch(request);
	}

	override async onStart(): Promise<void> {
		const now = Date.now();
		await this.ctx.storage.put({ [STARTED_AT]: now, [LAST_REQUEST]: now });
		await this.schedule(new Date(now + WATCHDOG_MS), "watchdog");
	}

	// Stops a running container once it has gone IDLE_LIMIT_MS without a new
	// request; otherwise checks again in WATCHDOG_MS.
	async watchdog(): Promise<void> {
		if (!this.ctx.container?.running) return;
		const last = (await this.ctx.storage.get<number>(LAST_REQUEST)) ?? 0;
		if (Date.now() - last > IDLE_LIMIT_MS) {
			await this.stop();
			return;
		}
		await this.schedule(new Date(Date.now() + WATCHDOG_MS), "watchdog");
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

		// The library can report a stop late (on the next wake), so a stop
		// never counts past the latest moment the watchdog allows.
		let end = now;
		if (stopped) {
			const last = (await this.ctx.storage.get<number>(LAST_REQUEST)) ?? since;
			end = Math.min(now, Math.max(since, last) + IDLE_LIMIT_MS + WATCHDOG_MS);
		}
		const spent = (Math.max(0, end - since) / 1000) * this.memoryGiB;
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
