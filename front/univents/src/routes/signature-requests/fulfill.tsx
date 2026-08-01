import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Check, Eraser, FileUp, PenLine } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import {
  denySignatureRequestFn,
  fulfillSignatureRequestFn,
  signatureRequestQueryOptions,
} from "@/features/signatures/api";
import type { SignatureRequestI } from "@/features/signatures/model";
import { SignatureCanvas } from "@/features/signatures/ui/SignatureCanvas";
import { SignatureRequestConfirmation } from "@/features/signatures/ui/SignatureRequestConfirmation";
import { SignatureRequestStatus } from "@/features/signatures/ui/SignatureRequestStatus";
import { uploadFile } from "@/features/storage/api";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import { Checkbox } from "@/shared/ui/shadcn/checkbox";
import { Input } from "@/shared/ui/shadcn/input";

export const Route = createFileRoute("/signature-requests/fulfill")({
  validateSearch: (search: Record<string, unknown>) => ({
    token: typeof search.token === "string" ? search.token : "",
  }),
  component: SignatureRequestPage,
});

type SignMethod = "draw" | "import";

function SignatureRequestPage() {
  const { token } = Route.useSearch();
  const requestId = token ? decodeTokenRequestId(token) : "";

  const requestQuery = useQuery({
    ...signatureRequestQueryOptions(requestId),
    enabled: Boolean(requestId),
  });

  const [signedAt, setSignedAt] = useState<string | null>(null);

  if (signedAt) {
    return <SignatureRequestConfirmation timestamp={signedAt} />;
  }

  if (!token || !requestId) {
    return (
      <SignatureRequestStatus
        title="Solicitação inválida"
        message="O link da solicitação de assinatura está incompleto ou inválido."
      />
    );
  }

  if (requestQuery.isLoading) {
    return (
      <SignatureRequestStatus
        title="Carregando solicitação"
        message="Aguarde enquanto buscamos os dados da assinatura."
        loading
      />
    );
  }

  if (requestQuery.isError) {
    return (
      <SignatureRequestStatus
        title="Não foi possível carregar a solicitação"
        message="Ocorreu um erro ao buscar a solicitação. Verifique o link e tente novamente."
      />
    );
  }

  if (!requestQuery.data) {
    return (
      <SignatureRequestStatus
        title="Solicitação não encontrada"
        message="Esta solicitação pode ter sido removida, cancelada ou já ter expirado."
      />
    );
  }

  if (requestQuery.data.status !== "pending") {
    const statusMessages = {
      completed: {
        title: "Solicitação já concluída",
        message: "Esta solicitação de assinatura já foi concluída.",
      },
      expired: {
        title: "Solicitação expirada",
        message: "O prazo desta solicitação de assinatura terminou.",
      },
      cancelled: {
        title: "Solicitação cancelada",
        message: "Esta solicitação de assinatura foi cancelada.",
      },
    } as const;
    const statusMessage = statusMessages[requestQuery.data.status];

    return (
      <SignatureRequestStatus
        title={statusMessage.title}
        message={statusMessage.message}
      />
    );
  }

  return (
    <SignatureForm
      token={token}
      request={requestQuery.data}
      onSigned={(timestamp) => setSignedAt(timestamp)}
    />
  );
}

function SignatureForm({
  token,
  request,
  onSigned,
}: {
  token: string;
  request: SignatureRequestI;
  onSigned: (timestamp: string) => void;
}) {
  const navigate = Route.useNavigate();
  const canvasRef = useRef<HTMLCanvasElement>(null);

  const [method, setMethod] = useState<SignMethod>("draw");
  const [file, setFile] = useState<File | null>(null);
  const [filePreview, setFilePreview] = useState<string | null>(null);
  const [dragOver, setDragOver] = useState(false);
  const [agreed, setAgreed] = useState(false);
  const [saving, setSaving] = useState(false);
  const [denying, setDenying] = useState(false);

  useEffect(() => {
    if (!file) {
      setFilePreview(null);
      return;
    }
    const previewUrl = URL.createObjectURL(file);
    setFilePreview(previewUrl);
    return () => URL.revokeObjectURL(previewUrl);
  }, [file]);

  const clearCanvas = () => {
    const canvas = canvasRef.current;
    const context = canvas?.getContext("2d");
    if (canvas && context) context.clearRect(0, 0, canvas.width, canvas.height);
  };

  const buildSignatureFile = async () => {
    if (method === "import") {
      if (!file) throw new Error("Selecione um arquivo de imagem");
      return file;
    }
    const canvas = canvasRef.current;
    if (!canvas) throw new Error("Área de assinatura indisponível");
    const blob = await new Promise<Blob | null>((resolve) =>
      canvas.toBlob(resolve, "image/png"),
    );
    if (!blob) throw new Error("Não foi possível gerar a assinatura");
    return new File([blob], "signature.png", { type: "image/png" });
  };

  const save = async () => {
    if (!token) return toast.error("Convite inválido ou incompleto");
    if (!agreed)
      return toast.error("Confirme que você tem autoridade para assinar");
    setSaving(true);
    try {
      const selected = await buildSignatureFile();
      const imageUrl = await uploadFile(selected, "signature-requests");
      const response = await fulfillSignatureRequestFn(token, imageUrl);
      if (!response.success)
        throw new Error(response.message || "Convite expirado");
      toast.success("Assinatura enviada com sucesso");
      onSigned(new Date().toISOString());
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
      await navigate({ to: "/" });
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

  const expiresLabel = request.expires_at
    ? new Date(request.expires_at).toLocaleDateString("pt-BR", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
      })
    : null;

  const onDropFile = (event: React.DragEvent<HTMLLabelElement>) => {
    event.preventDefault();
    setDragOver(false);
    const dropped = event.dataTransfer.files?.[0];
    if (dropped) setFile(dropped);
  };

  return (
    <main
      className={cn(
        "flex min-h-screen items-center justify-center",
        "bg-muted px-3 pt-6 pb-28 sm:px-4 sm:pt-10 sm:pb-24",
      )}
    >
      <section
        className={cn(
          "w-full min-w-0 max-w-2xl rounded-2xl",
          "border border-border bg-card shadow-sm",
        )}
      >
        <header
          className={cn(
            "flex flex-col gap-2 border-b border-border px-4 py-4",
            "sm:flex-row sm:items-start sm:justify-between sm:gap-4 sm:px-6 sm:py-5",
          )}
        >
          <div>
            <h1 className="text-base font-semibold text-card-foreground sm:text-lg">
              Solicitação de assinatura eletrônica
            </h1>
            <p className="mt-1 text-xs text-muted-foreground">
              ID do documento: {request.id}
            </p>
          </div>
          {expiresLabel && (
            <span
              className={cn(
                "w-fit whitespace-nowrap rounded-full bg-destructive/10",
                "px-3 py-1 text-xs font-medium text-destructive",
              )}
            >
              Expira em {expiresLabel}
            </span>
          )}
        </header>

        <div className="grid grid-cols-1 gap-4 px-4 py-4 sm:grid-cols-2 sm:px-6 sm:py-5">
          <Field label="Nome do signatário" value={request.signatory_name} />
          <Field
            label="Cargo do signatário"
            value={request.signatory_title ?? "—"}
          />
        </div>

        <div className="px-4 sm:px-6">
          <div className="flex gap-1 border-b border-border">
            <MethodTab
              active={method === "draw"}
              onClick={() => setMethod("draw")}
              icon={<PenLine className="size-4" />}
              label="Desenhar"
            />
            <MethodTab
              active={method === "import"}
              onClick={() => setMethod("import")}
              icon={<FileUp className="size-4" />}
              label="Importar"
            />
          </div>
        </div>

        <div className="px-4 py-4 sm:px-6 sm:py-5">
          {method === "draw" && (
            <div className="relative rounded-lg border border-border bg-muted/40">
              <SignatureCanvas
                ref={canvasRef}
                className="h-36 w-full sm:h-40"
              />
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="absolute right-2 top-2 sm:right-3 sm:top-3"
                onClick={clearCanvas}
              >
                <Eraser className="size-3.5" /> Limpar
              </Button>
            </div>
          )}

          {method === "import" && (
            <label
              htmlFor="signature-file"
              onDragOver={(event) => {
                event.preventDefault();
                setDragOver(true);
              }}
              onDragLeave={() => setDragOver(false)}
              onDrop={onDropFile}
              className={cn(
                "flex h-36 cursor-pointer flex-col items-center justify-center gap-2",
                "rounded-lg border border-dashed text-center text-sm sm:h-40",
                dragOver
                  ? "border-primary bg-primary/10 text-primary"
                  : "border-border bg-muted/40 text-muted-foreground hover:bg-muted",
              )}
            >
              {filePreview ? (
                <img
                  src={filePreview}
                  alt="Prévia da assinatura importada"
                  className="max-h-28 max-w-[85%] object-contain"
                />
              ) : (
                <>
                  <FileUp className="size-5" />
                  <span className="px-4">
                    Arraste uma imagem ou clique para enviar
                  </span>
                </>
              )}
              {file && (
                <span className="max-w-full truncate px-4 text-xs">
                  {file.name}
                </span>
              )}
              <Input
                id="signature-file"
                type="file"
                accept="image/png,image/jpeg,image/webp"
                className="hidden"
                onChange={(event) => setFile(event.target.files?.[0] ?? null)}
              />
            </label>
          )}
        </div>

        <div className="flex items-start gap-2 px-4 pb-4 sm:px-6 sm:pb-5">
          <Checkbox
            id="authority"
            checked={agreed}
            onCheckedChange={(value) => setAgreed(value === true)}
            className="mt-0.5 shrink-0"
          />
          <label
            htmlFor="authority"
            className="text-xs leading-relaxed text-muted-foreground"
          >
            Entendo que esta é uma assinatura digital juridicamente vinculante e
            que tenho autoridade para assinar este documento em nome da entidade
            representada.
          </label>
        </div>

        <footer
          className={cn(
            "flex  gap-3 border-t border-border p-4",
            "sm:flex-col sm:items-stretch sm:px-6 md:flex-row md:items-center",
            "md:justify-between gap-4",
          )}
        >
          <button
            type="button"
            onClick={() => void deny()}
            disabled={!token || denying || saving}
            className={cn(
              "text-sm font-medium text-muted-foreground",
              "hover:text-foreground disabled:opacity-50",
            )}
          >
            {denying ? "Cancelando..." : "Cancelar e voltar"}
          </button>
          <Button
            disabled={!token || saving || !agreed}
            onClick={() => void save()}
            className=""
          >
            <Check className="size-4" />{" "}
            {saving ? "Assinando..." : "Assinar documento"}
          </Button>
        </footer>
      </section>
    </main>
  );
}

function Field({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-[11px] font-medium uppercase tracking-wide text-muted-foreground">
        {label}
      </p>
      <p className="mt-0.5 text-sm font-medium text-card-foreground">{value}</p>
    </div>
  );
}

function MethodTab({
  active,
  onClick,
  icon,
  label,
}: {
  active: boolean;
  onClick: () => void;
  icon: React.ReactNode;
  label: string;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "flex items-center gap-1.5 border-b-2 px-3 py-2",
        "text-sm font-medium transition-colors",
        active
          ? "border-primary text-primary"
          : "border-transparent text-muted-foreground hover:text-foreground",
      )}
    >
      {icon} {label}
    </button>
  );
}

function decodeTokenRequestId(token: string) {
  try {
    const payload = token.split(".")[1];
    if (!payload) return "";
    const decoded = JSON.parse(
      atob(payload.replace(/-/g, "+").replace(/_/g, "/")),
    ) as {
      request_id?: string;
    };
    return decoded.request_id ?? "";
  } catch {
    return "";
  }
}
