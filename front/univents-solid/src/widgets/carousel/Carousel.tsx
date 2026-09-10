import ChevronLeftIcon from "~icons/lucide/chevron-left";
import ChevronRightIcon from "~icons/lucide/chevron-right";
import {
  For,
  Match,
  Switch,
  createEffect,
  createMemo,
  createSignal,
  flush,
  onSettled,
  untrack,
  type Accessor,
  type Element as SolidElement,
} from "solid-js";
import type { JSX } from "@solidjs/web";
type IconComponent = () => SolidElement;

const ChevronLeft = ChevronLeftIcon as unknown as IconComponent;
const ChevronRight = ChevronRightIcon as unknown as IconComponent;

function createContainerWidth(): [
  (element: HTMLDivElement) => void,
  Accessor<number>,
] {
  const [width, setWidth] = createSignal(0);
  let observer: ResizeObserver | undefined;

  const setRef = (element: HTMLDivElement) => {
    observer?.disconnect();

    const update = () => {
      setWidth(element.getBoundingClientRect().width);
    };

    update();

    if (typeof ResizeObserver === "undefined") return;

    observer = new ResizeObserver((entries) => {
      for (const entry of entries) {
        setWidth(entry.contentRect.width);
      }
    });

    observer.observe(element);
  };

  onSettled(() => {
    return () => observer?.disconnect();
  });

  return [setRef, width];
}

function createReducedMotion(): Accessor<boolean> {
  const [reduced, setReduced] = createSignal(false);

  onSettled(() => {
    if (typeof window === "undefined") return;

    const mql = window.matchMedia("(prefers-reduced-motion: reduce)");
    setReduced(mql.matches);

    const handler = (event: MediaQueryListEvent) => {
      setReduced(event.matches);
    };

    mql.addEventListener("change", handler);

    return () => {
      mql.removeEventListener("change", handler);
    };
  });

  return reduced;
}

type ArrowPosition = "overlay" | "outside" | "top" | "below";

interface CarouselProps<T> {
  items: T[];
  renderItem: (item: T, logicalIndex: number) => SolidElement;
  itemMinWidth?: number;
  itemMaxWidth?: number;
  itemsPerView?: number;
  gap?: number;
  autoPlay?: boolean;
  autoPlayInterval?: number;
  loop?: boolean;
  scrollBy?: "page" | number;
  showArrows?: boolean;
  showDots?: boolean;
  arrowPosition?: ArrowPosition;
  className?: string;
}

export function Carousel<T>(props: CarouselProps<T>) {
  const [setContainerRef, containerWidth] = createContainerWidth();
  const reducedMotion = createReducedMotion();

  const itemMinWidth = () => props.itemMinWidth ?? 220;
  const itemMaxWidth = () => props.itemMaxWidth ?? Infinity;
  const gap = () => props.gap ?? 16;
  const autoPlay = () => props.autoPlay ?? false;
  const autoPlayInterval = () => props.autoPlayInterval ?? 4000;
  const loop = () => props.loop ?? true;
  const scrollBy = () => props.scrollBy ?? "page";
  const showArrows = () => props.showArrows ?? true;
  const showDots = () => props.showDots ?? true;
  const arrowPosition = () => props.arrowPosition ?? "overlay";
  const className = () => props.className ?? "";

  const rawFit = () =>
    containerWidth() > 0
      ? Math.floor((containerWidth() + gap()) / (itemMinWidth() + gap()))
      : 1;

  const autoItemsPerView = () =>
    Math.max(1, Math.min(rawFit(), props.items.length || 1));

  const itemsPerView = () => {
    if (props.itemsPerView != null) {
      return Math.max(1, Math.min(props.itemsPerView, props.items.length || 1));
    }

    return autoItemsPerView();
  };

  const computedWidth = () =>
    itemsPerView() > 0
      ? (containerWidth() - gap() * (itemsPerView() - 1)) / itemsPerView()
      : 0;

  const itemWidth = () => Math.min(computedWidth(), itemMaxWidth());

  const step = () => {
    const value = scrollBy();

    return value === "page"
      ? itemsPerView()
      : Math.max(1, Math.min(value, itemsPerView()));
  };

  const maxIndex = () => Math.max(0, props.items.length - itemsPerView());
  const canSlide = () => props.items.length > itemsPerView();

  const effectiveLoop = () => loop() && canSlide();
  const buffer = () => props.items.length;

  const extendedItems = createMemo(() => {
    const currentItems = props.items;

    if (!effectiveLoop() || currentItems.length === 0) {
      return currentItems;
    }

    const count = buffer();
    const head = currentItems.slice(-count);
    const tail = currentItems.slice(0, count);

    return [...head, ...currentItems, ...tail];
  });

  const [index, setIndex] = createSignal(
    untrack(() => (effectiveLoop() ? buffer() : 0)),
  );
  const [transitionEnabled, setTransitionEnabled] = createSignal(true);

  let animating = false;

  let prevBuffer = untrack(buffer);
  let prevLoop = untrack(effectiveLoop);
  let prevItemsLen = untrack(() => props.items.length);

  createEffect(
    () => ({
      loop: effectiveLoop(),
      buffer: buffer(),
      itemsLength: props.items.length,
      maxIndex: maxIndex(),
    }),
    (next) => {
      const loopChanged = prevLoop !== next.loop;
      const bufferChanged = prevBuffer !== next.buffer;
      const lenChanged = prevItemsLen !== next.itemsLength;

      if (!loopChanged && !bufferChanged && !lenChanged) return;

      setTransitionEnabled(false);
      setIndex((previous) => {
        if (!next.loop) {
          return Math.min(previous, next.maxIndex);
        }

        if (next.itemsLength === 0) {
          return next.buffer;
        }

        const realIndex = prevLoop
          ? (((previous - prevBuffer) % next.itemsLength) + next.itemsLength) %
          next.itemsLength
          : Math.min(previous, Math.max(0, next.itemsLength - 1));

        return next.buffer + realIndex;
      });

      prevBuffer = next.buffer;
      prevLoop = next.loop;
      prevItemsLen = next.itemsLength;

      const raf = requestAnimationFrame(() => {
        setTransitionEnabled(true);
      });

      return () => cancelAnimationFrame(raf);
    },
  );

  const realIndex = () => {
    if (!effectiveLoop() || props.items.length === 0) {
      return index();
    }

    return (
      ((index() - buffer()) % props.items.length) + props.items.length
    ) % props.items.length;
  };

  const currentPage = () => Math.round(realIndex() / step());

  const totalPages = () =>
    effectiveLoop()
      ? Math.max(1, Math.ceil(props.items.length / step()))
      : Math.max(1, Math.ceil((maxIndex() + 1) / step()));

  const dotsVisible = () => showDots() && canSlide() && totalPages() > 1;

  const pages = createMemo(() =>
    Array.from({ length: totalPages() }, (_, page) => page),
  );

  let trackElement: HTMLDivElement | undefined;

  const jumpWithoutTransition = (target: number) => {
    setTransitionEnabled(false);
    flush();

    setIndex(target);
    flush();

    if (trackElement) {
      void trackElement.offsetWidth;
    }

    setTransitionEnabled(true);
    flush();
  };

  const recenterBeforeNavigation = (direction: "next" | "prev") => {
    if (!effectiveLoop() || props.items.length === 0) {
      return index();
    }

    const length = props.items.length;
    let current = index();
    if (direction === "next" && current >= buffer() * 2) {
      current -= length;
      jumpWithoutTransition(current);
    }
    if (direction === "prev" && current < buffer()) {
      current += length;
      jumpWithoutTransition(current);
    }

    return current;
  };

  const handleTransitionEnd: JSX.EventHandler<HTMLDivElement, TransitionEvent> = (
    event,
  ) => {
    if (event.target !== event.currentTarget || event.propertyName !== "transform") {
      return;
    }

    animating = false;
  };

  const goNext = () => {
    if (animating || !canSlide()) return;

    const current = effectiveLoop()
      ? recenterBeforeNavigation("next")
      : index();

    const raw = current + step();
    const target = effectiveLoop() ? raw : Math.min(raw, maxIndex());

    if (target === current) return;

    animating = !reducedMotion();
    setIndex(target);
  };

  const goPrev = () => {
    if (animating || !canSlide()) return;

    const current = effectiveLoop()
      ? recenterBeforeNavigation("prev")
      : index();

    const raw = current - step();
    const target = effectiveLoop() ? raw : Math.max(0, raw);

    if (target === current) return;

    animating = !reducedMotion();
    setIndex(target);
  };

  const goToPage = (page: number) => {
    if (animating || !canSlide()) return;

    const target = effectiveLoop()
      ? buffer() + page * step()
      : Math.min(page * step(), maxIndex());

    if (target === index()) return;

    animating = !reducedMotion();
    setIndex(target);
  };

  const [isHovering, setIsHovering] = createSignal(false);
  const [isVisible, setIsVisible] = createSignal(
    typeof document === "undefined" || !document.hidden,
  );

  onSettled(() => {
    if (typeof document === "undefined") return;

    const onVisibilityChange = () => {
      setIsVisible(!document.hidden);
    };

    document.addEventListener("visibilitychange", onVisibilityChange);

    return () => {
      document.removeEventListener("visibilitychange", onVisibilityChange);
    };
  });

  createEffect(
    () => ({
      enabled: autoPlay(),
      hovering: isHovering(),
      canSlide: canSlide(),
      visible: isVisible(),
      interval: autoPlayInterval(),
      index: index(),
      step: step(),
      loop: effectiveLoop(),
      maxIndex: maxIndex(),
      reducedMotion: reducedMotion(),
    }),
    (state) => {
      if (
        !state.enabled ||
        state.hovering ||
        !state.canSlide ||
        !state.visible
      ) {
        return;
      }

      const id = window.setTimeout(goNext, state.interval);
      return () => window.clearTimeout(id);
    },
  );

  const dragState = {
    dragging: false,
    startX: 0,
    delta: 0,
  };

  const [dragOffset, setDragOffset] = createSignal(0);
  const [isDragging, setIsDragging] = createSignal(false);

  const onPointerDown: JSX.EventHandler<HTMLDivElement, PointerEvent> = (
    event,
  ) => {
    if (!canSlide() || animating) return;

    event.currentTarget.setPointerCapture(event.pointerId);

    dragState.dragging = true;
    dragState.startX = event.clientX;
    dragState.delta = 0;

    setIsDragging(true);
  };

  const onPointerMove: JSX.EventHandler<HTMLDivElement, PointerEvent> = (
    event,
  ) => {
    if (!dragState.dragging) return;

    const delta = event.clientX - dragState.startX;
    dragState.delta = delta;
    setDragOffset(delta);
  };

  const endDrag = () => {
    if (!dragState.dragging) return;

    dragState.dragging = false;
    setIsDragging(false);

    const delta = dragState.delta;
    const threshold = (itemWidth() > 0 ? itemWidth() / 3 : 0) || 40;

    setDragOffset(0);
    flush();

    if (delta > threshold) goPrev();
    else if (delta < -threshold) goNext();
  };

  const translateX = () =>
    -(index() * (itemWidth() + gap())) + dragOffset();

  const arrowBtnBase =
    "flex h-9 w-9 shrink-0 items-center justify-center rounded-full border border-border shadow-md transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-40 disabled:cursor-not-allowed";

  const arrowBtnOverlay =
    "absolute top-1/2 -translate-y-1/2 z-10 bg-card/90 backdrop-blur text-foreground hover:bg-accent hover:text-accent-foreground";

  const arrowBtnStatic =
    "bg-card text-foreground hover:bg-accent hover:text-accent-foreground";

  const prevDisabled = () => !effectiveLoop() && index() <= 0;
  const nextDisabled = () => !effectiveLoop() && index() >= maxIndex();

  const PrevButton = (buttonProps: { variant: ArrowPosition }) => (
    <button
      type="button"
      aria-label="Slide anterior"
      onClick={goPrev}
      disabled={prevDisabled()}
      class={
        buttonProps.variant === "overlay"
          ? `left-2 ${arrowBtnOverlay} ${arrowBtnBase}`
          : `${arrowBtnStatic} ${arrowBtnBase}`
      }
    >
      <span class="flex h-5 w-5 items-center justify-center [&>svg]:h-5 [&>svg]:w-5">
        <ChevronLeft />
      </span>
    </button>
  );

  const NextButton = (buttonProps: { variant: ArrowPosition }) => (
    <button
      type="button"
      aria-label="Próximo slide"
      onClick={goNext}
      disabled={nextDisabled()}
      class={
        buttonProps.variant === "overlay"
          ? `right-2 ${arrowBtnOverlay} ${arrowBtnBase}`
          : `${arrowBtnStatic} ${arrowBtnBase}`
      }
    >
      <span class="flex h-5 w-5 items-center justify-center [&>svg]:h-5 [&>svg]:w-5">
        <ChevronRight />
      </span>
    </button>
  );

  const trackStyle = (): JSX.CSSProperties => ({
    gap: `${gap()}px`,
    transform: `translateX(${translateX()}px)`,
    transition:
      isDragging() || !transitionEnabled() || reducedMotion()
        ? "none"
        : "transform 450ms cubic-bezier(0.22, 1, 0.36, 1)",
    cursor: canSlide() ? (isDragging() ? "grabbing" : "grab") : "default",
    "touch-action": "pan-y",
  });

  const Track = () => (
    <div
      ref={setContainerRef}
      class="overflow-hidden"
      style={
        canSlide()
          ? {
            "mask-image":
              "linear-gradient(to right, transparent, black 32px, black calc(100% - 32px), transparent)",
            "-webkit-mask-image":
              "linear-gradient(to right, transparent, black 32px, black calc(100% - 32px), transparent)",
          }
          : undefined
      }
    >
      <div
        ref={(element) => {
          trackElement = element;
        }}
        class="flex select-none"
        style={trackStyle()}
        onTransitionEnd={handleTransitionEnd}
        onPointerDown={onPointerDown}
        onPointerMove={onPointerMove}
        onPointerUp={endDrag}
        onPointerCancel={endDrag}
      >
        <For each={extendedItems()} keyed={false}>
          {(item, itemIndex) => {
            const logicalIndex = () =>
              effectiveLoop()
                ? (((itemIndex - buffer()) % props.items.length) +
                  props.items.length) %
                props.items.length
                : itemIndex;

            return (
              <div
                class="flex-none"
                style={{
                  width:
                    itemWidth() > 0
                      ? `${itemWidth()}px`
                      : `${itemMinWidth()}px`,
                }}
              >
                {props.renderItem(item(), logicalIndex())}
              </div>
            );
          }}
        </For>
      </div>
    </div>
  );

  const Dots = () => (
    <>
      {dotsVisible() ? (
        <div
          class="flex justify-center gap-2"
          role="tablist"
          aria-label="Controles do carousel"
        >
          <For each={pages()}>
            {(page) => (
              <button
                type="button"
                role="tab"
                aria-selected={page === currentPage() ? "true" : "false"}
                aria-label={`Ir para página ${page + 1}`}
                onClick={() => goToPage(page)}
                class={`h-1.5 rounded-full transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-ring ${page === currentPage()
                    ? "w-6 bg-accent"
                    : "w-1.5 bg-border hover:bg-muted-foreground"
                  }`}
              />
            )}
          </For>
        </div>
      ) : null}
    </>
  );

  return (
    <Switch
      fallback={
        <div
          class={`relative w-full ${className()}`}
          onMouseEnter={() => setIsHovering(true)}
          onMouseLeave={() => setIsHovering(false)}
          role="region"
          aria-roledescription="carousel"
        >
          <Track />

          {showArrows() && canSlide() ? (
            <>
              <PrevButton variant="overlay" />
              <NextButton variant="overlay" />
            </>
          ) : null}

          {dotsVisible() ? (
            <div class="mt-4">
              <Dots />
            </div>
          ) : null}
        </div>
      }
    >
      <Match when={arrowPosition() === "outside"}>
        <div
          class={`w-full ${className()}`}
          onMouseEnter={() => setIsHovering(true)}
          onMouseLeave={() => setIsHovering(false)}
          role="region"
          aria-roledescription="carousel"
        >
          <div class="flex items-center gap-3">
            {showArrows() && canSlide() ? (
              <PrevButton variant="outside" />
            ) : null}

            <div class="min-w-0 flex-1">
              <Track />
            </div>

            {showArrows() && canSlide() ? (
              <NextButton variant="outside" />
            ) : null}
          </div>

          {dotsVisible() ? (
            <div class="mt-4">
              <Dots />
            </div>
          ) : null}
        </div>
      </Match>

      <Match when={arrowPosition() === "below"}>
        <div
          class={`w-full ${className()}`}
          onMouseEnter={() => setIsHovering(true)}
          onMouseLeave={() => setIsHovering(false)}
          role="region"
          aria-roledescription="carousel"
        >
          <Track />

          {(showArrows() && canSlide()) ||
            dotsVisible() ? (
            <div class="mt-4 flex items-center justify-center gap-4">
              {showArrows() && canSlide() ? (
                <PrevButton variant="below" />
              ) : null}

              <Dots />

              {showArrows() && canSlide() ? (
                <NextButton variant="below" />
              ) : null}
            </div>
          ) : null}
        </div>
      </Match>

      <Match when={arrowPosition() === "top"}>
        <div
          class={`w-full ${className()}`}
          onMouseEnter={() => setIsHovering(true)}
          onMouseLeave={() => setIsHovering(false)}
          role="region"
          aria-roledescription="carousel"
        >
          {(showArrows() && canSlide()) ||
            dotsVisible() ? (
            <div class="mb-4 flex items-center justify-center gap-4">
              {showArrows() && canSlide() ? (
                <PrevButton variant="below" />
              ) : null}

              <Dots />

              {showArrows() && canSlide() ? (
                <NextButton variant="below" />
              ) : null}
            </div>
          ) : null}

          <Track />
        </div>
      </Match>
    </Switch>
  );
}
