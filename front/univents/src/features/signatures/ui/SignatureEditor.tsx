import { useNavigate } from "@tanstack/react-router";
import { Eraser, PenLine, Upload } from "lucide-react";
import { useRef, useState } from "react";
import { toast } from "sonner";
import { useCreateSignatureMutation } from "@/features/signatures/api/mutations";
import { SignatureCanvas } from "@/features/signatures/ui/SignatureCanvas";
import { SignatureImageSelector } from "@/features/signatures/ui/SignatureImageSelector";
import { uploadFile } from "@/features/storage/api";
import { Button } from "@/shared/ui/shadcn/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/shadcn/card";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";

type Mode = "draw" | "upload";

export interface SignatureEditorProps {
  eventId: string;
  editionId: string;
}

export function SignatureEditor({ eventId, editionId }: SignatureEditorProps) {
  const navigate = useNavigate();
  const [signatoryName, setSignatoryName] = useState("");
  const [signatoryTitle, setSignatoryTitle] = useState("");
  const [signatoryEmail, setSignatoryEmail] = useState("");
  const [signatoryUserId, setSignatoryUserId] = useState("");
  const [mode, setMode] = useState<Mode>("draw");
  const [importedFile, setImportedFile] = useState<File | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  const clearCanvas = () => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext("2d");
    if (!canvas || !ctx) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);
  };

  const saveMutation = useCreateSignatureMutation();

  return (
    <div className="mx-auto max-w-6xl px-4">
      <div className="mb-6 space-y-1">
        <h1 className="text-2xl font-semibold">Nova assinatura</h1>
        <p className="max-w-2xl text-sm text-muted-foreground">
          Crie uma assinatura desenhando no canvas ou importando uma imagem
          pronta.
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,360px)_minmax(0,1fr)]">
        <Card className="h-fit">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-semibold">
              Configuração
            </CardTitle>
            <CardDescription className="text-xs">
              Dados do signatário e origem da assinatura.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-1.5">
              <Label className="text-xs text-muted-foreground">
                Nome do signatário
              </Label>
              <Input
                value={signatoryName}
                onChange={(e) => setSignatoryName(e.target.value)}
                placeholder="Nome completo"
              />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs text-muted-foreground">Cargo</Label>
              <Input
                value={signatoryTitle}
                onChange={(e) => setSignatoryTitle(e.target.value)}
                placeholder="Cargo ou função"
              />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs text-muted-foreground">E-mail</Label>
              <Input
                type="email"
                value={signatoryEmail}
                onChange={(e) => setSignatoryEmail(e.target.value)}
                placeholder="nome@exemplo.com"
              />
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs text-muted-foreground">
                ID do usuário
              </Label>
              <Input
                value={signatoryUserId}
                onChange={(e) => setSignatoryUserId(e.target.value)}
                placeholder="Opcional"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-xs text-muted-foreground">Modo</Label>
              <div className="grid grid-cols-2 gap-2">
                <Button
                  type="button"
                  variant={mode === "draw" ? "default" : "outline"}
                  className="h-9 gap-2"
                  onClick={() => setMode("draw")}
                >
                  <PenLine className="size-4" />
                  Desenhar
                </Button>
                <Button
                  type="button"
                  variant={mode === "upload" ? "default" : "outline"}
                  className="h-9 gap-2"
                  onClick={() => setMode("upload")}
                >
                  <Upload className="size-4" />
                  Importar
                </Button>
              </div>
            </div>

            <div className="rounded-2xl border bg-muted/20 p-3 text-xs text-muted-foreground">
              <p className="font-medium text-foreground">Dica</p>
              <p className="mt-1">
                Use uma imagem com fundo limpo para melhor legibilidade ou
                desenhe diretamente aqui.
              </p>
            </div>

            <Button
              type="button"
              className="h-9 w-full gap-2"
              onClick={async () => {
                try {
                  const trimmedName = signatoryName.trim();
                  if (!trimmedName) {
                    toast.error("Nome do signatário é obrigatório");
                    return;
                  }

                  let url: string | null = null;
                  if (mode === "draw") {
                    const canvas = canvasRef.current;
                    if (!canvas) throw new Error("Canvas indisponível");
                    const blob = await new Promise<Blob | null>((resolve) =>
                      canvas.toBlob(resolve, "image/png"),
                    );
                    if (!blob)
                      throw new Error("Falha ao gerar imagem da assinatura");
                    const file = new File(
                      [blob],
                      `${Date.now()}-signature.png`,
                      { type: "image/png" },
                    );
                    url = await uploadFile(
                      file,
                      `events/${eventId}/editions/${editionId}/signatures`,
                    );
                  } else if (importedFile) {
                    url = await uploadFile(
                      importedFile,
                      `events/${eventId}/editions/${editionId}/signatures`,
                    );
                  }

                  if (!url) {
                    toast.error("Selecione ou desenhe uma assinatura");
                    return;
                  }

                  const res = await saveMutation.mutateAsync({
                    eventId,
                    editionId,
                    data: {
                      signatory_name: trimmedName,
                      signatory_title: signatoryTitle.trim() || undefined,
                      signatory_email: signatoryEmail.trim() || undefined,
                      signatory_user_id: signatoryUserId.trim() || undefined,
                      image_url: url,
                    },
                  });

                  if (res.success) {
                    toast.success("Assinatura criada com sucesso");
                    void navigate({
                      to: "/admin/events/$eventId/editions/$editionId/signatures",
                      params: { eventId, editionId },
                    });
                    return;
                  }

                  toast.error(res.message || "Erro ao criar assinatura");
                } catch (error) {
                  toast.error(
                    error instanceof Error
                      ? error.message
                      : "Erro ao criar assinatura",
                  );
                }
              }}
              disabled={saveMutation.isPending}
            >
              Salvar assinatura
            </Button>
          </CardContent>
        </Card>

        <Card className="min-w-0">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-semibold">Prévia</CardTitle>
            <CardDescription className="text-xs">
              Veja o que vai ser salvo antes de concluir.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {mode === "upload" ? (
              <SignatureImageSelector
                file={importedFile}
                onChange={setImportedFile}
              />
            ) : (
              <div className="space-y-3">
                <div className="flex items-center justify-between gap-3">
                  <div className="space-y-0.5">
                    <p className="text-sm font-medium text-foreground">
                      Desenho
                    </p>
                    <p className="text-xs text-muted-foreground">
                      Use o mouse ou toque para assinar no quadro abaixo.
                    </p>
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={clearCanvas}
                  >
                    <Eraser className="size-4" />
                    Limpar
                  </Button>
                </div>
                <div className="rounded-2xl border bg-muted/10 p-2">
                  <SignatureCanvas
                    ref={canvasRef}
                    className="h-44 w-full min-w-full"
                  />
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
