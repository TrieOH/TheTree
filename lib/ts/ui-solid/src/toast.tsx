import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createSignal,
  onSettled,
  untrack,
} from "solid-js";
import { Portal } from "@solidjs/web";

export type ToastType =
  | "default"
  | "success"
  | "error"
  | "warning"
  | "info"
  | "loading";

export type ToastPosition =
  | "top-left"
  | "top-center"
  | "top-right"
  | "bottom-left"
  | "bottom-center"
  | "bottom-right";

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
export type ToastPayload = { message: string } & Omit<ToastOptions, "id">;

export type ToastRecord = {
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
let defaultPosition: ToastPosition = "bottom-right";

const closeHandlers = new Map<number, () => void>();

const [heightMap, setHeightMap] = createSignal<Map<number, number>>(
  new Map(),
  { ownedWrite: true },
);

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
    const overflow = atPosition.slice(
      0,
      atPosition.length - MAX_VISIBLE_PER_POSITION,
    );
    for (const item of overflow) closeHandlers.get(item.id)?.();
  }
}

function pushToast(
  title: string,
  type: ToastType,
  options: ToastOptions = {},
): number {
  const id = options.id ?? nextId++;
  const position = options.position ?? defaultPosition;

  const record: ToastRecord = {
    id,
    type,
    title,
    description: options.description,
    duration:
      options.duration ?? (type === "loading" ? Infinity : DEFAULT_DURATION),
    position,
    dismissible: options.dismissible ?? true,
    action: options.action,
    cancel: options.cancel,
    onDismiss: options.onDismiss,
    onAutoClose: options.onAutoClose,
  };

  setToasts((current) => [
    ...current.filter((item) => item.id !== id),
    record,
  ]);
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
    options?: Omit<ToastOptions, "description">,
  ) => Promise<T>;
}

const toastFn = ((
  input: string | ToastPayload,
  options?: ToastOptions,
) =>
  untrack(() => {
    if (typeof input === "string") {
      return pushToast(input, options?.type ?? "default", options);
    }
    const { message, ...rest } = input;
    return pushToast(message, rest.type ?? "default", rest);
  })) as ToastFn;

toastFn.success = (message, options) =>
  untrack(() => pushToast(message, "success", options));
toastFn.error = (message, options) =>
  untrack(() => pushToast(message, "error", options));
toastFn.warning = (message, options) =>
  untrack(() => pushToast(message, "warning", options));
toastFn.info = (message, options) =>
  untrack(() => pushToast(message, "info", options));
toastFn.message = (message, options) =>
  untrack(() => pushToast(message, "default", options));
toastFn.loading = (message, options) =>
  untrack(() =>
    pushToast(message, "loading", { duration: Infinity, ...options }),
  );

toastFn.dismiss = (id) => {
  if (id === undefined) {
    for (const close of closeHandlers.values()) close();
    return;
  }
  closeHandlers.get(id)?.();
};

toastFn.promise = (promise, messages, options) => {
  const initialId = untrack(() =>
    pushToast(messages.loading, "loading", {
      duration: Infinity,
      ...options,
    }),
  );

  promise
    .then((data) => {
      const text =
        typeof messages.success === "function"
          ? messages.success(data)
          : messages.success;
      untrack(() =>
        pushToast(text, "success", { ...options, id: initialId }),
      );
    })
    .catch((error: unknown) => {
      const text =
        typeof messages.error === "function"
          ? messages.error(error)
          : messages.error;
      untrack(() => pushToast(text, "error", { ...options, id: initialId }));
    });

  return promise;
};

export const toast: ToastFn = toastFn;

/* ------------------------------------------------------------------ */
/*  Inline SVG Icons                                                  */
/* ------------------------------------------------------------------ */

function CheckIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2.5"
      stroke-linecap="round"
      stroke-linejoin="round"
      class="size-3"
      aria-hidden="true"
    >
      <path d="M20 6 9 17l-5-5" />
    </svg>
  );
}

function ErrorIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2.5"
      stroke-linecap="round"
      stroke-linejoin="round"
      class="size-3"
      aria-hidden="true"
    >
      <path d="M18 6 6 18M6 6l12 12" />
    </svg>
  );
}

function WarningIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      class="size-3"
      aria-hidden="true"
    >
      <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z" />
      <path d="M12 9v4M12 17h.01" />
    </svg>
  );
}

function InfoIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      class="size-3"
      aria-hidden="true"
    >
      <circle cx="12" cy="12" r="10" />
      <path d="M12 16v-4M12 8h.01" />
    </svg>
  );
}

function LoadingIcon() {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      class="size-3 animate-spin"
      aria-hidden="true"
    >
      <path d="M21 12a9 9 0 1 1-6.219-8.56" />
    </svg>
  );
}

const ICON_STYLE: Record<
  Exclude<ToastType, "default">,
  { wrapper: string; icon: () => JSX.Element }
> = {
  success: {
    wrapper: "bg-primary text-primary-foreground",
    icon: CheckIcon,
  },
  error: { wrapper: "bg-destructive text-primary", icon: ErrorIcon },
  warning: { wrapper: "bg-accent text-accent-foreground", icon: WarningIcon },
  info: {
    wrapper: "bg-secondary text-secondary-foreground",
    icon: InfoIcon,
  },
  loading: {
    wrapper: "bg-muted text-muted-foreground",
    icon: LoadingIcon,
  },
};

const SWIPE_THRESHOLD = 60;
const DEFAULT_ITEM_HEIGHT = 52;

type ToastItemProps = {
  toast: ToastRecord;
  edge: "top" | "bottom";
  index: () => number;
  offset: () => number;
  scale: () => number;
  layerOpacity: () => number;
  interactive: () => boolean;
};

function ToastItem(props: ToastItemProps) {
  const toastItem = untrack(() => props.toast);

  const [mounted, setMounted] = createSignal(false);
  const [leaving, setLeaving] = createSignal(false);
  const [dragging, setDragging] = createSignal(false);
  const [dragX, setDragX] = createSignal(0);
  const [dragY, setDragY] = createSignal(0);

  let ref: HTMLDivElement | undefined;
  let timer: ReturnType<typeof setTimeout> | undefined;
  let remaining = toastItem.duration;
  let startedAt = 0;
  let pointerStartX = 0;
  let pointerStartY = 0;

  const startTimer = () => {
    if (!Number.isFinite(remaining) || remaining <= 0) return;
    startedAt = Date.now();
    timer = setTimeout(() => {
      toastItem.onAutoClose?.();
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
    toastItem.onDismiss?.();
    setTimeout(() => removeToast(toastItem.id), EXIT_ANIMATION_MS);
  }

  untrack(() => closeHandlers.set(toastItem.id, close));

  onSettled(() => {
    requestAnimationFrame(() => setMounted(true));
    if (Number.isFinite(toastItem.duration)) startTimer();

    let observer: ResizeObserver | undefined;
    if (ref) {
      setHeight(toastItem.id, ref.getBoundingClientRect().height);
      observer = new ResizeObserver((entries) => {
        const entry = entries[0];
        if (entry)
          setHeight(toastItem.id, entry.target.getBoundingClientRect().height);
      });
      observer.observe(ref);
    }

    return () => {
      closeHandlers.delete(toastItem.id);
      if (timer !== undefined) clearTimeout(timer);
      observer?.disconnect();
      clearHeight(toastItem.id);
    };
  });

  function onPointerDown(event: PointerEvent) {
    if (!toastItem.dismissible) return;
    if ((event.target as HTMLElement).closest("button")) return;
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
    const pos = toastItem.position;
    const pushedTowardEdge = pos.startsWith("top")
      ? dy < -SWIPE_THRESHOLD
      : pos.startsWith("bottom")
      ? dy > SWIPE_THRESHOLD
      : false;

    if (Math.abs(dx) > SWIPE_THRESHOLD || pushedTowardEdge) {
      close();
      return;
    }
    setDragX(0);
    setDragY(0);
    resumeTimer();
  }

  const enterOffset = () => {
    const pos = toastItem.position;
    if (pos.endsWith("right")) return { x: 90, y: 0 };
    if (pos.endsWith("left")) return { x: -90, y: 0 };
    return { x: 0, y: pos.startsWith("top") ? -24 : 24 };
  };

  const style = (): JSX.CSSProperties => {
    const sign = props.edge === "top" ? 1 : -1;
    const stackY = sign * props.offset();
    const zIndex = `${50 - props.index()}`;
    const pointerEvents = props.interactive() ? "auto" : "none";

    if (dragging()) {
      const distance = Math.max(Math.abs(dragX()), Math.abs(dragY()));
      return {
        transform: `translate(${dragX()}px, ${
          stackY + dragY()
        }px) scale(${props.scale()})`,
        opacity: `${Math.max(
          (1 - distance / 200) * props.layerOpacity(),
          0.35,
        )}`,
        transition: "none",
        "touch-action": "none",
        cursor: "grabbing",
        "z-index": zIndex,
        "pointer-events": pointerEvents,
      };
    }
    if (leaving()) {
      const off = enterOffset();
      return {
        transform: `translate(${dragX() || off.x}px, ${
          stackY + (dragY() || off.y)
        }px) scale(0.96)`,
        opacity: "0",
        transition: `transform ${EXIT_ANIMATION_MS}ms ease, opacity ${EXIT_ANIMATION_MS}ms ease`,
        "z-index": zIndex,
        "pointer-events": "none",
      };
    }
    if (!mounted()) {
      const off = enterOffset();
      return {
        transform: `translate(${off.x}px, ${stackY + off.y}px) scale(0.96)`,
        opacity: "0",
        "z-index": zIndex,
        "pointer-events": "none",
      };
    }
    return {
      transform: `translate(0, ${stackY}px) scale(${props.scale()})`,
      opacity: `${props.layerOpacity()}`,
      transition:
        "transform 260ms cubic-bezier(0.22, 1, 0.36, 1), opacity 200ms ease, scale 260ms ease",
      "z-index": zIndex,
      "pointer-events": pointerEvents,
    };
  };

  const iconMeta = () =>
    toastItem.type === "default" ? undefined : ICON_STYLE[toastItem.type];

  return (
    <div
      ref={(element) => {
        ref = element;
      }}
      role={
        toastItem.type === "error" || toastItem.type === "warning"
          ? "alert"
          : "status"
      }
      aria-live="polite"
      style={style()}
      onPointerDown={onPointerDown}
      onPointerMove={onPointerMove}
      onPointerUp={onPointerUp}
      onPointerCancel={onPointerUp}
      onMouseEnter={pauseTimer}
      onMouseLeave={resumeTimer}
      class={`pointer-events-auto absolute inset-x-0 ${
        props.edge === "top" ? "top-0" : "bottom-0"
      } flex items-start gap-2.5 rounded-lg border border-border bg-card px-3.5 py-3 text-[13px] text-card-foreground shadow-md shadow-black/6`}
    >
      <Show when={iconMeta()}>
        {(meta) => (
          <div
            class={`flex h-5 w-5 shrink-0 items-center justify-center rounded-full p-1 ${
              meta().wrapper
            }`}
          >
            {meta().icon()}
          </div>
        )}
      </Show>

      <div class="min-w-0 flex-1 mt-px">
        <p class="font-medium leading-snug">{toastItem.title}</p>
        <Show when={toastItem.description}>
          <p class="mt-0.5 text-xs leading-snug text-muted-foreground">
            {toastItem.description}
          </p>
        </Show>
        <Show when={toastItem.action || toastItem.cancel}>
          <div class="mt-2.5 flex gap-2">
            <Show when={toastItem.action}>
              {(action) => (
                <button
                  type="button"
                  onClick={() => {
                    action().onClick(toastItem.id);
                    close();
                  }}
                  class="rounded-md bg-primary px-2 py-1 text-[11px] font-medium text-primary-foreground transition-opacity hover:opacity-90 cursor-pointer"
                >
                  {action().label}
                </button>
              )}
            </Show>
            <Show when={toastItem.cancel}>
              {(cancel) => (
                <button
                  type="button"
                  onClick={() => {
                    cancel().onClick(toastItem.id);
                    close();
                  }}
                  class="rounded-md border border-border px-2 py-1 text-[11px] font-medium text-muted-foreground transition-colors hover:bg-muted cursor-pointer"
                >
                  {cancel().label}
                </button>
              )}
            </Show>
          </div>
        </Show>
      </div>

      <Show when={toastItem.dismissible}>
        <button
          type="button"
          aria-label="Fechar notificação"
          onClick={close}
          class="-mr-1 -mt-0.5 shrink-0 rounded-md p-1 text-muted-foreground/70 transition-colors hover:bg-muted hover:text-foreground cursor-pointer"
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
  "top-left":
    "top-[calc(env(safe-area-inset-top)+1rem)] left-[calc(env(safe-area-inset-left)+0.75rem)]",
  "top-center":
    "top-[calc(env(safe-area-inset-top)+1rem)] left-1/2 -translate-x-1/2",
  "top-right":
    "top-[calc(env(safe-area-inset-top)+1rem)] right-[calc(env(safe-area-inset-right)+0.75rem)]",
  "bottom-left":
    "bottom-[calc(env(safe-area-inset-bottom)+1rem)] left-[calc(env(safe-area-inset-left)+0.75rem)]",
  "bottom-center":
    "bottom-[calc(env(safe-area-inset-bottom)+1rem)] left-1/2 -translate-x-1/2",
  "bottom-right":
    "bottom-[calc(env(safe-area-inset-bottom)+1rem)] right-[calc(env(safe-area-inset-right)+0.75rem)]",
};

function PositionGroup(props: {
  position: ToastPosition;
  items: () => ToastRecord[];
}) {
  const position = untrack(() => props.position);

  const [hovered, setHovered] = createSignal(false);
  const [coarsePointer, setCoarsePointer] = createSignal(false);

  onSettled(() => {
    if (typeof window !== "undefined" && window.matchMedia) {
      setCoarsePointer(
        window.matchMedia("(hover: none), (pointer: coarse)").matches,
      );
    }
  });

  const expanded = () => hovered() || coarsePointer();
  const edge: "top" | "bottom" = position.startsWith("top") ? "top" : "bottom";

  const ordered = () => [...props.items()].reverse();

  const heightOf = (id: number) =>
    heightMap().get(id) ?? DEFAULT_ITEM_HEIGHT;

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
      style={{
        height: `${expanded() ? totalHeight() : frontHeight()}px`,
        transition: "height 220ms ease",
      }}
      class={`pointer-events-none fixed z-100 w-[calc(100%-1.5rem)] max-w-sm ${POSITION_CLASS[position]}`}
    >
      <For each={ordered()}>
        {(item, index) => (
          <ToastItem
            toast={item}
            edge={edge}
            index={index}
            offset={() => offsetFor(index())}
            scale={() =>
              expanded()
                ? 1
                : Math.max(1 - index() * COLLAPSE_SCALE_STEP, 0.9)
            }
            layerOpacity={() =>
              expanded() || index() < VISIBLE_STACK_DEPTH ? 1 : 0
            }
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
  "top-left",
  "top-center",
  "top-right",
  "bottom-left",
  "bottom-center",
  "bottom-right",
];

export function Toaster(props: { position?: ToastPosition }): JSX.Element {
  const initialPosition = untrack(() => props.position);
  if (initialPosition) defaultPosition = initialPosition;

  const atPosition = (position: ToastPosition) => () =>
    toasts().filter((item) => item.position === position);

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
