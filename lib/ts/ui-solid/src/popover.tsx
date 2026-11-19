import { Portal, type JSX } from "@solidjs/web";
import { Show, createEffect, createSignal, createUniqueId } from "solid-js";
import { cn } from "./lib/cn";
import { shapeClass, type Shape } from "./lib/shape";

export interface PopoverTriggerProps {
  ref: (element: HTMLButtonElement) => void;
  onClick: () => void;
  onKeyDown: (event: KeyboardEvent) => void;
  "aria-expanded": "true" | "false";
  "aria-controls": string;
  "aria-haspopup": "dialog" | "menu";
  disabled?: boolean;
}
export interface PopoverProps {
  trigger: (props: PopoverTriggerProps) => JSX.Element;
  children: JSX.Element | ((close: () => void) => JSX.Element);
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  disabled?: boolean;
  ariaLabel: string;
  role?: "dialog" | "menu";
  align?: "start" | "end";
  width?: number;
  shape?: Shape;
  class?: string;
}
const focusable =
  'button:not([disabled]),input:not([disabled]),textarea:not([disabled]),a[href],[tabindex="0"]';
/** Anchored, non-modal surface with viewport collision handling and focus restoration. */
export function Popover(props: PopoverProps) {
  const id = createUniqueId();
  const [internal, setInternal] = createSignal(false);
  const [anchor, setAnchor] = createSignal<HTMLButtonElement>();
  const [panel, setPanel] = createSignal<HTMLDivElement>();
  const [entry, setEntry] = createSignal<"first" | "last">("first");
  const [position, setPosition] = createSignal({ top: 0, left: 0 });
  const open = () => !props.disabled && (props.open ?? internal());
  const change = (value: boolean) => {
    setInternal(value);
    props.onOpenChange?.(value);
  };
  createEffect(
    () => props.disabled,
    (disabled) => {
      if (disabled) change(false);
    },
  );
  const close = () => {
    change(false);
    anchor()?.focus({ preventScroll: true });
  };
  createEffect(
    () => ({
      open: open(),
      anchor: anchor(),
      panel: panel(),
      width: props.width ?? 280,
      align: props.align,
      entry: entry(),
      role: props.role,
    }),
    ({ open, anchor, panel, width, align, entry, role }) => {
      if (!open || !anchor || !panel) return;
      // The browser top layer escapes transformed or clipped dialog/preview ancestors
      // while keeping the surface inside its owner for keyboard focus containment.
      if (typeof panel.showPopover === "function" && panel.isConnected)
        panel.showPopover();
      const reposition = () => {
        const a = anchor.getBoundingClientRect(),
          p = panel.getBoundingClientRect();
        const below = innerHeight - a.bottom - 8,
          above = a.top - 8;
        const top =
          below < p.height + 6 && above > below
            ? a.top - p.height - 6
            : a.bottom + 6;
        setPosition({
          top: Math.max(8, Math.min(top, innerHeight - p.height - 8)),
          left: Math.max(
            8,
            Math.min(
              align === "end"
                ? a.right - Math.min(width, innerWidth - 16)
                : a.left,
              innerWidth - p.width - 8,
            ),
          ),
        });
      };
      reposition();
      const candidates = [...panel.querySelectorAll<HTMLElement>(focusable)];
      (
        panel.querySelector<HTMLElement>("[autofocus], [data-autofocus]") ??
        (role === "menu" && entry === "last"
          ? candidates.at(-1)
          : candidates[0]) ??
        panel
      ).focus({ preventScroll: true });
      const outside = (event: PointerEvent) => {
        const target = event.target as Node;
        if (!anchor.contains(target) && !panel.contains(target)) change(false);
      };
      const keyboard = (event: KeyboardEvent) => {
        if (event.key === "Escape") {
          event.preventDefault();
          event.stopPropagation();
          close();
        }
      };
      const focus = (event: FocusEvent) => {
        const target = event.target as Node;
        if (!anchor.contains(target) && !panel.contains(target)) change(false);
      };
      const observer =
        typeof ResizeObserver !== "undefined"
          ? new ResizeObserver(reposition)
          : undefined;
      observer?.observe(panel);
      observer?.observe(anchor);
      document.addEventListener("pointerdown", outside);
      document.addEventListener("focusin", focus);
      panel.addEventListener("keydown", keyboard);
      window.addEventListener("resize", reposition);
      window.addEventListener("scroll", reposition, true);
      return () => {
        if (typeof panel.hidePopover === "function") panel.hidePopover();
        observer?.disconnect();
        document.removeEventListener("pointerdown", outside);
        document.removeEventListener("focusin", focus);
        panel.removeEventListener("keydown", keyboard);
        window.removeEventListener("resize", reposition);
        window.removeEventListener("scroll", reposition, true);
      };
    },
  );
  return (
    <>
      {props.trigger({
        ref: setAnchor,
        get "aria-expanded"() {
          return open() ? "true" : "false";
        },
        "aria-controls": id,
        get "aria-haspopup"() {
          return props.role ?? "dialog";
        },
        get disabled() {
          return props.disabled;
        },
        onClick: () => {
          setEntry("first");
          change(!open());
        },
        onKeyDown: (event) => {
          if (event.key === "ArrowDown" || event.key === "ArrowUp") {
            event.preventDefault();
            setEntry(event.key === "ArrowUp" ? "last" : "first");
            change(true);
          }
        },
      })}
      <Show when={open()}>
        <Portal mount={anchor()?.closest('[role="dialog"]') ?? undefined}>
          <div
            ref={setPanel}
            popover="manual"
            id={id}
            role={props.role ?? "dialog"}
            aria-label={props.ariaLabel}
            tabindex={-1}
            class={cn(
              "fixed z-100 max-w-[calc(100vw-1rem)] overflow-y-auto overscroll-contain rounded-xl border border-border bg-popover p-3 text-popover-foreground shadow-lg outline-none",
              shapeClass(props.shape, "surface"),
              props.class,
            )}
            style={{
              margin: "0",
              inset: "auto",
              width: `${props.width ?? 280}px`,
              "max-height": "calc(100dvh - 1rem)",
              top: `${position().top}px`,
              left: `${position().left}px`,
            }}
          >
            {typeof props.children === "function"
              ? props.children(close)
              : props.children}
          </div>
        </Portal>
      </Show>
    </>
  );
}
