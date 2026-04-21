import TextField from "@/shared/ui/form/TextField";
import IconSelector from "@/shared/ui/form/IconSelector";
import ColorSelector from "@/shared/ui/form/ColorSelector";
import { createFormHook, createFormHookContexts } from "@tanstack/react-form";

export const { fieldContext, formContext, useFieldContext } = createFormHookContexts();
export const { useAppForm } = createFormHook({
  fieldComponents: {
    TextField,
    IconSelector,
    ColorSelector
  },
  formComponents: { },
  fieldContext,
  formContext
});
