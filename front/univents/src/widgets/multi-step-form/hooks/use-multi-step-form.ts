import { useEffect, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import type { DefaultValues, FieldValues, UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import type { ZodType } from "zod";
import type { FieldConfig, StepConfig } from "../model/types";
import { isFieldVisible } from "../model/visibility";

export interface MultiStepFormController<TInput extends FieldValues, TOutput = TInput> {
  form: UseFormReturn<TInput, unknown, TOutput>;
  steps: StepConfig<TInput>[];
  currentStep: StepConfig<TInput>;
  stepIndex: number;
  isFirstStep: boolean;
  isLastStep: boolean;
  visibleFields: FieldConfig<TInput>[];
  goNext: () => Promise<void>;
  goBack: () => void;
  goToStep: (index: number) => void;
  isSubmitting: boolean;
  /** True once any field differs from its default/edited value. */
  isDirty: boolean;
  /**
   * Whether the last-step submit action should be enabled right now.
   * Always `true` unless `requireDirtyToSubmit` was set and nothing has changed yet.
   */
  canSubmit: boolean;
}

export interface UseMultiStepFormOptions<TInput extends FieldValues, TOutput> {
  /** A zod schema whose *input* shape drives the form fields and whose
   * *output* shape (after `.transform()`s) is handed to `onSubmit`. */
  schema: ZodType<TOutput, TInput>;
  steps: StepConfig<TInput>[];
  /** Initial values for a brand-new (create) form. */
  defaultValues: DefaultValues<TInput>;
  /**
   * Edit mode: pass the record being edited here, already mapped to the
   * form's *input* shape (e.g. `null` -> `""` for text fields). Unlike
   * `defaultValues`, RHF re-syncs the form whenever this reference
   * changes — e.g. once an async fetch resolves, or when the user picks
   * a different record to edit in the same modal instance. Omit/leave
   * `undefined` for create mode.
   */
  values?: TInput;
  /**
   * When true (typically: editing an existing record), the submit
   * button on the last step stays disabled until at least one field
   * has actually changed. Has no effect on create forms.
   */
  requireDirtyToSubmit?: boolean;
  onSubmit: (values: TOutput) => void | Promise<void>;
  /**
   * When errors show up. "onChange" gives instant feedback on every
   * keystroke (revalidation is always "onChange" regardless of this
   * setting, once a field has been validated once). Defaults to
   * "onChange" for immediate feedback; use "onTouched" if you'd rather
   * wait for the first blur before nagging the user.
   */
  validationMode?: "onChange" | "onTouched" | "onBlur";
}

export function useMultiStepForm<TInput extends FieldValues, TOutput>({
  schema,
  steps,
  defaultValues,
  values,
  requireDirtyToSubmit = false,
  onSubmit,
  validationMode = "onChange",
}: UseMultiStepFormOptions<TInput, TOutput>): MultiStepFormController<TInput, TOutput> {
  const form = useForm<TInput, unknown, TOutput>({
    resolver: zodResolver(schema),
    defaultValues,
    values,
    mode: validationMode,
  });

  const [stepIndex, setStepIndex] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Whenever the edited record changes (a new `values` reference lands),
  // jump the wizard back to the first step instead of leaving the user
  // stranded on step 3 of a completely different record.
  useEffect(() => {
    setStepIndex(0);
  }, [values]);

  const watchedValues = form.watch();
  const currentStep = steps[stepIndex];
  const isFirstStep = stepIndex === 0;
  const isLastStep = stepIndex === steps.length - 1;

  const visibleFields = useMemo(
    () => currentStep.fields.filter((field) => isFieldVisible(field.visibleIf, watchedValues)),
    [currentStep, watchedValues],
  );

  const isDirty = form.formState.isDirty;
  const canSubmit = !requireDirtyToSubmit || isDirty;

  const handleSubmit = form.handleSubmit(async (submittedValues) => {
    setIsSubmitting(true);
    try {
      await onSubmit(submittedValues);
    } finally {
      setIsSubmitting(false);
    }
  });

  const goNext = async () => {
    const fieldNamesToValidate = visibleFields
      .filter((field) => field.kind === "text")
      .map((field) => field.name);

    const isStepValid =
      fieldNamesToValidate.length > 0 ? await form.trigger(fieldNamesToValidate) : true;

    if (!isStepValid) return;

    if (isLastStep) {
      if (!canSubmit) return;
      await handleSubmit();
      return;
    }

    setStepIndex((index) => Math.min(index + 1, steps.length - 1));
  };

  const goBack = () => {
    setStepIndex((index) => Math.max(index - 1, 0));
  };

  const goToStep = (index: number) => {
    setStepIndex(Math.min(Math.max(index, 0), steps.length - 1));
  };

  return {
    form,
    steps,
    currentStep,
    stepIndex,
    isFirstStep,
    isLastStep,
    visibleFields,
    goNext,
    goBack,
    goToStep,
    isSubmitting,
    isDirty,
    canSubmit,
  };
}