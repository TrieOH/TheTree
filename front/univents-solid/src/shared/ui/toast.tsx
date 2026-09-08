import { For, createSignal, onCleanup } from 'solid-js';

type Toast = { id: number; type: 'success' | 'error' | 'info'; message: string };
const [items, setItems] = createSignal<Toast[]>([]);
let nextId = 0;

function show(type: Toast['type'], message: string) {
  const id = nextId++;
  setItems((current) => [...current, { id, type, message }]);
  setTimeout(() => setItems((current) => current.filter((item) => item.id !== id)), 4000);
}

export const toast = Object.assign(
  (input: string | Omit<Toast, 'id'>) =>
    typeof input === 'string'
      ? show('success', input)
      : show(input.type, input.message),
  {
    success: (message: string) => show('success', message),
    error: (message: string) => show('error', message),
    info: (message: string) => show('info', message),
  },
);

export function Toaster() {
  onCleanup(() => setItems([]));
  return (
    <div class="fixed right-4 top-4 z-50 grid max-w-sm gap-2">
      <For each={items()}>
        {(item) => (
          <div class={`rounded-lg border px-4 py-3 text-sm shadow-lg ${item.type === 'error' ? 'border-destructive/40 bg-destructive text-destructive-foreground' : 'border-emerald-500/40 bg-emerald-600 text-white'}`}>
            {item.message}
          </div>
        )}
      </For>
    </div>
  );
}
