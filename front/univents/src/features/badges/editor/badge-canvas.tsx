import { Maximize, Minus, Plus } from "lucide-react";
import QRCode from "qrcode";
import type { CSSProperties, PointerEvent as ReactPointerEvent } from "react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useElementSize } from "../../certifications/editor/hooks/use-element-size";
import {
  type TextElementAdapter,
  TextElementView,
} from "../../certifications/editor/ui/elements/text-element-view";
import type { BadgeElement, BadgeTemplateCreate } from "../model";
import { BadgeElementFrame } from "./badge-element-frame";

const CHECKERBOARD: CSSProperties = {
  backgroundColor: "#e5e7eb",
  backgroundImage:
    "linear-gradient(45deg, #d1d5db 25%, transparent 25%), linear-gradient(-45deg, #d1d5db 25%, transparent 25%), linear-gradient(45deg, transparent 75%, #d1d5db 75%), linear-gradient(-45deg, transparent 75%, #d1d5db 75%)",
  backgroundSize: "16px 16px",
  backgroundPosition: "0 0, 0 8px, 8px -8px, -8px 0",
};

function QrPreview({
  element,
}: {
  element: Extract<BadgeElement, { type: "qr" }>;
}) {
  const matrix = useMemo(
    () => QRCode.create("https://univents.app/check-in/example").modules,
    [],
  );
  const margin = 2;
  const size = matrix.size + margin * 2;
  const modules = [];
  for (let row = 0; row < matrix.size; row++) {
    for (let column = 0; column < matrix.size; column++) {
      if (!matrix.get(row, column)) continue;
      const x = column + margin;
      const y = row + margin;
      modules.push(
        element.style === "dots" ? (
          <circle
            key={`${row}-${column}`}
            cx={x + 0.5}
            cy={y + 0.5}
            r="0.42"
            fill={element.foreground}
          />
        ) : (
          <rect
            key={`${row}-${column}`}
            x={x}
            y={y}
            width="1"
            height="1"
            rx={element.style === "rounded" ? 0.28 : 0}
            fill={element.foreground}
          />
        ),
      );
    }
  }
  return (
    <svg
      role="img"
      aria-label="QR Code de check-in"
      className="size-full"
      viewBox={`0 0 ${size} ${size}`}
      style={{ background: element.background }}
      shapeRendering={
        element.style === "square" ? "crispEdges" : "geometricPrecision"
      }
    >
      {modules}
    </svg>
  );
}

function BadgeElementView({
  element,
  editing,
  adapter,
}: {
  element: BadgeElement;
  editing: boolean;
  adapter: TextElementAdapter;
}) {
  if (element.type === "image")
    return (
      <img
        src={element.src}
        alt=""
        draggable={false}
        className="size-full"
        style={{
          objectFit: element.fit,
          borderRadius: element.radius,
          opacity: element.opacity,
        }}
      />
    );
  if (element.type === "qr") return <QrPreview element={element} />;
  return (
    <TextElementView element={element} editing={editing} adapter={adapter} />
  );
}

interface BadgeCanvasProps {
  design: BadgeTemplateCreate["design_data"];
  selectedId: string | null;
  onSelect: (id: string | null) => void;
  onChangeElement: (id: string, changes: Partial<BadgeElement>) => void;
  onDeleteElement: (id: string) => void;
  textAdapter: TextElementAdapter;
}

export function BadgeCanvas({
  design,
  selectedId,
  onSelect,
  onChangeElement,
  onDeleteElement,
  textAdapter,
}: BadgeCanvasProps) {
  const { ref: stageRef, size: viewport } = useElementSize<HTMLDivElement>();
  const [zoom, setZoom] = useState(1);
  const [editingId, setEditingId] = useState<string | null>(null);
  const effectiveTextAdapter = useMemo(
    () => ({ ...textAdapter, stopEditing: () => setEditingId(null) }),
    [textAdapter],
  );
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const panStart = useRef<{
    pointerX: number;
    pointerY: number;
    x: number;
    y: number;
  } | null>(null);

  const fitScale = useMemo(() => {
    if (!viewport.width || !viewport.height) return 0.62;
    return Math.min(
      (viewport.width - 96) / design.canvas.width,
      (viewport.height - 96) / design.canvas.height,
      2,
    );
  }, [design.canvas, viewport]);
  const scale = fitScale * zoom;

  useEffect(() => {
    const stage = stageRef.current;
    if (!stage) return;
    const wheel = (event: WheelEvent) => {
      event.preventDefault();
      setZoom((value) =>
        Math.min(4, Math.max(0.25, value * (event.deltaY < 0 ? 1.1 : 0.9))),
      );
    };
    stage.addEventListener("wheel", wheel, { passive: false });
    return () => stage.removeEventListener("wheel", wheel);
  }, [stageRef]);

  function startPan(event: ReactPointerEvent<HTMLDivElement>) {
    const target = event.target as HTMLElement;
    if (target.closest("[data-badge-element], [data-canvas-controls]")) return;
    onSelect(null);
    panStart.current = {
      pointerX: event.clientX,
      pointerY: event.clientY,
      ...pan,
    };
    event.currentTarget.setPointerCapture(event.pointerId);
  }

  return (
    <main
      ref={stageRef}
      className="relative min-w-90 flex-1 touch-none overflow-hidden"
      style={CHECKERBOARD}
      onPointerDown={startPan}
      onPointerMove={(event) => {
        const start = panStart.current;
        if (!start) return;
        setPan({
          x: start.x + event.clientX - start.pointerX,
          y: start.y + event.clientY - start.pointerY,
        });
      }}
      onPointerUp={() => {
        panStart.current = null;
      }}
      onPointerCancel={() => {
        panStart.current = null;
      }}
    >
      <div
        data-canvas-controls
        className="absolute bottom-4 left-1/2 z-50 flex -translate-x-1/2 items-center gap-1 rounded-lg border border-muted bg-card p-1 shadow-lg"
      >
        <button
          type="button"
          className="rounded p-1.5 hover:bg-muted"
          aria-label="Diminuir zoom"
          onClick={() => setZoom((value) => Math.max(0.25, value / 1.2))}
        >
          <Minus className="size-4" />
        </button>
        <span className="w-14 text-center text-xs tabular-nums">
          {Math.round(scale * 100)}%
        </span>
        <button
          type="button"
          className="rounded p-1.5 hover:bg-muted"
          aria-label="Aumentar zoom"
          onClick={() => setZoom((value) => Math.min(4, value * 1.2))}
        >
          <Plus className="size-4" />
        </button>
        <button
          type="button"
          className="rounded p-1.5 hover:bg-muted"
          aria-label="Ajustar à tela"
          onClick={() => {
            setZoom(1);
            setPan({ x: 0, y: 0 });
          }}
        >
          <Maximize className="size-4" />
        </button>
      </div>

      <div
        className="absolute left-1/2 top-1/2 shadow-2xl ring-1 ring-muted"
        style={{
          width: design.canvas.width * scale,
          height: design.canvas.height * scale,
          transform: `translate(calc(-50% + ${pan.x}px), calc(-50% + ${pan.y}px))`,
        }}
      >
        <div
          className="relative overflow-hidden"
          style={{
            width: design.canvas.width,
            height: design.canvas.height,
            transform: `scale(${scale})`,
            transformOrigin: "top left",
            backgroundColor: design.backgroundColor,
            backgroundImage: design.background
              ? `url(${design.background})`
              : undefined,
            backgroundSize: "cover",
            backgroundPosition: "center",
          }}
        >
          {design.elements.map((element, index) => (
            <BadgeElementFrame
              key={element.id}
              bounds={element}
              scale={scale}
              zIndex={index + 1}
              canvas={design.canvas}
              selected={selectedId === element.id}
              editing={editingId === element.id}
              onSelect={() => onSelect(element.id)}
              onDoubleClick={
                element.type === "text"
                  ? () => setEditingId(element.id)
                  : undefined
              }
              onChangeBounds={(bounds) => onChangeElement(element.id, bounds)}
              onDelete={
                element.type === "qr"
                  ? undefined
                  : () => onDeleteElement(element.id)
              }
            >
              <BadgeElementView
                element={element}
                editing={editingId === element.id}
                adapter={effectiveTextAdapter}
              />
            </BadgeElementFrame>
          ))}
        </div>
      </div>
    </main>
  );
}
