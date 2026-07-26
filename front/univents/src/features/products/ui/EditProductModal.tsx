import { useMemo } from "react";

import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import {
  type ProductI,
  type ProductPatchInputI,
  type ProductPatchOutputI,
  productPatchSchema,
} from "../model";
import { createProductPatchFormSteps } from "../model/product-form-steps";

interface EditProductModalProps {
  open: boolean;
  product: ProductI;
  onOpenChange: (open: boolean) => void;
  onUpdate: (
    values: ProductPatchOutputI,
  ) => Promise<ProductI | null | boolean> | ProductI | null | boolean;
}

export function EditProductModal({
  open,
  product,
  onOpenChange,
  onUpdate,
}: EditProductModalProps) {
  const steps = useMemo(() => createProductPatchFormSteps(), []);
  const defaultValues = useMemo<ProductPatchInputI>(
    () => ({
      vendor_code: product.vendor_code,
      requires_registration: product.requires_registration ?? false,
    }),
    [product],
  );
  const controller = useMultiStepForm({
    schema: productPatchSchema,
    steps,
    defaultValues,
    resetOnSuccessValues: defaultValues,
    onSubmit: async (values) => Boolean(await onUpdate(values)),
    onSubmitSuccess: () => onOpenChange(false),
  });

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title="Editar produto"
      controller={controller}
      submitLabel="Salvar alterações"
    />
  );
}
