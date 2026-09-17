// Secrets (set via `wrangler secret put`) don't appear in wrangler's
// generated Env — declare them here; ambient interfaces merge.
interface Env {
  /** ntfy topic URL that receives dead/recovery pings. */
  NTFY_URL?: string;
}
