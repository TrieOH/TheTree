import type { JSX } from '@solidjs/web';
import { For, Show, createSignal, onSettled, untrack } from 'solid-js';
import { Portal } from '@solidjs/web';
import AlertTriangle from '~icons/lucide/triangle-alert';
import Check from '~icons/lucide/check';
import Info from '~icons/lucide/info';
import LoaderCircle from '~icons/lucide/loader-circle';
import X from '~icons/lucide/x';

type IconComponent = () => JSX.Element;

const asIcon = (icon: unknown) => icon as IconComponent;
const CheckIcon = asIcon(Check);
const ErrorIcon = asIcon(X);
const WarningIcon = asIcon(AlertTriangle);
const InfoIcon = asIcon(Info);
const LoadingIcon = asIcon(LoaderCircle);

export type ToastType = 'default' | 'success' | 'error' | 'warning' | 'info' | 'loading';

export type ToastPosition =
  | 'top-left'
  | 'top-center'
  | 'top-right'
  | 'bottom-left'
  | 'bottom-center'
  | 'bottom-right';

export type ToastAction = {
  label: string;
  onClick: (id: number) => void;
};

export type ToastOptions = {
  id?: number;
  type?: ToastType;
  description?: string;
  /** Duration in milliseconds; use `Infinity` to keep the toast open. */
  duration?: number;
  position?: ToastPosition;
  /** Shows the close button and enables swipe dismissal. */
  dismissible?: boolean;
  action?: ToastAction;
  cancel?: ToastAction;
  onDismiss?: () => void;
  onAutoClose?: () => void;
};

/** Alternative form: `toast({ type: 'error', message: '...' })`. */
export type ToastPayload = { message: string } & Omit<ToastOptions, 'id'>;

type ToastRecord = {
  id: number;
  type: ToastType;
  title: string;
  description?: string;
  duration: number;
  position: ToastPosition;
  dismissible: boolean;
  action?: ToastAction;
  cancel?: ToastAction;
  onDismiss?: () => void;
  onAutoClose?: () => void;
};

/* ------------------------------------------------------------------ */
/*  Global state                                                        */
/* ------------------------------------------------------------------ */

const DEFAULT_DURATION = 4000;
const MAX_VISIBLE_PER_POSITION = 5;
const EXIT_ANIMATION_MS = 200;

const [toasts, setToasts] = createSignal<ToastRecord[]>([], { ownedWrite: true });
let nextId = 1;
let defaultPosition: ToastPosition = 'bottom-right';

const closeHandlers = new Map<number, () => void>();

const [heightMap, setHeightMap] = createSignal<Map<number, number>>(new Map(), { ownedWrite: true });

function setHeight(id: number, height: number) {
  setHeightMap((current) => {
    if (current.get(id) === height) return current;
    const next = new Map(current);
    next.set(id, height);
    return next;
  });
}

function clearHeight(id: number) {
  setHeightMap((current) => {
    if (!current.has(id)) return current;
    const next = new Map(current);
    next.delete(id);
    return next;
  });
}

function removeToast(id: number) {
  setToasts((current) => current.filter((item) => item.id !== id));
  closeHandlers.delete(id);
}

function enforcePositionLimit(position: ToastPosition) {
  const atPosition = toasts().filter((item) => item.position === position);
  if (atPosition.length > MAX_VISIBLE_PER_POSITION) {
    const overflow = atPosition.slice(0, atPosition.length - MAX_VISIBLE_PER_POSITION);
    for (const item of overflow) closeHandlers.get(item.id)?.();
  }
}

function pushToast(title: string, type: ToastType, options: ToastOptions = {}): number {
  const id = options.id ?? nextId++;
  const position = options.position ?? defaultPosition;

  const record: ToastRecord = {
    id,
    type,
    title,
    description: options.description,
    duration: options.duration ?? (type === 'loading' ? Infinity : DEFAULT_DURATION),
    position,
    dismissible: options.dismissible ?? true,
    action: options.action,
    cancel: options.cancel,
    onDismiss: options.onDismiss,
    onAutoClose: options.onAutoClose,
  };

  setToasts((current) => [...current.filter((item) => item.id !== id), record]);
  enforcePositionLimit(position);
  return id;
}

/* ------------------------------------------------------------------ */
/*  Public API: toast(...)                                              */
/* ------------------------------------------------------------------ */

export interface ToastFn {
  (input: string | ToastPayload, options?: ToastOptions): number;
  success: (message: string, options?: ToastOptions) => number;
  error: (message: string, options?: ToastOptions) => number;
  warning: (message: string, options?: ToastOptions) => number;
  info: (message: string, options?: ToastOptions) => number;
  loading: (message: string, options?: ToastOptions) => number;
  message: (message: string, options?: ToastOptions) => number;
  dismiss: (id?: number) => void;
  promise: <T>(
    promise: Promise<T>,
    messages: {
      loading: string;
      success: string | ((data: T) => string);
      error: string | ((error: unknown) => string);
    },
    options?: Omit<ToastOptions, 'description'>,
  ) => Promise<T>;
}

const toastFn = ((input: string | ToastPayload, options?: ToastOptions) =>
  untrack(() => {
    if (typeof input === 'string') {
      return pushToast(input, options?.type ?? 'default', options);
    }
    const { message, ...rest } = input;
    return pushToast(message, rest.type ?? 'default', rest);
  })) as ToastFn;

toastFn.success = (message, options) => untrack(() => pushToast(message, 'success', options));
toastFn.error = (message, options) => untrack(() => pushToast(message, 'error', options));
toastFn.warning = (message, options) => untrack(() => pushToast(message, 'warning', options));
toastFn.info = (message, options) => untrack(() => pushToast(message, 'info', options));
toastFn.message = (message, options) => untrack(() => pushToast(message, 'default', options));
toastFn.loading = (message, options) =>
  untrack(() => pushToast(message, 'loading', { duration: Infinity, ...options }));

toastFn.dismiss = (id) => {
  if (id === undefined) {
    for (const close of closeHandlers.values()) close();
    return;
  }
  closeHandlers.get(id)?.();
};

toastFn.promise = (promise, messages, options) => {
  const initialId = untrack(() => pushToast(messages.loading, 'loading', { duration: Infinity, ...options }));

  promise
    .then((data) => {
      const text = typeof messages.success === 'function' ? messages.success(data) : messages.success;
      untrack(() => pushToast(text, 'success', { ...options, id: initialId }));
    })
    .catch((error: unknown) => {
      const text = typeof messages.error === 'function' ? messages.error(error) : messages.error;
      untrack(() => pushToast(text, 'error', { ...options, id: initialId }));
    });

  return promise;
};

/**
 * `toast('mensagem')`
 * `toast('mensagem', { type: 'error' })`
 * `toast({ type: 'error', message: 'message' })`
 * `toast.success('salvo!')`
 */
export const toast: ToastFn = toastFn;

const ICON_STYLE: Record<Exclude<ToastType, 'default'>, { wrapper: string; icon: () => JSX.Element }> = {
  success: { wrapper: 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400', icon: CheckIcon },
  error: { wrapper: 'bg-destructive/15 text-destructive', icon: ErrorIcon },
  warning: { wrapper: 'bg-accent/15 text-accent', icon: WarningIcon },
  info: { wrapper: 'bg-secondary/25 text-secondary-foreground', icon: InfoIcon },
  loading: { wrapper: 'bg-muted text-muted-foreground animate-spin', icon: LoadingIcon },
};

const SWIPE_THRESHOLD = 60;
const DEFAULT_ITEM_HEIGHT = 52;

type ToastItemProps = {
  toast: ToastRecord;
  edge: 'top' | 'bottom';
  index: () => number;
  offset: () => number;
  scale: () => number;
  layerOpacity: () => number;
  interactive: () => boolean;
};

function ToastItem(props: ToastItemProps) {
  const toast = untrack(() => props.toast);

  const [mounted, setMounted] = createSignal(false);
  const [leaving, setLeaving] = createSignal(false);
  const [dragging, setDragging] = createSignal(false);
  const [dragX, setDragX] = createSignal(0);
  const [dragY, setDragY] = createSignal(0);

  let ref: HTMLDivElement | undefined;
  let timer: ReturnType<typeof setTimeout> | undefined;
  let remaining = toast.duration;
  let startedAt = 0;
  let pointerStartX = 0;
  let pointerStartY = 0;

  const startTimer = () => {
    if (!Number.isFinite(remaining) || remaining <= 0) return;
    startedAt = Date.now();
    timer = setTimeout(() => {
      toast.onAutoClose?.();
      close();
    }, remaining);
  };

  const pauseTimer = () => {
    if (timer === undefined) return;
    clearTimeout(timer);
    remaining -= Date.now() - startedAt;
    timer = undefined;
  };

  const resumeTimer = () => {
    if (timer !== undefined || dragging()) return;
    startTimer();
  };

  function close() {
    if (leaving()) return;
    if (timer !== undefined) clearTimeout(timer);
    setLeaving(true);
    toast.onDismiss?.();
    setTimeout(() => removeToast(toast.id), EXIT_ANIMATION_MS);
  }

  untrack(() => closeHandlers.set(toast.id, close));

  onSettled(() => {
    requestAnimationFrame(() => setMounted(true));
    if (Number.isFinite(toast.duration)) startTimer();

    let observer: ResizeObserver | undefined;
    if (ref) {
      setHeight(toast.id, ref.getBoundingClientRect().height);
      observer = new ResizeObserver((entries) => {
        const entry = entries[0];
        if (entry) setHeight(toast.id, entry.target.getBoundingClientRect().height);
      });
      observer.observe(ref);
    }

    return () => {
      closeHandlers.delete(toast.id);
      if (timer !== undefined) clearTimeout(timer);
      observer?.disconnect();
      clearHeight(toast.id);
    };
  });

  function onPointerDown(event: PointerEvent) {
    if (!toast.dismissible) return;
    if ((event.target as HTMLElement).closest('button')) return;
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
    pointerStartX = event.clientX;
    pointerStartY = event.clientY;
    setDragging(true);
    pauseTimer();
  }

  function onPointerMove(event: PointerEvent) {
    if (!dragging()) return;
    setDragX(event.clientX - pointerStartX);
    setDragY(event.clientY - pointerStartY);
  }

  function onPointerUp() {
    if (!dragging()) return;
    setDragging(false);

    const dx = dragX();
    const dy = dragY();
    const pos = toast.position;
    const pushedTowardEdge = pos.startsWith('top') ? dy < -SWIPE_THRESHOLD : pos.startsWith('bottom') ? dy > SWIPE_THRESHOLD : false;

    if (Math.abs(dx) > SWIPE_THRESHOLD || pushedTowardEdge) {
      close();
      return;
    }
    setDragX(0);
    setDragY(0);
    resumeTimer();
  }

  const enterOffset = () => {
    const pos = toast.position;
    if (pos.endsWith('right')) return { x: 90, y: 0 };
    if (pos.endsWith('left')) return { x: -90, y: 0 };
    return { x: 0, y: pos.startsWith('top') ? -24 : 24 };
  };

  const style = (): JSX.CSSProperties => {
    const sign = props.edge === 'top' ? 1 : -1;
    const stackY = sign * props.offset();
    const zIndex = `${50 - props.index()}`;
    const pointerEvents = props.interactive() ? 'auto' : 'none';

    if (dragging()) {
      const distance = Math.max(Math.abs(dragX()), Math.abs(dragY()));
      return {
        transform: `translate(${dragX()}px, ${stackY + dragY()}px) scale(${props.scale()})`,
        opacity: `${Math.max((1 - distance / 200) * props.layerOpacity(), 0.35)}`,
        transition: 'none',
        'touch-action': 'none',
        cursor: 'grabbing',
        'z-index': zIndex,
        'pointer-events': pointerEvents,
      };
    }
    if (leaving()) {
      const off = enterOffset();
      return {
        transform: `translate(${dragX() || off.x}px, ${stackY + (dragY() || off.y)}px) scale(0.96)`,
        opacity: '0',
        transition: `transform ${EXIT_ANIMATION_MS}ms ease, opacity ${EXIT_ANIMATION_MS}ms ease`,
        'z-index': zIndex,
        'pointer-events': 'none',
      };
    }
    if (!mounted()) {
      const off = enterOffset();
      return {
        transform: `translate(${off.x}px, ${stackY + off.y}px) scale(0.96)`,
        opacity: '0',
        'z-index': zIndex,
        'pointer-events': 'none',
      };
    }
    return {
      transform: `translate(0, ${stackY}px) scale(${props.scale()})`,
      opacity: `${props.layerOpacity()}`,
      transition: 'transform 260ms cubic-bezier(0.22, 1, 0.36, 1), opacity 200ms ease, scale 260ms ease',
      'z-index': zIndex,
      'pointer-events': pointerEvents,
    };
  };

  const iconMeta = () => (toast.type === 'default' ? undefined : ICON_STYLE[toast.type]);

  return (
    <div
      ref={(element) => { ref = element; }}
      role={toast.type === 'error' || toast.type === 'warning' ? 'alert' : 'status'}
      aria-live="polite"
      style={style()}
      onPointerDown={onPointerDown}
      onPointerMove={onPointerMove}
      onPointerUp={onPointerUp}
      onPointerCancel={onPointerUp}
      onMouseEnter={pauseTimer}
      onMouseLeave={resumeTimer}
      class={`pointer-events-auto absolute inset-x-0 ${props.edge === 'top' ? 'top-0' : 'bottom-0'} flex items-start gap-2.5 rounded-lg border border-border bg-card px-3.5 py-3 text-[13px] text-card-foreground shadow-md shadow-black/6`}
    >
      <Show when={iconMeta()}>
        {(meta) => (
          <div class={`mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full p-1 ${meta().wrapper}`}>
            {meta().icon()}
          </div>
        )}
      </Show>

      <div class="min-w-0 flex-1">
        <p class="font-medium leading-snug">{toast.title}</p>
        <Show when={toast.description}>
          <p class="mt-0.5 text-xs leading-snug text-muted-foreground">{toast.description}</p>
        </Show>
        <Show when={toast.action || toast.cancel}>
          <div class="mt-2.5 flex gap-2">
            <Show when={toast.action}>
              {(action) => (
                <button
                  type="button"
                  onClick={() => {
                    action().onClick(toast.id);
                    close();
                  }}
                  class="rounded-md bg-primary px-2 py-1 text-[11px] font-medium text-primary-foreground transition-opacity hover:opacity-90"
                >
                  {action().label}
                </button>
              )}
            </Show>
            <Show when={toast.cancel}>
              {(cancel) => (
                <button
                  type="button"
                  onClick={() => {
                    cancel().onClick(toast.id);
                    close();
                  }}
                  class="rounded-md border border-border px-2 py-1 text-[11px] font-medium text-muted-foreground transition-colors hover:bg-muted"
                >
                  {cancel().label}
                </button>
              )}
            </Show>
          </div>
        </Show>
      </div>

      <Show when={toast.dismissible}>
        <button
          type="button"
          aria-label="Fechar notificação"
          onClick={close}
          class="-mr-1 -mt-0.5 shrink-0 rounded-md p-1 text-muted-foreground/70 transition-colors hover:bg-muted hover:text-foreground"
        >
          <ErrorIcon />
        </button>
      </Show>
    </div>
  );
}

/* ------------------------------------------------------------------ */
/*  Position group: compact stack that expands on hover                 */
/* ------------------------------------------------------------------ */

const STACK_GAP = 8;
const COLLAPSE_GAP = 10;
const COLLAPSE_SCALE_STEP = 0.055;
const VISIBLE_STACK_DEPTH = 3;

const POSITION_CLASS: Record<ToastPosition, string> = {
  'top-left': 'top-[calc(env(safe-area-inset-top)+1rem)] left-[calc(env(safe-area-inset-left)+0.75rem)]',
  'top-center': 'top-[calc(env(safe-area-inset-top)+1rem)] left-1/2 -translate-x-1/2',
  'top-right': 'top-[calc(env(safe-area-inset-top)+1rem)] right-[calc(env(safe-area-inset-right)+0.75rem)]',
  'bottom-left': 'bottom-[calc(env(safe-area-inset-bottom)+1rem)] left-[calc(env(safe-area-inset-left)+0.75rem)]',
  'bottom-center': 'bottom-[calc(env(safe-area-inset-bottom)+1rem)] left-1/2 -translate-x-1/2',
  'bottom-right': 'bottom-[calc(env(safe-area-inset-bottom)+1rem)] right-[calc(env(safe-area-inset-right)+0.75rem)]',
};

function PositionGroup(props: { position: ToastPosition; items: () => ToastRecord[] }) {
  const position = untrack(() => props.position);

  const [hovered, setHovered] = createSignal(false);
  const [coarsePointer, setCoarsePointer] = createSignal(false);

  onSettled(() => {
    if (typeof window !== 'undefined' && window.matchMedia) {
      setCoarsePointer(window.matchMedia('(hover: none), (pointer: coarse)').matches);
    }
  });

  // Touch devices have no hover, so keep the stack expanded.
  const expanded = () => hovered() || coarsePointer();
  const edge: 'top' | 'bottom' = position.startsWith('top') ? 'top' : 'bottom';

  // Newest first; index 0 is closest to the screen edge.
  const ordered = () => [...props.items()].reverse();

  const heightOf = (id: number) => heightMap().get(id) ?? DEFAULT_ITEM_HEIGHT;

  const frontHeight = () => {
    const list = ordered();
    return list.length ? heightOf(list[0].id) : 0;
  };

  const totalHeight = () => {
    const list = ordered();
    if (!list.length) return 0;
    let sum = 0;
    for (const item of list) sum += heightOf(item.id);
    return sum + STACK_GAP * (list.length - 1);
  };

  const offsetFor = (index: number) => {
    if (!expanded()) return index * COLLAPSE_GAP;
    const list = ordered();
    let sum = 0;
    for (let i = 0; i < index; i++) sum += heightOf(list[i].id) + STACK_GAP;
    return sum;
  };

  return (
    <div
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
      style={{ height: `${expanded() ? totalHeight() : frontHeight()}px`, transition: 'height 220ms ease' }}
      class={`pointer-events-none fixed z-100 w-[calc(100%-1.5rem)] max-w-sm ${POSITION_CLASS[position]}`}
    >
      <For each={ordered()}>
        {(item, index) => (
          <ToastItem
            toast={item}
            edge={edge}
            index={index}
            offset={() => offsetFor(index())}
            scale={() => (expanded() ? 1 : Math.max(1 - index() * COLLAPSE_SCALE_STEP, 0.9))}
            layerOpacity={() => (expanded() || index() < VISIBLE_STACK_DEPTH ? 1 : 0)}
            interactive={() => expanded() || index() === 0}
          />
        )}
      </For>
    </div>
  );
}

/* ------------------------------------------------------------------ */
/*  Toaster                                                             */
/* ------------------------------------------------------------------ */

const POSITIONS: ToastPosition[] = [
  'top-left',
  'top-center',
  'top-right',
  'bottom-left',
  'bottom-center',
  'bottom-right',
];

export function Toaster(props: { position?: ToastPosition }) {
  const initialPosition = untrack(() => props.position);
  if (initialPosition) defaultPosition = initialPosition;

  const atPosition = (position: ToastPosition) => () => toasts().filter((item) => item.position === position);

  return (
    <Portal>
      <For each={POSITIONS}>
        {(position) => (
          <Show when={atPosition(position)().length > 0}>
            <PositionGroup position={position} items={atPosition(position)} />
          </Show>
        )}
      </For>
    </Portal>
  );
}