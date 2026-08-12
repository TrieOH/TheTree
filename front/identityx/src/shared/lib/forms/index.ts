import { createFormHook, createFormHookContexts } from "@tanstack/react-form";
import DateField from "@/shared/ui/form/DateField";
import MultiOptionPicker from "@/shared/ui/form/MultiOptionPicker";
import OptionPicker from "@/shared/ui/form/OptionPicker";
import SelectField from "@/shared/ui/form/SelectField";
import TextField from "@/shared/ui/form/TextField";

export const { fieldContext, formContext, useFieldContext } =
  createFormHookContexts();
export const { useAppForm } = createFormHook({
  fieldComponents: {
    TextField,
    OptionPicker,
    MultiOptionPicker,
    DateField,
    SelectField,
  },
  formComponents: {},
  fieldContext,
  formContext,
});
