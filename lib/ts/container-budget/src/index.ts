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

// COLD_STARTS is optional: an Analytics Engine dataset that gets one data
// point per request that had to start the container. The dev portal reads it.
export type BudgetEnv = { BUDGET: KVNamespace; COLD_STARTS?: AnalyticsEngineDataset };

// The dataset the API Workers bind as COLD_STARTS.
export const COLD_START_DATASET = "container_cold_starts";

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

	// The start in progress: when its first request arrived, and when the
	// port answered (onStart). Requests that arrive during the same start
	// wait on it too but do not count it again.
	private boot?: { since: number; readyAt?: number };

	override async fetch(request: Request): Promise<Response> {
		// Asleep when the request arrived: this request pays for the start,
		// unless another request already started it.
		const asleep = !this.ctx.container?.running;
		const cold = asleep && this.boot === undefined;
		// Not running with a start on record: it stopped and the library has
		// not reported it yet. Close that run before starting a new one.
		await this.flush(asleep);
		if (await this.overBudget()) {
			return new Response("Paused: this month's compute budget is spent.", {
				status: 503,
				headers: { "Retry-After": "3600" },
			});
		}
		await this.ctx.storage.put(LAST_REQUEST, Date.now());
		if (!cold) return super.fetch(request);

		// Time from the request to the container's first answer, which is what
		// the caller waited: the boot plus that one response.
		const boot: { since: number; readyAt?: number } = { since: Date.now() };
		this.boot = boot;
		let status = 0;
		try {
			const response = await super.fetch(request);
			status = response.status;
			return response;
		} finally {
			this.boot = undefined;
			this.recordColdStart(boot, Date.now(), status);
		}
	}

	// One data point per cold start: index and blob1 the app, blob2 "ok" or
	// "failed" (a 5xx or no answer), double1 the milliseconds to the first
	// answer, double2 the status (0 when the container never answered),
	// double3 the milliseconds until the port answered (0 when it never did).
	// double3 is Cloudflare's start plus the process's own boot, which the
	// process traces as its `boot` span; the rest of double1 is the request.
	private recordColdStart(boot: { since: number; readyAt?: number }, end: number, status: number): void {
		try {
			this.env.COLD_STARTS?.writeDataPoint({
				indexes: [this.app],
				blobs: [this.app, status === 0 || status >= 500 ? "failed" : "ok"],
				doubles: [end - boot.since, status, boot.readyAt === undefined ? 0 : boot.readyAt - boot.since],
			});
		} catch {
			// Metrics never get in the way of the request.
		}
	}

	override async onStart(): Promise<void> {
		const now = Date.now();
		// The library calls this once the port answers.
		if (this.boot) this.boot.readyAt = now;
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

export { isNoise, isRoute, notFound } from "./noise";
