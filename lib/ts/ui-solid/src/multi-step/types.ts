import type { JSX } from "@solidjs/web";

export type MultiStepFieldKind =
  | "text"
  | "email"
  | "password"
  | "number"
  | "tel"
  | "url"
  | "textarea"
  | "custom";

export interface BaseMultiStepField<T = Record<string, unknown>> {
  name: keyof T & string;
  label: JSX.Element;
  placeholder?: string;
  hint?: JSX.Element;
  required?: boolean;
  disabled?: boolean;
  layout?: "full" | "half";
}

export interface TextMultiStepField<T = Record<string, unknown>>
  extends BaseMultiStepField<T> {
  kind?: "text" | "email" | "password" | "number" | "tel" | "url";
}

export interface TextareaMultiStepField<T = Record<string, unknown>>
  extends BaseMultiStepField<T> {
  kind: "textarea";
  rows?: number;
}

export interface CustomMultiStepField<T = Record<string, unknown>>
  extends BaseMultiStepField<T> {
  kind: "custom";
  render: (args: {
    value: unknown;
    onChange: (value: unknown) => void;
    error?: string;
    ids: { id: string; describedBy?: string };
  }) => JSX.Element;
}

export type MultiStepFieldConfig<T = Record<string, unknown>> =
  | TextMultiStepField<T>
  | TextareaMultiStepField<T>
  | CustomMultiStepField<T>;

export interface MultiStepFieldsProps<T = Record<string, unknown>> {
  fields: MultiStepFieldConfig<T>[];
  values?: T;
  errors?: Partial<Record<keyof T & string, string>>;
  onChange?: (key: keyof T & string, value: unknown) => void;
  class?: string;
}

export interface MultiStepSummaryItem {
  label: string;
  value: unknown | (() => unknown);
  fullWidth?: boolean;
  mono?: boolean;
  href?: string | (() => string | undefined);
}

export interface MultiStepSummaryConfig {
  title?: string;
  badge?: string;
  items: MultiStepSummaryItem[];
  extra?: () => JSX.Element;
}

export interface MultiStepRenderContext<T = Record<string, unknown>> {
  currentStep: number;
  step: MultiStepItem<T>;
  isFirst: boolean;
  isLast: boolean;
  goNext: () => void;
  goBack: () => void;
  goTo: (index: number) => void;
}

export interface MultiStepItem<T = Record<string, unknown>> {
  id: string;
  title: string;
  description?: string;
  fields?: MultiStepFieldConfig<T>[];
  summary?: MultiStepSummaryConfig;
  render?: (ctx: MultiStepRenderContext<T>) => JSX.Element;
}

export interface MultiStepStepperProps {
  steps: MultiStepItem<any>[];
  currentStep: number;
  onStepClick?: (index: number) => void;
  class?: string;
}
