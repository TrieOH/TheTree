import type { JSX } from "@solidjs/web";
import { Show, createSignal } from "solid-js";
import Loader2Icon from "~icons/lucide/loader-2";
import MailIcon from "~icons/lucide/mail";

import { Button, Dialog, Input, Label } from "@trieoh/ui-solid";
import { toast } from "@/shared/ui/toast";
import { useCreateSignatureRequestMutation } from "../api/mutations";

const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface CreateSignatureRequestModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  editionId: string;
  onSuccess?: () => void;
}

export function CreateSignatureRequestModal(
  props: CreateSignatureRequestModalProps,
): JSX.Element {
  const [signatoryName, setSignatoryName] = createSignal("");
  const [signatoryTitle, setSignatoryTitle] = createSignal("");
  const [signatoryEmail, setSignatoryEmail] = createSignal("");
  const [expiresInDays, setExpiresInDays] = createSignal<number>(7);
  const [isSubmitting, setIsSubmitting] = createSignal(false);

  const createRequestMutation = useCreateSignatureRequestMutation();

  const resetForm = () => {
    setSignatoryName("");
    setSignatoryTitle("");
    setSignatoryEmail("");
    setExpiresInDays(7);
  };

  const handleSubmit = async () => {
    const name = signatoryName().trim();
    const email = signatoryEmail().trim();

    if (!name || name.length < 2) {
      toast.error("Informe o nome do signatário (mínimo 2 caracteres).");
      return;
    }

    if (!email || !email.includes("@")) {
      toast.error("Informe um e-mail válido para envio do convite.");
      return;
    }

    setIsSubmitting(true);
    try {
      await createRequestMutation.mutateAsync({
        editionId: props.editionId,
        data: {
          signatory_name: name,
          signatory_title: signatoryTitle().trim() || undefined,
          signatory_email: email,
          expires_in_days: expiresInDays(),
        },
      });

      toast.success("Convite de assinatura enviado com sucesso!");
      resetForm();
      props.onOpenChange(false);
      props.onSuccess?.();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Erro ao enviar convite.";
      toast.error(msg);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={(open) => {
        if (!open) resetForm();
        props.onOpenChange(open);
      }}
      title="Enviar convite de assinatura"
      description="O signatário receberá um link por e-mail para assinar digitalmente através de uma interface segura."
      footer={
        <div class="flex w-full items-center justify-end gap-2 pt-3 border-t border-border">
          <Button
            variant="outline"
            disabled={isSubmitting()}
            onClick={() => props.onOpenChange(false)}
          >
            Cancelar
          </Button>
          <Button
            disabled={isSubmitting()}
            onClick={handleSubmit}
            class="min-w-28 gap-2"
          >
            <Show when={isSubmitting()} fallback={
              <>
                <Mail class="size-4" />
                <span>Enviar convite</span>
              </>
            }>
              <Loader2 class="size-4 animate-spin" />
              <span>Enviando...</span>
            </Show>
          </Button>
        </div>
      }
    >
      <div class="space-y-4 py-2 text-left">
        <div class="space-y-1.5">
          <Label for="req-name" class="text-xs font-medium">
            Nome do signatário <span class="text-destructive">*</span>
          </Label>
          <Input
            id="req-name"
            placeholder="Ex.: Prof. Dr. Carlos Eduardo"
            value={signatoryName()}
            onInput={(e) => setSignatoryName((e.target as HTMLInputElement).value)}
            disabled={isSubmitting()}
          />
        </div>

        <div class="space-y-1.5">
          <Label for="req-email" class="text-xs font-medium">
            E-mail do signatário <span class="text-destructive">*</span>
          </Label>
          <Input
            id="req-email"
            type="email"
            placeholder="carlos.eduardo@universidade.edu.br"
            value={signatoryEmail()}
            onInput={(e) => setSignatoryEmail((e.target as HTMLInputElement).value)}
            disabled={isSubmitting()}
          />
          <p class="mt-1 text-[11px] text-muted-foreground">
            O link de assinatura individual será encaminhado para esta caixa postal.
          </p>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1.5">
            <Label for="req-title" class="text-xs font-medium">
              Cargo / Função
            </Label>
            <Input
              id="req-title"
              placeholder="Ex.: Reitor / Diretor"
              value={signatoryTitle()}
              onInput={(e) => setSignatoryTitle((e.target as HTMLInputElement).value)}
              disabled={isSubmitting()}
            />
          </div>

          <div class="space-y-1.5">
            <Label for="req-expires" class="text-xs font-medium">
              Validade do link (dias)
            </Label>
            <Input
              id="req-expires"
              type="number"
              min="1"
              max="365"
              value={String(expiresInDays())}
              onInput={(e) => {
                const val = parseInt((e.target as HTMLInputElement).value, 10);
                setExpiresInDays(isNaN(val) ? 7 : Math.max(1, Math.min(365, val)));
              }}
              disabled={isSubmitting()}
            />
          </div>
        </div>
      </div>
    </Dialog>
  );
}
