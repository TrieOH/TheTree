import {
  createEffect,
  createMemo,
  createSignal,
  For,
  onSettled,
} from "solid-js";
import type { JSX } from "@solidjs/web";

export interface VirtualizedContainerProps<T> {
  items: readonly T[];
  children: (item: T, index: () => number) => JSX.Element;
  class?: string;
  /** Responsive column count, based on the container width in pixels. */
  columns?: number | ((width: number) => number);
  /** Initial row height in pixels; mounted rows are measured automatically. */
  estimatedRowHeight?: number | ((width: number) => number);
  gap?: number | ((width: number) => number);
  /** Extra rows mounted above and below the viewport. Defaults to 2. */
  overscan?: number;
}

/** Window-scrolling list/grid. Keep editable item state outside mounted rows. */
export function VirtualizedContainer<T>(props: VirtualizedContainerProps<T>) {
  let container!: HTMLDivElement;
  const [viewport, setViewport] = createSignal({ width: 0, top: 0, height: 0 });
  const [measurements, setMeasurements] = createSignal<
    ReadonlyMap<number, number>
  >(new Map());
  const width = createMemo(() => viewport().width);
  const resolve = (
    value: number | ((width: number) => number) | undefined,
    fallback: number,
  ) => (typeof value === "function" ? value(width()) : (value ?? fallback));
  const columns = () => Math.max(1, Math.floor(resolve(props.columns, 1)));
  const gap = () => Math.max(0, resolve(props.gap, 0));
  const offsets = createMemo(() => {
    // Rebuild measurements when the collection or responsive geometry changes.
    const items = props.items;
    const measuredWidth = width();
    const count = columns();
    const measured = measurements();
    const estimate = Math.max(1, resolve(props.estimatedRowHeight, 100));
    const result = [0];
    for (let row = 0; row < Math.ceil(items.length / count); row++) {
      const height = measured.get(row) ?? estimate;
      result.push(result[row] + height + gap());
    }
    return { positions: result, items, width: measuredWidth, count };
  });
  const visible = createMemo(() => {
    const { positions } = offsets();
    const { top, height } = viewport();
    const find = (target: number) => {
      let low = 0;
      let high = positions.length - 1;
      while (low < high) {
        const middle = Math.floor((low + high + 1) / 2);
        if (positions[middle] <= target) low = middle;
        else high = middle - 1;
      }
      return low;
    };
    const extra = Math.max(0, Math.floor(props.overscan ?? 2));
    const first = Math.max(0, find(Math.max(0, top)) - extra);
    const end = Math.min(
      positions.length - 1,
      find(Math.max(0, top + height)) + extra + 1,
    );
    return Array.from(
      { length: Math.max(0, end - first) },
      (_, i) => first + i,
    );
  });

  onSettled(() => {
    let frame = 0;
    const update = () => {
      frame = 0;
      const rect = container.getBoundingClientRect();
      setViewport({
        width: rect.width,
        top: -rect.top,
        height: window.innerHeight,
      });
    };
    const schedule = () => {
      if (!frame) frame = requestAnimationFrame(update);
    };
    const observer = new ResizeObserver(schedule);
    observer.observe(container);
    window.addEventListener("scroll", schedule, { passive: true });
    window.addEventListener("resize", schedule);
    update();
    return () => {
      cancelAnimationFrame(frame);
      observer.disconnect();
      window.removeEventListener("scroll", schedule);
      window.removeEventListener("resize", schedule);
    };
  });

  // Reset cached heights for filtering and responsive reflow.
  createEffect(
    () => [props.items, width(), columns(), gap()],
    () => {
      setMeasurements(new Map());
    },
  );

  function Row(rowProps: { row: number }) {
    let element!: HTMLDivElement;
    createEffect(
      () => ({
        collection: props.items,
        measuredWidth: width(),
        count: columns(),
      }),
      ({ collection, measuredWidth }) => {
        const observer = new ResizeObserver(() => {
          if (collection !== props.items || measuredWidth !== width()) return;
          const height = element.getBoundingClientRect().height;
          setMeasurements((previous) => {
            if (previous.get(rowProps.row) === height) return previous;
            const next = new Map(previous);
            next.set(rowProps.row, height);
            return next;
          });
        });
        observer.observe(element);
        return () => observer.disconnect();
      },
    );
    return (
      <div
        ref={(el) => {
          element = el;
        }}
        style={{
          display: "grid",
          "grid-template-columns": `repeat(${columns()}, minmax(0, 1fr))`,
          gap: `${gap()}px`,
        }}
      >
        <For
          each={props.items.slice(
            rowProps.row * columns(),
            (rowProps.row + 1) * columns(),
          )}
        >
          {(item, index) =>
            props.children(item, () => rowProps.row * columns() + index())
          }
        </For>
      </div>
    );
  }

  return (
    <div
      ref={(el) => {
        container = el;
      }}
      class={props.class}
      style={{
        position: "relative",
        height: `${Math.max(0, offsets().positions.at(-1)! - gap())}px`,
        "overflow-anchor": "none",
      }}
    >
      <For each={visible()}>
        {(row) => (
          <div
            style={{
              position: "absolute",
              top: `${offsets().positions[row]}px`,
              left: 0,
              right: 0,
            }}
          >
            <Row row={row} />
          </div>
        )}
      </For>
    </div>
  );
}
