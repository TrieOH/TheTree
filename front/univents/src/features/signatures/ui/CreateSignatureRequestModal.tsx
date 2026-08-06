import { useMemo } from "react";
import { toast } from "sonner";
import { createSignatureRequestFn } from "@/features/signatures/api";
import { findActorIdByEmailServerFn } from "@/features/signatures/api/actor-lookup";
import {
  type SignatureRequestCreateInputI,
  type SignatureRequestCreateOutputI,
  signatureRequestCreateSchema,
} from "@/features/signatures/model";
import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import { createSignatureRequestFormSteps } from "../model/signature-request-form-steps";

interface CreateSignatureRequestModalProps {
  editionId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onCreated: () => Promise<void>;
}

const defaults: SignatureRequestCreateInputI = {
  signatory_name: "",
  signatory_title: "",
  signatory_email: "",
  signatory_user_id: "",
  expires_in_days: 7,
};

export function CreateSignatureRequestModal({
  editionId,
  open,
  onOpenChange,
  onCreated,
}: CreateSignatureRequestModalProps) {
  const steps = useMemo(() => createSignatureRequestFormSteps(), []);
  const controller = useMultiStepForm<
    SignatureRequestCreateInputI,
    SignatureRequestCreateOutputI
  >({
    schema: signatureRequestCreateSchema,
    steps,
    defaultValues: defaults,
    onSubmit: async (values) => {
      try {
        await createSignatureRequestFn(editionId, {
          ...values,
          signatory_title: values.signatory_title || undefined,
          signatory_user_id: await findActorIdByEmailServerFn({
            data: { email: values.signatory_email },
          }),
        });
        await onCreated();
        return true;
      } catch (error) {
        toast.error(
          error instanceof Error ? error.message : "Erro ao criar convite",
        );
        return false;
      }
    },
    onSubmitSuccess: () => onOpenChange(false),
    resetOnSuccessValues: defaults,
  });

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Novo convite de assinatura"
      controller={controller}
      submitLabel="Enviar convite"
    />
  );
}
