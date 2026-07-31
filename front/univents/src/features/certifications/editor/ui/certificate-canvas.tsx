import { Maximize, Minus, Plus } from "lucide-react";
import type { CSSProperties, PointerEvent as ReactPointerEvent } from "react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useElementSize } from "../hooks/use-element-size";
import { certificateEditorActions, useCertificateEditorState } from "../store";
import { CertificateElementFrame } from "./elements/certificate-element-frame";
import { CertificateElementView } from "./elements/certificate-element-view";

export const CERTIFICATE_CANVAS_DOM_ID = "certificate-canvas-root";

const MIN_ZOOM = 0.25;
const MAX_ZOOM = 4;

const CHECKERBOARD_STYLE: CSSProperties = {
  backgroundColor: "#e5e7eb",
  backgroundImage:
    "linear-gradient(45deg, #d1d5db 25%, transparent 25%), linear-gradient(-45deg, #d1d5db 25%, transparent 25%), linear-gradient(45deg, transparent 75%, #d1d5db 75%), linear-gradient(-45deg, transparent 75%, #d1d5db 75%)",
  backgroundSize: "16px 16px",
  backgroundPosition: "0 0, 0 8px, 8px -8px, -8px 0",
};

export function CertificateCanvas() {
  const elements = useCertificateEditorState(
    (state) => state.draft.design_data.elements,
  );
  const backgroundUrl = useCertificateEditorState(
    (state) => state.draft.design_data.background,
  );
  const canvas = useCertificateEditorState((state) => state.canvas);
  const selectedElementId = useCertificateEditorState(
    (state) => state.selectedElementId,
  );
  const editingElementId = useCertificateEditorState(
    (state) => state.editingElementId,
  );

  const { ref: stageRef, size: viewport } = useElementSize<HTMLDivElement>();
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const [panning, setPanning] = useState(false);
  const panStart = useRef<{
    pointerX: number;
    pointerY: number;
    x: number;
    y: number;
  } | null>(null);
  const lastFitScale = useRef(1);

  const fitScale = useMemo(() => {
    if (viewport.width <= 0 || viewport.height <= 0) {
      return lastFitScale.current;
    }

    const padding = 96;
    const availableWidth = Math.max(120, viewport.width - padding);
    const availableHeight = Math.max(120, viewport.height - padding);
    const nextScale = Math.min(
      availableWidth / canvas.width,
      availableHeight / canvas.height,
      2,
    );

    if (Number.isFinite(nextScale) && nextScale > 0) {
      lastFitScale.current = nextScale;
    }
    return lastFitScale.current;
  }, [canvas.height, canvas.width, viewport.height, viewport.width]);

  const scale = fitScale * zoom;

  useEffect(() => {
    const stage = stageRef.current;
    if (!stage) return;

    const handleWheel = (event: WheelEvent) => {
      event.preventDefault();
      const factor = event.deltaY < 0 ? 1.1 : 0.9;
      setZoom((value) =>
        Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, value * factor)),
      );
    };

    stage.addEventListener("wheel", handleWheel, { passive: false });
    return () => stage.removeEventListener("wheel", handleWheel);
  }, [stageRef]);

  function startPan(event: ReactPointerEvent<HTMLDivElement>) {
    const target = event.target as HTMLElement;
    if (target.closest("[data-certificate-element], [data-canvas-controls]")) {
      return;
    }

    certificateEditorActions.selectElement(null);
    certificateEditorActions.stopEditing();
    panStart.current = {
      pointerX: event.clientX,
      pointerY: event.clientY,
      ...pan,
    };
    setPanning(true);
    event.currentTarget.setPointerCapture(event.pointerId);
  }

  function finishPan(event: ReactPointerEvent<HTMLDivElement>) {
    panStart.current = null;
    setPanning(false);
    if (event.currentTarget.hasPointerCapture(event.pointerId)) {
      event.currentTarget.releasePointerCapture(event.pointerId);
    }
  }

  const canvasStyle: CSSProperties = {
    width: canvas.width,
    height: canvas.height,
    transform: `scale(${scale})`,
    transformOrigin: "top left",
    backgroundColor: "#ffffff",
    ...(backgroundUrl
      ? {
          backgroundImage: `url(${backgroundUrl})`,
          backgroundSize: "cover",
          backgroundPosition: "center",
          backgroundRepeat: "no-repeat",
        }
      : {}),
  };

  return (
    <div
      ref={stageRef}
      className={
        "relative flex min-h-0 min-w-0 flex-1 touch-none items-center justify-center overflow-hidden p-12 " +
        (panning ? "cursor-grabbing" : "cursor-grab")
      }
      style={CHECKERBOARD_STYLE}
      onPointerDown={startPan}
      onPointerMove={(event) => {
        const start = panStart.current;
        if (!start) return;
        setPan({
          x: start.x + event.clientX - start.pointerX,
          y: start.y + event.clientY - start.pointerY,
        });
      }}
      onPointerUp={finishPan}
      onPointerCancel={finishPan}
    >
      <div
        data-canvas-controls
        className="absolute bottom-4 left-1/2 z-50 flex -translate-x-1/2 items-center gap-1 rounded-lg border border-border bg-card p-1 text-card-foreground shadow-lg"
      >
        <button
          type="button"
          className="rounded p-1.5 hover:bg-muted"
          title="Diminuir zoom"
          aria-label="Diminuir zoom"
          onClick={() => setZoom((value) => Math.max(MIN_ZOOM, value / 1.2))}
        >
          <Minus className="size-4" />
        </button>
        <span className="w-14 text-center text-xs tabular-nums">
          {Math.round(scale * 100)}%
        </span>
        <button
          type="button"
          className="rounded p-1.5 hover:bg-muted"
          title="Aumentar zoom"
          aria-label="Aumentar zoom"
          onClick={() => setZoom((value) => Math.min(MAX_ZOOM, value * 1.2))}
        >
          <Plus className="size-4" />
        </button>
        <button
          type="button"
          className="rounded p-1.5 hover:bg-muted"
          title="Ajustar à tela e centralizar"
          aria-label="Ajustar à tela e centralizar"
          onClick={() => {
            setZoom(1);
            setPan({ x: 0, y: 0 });
          }}
        >
          <Maximize className="size-4" />
        </button>
      </div>

      <div
        className="relative shrink-0 shadow-2xl ring-1 ring-black/5"
        style={{
          width: canvas.width * scale,
          height: canvas.height * scale,
          transform: `translate(${pan.x}px, ${pan.y}px)`,
        }}
      >
        <div
          id={CERTIFICATE_CANVAS_DOM_ID}
          className="relative overflow-hidden"
          style={canvasStyle}
        >
          {elements.map((element, elementIndex) => (
            <CertificateElementFrame
              key={element.id}
              type={element.type}
              bounds={element}
              scale={scale}
              zIndex={elementIndex + 1}
              canvas={canvas}
              selected={selectedElementId === element.id}
              editing={editingElementId === element.id}
              deletable={element.type !== "hash"}
              onSelect={() => {
                if (editingElementId && editingElementId !== element.id) {
                  certificateEditorActions.stopEditing();
                }
                certificateEditorActions.selectElement(element.id);
              }}
              onDoubleClick={
                element.type === "text"
                  ? () => certificateEditorActions.startEditing(element.id)
                  : undefined
              }
              onChangeBounds={(bounds) =>
                certificateEditorActions.updateElementBounds(element.id, bounds)
              }
              onDelete={() =>
                certificateEditorActions.removeElement(element.id)
              }
            >
              <CertificateElementView
                element={element}
                editing={editingElementId === element.id}
              />
            </CertificateElementFrame>
          ))}
        </div>
      </div>
    </div>
  );
}
