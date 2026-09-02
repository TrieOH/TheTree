import { useMemo } from "react";
import { useCreateSignatureRequestMutation } from "@/features/signatures/api/mutations";
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
}: CreateSignatureRequestModalProps) {
  const steps = useMemo(() => createSignatureRequestFormSteps(), []);
  const createRequest = useCreateSignatureRequestMutation();
  const controller = useMultiStepForm<
    SignatureRequestCreateInputI,
    SignatureRequestCreateOutputI
  >({
    schema: signatureRequestCreateSchema,
    steps,
    defaultValues: defaults,
    onSubmit: async (values) => {
      await createRequest.mutateAsync({ editionId, data: values });
      return true;
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
