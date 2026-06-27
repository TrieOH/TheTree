import OptionPicker from "@/shared/ui/form/OptionPicker";
import TextField from "@/shared/ui/form/TextField";
import { createFormHook, createFormHookContexts } from "@tanstack/react-form";

export const { fieldContext, formContext, useFieldContext } = createFormHookContexts();
export const { useAppForm } = createFormHook({
  fieldComponents: { TextField, OptionPicker },
  formComponents: {},
  fieldContext,
  formContext
});
