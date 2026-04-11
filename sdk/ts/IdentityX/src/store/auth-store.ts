type AuthState = {
  isAuthenticated: boolean;
  isInitializing: boolean;
};

let _state: AuthState = {
  isAuthenticated: false,
  isInitializing: true,
};
const _listeners = new Set<() => void>();

const notify = () => _listeners.forEach((l) => l());

export const authStore = {
  subscribe: (cb: () => void) => {
    _listeners.add(cb);
    return () => _listeners.delete(cb);
  },
  getSnapshot: () => _state,
  getServerSnapshot: () => _state,
  set: (partial: Partial<AuthState>) => {
    _state = { ..._state, ...partial };
    notify();
  },
  reset: () => {
    _state = {
      isAuthenticated: false,
      isInitializing: false,
    };
    notify();
  },
};

// Sync between tabs
if (typeof window !== "undefined") {
  window.addEventListener("storage", (event) => {
    if (event.key === "trieoh_access_expiry") {
      if (!event.newValue) authStore.reset();
      else {
        const expiry = parseInt(event.newValue, 10);
        const isAuthenticated = !isNaN(expiry) && expiry > Date.now();
        authStore.set({ isAuthenticated, isInitializing: false });
      }
    }
  });
}
