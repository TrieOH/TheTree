type AuthState = {
  isAuthenticated: boolean;
};

let _state: AuthState = { isAuthenticated: false };
const _listeners = new Set<() => void>();

const notify = () => _listeners.forEach(l => l());

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
    _state = { isAuthenticated: false };
    notify();
  },
};
