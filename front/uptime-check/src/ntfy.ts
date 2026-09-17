/**
 * Post a notification to the ntfy topic. Degraded never calls this —
 * only dead (loud) and recovery (quieter) do.
 */
export async function notify(
  env: Env,
  title: string,
  body: string,
  priority: string,
  tags: string[],
): Promise<void> {
  if (!env.NTFY_URL) {
    console.error("NTFY_URL is not set; ping dropped:", title);
    return;
  }
  try {
    const res = await fetch(env.NTFY_URL, {
      method: "POST",
      headers: {
        Title: title,
        Priority: priority,
        Tags: tags.join(","),
      },
      body,
    });
    if (!res.ok) {
      console.error(`ntfy responded HTTP ${res.status} for "${title}"`);
    }
  } catch (e) {
    // Never let a notification failure break the check run.
    console.error("ntfy post failed:", e instanceof Error ? e.message : e);
  }
}
