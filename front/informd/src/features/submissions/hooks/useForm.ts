import { useState, useCallback, useMemo, useEffect } from "react";
import { toast } from "sonner";
import type { AnswerCreateI } from "../model";
import { useSuspenseQueries } from "@tanstack/react-query";
import { allFormsAnswerableQueryOptions, submitFormFn } from "../api";
import type { FieldAnswerable, StepAnswerable } from "@trieoh/informd-models";

interface FormState {
  currentStepIndex: number;
  formData: Record<string, unknown>;
  errors: Record<string, string>;
  submitting: boolean;
  submitted: boolean;
}

export function useForm(formId: string) {
  const [state, setState] = useState<FormState>({
    currentStepIndex: 0,
    formData: {},
    errors: {},
    submitting: false,
    submitted: false,
  });

  const { data: answerableForm } = useSuspenseQueries({
    queries: [allFormsAnswerableQueryOptions(formId)],
    combine: (result) => ({ data: result[0].data })
  })

  const steps = useMemo(() =>
    [...answerableForm.steps].sort((a, b) => a.step.position_hint - b.step.position_hint),
    [answerableForm.steps]
  );

  const fieldsByStep = useMemo(() => {
    const map: Record<string, FieldAnswerable[]> = {};
    steps.forEach((s) => (map[s.step.id] = s.fields));
    return map;
  }, [steps]);

  // Set default values
  useEffect(() => {
    if (Object.keys(state.formData).length > 0 || steps.length === 0) return;

    const defaults: Record<string, unknown> = {};
    steps.flatMap(s => s.fields).forEach(f => {
      let val = f.field.default_value;
      if (val && typeof val === 'object' && val !== null && 'value' in val) {
        val = (val as { value: unknown }).value;
      }

      if (f.field.type === 'select') {
        defaults[f.field.id] = Array.isArray(val) ? val : (val ? [String(val)] : []);
      } else if (f.field.type === 'bool') {
        defaults[f.field.id] = val ?? false;
      } else if (val !== undefined && val !== null && val !== "") {
        defaults[f.field.id] = val;
      }
    });

    if (Object.keys(defaults).length > 0) {
      setState(prev => ({ ...prev, formData: defaults }));
    }
  }, [steps]);

  const setFieldValue = useCallback((fieldId: string, value: unknown) => {
    setState((prev) => ({
      ...prev,
      formData: { ...prev.formData, [fieldId]: value },
      errors: { ...prev.errors, [fieldId]: "" },
    }));
  }, []);

  const currentStep = steps[state.currentStepIndex] as StepAnswerable | undefined;
  const lastStepId = steps.length > 0 ? steps[steps.length - 1].step.id : null;
  const lastStepHasFields = lastStepId ? fieldsByStep[lastStepId].length > 0 : false;
  const totalStepsCount = steps.length + (lastStepHasFields ? 1 : 0);
  const isReviewStep = state.currentStepIndex >= steps.length || (currentStep && fieldsByStep[currentStep.step.id].length === 0);

  const validateStep = useCallback((): boolean => {
    if (isReviewStep || !currentStep) return true;

    const stepFields = fieldsByStep[currentStep.step.id];
    const newErrors: Record<string, string> = {};

    stepFields.forEach(f => {
      const val = state.formData[f.field.id];
      const isEmpty = val === undefined || val === "" || val === null || (Array.isArray(val) && val.length === 0);

      if (f.field.required && isEmpty) {
        newErrors[f.field.id] = "This field is required";
      } else if (!isEmpty && typeof val === "string") {
        if (f.field.type === "email" && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val)) {
          newErrors[f.field.id] = "Invalid email";
        } else if (f.field.type === "url") {
          try { new URL(val) } catch { newErrors[f.field.id] = "Invalid URL" }
        }
      }
    });

    setState((prev) => ({ ...prev, errors: newErrors }));
    return Object.keys(newErrors).length === 0;
  }, [currentStep, fieldsByStep, state.formData, isReviewStep]);

  const goNext = useCallback(() => {
    if (!validateStep()) {
      toast.error("Please fill in all required fields.");
      return;
    }
    setState((prev) => ({ ...prev, currentStepIndex: Math.min(prev.currentStepIndex + 1, totalStepsCount - 1) }));
  }, [validateStep, totalStepsCount]);

  const goBack = useCallback(() => {
    setState((prev) => ({ ...prev, currentStepIndex: Math.max(prev.currentStepIndex - 1, 0) }));
  }, []);

  const submit = useCallback(async () => {
    if (!validateStep()) {
      toast.error("Please fill in all required fields.");
      return;
    }
    setState((prev) => ({ ...prev, submitting: true }));

    try {
      const allFields = steps.flatMap(s => s.fields);
      const emailField = allFields.find((f) => f.field.type === "email");
      const email = emailField ? (state.formData[emailField.field.id] as string) : undefined;

      const answers: AnswerCreateI[] = allFields
        .filter((f) => {
          const val = state.formData[f.field.id];
          return val !== undefined && val !== "" && val !== null;
        })
        .map((f) => ({
          field_id: f.field.id,
          answer: JSON.stringify(state.formData[f.field.id]),
        }));

      const res = await submitFormFn(formId, { email, answers });
      if (!res.success) throw new Error(res.message || "Error submitting form");

      setState((prev) => ({ ...prev, submitting: false, submitted: true }));
      toast.success("Form submitted successfully!");
    } catch (err: unknown) {
      setState((prev) => ({ ...prev, submitting: false }));
      const message = err instanceof Error ? err.message : String(err);
      toast.error(message);
    }
  }, [formId, state.formData, steps, validateStep]);

  return {
    ...state,
    steps,
    fields: fieldsByStep,
    isReviewStep,
    totalStepsCount,
    loading: false,
    setFieldValue,
    goNext,
    goBack,
    submit,
  };
}
