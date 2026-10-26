import type { JSX } from "@solidjs/web";
import { createSignal, untrack } from "solid-js";
import SendIcon from "~icons/lucide/send";

import {
  MultiStepDialog,
  type MultiStepItem,
} from "@trieoh/ui-solid";
import { toast } from "@/shared/ui/toast";
import { getFieldErrors } from "@/shared/lib/effect-schema";
import { useCreateSignatureRequestMutation } from "../api/mutations";
import { signatureRequestStep1Schema } from "../model";

const Send = SendIcon as unknown as (props: { class?: string }) => JSX.Element;

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
    let nextErrors: Partial<Record<keyof CreateSignatureRequestValues, string>> = {};

    if (stepIndex === 0) {
      nextErrors = getFieldErrors<CreateSignatureRequestValues>(
        signatureRequestStep1Schema.safeParse(current),
      );
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

      toast.success("Convite enviado com sucesso!");
      resetForm();
      props.onOpenChange(false);
      props.onSuccess?.();
      return true;
    } catch {
      toast.error("Não foi possível enviar o convite de assinatura.");
      return false;
    } finally {
      setIsSubmitting(false);
    }
  };

  const steps: MultiStepItem<CreateSignatureRequestValues>[] = [
    {
      id: "destinatario",
      title: "Destinatário",
      description: "Dados do signatário e prazo para resposta",
      fields: [
        {
          name: "signatory_name",
          label: "Nome do signatário",
          required: true,
          placeholder: "Ex: Profa. Dra. Helena Silveira",
          hint: "Nome da autoridade convidada a assinar os documentos.",
        },
        {
          name: "signatory_email",
          label: "E-mail para envio",
          kind: "email",
          required: true,
          placeholder: "helena.silveira@universidade.edu",
          hint: "O link com token exclusivo e temporário será encaminhado para este endereço.",
        },
        {
          name: "signatory_title",
          label: "Cargo / Função",
          required: false,
          placeholder: "Ex: Pró-Reitora de Pós-Graduação",
          hint: "Cargo ou titulação exibida no certificado abaixo da assinatura.",
        },
        {
          name: "expires_in_days",
          label: "Validade do link (dias)",
          kind: "number",
          required: true,
          min: "1",
          max: "365",
          step: "1",
          placeholder: "7",
          hint: "Tempo em dias corridos até o link de assinatura expirar (1 a 365 dias).",
        },
      ],
    },
    {
      id: "resumo",
      title: "Resumo",
      description: "Revise os dados antes de disparar o convite",
      summary: {
        title: "Dados do convite",
        badge: () => "Pronto para enviar",
        items: [
          {
            label: "Signatário convidado",
            value: () => values().signatory_name,
          },
          {
            label: "E-mail de destino",
            value: () => values().signatory_email,
          },
          {
            label: "Cargo / Função",
            value: () => values().signatory_title || "Não informado",
          },
          {
            label: "Prazo de expiração",
            value: () => `${values().expires_in_days} dias`,
          },
        ],
        extra: () => (
          <div class="relative overflow-hidden rounded-xl border border-border/80 bg-card p-4 shadow-xs">
            <div class="flex items-start gap-3">
              <div class="rounded-lg p-2.5 bg-primary/10 text-primary shrink-0 mt-0.5">
                <Send class="size-4" />
              </div>
              <div class="space-y-1 min-w-0">
                <h4 class="text-sm font-semibold text-foreground">
                  Como funciona o convite externo?
                </h4>
                <p class="text-xs text-muted-foreground leading-relaxed">
                  Um link criptografado com token único será gerado e enviado ao
                  e-mail informado. O convidado poderá desenhar ou subir sua
                  assinatura sem necessidade de login. Assim que concluído, a
                  assinatura ficará disponível automaticamente para a edição.
                </p>
              </div>
            </div>
          </div>
        ),
      },
    },
  ];

  return (
    <MultiStepDialog<CreateSignatureRequestValues>
      open={props.open}
      onOpenChange={handleOpenChange}
      title="Enviar convite de assinatura"
      description="Gere um convite seguro com token para que uma autoridade envie sua assinatura digital."
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
