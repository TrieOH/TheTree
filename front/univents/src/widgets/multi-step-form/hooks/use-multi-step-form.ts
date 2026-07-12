import { useEffect, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import type { DefaultValues, FieldValues, UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import type { ZodType } from "zod";
import { toast } from "sonner";
import type { FieldConfig, StepConfig } from "../model/types";
import { isFieldVisible } from "../model/visibility";
import { flushImageUploadTasks } from "./use-image-upload-queue";

export interface MultiStepFormController<TInput extends FieldValues, TOutput = TInput> {
  form: UseFormReturn<TInput, unknown, TOutput>;
  steps: StepConfig<TInput>[];
  currentStep: StepConfig<TInput>;
  stepIndex: number;
  isFirstStep: boolean;
  isLastStep: boolean;
  visibleFields: FieldConfig<TInput>[];
  goNext: () => Promise<boolean>;
  goBack: () => void;
  goToStep: (index: number) => void;
  isSubmitting: boolean;
  isProcessingUploads: boolean;
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
  onSubmit: (values: TOutput) => boolean | void | Promise<boolean | void>;
  /** Optional values to restore after a successful submit (useful for create flows). */
  resetOnSuccessValues?: DefaultValues<TInput>;
  /** Called after a successful submit, after any optional reset. */
  onSubmitSuccess?: () => void | Promise<void>;
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
  resetOnSuccessValues,
  onSubmitSuccess,
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
  const [isProcessingUploads, setIsProcessingUploads] = useState(false);

  // Whenever the edited record changes (a new `values` reference lands),
  // jump the wizard back to the first step instead of leaving the user
  // stranded on step 'n' of a completely different record.
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

  const handleSubmit = async (): Promise<boolean> => {
    let didSucceed = false;

    await form.handleSubmit(async (submittedValues) => {
      setIsSubmitting(true);
      const result = await onSubmit(submittedValues);
      setIsSubmitting(false);

      if (result === false) return;

      if (resetOnSuccessValues) {
        form.reset(resetOnSuccessValues);
        setStepIndex(0);
      }

      await onSubmitSuccess?.();
      didSucceed = true;
    })().catch(() => {
      setIsSubmitting(false);
    });

    return didSucceed;
  };

  const goNext = async (): Promise<boolean> => {
    const fieldNamesToValidate = visibleFields
      .filter((field) => field.kind !== "custom")
      .map((field) => field.name);

    const isStepValid =
      fieldNamesToValidate.length > 0 ? await form.trigger(fieldNamesToValidate) : true;

    if (!isStepValid) return false;

    if (isLastStep) {
      if (!canSubmit) return false;
      try {
        setIsProcessingUploads(true);
        await flushImageUploadTasks();
      } catch (error) {
        toast.error(error instanceof Error ? error.message : "Falha ao processar imagens");
        setIsProcessingUploads(false);
        return false;
      }
      setIsProcessingUploads(false);
      return await handleSubmit();
    }

    setStepIndex((index) => Math.min(index + 1, steps.length - 1));
    return true;
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
    isProcessingUploads,
    isDirty,
    canSubmit,
  };
}
