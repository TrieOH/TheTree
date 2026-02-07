import type { AnyFormApi, FormAsyncValidateOrFn, FormValidateOrFn } from "@tanstack/react-form";

export interface CrudFormConfig<TFormData> {
  defaultValues: TFormData;
  validators?: {
    onChange?: FormValidateOrFn<TFormData>;
    onBlur?: FormValidateOrFn<TFormData>;
    onSubmit?: FormValidateOrFn<TFormData>;
    onChangeAsync?: FormAsyncValidateOrFn<TFormData>;
    onBlurAsync?: FormAsyncValidateOrFn<TFormData>;
    onSubmitAsync?: FormAsyncValidateOrFn<TFormData>;
  };
  onSubmit: (props: { value: TFormData, formApi: AnyFormApi }) => Promise<void> | void;
}

export interface FieldConfig {
  name: string;
  label: string;
  placeholder: string;
  autoComplete?: string;
}