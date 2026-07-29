import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Check, Eraser, PenLine, Upload } from "lucide-react";
import { useRef, useState } from "react";
import { toast } from "sonner";
import {
  denySignatureRequestFn,
  fulfillSignatureRequestFn,
  signatureRequestQueryOptions,
} from "@/features/signatures/api";
import { uploadFile } from "@/features/storage/api";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";

export const Route = createFileRoute("/signature-requests/fulfill")({
  validateSearch: (search: Record<string, unknown>) => ({
    token: typeof search.token === "string" ? search.token : "",
  }),
  component: SignatureRequestPage,
});

const WIDTH = 1200;
const HEIGHT = 420;

function SignatureRequestPage() {
  const { token } = Route.useSearch();
  const requestId = token ? decodeTokenRequestId(token) : "";
  const requestQuery = useQuery({
    ...signatureRequestQueryOptions(requestId),
    enabled: Boolean(requestId),
  });
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const drawingRef = useRef(false);
  const lastPointRef = useRef<{ x: number; y: number } | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const [saving, setSaving] = useState(false);
  const [denying, setDenying] = useState(false);

  const point = (event: React.PointerEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current;
    if (!canvas) return null;
    const rect = canvas.getBoundingClientRect();
    return {
      x: (event.clientX - rect.left) * (canvas.width / rect.width),
      y: (event.clientY - rect.top) * (canvas.height / rect.height),
    };
  };

  const clear = () => {
    const canvas = canvasRef.current;
    const context = canvas?.getContext("2d");
    if (canvas && context) context.clearRect(0, 0, canvas.width, canvas.height);
  };

  const save = async () => {
    if (!token) return toast.error("Convite inválido ou incompleto");
    setSaving(true);
    try {
      let selected = file;
      if (!selected) {
        const canvas = canvasRef.current;
        if (!canvas) throw new Error("Área de assinatura indisponível");
        const blob = await new Promise<Blob | null>((resolve) =>
          canvas.toBlob(resolve, "image/png"),
        );
        if (!blob) throw new Error("Não foi possível gerar a assinatura");
        selected = new File([blob], "signature.png", { type: "image/png" });
      }
      const imageUrl = await uploadFile(selected, "signature-requests");
      const response = await fulfillSignatureRequestFn(token, imageUrl);
      if (!response.success)
        throw new Error(response.message || "Convite expirado");
      toast.success("Assinatura enviada com sucesso");
      clear();
      setFile(null);
    } catch (error) {
      toast.error(
        error instanceof Error
          ? error.message
          : "Não foi possível enviar a assinatura",
      );
    } finally {
      setSaving(false);
    }
  };

  const deny = async () => {
    if (!token) return;
    setDenying(true);
    try {
      const response = await denySignatureRequestFn(
        token,
        "Recusada pelo signatário",
      );
      if (!response.success)
        throw new Error(
          response.message || "Não foi possível recusar o convite",
        );
      toast.success("Solicitação recusada");
    } catch (error) {
      toast.error(
        error instanceof Error
          ? error.message
          : "Não foi possível recusar o convite",
      );
    } finally {
      setDenying(false);
    }
  };

  return (
    <main className="min-h-screen bg-muted/30 px-4 py-10 sm:py-16">
      <section className="mx-auto max-w-3xl space-y-6">
        <header className="space-y-3">
          <div className="mb-2 flex size-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <PenLine className="size-5" />
          </div>
          <h1 className="text-2xl font-semibold tracking-tight">
            Assine o convite
          </h1>
          <p className="text-sm text-muted-foreground">
            Desenhe sua assinatura no quadro abaixo ou envie uma imagem pronta.
            Este convite é pessoal e válido somente para esta solicitação.
          </p>
          {requestQuery.data && (
            <p className="text-sm text-muted-foreground">
              Solicitação para{" "}
              <strong>{requestQuery.data.signatory_name}</strong>
              {requestQuery.data.signatory_email
                ? ` (${requestQuery.data.signatory_email})`
                : ""}
              .
            </p>
          )}
        </header>
        <div className="space-y-5">
          {!token && (
            <p className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
              Link de convite inválido.
            </p>
          )}
          <div className="rounded-xl border bg-white p-2">
            <canvas
              ref={(canvas) => {
                canvasRef.current = canvas;
                if (canvas) {
                  canvas.width = WIDTH;
                  canvas.height = HEIGHT;
                }
              }}
              className="h-48 w-full touch-none rounded-lg"
              onPointerDown={(event) => {
                const current = point(event);
                if (!current) return;
                drawingRef.current = true;
                lastPointRef.current = current;
              }}
              onPointerMove={(event) => {
                if (!drawingRef.current) return;
                const current = point(event);
                const previous = lastPointRef.current;
                const context = canvasRef.current?.getContext("2d");
                if (!current || !previous || !context) return;
                context.lineWidth = 4;
                context.lineCap = "round";
                context.strokeStyle = "#111827";
                context.beginPath();
                context.moveTo(previous.x, previous.y);
                context.lineTo(current.x, current.y);
                context.stroke();
                lastPointRef.current = current;
              }}
              onPointerUp={() => {
                drawingRef.current = false;
                lastPointRef.current = null;
              }}
              onPointerLeave={() => {
                drawingRef.current = false;
                lastPointRef.current = null;
              }}
            />
          </div>
          <div className="flex flex-wrap gap-2">
            <Button type="button" variant="outline" onClick={clear}>
              <Eraser className="size-4" /> Limpar
            </Button>
            <label
              htmlFor="signature-file"
              className="inline-flex h-9 cursor-pointer items-center gap-2 rounded-md border px-4 text-sm font-medium hover:bg-muted"
            >
              <Upload className="size-4" /> Enviar imagem
              <Input
                id="signature-file"
                type="file"
                accept="image/png,image/jpeg,image/webp"
                className="hidden"
                onChange={(event) => setFile(event.target.files?.[0] ?? null)}
              />
            </label>
            {file && (
              <span className="self-center text-sm text-muted-foreground">
                {file.name}
              </span>
            )}
          </div>
          <Button
            className="w-full"
            disabled={!token || saving}
            onClick={() => void save()}
          >
            <Check className="size-4" />{" "}
            {saving ? "Enviando..." : "Confirmar assinatura"}
          </Button>
          <Button
            className="w-full"
            variant="ghost"
            disabled={!token || denying || saving}
            onClick={() => void deny()}
          >
            {denying ? "Recusando..." : "Recusar solicitação"}
          </Button>
        </div>
      </section>
    </main>
  );
}

function decodeTokenRequestId(token: string) {
  try {
    const payload = token.split(".")[1];
    if (!payload) return "";
    const decoded = JSON.parse(
      atob(payload.replace(/-/g, "+").replace(/_/g, "/")),
    ) as { request_id?: string };
    return decoded.request_id ?? "";
  } catch {
    return "";
  }
}
