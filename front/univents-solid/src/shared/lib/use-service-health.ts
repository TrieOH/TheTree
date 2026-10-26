import { createSignal, onCleanup, untrack } from "solid-js";
import { serviceHealth } from "./service-health";

/** The live health snapshot as a signal, unsubscribed with the owner. */
export function useServiceHealth() {
  const [snapshot, setSnapshot] = createSignal(
    untrack(() => serviceHealth.getSnapshot()),
    { ownedWrite: true },
  );
  onCleanup(
    serviceHealth.subscribe(() => setSnapshot(serviceHealth.getSnapshot())),
  );
  return snapshot;
}
