import DateField from "@/shared/ui/form/DateField";
import OptionPicker from "@/shared/ui/form/OptionPicker";
import MultiOptionPicker from "@/shared/ui/form/MultiOptionPicker";
import TextField from "@/shared/ui/form/TextField";
import { createFormHook, createFormHookContexts } from "@tanstack/react-form";

export const { fieldContext, formContext, useFieldContext } = createFormHookContexts();
export const { useAppForm } = createFormHook({
  fieldComponents: { TextField, OptionPicker, MultiOptionPicker, DateField },
  formComponents: {},
  fieldContext,
  formContext
});
