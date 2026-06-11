import { useState, useCallback, useMemo, useEffect } from "react";
import { toast } from "sonner";
import type { AnswerCreateI, SubmitRequestI } from "../model";
import { useSuspenseQueries } from "@tanstack/react-query";
import { allFormsAnswerableQueryOptions, submitFormFn } from "../api";
import type { FieldAnswerable } from "@trieoh/informd-models";

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
  const steps = useMemo(() => answerableForm.steps.sort((a, b) => a.step.position_hint - b.step.position_hint), [answerableForm]);

  const fields = useMemo(() => {
    const map: Record<string, FieldAnswerable[]> = {};
    steps.forEach((answerableStep) => {
      map[answerableStep.step.id] = answerableStep.fields;
    });
    return map;
  }, [steps]);

  // const allFields = useMemo(() => Object.values(fields).flat(), [fields]);
  // const selectFields = useMemo(() => allFields.filter(f => f.field.type === 'select'), [allFields]);

  // Enrich fields with their select configs
  // const enrichedFields = useMemo(() => {
  //   const enrichedMap: Record<string, FieldAnswerable[]> = { ...fields };
  //   return enrichedMap;
  // }, [fields]);

  useEffect(() => {
    if (answerableForm.steps.length === 0) return;

    const defaults: Record<string, unknown> = {};
    Object.values(fields).flat().forEach(answerableField => {
      let defaultValue: unknown;
      if (typeof answerableField.field.default_value === 'object' && answerableField.field.default_value !== null && 'value' in answerableField.field.default_value) {
        defaultValue = (answerableField.field.default_value as { value: unknown }).value;
      } else defaultValue = answerableField.field.default_value;

      if (defaultValue !== undefined && defaultValue !== null && defaultValue !== "") {
        // Standardize select field defaults to array if they aren't already
        if (answerableField.field.type === 'select' && !Array.isArray(defaultValue)) {
          defaults[answerableField.field.id] = [String(defaultValue)];
        } else defaults[answerableField.field.id] = defaultValue;
      } else if (answerableField.field.type === 'select') {
        // Ensure select fields always start as arrays
        defaults[answerableField.field.id] = [];
      } else if (answerableField.field.type === 'bool') {
        defaults[answerableField.field.id] = false;
      }
    });

    if (Object.keys(defaults).length > 0) {
      setState(prev => {
        // Only apply if formData is completely empty to avoid overwriting during session
        if (Object.keys(prev.formData).length === 0) return { ...prev, formData: defaults };
        return prev;
      });
    }
  }, [fields, steps.length]);

  const setFieldValue = useCallback((fieldId: string, value: unknown) => {
    setState((prev) => ({
      ...prev,
      formData: { ...prev.formData, [fieldId]: value },
      errors: { ...prev.errors, [fieldId]: "" },
    }));
  }, []);

  const lastStepHasFields = useMemo(() => {
    if (steps.length === 0) return false;
    const lastStepId = steps[steps.length - 1].step.id;
    return fields[lastStepId].length > 0;
  }, [steps, fields]);

  const totalStepsCount = useMemo(() => {
    return steps.length + (lastStepHasFields ? 1 : 0);
  }, [steps.length, lastStepHasFields]);

  const isReviewStep = useMemo(() => {
    if (state.currentStepIndex >= steps.length) return true;
    const currentStep = steps[state.currentStepIndex];
    return fields[currentStep.step.id].length === 0;
  }, [state.currentStepIndex, steps, fields]);

  const validateStep = useCallback((): boolean => {
    if (isReviewStep) return true;
    const step = steps[state.currentStepIndex];
    const stepFields = fields[step.step.id];
    const newErrors: Record<string, string> = {};
    let valid = true;

    for (const field of stepFields) {
      if (!field.field.required) continue;

      const value = state.formData[field.field.id];

      const isEmpty =
        value === undefined ||
        value === "" ||
        value === null ||
        (Array.isArray(value) && value.length === 0);

      if (isEmpty) {
        newErrors[field.field.id] = "This field is required";
        valid = false;
        continue;
      }

      if (field.field.type === "email" && typeof value === "string") {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(value)) {
          newErrors[field.field.id] = "Invalid email";
          valid = false;
        }
      }

      if (field.field.type === "url" && typeof value === "string") {
        try {
          new URL(value);
        } catch {
          newErrors[field.field.id] = "Invalid URL";
          valid = false;
        }
      }
    }

    setState((prev) => ({ ...prev, errors: newErrors }));
    return valid;
  }, [state.currentStepIndex, steps, fields, isReviewStep, state.formData]);

  const goNext = useCallback(() => {
    if (!validateStep()) {
      toast.error("Please fill in all required fields.");
      return;
    }

    setState((prev) => ({
      ...prev,
      currentStepIndex: Math.min(prev.currentStepIndex + 1, totalStepsCount - 1),
    }));
  }, [validateStep, totalStepsCount]);

  const goBack = useCallback(() => {
    setState((prev) => ({
      ...prev,
      currentStepIndex: Math.max(prev.currentStepIndex - 1, 0),
    }));
  }, []);

  const submit = useCallback(async () => {
    if (!validateStep()) {
      toast.error("Please fill in all required fields.");
      return;
    }

    setState((prev) => ({ ...prev, submitting: true }));

    try {
      const allFieldsList = Object.values(fields).flat();
      const emailField = allFieldsList.find((f) => f.field.type === "email");
      const email = emailField ? (state.formData[emailField.field.id] as string) : undefined;

      const answers: AnswerCreateI[] = allFieldsList
        .filter((f) => {
          const val = state.formData[f.field.id];
          return val !== undefined && val !== "" && val !== null;
        })
        .map((f) => ({
          field_id: f.field.id,
          answer: JSON.stringify(state.formData[f.field.id]),
        }));

      const request: SubmitRequestI = {
        email,
        answers,
      };

      const res = await submitFormFn(formId, request);

      if (!res.success) {
        const message = res.message || "Error submitting form";
        setState((prev) => ({ ...prev, submitting: false }));
        toast.error(message);
        return;
      }

      setState((prev) => ({ ...prev, submitting: false, submitted: true }));
      toast.success("Form submitted successfully!");
    } catch (err: unknown) {
      console.error("Form submission error:", err);
      let message = ""
      if (err instanceof Error) message = err.message || "Error submitting form";
      else message = String(err)

      setState((prev) => ({ ...prev, submitting: false }));
      toast.error(message);
    }
  }, [formId, state.formData, fields, validateStep]);

  return {
    ...state,
    steps,
    fields: fields,
    isReviewStep,
    totalStepsCount,
    loading: false, // Handled by Suspense
    setFieldValue,
    goNext,
    goBack,
    submit,
  };
}