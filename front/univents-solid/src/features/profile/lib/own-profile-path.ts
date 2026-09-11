import { createSignal } from "solid-js";

/**
 * Canonical `/profile/<handle>` path of the viewer's own profile, published by
 * the profile route once it knows which actor it loaded. `null` while unknown
 * or while viewing somebody else.
 *
 * Read imperatively (from a click handler) to tell "my profile" apart from
 * someone else's — the route id is the same for both, and the URL segment can
 * be either a handle or an actor id.
 */
const [ownProfilePath, setOwnProfilePath] = createSignal<string | null>(null);

export { ownProfilePath, setOwnProfilePath };
