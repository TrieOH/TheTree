import type { AnyFormApi, FormAsyncValidateOrFn, FormValidateOrFn } from "@tanstack/react-form";

export interface RuleStatus {
  message: string
  passed: boolean
}

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

export type FieldType = "text" | "select" | "icon" | "color";

export interface FieldOption {
  label: string;
  value: string;
}

export interface FieldConfig {
  name: string;
  label: string;
  placeholder?: string;
  type?: FieldType;
  autoComplete?: string;
  required?: boolean;
  getRulesStatus?: (value: unknown) => RuleStatus[]
  options?: FieldOption[];
}
