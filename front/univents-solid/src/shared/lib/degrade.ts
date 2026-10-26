import { classifyError } from "./service-error";

/**
 * Resolves to `fallback` when optional data can't be loaded.
 *
 * For personal data (the ticket a user holds, their registrations) fetched next
 * to public data: those calls go through the session, so an auth outage fails
 * them while the public catalog is perfectly healthy. They must not take the
 * whole page down; the page just renders without them.
 */
export function degradeTo<T>(work: Promise<T> | T, fallback: T): Promise<T> {
  return Promise.resolve(work).catch(() => fallback);
}

/**
 * Like {@link degradeTo}, but only for definitive answers (a 403 or 404 for an
 * optional list). A transient failure is rethrown: swallowing it would cache
 * "there is nothing" as a success, and the page would show an empty state for a
 * backend that was only briefly unreachable.
 */
export function degradeUnlessTransient<T>(work: Promise<T> | T, fallback: T): Promise<T> {
  return Promise.resolve(work).catch((error: unknown) => {
    if (classifyError(error).transient) throw error;
    return fallback;
  });
}
