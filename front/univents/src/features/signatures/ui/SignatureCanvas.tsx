import {
  forwardRef,
  useImperativeHandle,
  useLayoutEffect,
  useRef,
} from "react";
import { cn } from "@/shared/lib/utils";

export const SIGNATURE_CANVAS_WIDTH = 1200;
export const SIGNATURE_CANVAS_HEIGHT = 420;

export const SignatureCanvas = forwardRef<
  HTMLCanvasElement,
  {
    className?: string;
  }
>(({ className }, forwardedRef) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useImperativeHandle(
    forwardedRef,
    () => canvasRef.current as HTMLCanvasElement,
    [],
  );

  useLayoutEffect(() => {
    const canvas = canvasRef.current;
    const context = canvas?.getContext("2d");
    if (!canvas || !context) return;
    canvas.width = SIGNATURE_CANVAS_WIDTH;
    canvas.height = SIGNATURE_CANVAS_HEIGHT;
    context.lineWidth = 3;
    context.lineCap = "round";
    context.lineJoin = "round";
    context.strokeStyle = "#111827";
  }, []);

  const point = (event: React.PointerEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current;
    if (!canvas) return null;
    const rect = canvas.getBoundingClientRect();
    return {
      x: (event.clientX - rect.left) * (canvas.width / rect.width),
      y: (event.clientY - rect.top) * (canvas.height / rect.height),
    };
  };

  const drawing = useRef(false);
  const lastPoint = useRef<{ x: number; y: number } | null>(null);
  const stopDrawing = () => {
    drawing.current = false;
    lastPoint.current = null;
  };

  return (
    <canvas
      ref={canvasRef}
      className={cn("touch-none rounded-xl bg-white", className)}
      style={{
        aspectRatio: `${SIGNATURE_CANVAS_WIDTH} / ${SIGNATURE_CANVAS_HEIGHT}`,
      }}
      onPointerDown={(event) => {
        const current = point(event);
        if (!current) return;
        drawing.current = true;
        lastPoint.current = current;
      }}
      onPointerMove={(event) => {
        if (!drawing.current) return;
        const current = point(event);
        const previous = lastPoint.current;
        const context = canvasRef.current?.getContext("2d");
        if (!current || !previous || !context) return;
        context.beginPath();
        context.moveTo(previous.x, previous.y);
        context.lineTo(current.x, current.y);
        context.stroke();
        lastPoint.current = current;
      }}
      onPointerUp={stopDrawing}
      onPointerLeave={stopDrawing}
      onPointerCancel={stopDrawing}
    />
  );
});
SignatureCanvas.displayName = "SignatureCanvas";
