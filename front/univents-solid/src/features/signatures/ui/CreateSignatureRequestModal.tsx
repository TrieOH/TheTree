import type { JSX } from "@solidjs/web";
import { createSignal, untrack } from "solid-js";
import CalendarIcon from "~icons/lucide/calendar";
import MailIcon from "~icons/lucide/mail";
import SendIcon from "~icons/lucide/send";
import ShieldCheckIcon from "~icons/lucide/shield-check";

import {
  MultiStepDialog,
  type MultiStepItem,
} from "@trieoh/ui-solid";
import { toast } from "@/shared/ui/toast";
import { useCreateSignatureRequestMutation } from "../api/mutations";

const Calendar = CalendarIcon as unknown as (props: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;
const Send = SendIcon as unknown as (props: { class?: string }) => JSX.Element;
const ShieldCheck = ShieldCheckIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface CreateSignatureRequestModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  editionId: string;
  onSuccess?: () => void;
}

export interface CreateSignatureRequestValues {
  signatory_name: string;
  signatory_email: string;
  signatory_title: string;
  expires_in_days: string | number;
}

const emptyValues: CreateSignatureRequestValues = {
  signatory_name: "",
  signatory_email: "",
  signatory_title: "",
  expires_in_days: 7,
};

export function CreateSignatureRequestModal(
  props: CreateSignatureRequestModalProps,
): JSX.Element {
  const [currentStep, setCurrentStep] = createSignal(0);
  let currentValues: CreateSignatureRequestValues = { ...emptyValues };
  const [values, setValues] = createSignal<CreateSignatureRequestValues>(
    untrack(() => currentValues),
  );
  const [errors, setErrors] = createSignal<
    Partial<Record<keyof CreateSignatureRequestValues, string>>
  >({});
  const [isSubmitting, setIsSubmitting] = createSignal(false);

  const createRequestMutation = useCreateSignatureRequestMutation();

  const resetForm = () => {
    currentValues = { ...emptyValues };
    setValues({ ...emptyValues });
    setErrors({});
    setCurrentStep(0);
  };

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      resetForm();
    }
    props.onOpenChange(open);
  };

  const handleChange = (
    key: keyof CreateSignatureRequestValues & string,
    value: unknown,
  ) => {
    currentValues = { ...currentValues, [key]: value };
    setValues(currentValues);
    setErrors((prev) => {
      if (!prev[key as keyof CreateSignatureRequestValues]) return prev;
      const nextErrors = { ...prev };
      delete nextErrors[key as keyof CreateSignatureRequestValues];
      return nextErrors;
    });
  };

  const validateStep = (stepIndex: number): boolean => {
    const current = currentValues;
    const nextErrors: Partial<Record<keyof CreateSignatureRequestValues, string>> =
      {};

    if (stepIndex === 0) {
      const name = current.signatory_name.trim();
      if (name.length < 2) {
        nextErrors.signatory_name =
          "O nome do signatário deve ter pelo menos 2 caracteres.";
      }

      const email = current.signatory_email.trim();
      if (!email || !email.includes("@") || !email.includes(".")) {
        nextErrors.signatory_email =
          "Informe um e-mail válido para envio do convite.";
      }

      const days = Number(current.expires_in_days);
      if (isNaN(days) || !Number.isInteger(days) || days < 1 || days > 365) {
        nextErrors.expires_in_days =
          "A validade deve ser um número inteiro entre 1 e 365 dias.";
      }
    }

    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors);
      return false;
    }

    setErrors({});
    return true;
  };

  const handleSubmit = async (): Promise<boolean> => {
    const current = currentValues;
    const name = current.signatory_name.trim();
    const email = current.signatory_email.trim();
    const days = Number(current.expires_in_days) || 7;

    setIsSubmitting(true);
    try {
      await createRequestMutation.mutateAsync({
        editionId: props.editionId,
        data: {
          signatory_name: name,
          signatory_title: current.signatory_title.trim() || undefined,
          signatory_email: email,
          expires_in_days: days,
        },
      });

      toast.success("Convite de assinatura enviado com sucesso!");
      resetForm();
      props.onOpenChange(false);
      props.onSuccess?.();
      return true;
    } catch (err: unknown) {
      const msg =
        err instanceof Error ? err.message : "Erro ao enviar convite.";
      toast.error(msg);
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  const steps: MultiStepItem<CreateSignatureRequestValues>[] = [
    {
      id: "destinatario",
      title: "Destinatário",
      description: "Identificação da autoridade que receberá o convite e prazo",
      fields: [
        {
          name: "signatory_name",
          label: "Nome do signatário",
          required: true,
          placeholder: "Ex.: Prof. Dr. Carlos Eduardo",
          hint: "Nome completo da autoridade signatária.",
        },
        {
          name: "signatory_email",
          label: "E-mail para envio",
          kind: "email",
          required: true,
          placeholder: "carlos.eduardo@universidade.edu.br",
          hint: "O link com token seguro de assinatura será enviado para esta caixa postal.",
        },
        {
          name: "signatory_title",
          label: "Cargo / Função",
          layout: "half",
          placeholder: "Ex.: Reitor, Coordenador Científico",
          hint: "Cargo oficial no evento (opcional).",
        },
        {
          name: "expires_in_days",
          label: "Validade do link (dias)",
          kind: "number",
          step: "1",
          min: "1",
          max: "365",
          layout: "half",
          required: true,
          placeholder: "7",
          hint: "Prazo para expiração do link.",
        },
      ],
    },
    {
      id: "resumo",
      title: "Resumo",
      description: "Revise os dados antes de disparar o convite por e-mail",
      summary: {
        title: "Dados do convite",
        badge: () => "Pronto para enviar",
        items: [
          {
            label: "Nome do signatário",
            value: () => values().signatory_name || "Não informado",
          },
          {
            label: "E-mail de destino",
            value: () => values().signatory_email || "Não informado",
          },
          {
            label: "Cargo / Função",
            value: () => values().signatory_title || "Não informado",
          },
          {
            label: "Validade do link",
            value: () => `${values().expires_in_days || 7} dias`,
          },
        ],
        extra: () => (
          <div class="relative overflow-hidden rounded-xl border border-border/80 bg-card p-5 shadow-xs">
            <div class="flex items-center justify-between pb-3 border-b border-border/60">
              <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider flex items-center gap-1.5">
                <Mail class="size-3.5" />
                Prévia da Notificação
              </span>
              <span class="inline-flex items-center gap-1 rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-semibold text-primary">
                <ShieldCheck class="size-3.5" />
                Link Autenticado
              </span>
            </div>

            <div class="mt-4 space-y-3 rounded-lg border border-border/60 bg-muted/20 p-4 text-xs">
              <div class="flex items-center justify-between border-b border-border/40 pb-2 text-muted-foreground">
                <span>Para: <strong class="text-foreground font-medium">{values().signatory_email || "email@instituicao.org"}</strong></span>
                <span class="inline-flex items-center gap-1 font-mono text-[11px]">
                  <Calendar class="size-3" />
                  Expira em {values().expires_in_days || 7}d
                </span>
              </div>

              <div>
                <p class="font-semibold text-foreground text-sm">
                  Convite para Assinatura de Certificados Digitais
                </p>
                <p class="mt-1 text-muted-foreground leading-relaxed">
                  Olá, <strong>{values().signatory_name || "Signatário(a)"}</strong>. Você foi convidado(a) para assinar digitalmente os certificados oficiais desta edição.
                </p>
              </div>

              <div class="pt-2">
                <div class="inline-flex items-center gap-1.5 rounded-lg bg-primary px-3 py-1.5 text-xs font-semibold text-primary-foreground shadow-xs">
                  <Send class="size-3.5" />
                  <span>Acessar e Assinar Documentos</span>
                </div>
              </div>
            </div>
          </div>
        ),
      },
    },
  ];

  return (
    <MultiStepDialog
      open={props.open}
      onOpenChange={handleOpenChange}
      title="Enviar convite de assinatura"
      description="O signatário receberá um link seguro por e-mail para assinar digitalmente."
      steps={steps}
      currentStep={currentStep()}
      onStepChange={setCurrentStep}
      values={values()}
      errors={errors()}
      onChange={handleChange}
      loading={isSubmitting()}
      submitLabel={isSubmitting() ? "Enviando..." : "Enviar convite"}
      onBeforeNext={validateStep}
      onSubmit={handleSubmit}
    />
  );
}

export { CreateSignatureRequestModal as ManageSignatureRequestDialog };
