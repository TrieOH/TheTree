import TextField from "@/shared/ui/form/TextField";
import { createFormHook, createFormHookContexts } from "@tanstack/react-form";

export const { fieldContext, formContext, useFieldContext } = createFormHookContexts();
export const { useAppForm } = createFormHook({
  fieldComponents: {
    TextField
  },
  formComponents: { },
  fieldContext,
  formContext
});
