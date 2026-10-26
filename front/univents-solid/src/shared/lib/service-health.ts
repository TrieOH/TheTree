import {
  CIRCUIT_OPEN_ERROR_ID,
  unavailableEnvelope,
  type ServiceFailure,
  type ServiceName,
} from "./service-error";

/**
 * Tracks whether each backend the app talks to is answering, and acts as a
 * per-service circuit breaker.
 *
 * The services are independent microservices, so one being down must not stall
 * the others: after a few consecutive transient failures a service's circuit
 * opens and reads against it fail immediately (no 10-second wait per request)
 * until a cool-down passes, at which point the next request is allowed through
 * as a probe. A single success closes the circuit again.
 *
 * Plain TypeScript with `subscribe`/`getSnapshot`, like the SDK's `authStore`:
 * the query client and the transports live outside any component, so they can't
 * reach a context. `useServiceHealth` bridges it into a signal.
 */

export const DEGRADED_AFTER = 2;
export const OPEN_AFTER = 3;
export const COOL_DOWN_MS = 15_000;

export interface ServiceStatus {
  /** Consecutive transient failures since the last success. */
  failures: number;
  /** Epoch ms until which the circuit rejects reads; 0 when closed. */
  openUntil: number;
  lastFailure?: ServiceFailure;
}

export type ServiceHealthSnapshot = Readonly<Record<ServiceName, ServiceStatus>>;

const SERVICES: readonly ServiceName[] = ["api", "identityx"];

const initial = (): ServiceHealthSnapshot => ({
  api: { failures: 0, openUntil: 0 },
  identityx: { failures: 0, openUntil: 0 },
});

/** Raised by {@link guard} while a service's circuit is open. */
export class CircuitOpenError extends Error {
  readonly envelope = unavailableEnvelope(CIRCUIT_OPEN_ERROR_ID);

  constructor(readonly service: ServiceName) {
    super(`Service "${service}" is unavailable (circuit open)`);
    this.name = "CircuitOpenError";
  }
}

export function createServiceHealth(now: () => number = Date.now) {
  let snapshot = initial();
  const listeners = new Set<() => void>();
  const probeListeners = new Set<(service: ServiceName) => void>();
  const recoverListeners = new Set<(service: ServiceName) => void>();
  const retryListeners = new Set<() => void>();
  const probeTimers = new Map<ServiceName, ReturnType<typeof setTimeout>>();

  const publish = (next: ServiceHealthSnapshot) => {
    snapshot = next;
    listeners.forEach((listener) => listener());
  };

  const patch = (service: ServiceName, status: ServiceStatus) =>
    publish({ ...snapshot, [service]: status });

  const clearProbe = (service: ServiceName) => {
    const timer = probeTimers.get(service);
    if (timer === undefined) return;
    clearTimeout(timer);
    probeTimers.delete(service);
  };

  const scheduleProbe = (service: ServiceName) => {
    clearProbe(service);
    probeTimers.set(
      service,
      setTimeout(() => {
        probeTimers.delete(service);
        // Half-open: let the next request through, and tell whoever owns the
        // queries to retry the ones that failed, so the page heals on its own.
        probeListeners.forEach((listener) => listener(service));
      }, COOL_DOWN_MS),
    );
  };

  const halfOpenAll = () => {
    const next = { ...snapshot };
    let changed = false;
    for (const service of SERVICES) {
      if (next[service].openUntil === 0) continue;
      clearProbe(service);
      next[service] = { ...next[service], openUntil: 0 };
      changed = true;
    }
    if (changed) publish(next);
  };

  return {
    subscribe(listener: () => void) {
      listeners.add(listener);
      return () => listeners.delete(listener);
    },

    getSnapshot: () => snapshot,

    /** Called when a service's cool-down ends and it should be probed again. */
    onProbe(listener: (service: ServiceName) => void) {
      probeListeners.add(listener);
      return () => probeListeners.delete(listener);
    },

    /** Called when a previously struggling or failing service answers successfully again. */
    onRecover(listener: (service: ServiceName) => void) {
      recoverListeners.add(listener);
      return () => recoverListeners.delete(listener);
    },

    /** Called when user explicitly requests a retry across components. */
    onRequestRetry(listener: () => void) {
      retryListeners.add(listener);
      return () => retryListeners.delete(listener);
    },

    isOpen(service: ServiceName) {
      return snapshot[service].openUntil > now();
    },

    /** Throws {@link CircuitOpenError} instead of letting a doomed request wait. */
    guard(service: ServiceName) {
      if (snapshot[service].openUntil > now()) throw new CircuitOpenError(service);
    },

    success(service: ServiceName) {
      clearProbe(service);
      const current = snapshot[service];
      const hadFailed = current.failures > 0 || current.openUntil > 0;
      if (current.failures === 0 && current.openUntil === 0) return;
      patch(service, { failures: 0, openUntil: 0 });
      if (hadFailed) {
        recoverListeners.forEach((listener) => listener(service));
      }
    },

    /** Only transient failures count: a 404 or 422 says nothing about health. */
    failure(service: ServiceName, failure: ServiceFailure) {
      if (!failure.transient) return;
      const failures = snapshot[service].failures + 1;
      const opens = failures >= OPEN_AFTER;
      patch(service, {
        failures,
        openUntil: opens ? now() + COOL_DOWN_MS : 0,
        lastFailure: failure,
      });
      if (opens) scheduleProbe(service);
    },

    /**
     * An explicit "try again" from the user: let the next requests through now
     * instead of making them wait out the cool-down. A failure re-opens the
     * circuit, so this can't turn into a request storm.
     */
    halfOpenAll,

    /**
     * Halfs-open all circuits and notifies listeners (such as QueryClient)
     * to refetch all active failed queries.
     */
    requestRetry() {
      halfOpenAll();
      retryListeners.forEach((listener) => listener());
    },

    reset() {
      probeTimers.forEach((timer) => clearTimeout(timer));
      probeTimers.clear();
      publish(initial());
    },
  };
}

export type ServiceHealth = ReturnType<typeof createServiceHealth>;

export const serviceHealth = createServiceHealth();

export function degradedServices(snapshot: ServiceHealthSnapshot): ServiceName[] {
  return SERVICES.filter((service) => snapshot[service].failures >= DEGRADED_AFTER);
}
