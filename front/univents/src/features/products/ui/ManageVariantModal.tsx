import { useMemo } from "react";
import { useMultiStepForm } from "@/widgets/multi-step-form/hooks/use-multi-step-form";
import { MultiStepFormModal } from "@/widgets/multi-step-form/ui/multi-step-form-modal";
import {
  type VariantCreateInputI,
  type VariantCreateOutputI,
  type VariantI,
  variantCreateSchema,
} from "../model";
import { createVariantFormSteps } from "../model/product-form-steps";

export interface ManageVariantModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  productId: string;
  variant?: VariantI;
  onCreate?: (
    values: VariantCreateOutputI,
  ) => Promise<VariantI | null | boolean>;
  onUpdate?: (
    id: string,
    values: VariantCreateOutputI,
  ) => Promise<VariantI | null | boolean>;
}

const emptyDefaultValues: VariantCreateInputI = {
  vendor_code: "",
  name: "",
  description: "",
  price: 0,
  stock: null,
};

function toFormValues(variant: VariantI): VariantCreateInputI {
  return {
    vendor_code: variant.vendor_code,
    name: variant.name,
    description: variant.description ?? "",
    price: variant.price,
    stock: variant.stock,
  };
}

export function ManageVariantModal({
  open,
  onOpenChange,
  variant,
  onCreate,
  onUpdate,
}: ManageVariantModalProps) {
  const isEditing = Boolean(variant);
  const editValues = useMemo(
    () => (variant ? toFormValues(variant) : undefined),
    [variant],
  );
  const steps = useMemo(() => createVariantFormSteps(), []);

  const controller = useMultiStepForm({
    schema: variantCreateSchema,
    steps,
    defaultValues: emptyDefaultValues,
    values: editValues,
    requireDirtyToSubmit: isEditing,
    onSubmit: async (values): Promise<boolean> => {
      const result = variant
        ? await onUpdate?.(variant.id, values)
        : await onCreate?.(values);

      return Boolean(result);
    },
    onSubmitSuccess: () => onOpenChange(false),
  });

  return (
    <MultiStepFormModal
      open={open}
      onOpenChange={onOpenChange}
      title={isEditing ? "Editar variação" : "Criar variação"}
      controller={controller}
      submitLabel={isEditing ? "Salvar alterações" : "Criar variação"}
    />
  );
}
