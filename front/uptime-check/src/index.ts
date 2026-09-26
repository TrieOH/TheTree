/**
 * The uptime worker: a 5-min cron that checks the public health endpoints
 * from outside the box, a daily cron that compacts the retention ladder
 * (5-min -> hourly -> daily averages), and a fetch handler that serves the
 * status page from KV.
 */

import { runChecks } from "./check";
import { runCompaction } from "./compact";
import { CRON } from "./config";
import { renderPage } from "./page";

export default {
  async scheduled(
    controller: ScheduledController,
    env: Env,
    _ctx: ExecutionContext,
  ) {
    console.log(
      `[scheduled] fired cron=${controller.cron} at=${new Date().toISOString()}`,
    );
    if (controller.cron === CRON.checks) {
      await runChecks(env);
    } else if (controller.cron === CRON.compaction) {
      await runCompaction(env);
    } else {
      console.error(`[scheduled] unexpected cron: ${controller.cron}`);
    }
    console.log(`[scheduled] finished cron=${controller.cron}`);
  },

  async fetch(req: Request, env: Env): Promise<Response> {
    const { pathname } = new URL(req.url);
    // TEMP: manual prod trigger while diagnosing crons — REMOVE ME
    if (pathname === "/__trigger") {
      console.log("[trigger] manual check requested");
      await runChecks(env);
      return new Response("triggered\n");
    }
    return renderPage(env);
  },
} satisfies ExportedHandler<Env>;
