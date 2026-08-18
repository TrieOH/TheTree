import { ChevronLeft, ChevronRight } from "lucide-react";
import {
  type CSSProperties,
  type PointerEvent,
  type ReactNode,
  type RefObject,
  useCallback,
  useEffect,
  useLayoutEffect,
  useMemo,
  useRef,
  useState,
} from "react";

/* ------------------------------------------------------------------ */
/*  Hook: largura do container (ResizeObserver)                       */
/* ------------------------------------------------------------------ */
function useContainerWidth(): [RefObject<HTMLDivElement | null>, number] {
  const ref = useRef<HTMLDivElement | null>(null);
  const [width, setWidth] = useState(0);

  useLayoutEffect(() => {
    const el = ref.current;
    if (!el) return;
    const update = () => setWidth(el.getBoundingClientRect().width);
    update();
    const ro = new ResizeObserver((entries) => {
      for (const entry of entries) setWidth(entry.contentRect.width);
    });
    ro.observe(el);
    return () => ro.disconnect();
  }, []);

  return [ref, width];
}

/* ------------------------------------------------------------------ */
/*  Hook: prefers-reduced-motion                                      */
/* ------------------------------------------------------------------ */
function useReducedMotion(): boolean {
  const [reduced, setReduced] = useState(false);
  useEffect(() => {
    const mql = window.matchMedia("(prefers-reduced-motion: reduce)");
    setReduced(mql.matches);
    const handler = (e: MediaQueryListEvent) => setReduced(e.matches);
    mql.addEventListener("change", handler);
    return () => mql.removeEventListener("change", handler);
  }, []);
  return reduced;
}

/* ------------------------------------------------------------------ */
/*  Tipos do Carousel                                                 */
/* ------------------------------------------------------------------ */
type ArrowPosition = "overlay" | "outside" | "below";

interface CarouselProps<T> {
  items: T[];
  renderItem: (item: T, logicalIndex: number) => ReactNode;
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

/* ------------------------------------------------------------------ */
/*  Carousel genérico                                                 */
/* ------------------------------------------------------------------ */
export function Carousel<T>({
  items,
  renderItem,
  itemMinWidth = 220,
  itemMaxWidth = Infinity,
  itemsPerView: forcedItemsPerView,
  gap = 16,
  autoPlay = false,
  autoPlayInterval = 4000,
  loop = true,
  scrollBy = "page",
  showArrows = true,
  showDots = true,
  arrowPosition = "overlay",
  className = "",
}: CarouselProps<T>) {
  const [containerRef, containerWidth] = useContainerWidth();
  const reducedMotion = useReducedMotion();

  /* ---------- quantos cabem? ---------- */
  const rawFit =
    containerWidth > 0
      ? Math.floor((containerWidth + gap) / (itemMinWidth + gap))
      : 1;

  const autoItemsPerView = Math.max(1, Math.min(rawFit, items.length || 1));

  const itemsPerView =
    forcedItemsPerView != null
      ? Math.max(1, Math.min(forcedItemsPerView, items.length || 1))
      : autoItemsPerView;

  const computedWidth =
    itemsPerView > 0
      ? (containerWidth - gap * (itemsPerView - 1)) / itemsPerView
      : 0;

  const itemWidth = Math.min(computedWidth, itemMaxWidth);

  const step =
    scrollBy === "page"
      ? itemsPerView
      : Math.max(1, Math.min(scrollBy, itemsPerView));

  const maxIndex = Math.max(0, items.length - itemsPerView);
  const canSlide = items.length > itemsPerView;

  /* ---------- loop infinito ---------- */
  const effectiveLoop = loop && canSlide;
  // Keep a complete list on both sides. A partial buffer can leave empty
  // space when the last visible item wraps to the first item.
  const buffer = items.length;

  const extendedItems = useMemo(() => {
    if (!effectiveLoop || items.length === 0) return items;
    const head = items.slice(-buffer);
    const tail = items.slice(0, buffer);
    return [...head, ...items, ...tail];
  }, [items, effectiveLoop, buffer]);

  /* ---------- estado de índice ---------- */
  const [index, setIndex] = useState(effectiveLoop ? buffer : 0);
  const [transitionEnabled, setTransitionEnabled] = useState(true);
  const animatingRef = useRef(false);

  const prevBufferRef = useRef(buffer);
  const prevLoopRef = useRef(effectiveLoop);
  const prevItemsLenRef = useRef(items.length);

  useEffect(() => {
    const loopChanged = prevLoopRef.current !== effectiveLoop;
    const bufferChanged = prevBufferRef.current !== buffer;
    const lenChanged = prevItemsLenRef.current !== items.length;

    if (!loopChanged && !bufferChanged && !lenChanged) return;

    setTransitionEnabled(false);
    setIndex((prev) => {
      if (!effectiveLoop) return Math.min(prev, maxIndex);
      const realIndex = prevLoopRef.current
        ? (((prev - prevBufferRef.current) % items.length) + items.length) %
          items.length
        : Math.min(prev, Math.max(0, items.length - 1));
      return buffer + realIndex;
    });

    prevBufferRef.current = buffer;
    prevLoopRef.current = effectiveLoop;
    prevItemsLenRef.current = items.length;

    const raf = requestAnimationFrame(() => setTransitionEnabled(true));
    return () => cancelAnimationFrame(raf);
  }, [effectiveLoop, buffer, items.length, maxIndex]);

  /* ---------- índice lógico ---------- */
  const realIndex = effectiveLoop
    ? (((index - buffer) % items.length) + items.length) % items.length
    : index;

  const currentPage = Math.round(realIndex / step);
  const totalPages = effectiveLoop
    ? Math.max(1, Math.ceil(items.length / step))
    : Math.max(1, Math.ceil((maxIndex + 1) / step));

  /* ---------- navegação ---------- */
  const handleTransitionEnd = useCallback(() => {
    animatingRef.current = false;
    if (!effectiveLoop) return;

    setIndex((i) => {
      const normalized =
        buffer +
        ((((i - buffer) % items.length) + items.length) % items.length);
      if (normalized === i) return i;
      setTransitionEnabled(false);
      requestAnimationFrame(() => setTransitionEnabled(true));
      return normalized;
    });
  }, [effectiveLoop, buffer, items.length]);

  const goNext = useCallback(() => {
    if (animatingRef.current || !canSlide) return;
    const raw = index + step;
    const target = effectiveLoop ? raw : Math.min(raw, maxIndex);
    if (target === index) return;
    animatingRef.current = !reducedMotion;
    setIndex(target);
  }, [index, step, effectiveLoop, maxIndex, canSlide, reducedMotion]);

  const goPrev = useCallback(() => {
    if (animatingRef.current || !canSlide) return;
    const raw = index - step;
    const target = effectiveLoop ? raw : Math.max(0, raw);
    if (target === index) return;
    animatingRef.current = !reducedMotion;
    setIndex(target);
  }, [index, step, effectiveLoop, canSlide, reducedMotion]);

  const goToPage = useCallback(
    (p: number) => {
      if (animatingRef.current || !canSlide) return;
      const target = effectiveLoop
        ? buffer + p * step
        : Math.min(p * step, maxIndex);
      if (target === index) return;
      animatingRef.current = !reducedMotion;
      setIndex(target);
    },
    [index, step, effectiveLoop, buffer, maxIndex, canSlide, reducedMotion],
  );

  /* ---------- autoplay + visibilidade ---------- */
  const [isHovering, setIsHovering] = useState(false);
  const [isVisible, setIsVisible] = useState(
    typeof document === "undefined" || !document.hidden,
  );

  useEffect(() => {
    const onVis = () => setIsVisible(!document.hidden);
    document.addEventListener("visibilitychange", onVis);
    return () => document.removeEventListener("visibilitychange", onVis);
  }, []);

  useEffect(() => {
    if (!autoPlay || isHovering || !canSlide || !isVisible) return;
    const id = setTimeout(goNext, autoPlayInterval);
    return () => clearTimeout(id);
  }, [autoPlay, isHovering, canSlide, isVisible, goNext, autoPlayInterval]);

  /* ---------- drag / swipe ---------- */
  const dragState = useRef({ dragging: false, startX: 0, delta: 0 });
  const [dragOffset, setDragOffset] = useState(0);
  const [isDragging, setIsDragging] = useState(false);

  const onPointerDown = (e: PointerEvent<HTMLDivElement>) => {
    if (!canSlide || animatingRef.current) return;
    e.currentTarget.setPointerCapture(e.pointerId);
    dragState.current = { dragging: true, startX: e.clientX, delta: 0 };
    setIsDragging(true);
  };

  const onPointerMove = (e: PointerEvent<HTMLDivElement>) => {
    if (!dragState.current.dragging) return;
    const delta = e.clientX - dragState.current.startX;
    dragState.current.delta = delta;
    setDragOffset(delta);
  };

  const endDrag = useCallback(() => {
    if (!dragState.current.dragging) return;
    dragState.current.dragging = false;
    setIsDragging(false);
    const delta = dragState.current.delta;
    const threshold = (itemWidth > 0 ? itemWidth / 3 : 0) || 40;
    setDragOffset(0);
    if (delta > threshold) goPrev();
    else if (delta < -threshold) goNext();
  }, [itemWidth, goPrev, goNext]);

  const translateX = -(index * (itemWidth + gap)) + dragOffset;

  /* ---------- setas ---------- */
  const arrowBtnBase =
    "flex h-9 w-9 shrink-0 items-center justify-center rounded-full border border-border shadow-md transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-40 disabled:cursor-not-allowed";

  const arrowBtnOverlay =
    "absolute top-1/2 -translate-y-1/2 z-10 bg-card/90 backdrop-blur text-foreground hover:bg-accent hover:text-accent-foreground";

  const arrowBtnStatic =
    "bg-card text-foreground hover:bg-accent hover:text-accent-foreground";

  const prevDisabled = !effectiveLoop && index <= 0;
  const nextDisabled = !effectiveLoop && index >= maxIndex;

  const PrevButton = ({ variant }: { variant: ArrowPosition }) => (
    <button
      type="button"
      aria-label="Slide anterior"
      onClick={goPrev}
      disabled={prevDisabled}
      className={
        variant === "overlay"
          ? `left-2 ${arrowBtnOverlay} ${arrowBtnBase}`
          : `${arrowBtnStatic} ${arrowBtnBase}`
      }
    >
      <ChevronLeft className="h-5 w-5" />
    </button>
  );

  const NextButton = ({ variant }: { variant: ArrowPosition }) => (
    <button
      type="button"
      aria-label="Próximo slide"
      onClick={goNext}
      disabled={nextDisabled}
      className={
        variant === "overlay"
          ? `right-2 ${arrowBtnOverlay} ${arrowBtnBase}`
          : `${arrowBtnStatic} ${arrowBtnBase}`
      }
    >
      <ChevronRight className="h-5 w-5" />
    </button>
  );

  /* ---------- track ---------- */
  const trackStyle: CSSProperties = {
    gap: `${gap}px`,
    transform: `translateX(${translateX}px)`,
    transition:
      isDragging || !transitionEnabled || reducedMotion
        ? "none"
        : "transform 450ms cubic-bezier(0.22, 1, 0.36, 1)",
    cursor: canSlide ? (isDragging ? "grabbing" : "grab") : "default",
    touchAction: "pan-y",
  };

  const track = (
    <div
      ref={containerRef}
      className="overflow-hidden"
      style={
        canSlide
          ? {
              maskImage:
                "linear-gradient(to right, transparent, black 32px, black calc(100% - 32px), transparent)",
              WebkitMaskImage:
                "linear-gradient(to right, transparent, black 32px, black calc(100% - 32px), transparent)",
            }
          : undefined
      }
    >
      <div
        className="flex select-none"
        style={trackStyle}
        onTransitionEnd={handleTransitionEnd}
        onPointerDown={onPointerDown}
        onPointerMove={onPointerMove}
        onPointerUp={endDrag}
        onPointerCancel={endDrag}
      >
        {extendedItems.map((item, i) => {
          const logicalIndex = effectiveLoop
            ? (((i - buffer) % items.length) + items.length) % items.length
            : i;
          return (
            <div
              key={`${i}-${logicalIndex}`}
              className="flex-none"
              style={{
                width: itemWidth > 0 ? `${itemWidth}px` : `${itemMinWidth}px`,
              }}
            >
              {renderItem(item, logicalIndex)}
            </div>
          );
        })}
      </div>
    </div>
  );

  /* ---------- dots ---------- */
  const dots =
    showDots && canSlide && totalPages > 1 ? (
      <div
        className="flex justify-center gap-2"
        role="tablist"
        aria-label="Controles do carousel"
      >
        {Array.from({ length: totalPages }).map((_, p) => (
          <button
            key={`page-${p + 1}`}
            type="button"
            role="tab"
            aria-selected={p === currentPage}
            aria-label={`Ir para página ${p + 1}`}
            onClick={() => goToPage(p)}
            className={`h-1.5 rounded-full transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-ring ${
              p === currentPage
                ? "w-6 bg-accent"
                : "w-1.5 bg-border hover:bg-muted-foreground"
            }`}
          />
        ))}
      </div>
    ) : null;

  /* ---------- layouts ---------- */
  const wrapperProps = {
    className: `w-full ${className}`,
    onMouseEnter: () => setIsHovering(true),
    onMouseLeave: () => setIsHovering(false),
    role: "region" as const,
    "aria-roledescription": "carousel" as const,
  };

  if (arrowPosition === "outside") {
    return (
      <div {...wrapperProps}>
        <div className="flex items-center gap-3">
          {showArrows && canSlide && <PrevButton variant="outside" />}
          <div className="min-w-0 flex-1">{track}</div>
          {showArrows && canSlide && <NextButton variant="outside" />}
        </div>
        {dots && <div className="mt-4">{dots}</div>}
      </div>
    );
  }

  if (arrowPosition === "below") {
    return (
      <div {...wrapperProps}>
        {track}
        {(showArrows && canSlide) || dots ? (
          <div className="mt-4 flex items-center justify-center gap-4">
            {showArrows && canSlide && <PrevButton variant="below" />}
            {dots}
            {showArrows && canSlide && <NextButton variant="below" />}
          </div>
        ) : null}
      </div>
    );
  }

  // overlay (padrão)
  return (
    <div {...wrapperProps} className={`relative w-full ${className}`}>
      {track}
      {showArrows && canSlide && (
        <>
          <PrevButton variant="overlay" />
          <NextButton variant="overlay" />
        </>
      )}
      {dots && <div className="mt-4">{dots}</div>}
    </div>
  );
}
