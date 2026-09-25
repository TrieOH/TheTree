import type { JSX } from "@solidjs/web";

export type MultiStepFieldKind =
  | "text"
  | "email"
  | "password"
  | "number"
  | "url"
  | "tel"
  | "date"
  | "textarea"
  | "custom";

export interface BaseMultiStepField<T = Record<string, unknown>> {
  name: keyof T & string;
  label: string;
  placeholder?: string;
  hint?: string;
  required?: boolean;
  disabled?: boolean;
  layout?: "full" | "half";
}

export interface TextMultiStepField<T = Record<string, unknown>>
  extends BaseMultiStepField<T> {
  kind?: "text" | "email" | "password" | "number" | "url" | "tel" | "date";
  autocomplete?: string;
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

export type MultiStepField<T = Record<string, unknown>> = MultiStepFieldConfig<T>;

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
  href?: string | (() => string | undefined);
  fullWidth?: boolean;
  mono?: boolean;
}

export interface MultiStepSummaryConfig {
  title?: string;
  badge?: string;
  items: MultiStepSummaryItem[];
  extra?: () => JSX.Element;
}

export interface MultiStepRenderContext<T = Record<string, unknown>> {
  readonly currentStep: number;
  readonly step: MultiStepItem<T>;
  readonly isFirst: boolean;
  readonly isLast: boolean;
  values: () => T;
  errors: () => Partial<Record<keyof T & string, string>>;
  next: () => void;
  prev: () => void;
  cancel: () => void;
  onChange: (key: keyof T & string, value: unknown) => void;
  goNext: () => void;
  goBack: () => void;
  goTo: (step: number) => void;
}

export interface MultiStepItem<T = Record<string, unknown>> {
  id: string;
  title: string;
  description?: string;
  fields?: MultiStepFieldConfig<T>[];
  summary?: MultiStepSummaryConfig;
  render?: (context: MultiStepRenderContext<T>) => JSX.Element;
}

export interface MultiStepStepperProps<T = Record<string, unknown>> {
  steps: MultiStepItem<T>[];
  currentStep: number;
  onStepClick?: (stepIndex: number) => void;
  class?: string;
}

export interface MultiStepDialogProps<T = Record<string, unknown>> {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description?: string;
  steps: MultiStepItem<T>[];
  currentStep?: number;
  onStepChange?: (step: number) => void;
  values?: T;
  errors?: Partial<Record<keyof T & string, string>>;
  onChange?: (key: keyof T & string, value: unknown) => void;
  formId?: string;
  loading?: boolean;
  submitLabel?: string;
  nextLabel?: string;
  cancelLabel?: string;
  onBeforeNext?: (currentStep: number) => boolean | Promise<boolean>;
  onFormSubmit?: () => boolean | Promise<boolean>;
  onSubmit?: () => boolean | Promise<boolean>;
  fixedHeight?: boolean;
  class?: string;
  contentClass?: string;
  /**
   * Optional custom children. If omitted, `step.fields` and `step.summary` are rendered automatically.
   */
  children?: (context: MultiStepRenderContext<T>) => JSX.Element;
  /**
   * Optional custom footer. If omitted, standard Back / Continue buttons are rendered.
   */
  footer?: (context: MultiStepRenderContext<T>) => JSX.Element;
}
