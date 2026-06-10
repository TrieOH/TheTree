import { useState, useCallback, useMemo, useEffect } from "react";
import type { FieldI } from "#/features/fields/model";
import { toast } from "sonner";
import type { AnswerCreateI, SubmitRequestI } from "../model";
import { useSuspenseQueries } from "@tanstack/react-query";
import { allFormsStepsQueryOptions } from "#/features/steps/api";
import { allStepsFieldsQueryOptions, allSelectConfigsQueryOptions } from "#/features/fields/api";
import { submitFormFn } from "../api";

interface FormState {
  currentStepIndex: number;
  formData: Record<string, unknown>;
  errors: Record<string, string>;
  submitting: boolean;
  submitted: boolean;
}

export function useForm(formId: string, namespaceId?: string) {
  const [state, setState] = useState<FormState>({
    currentStepIndex: 0,
    formData: {},
    errors: {},
    submitting: false,
    submitted: false,
  });

  // 1. Fetch Steps
  const { data: steps } = useSuspenseQueries({
    queries: [allFormsStepsQueryOptions(formId, namespaceId)],
    combine: (results) => ({
      data: results[0].data.sort((a, b) => a.position_hint - b.position_hint)
    })
  });

  // 2. Fetch Fields for each step
  const fieldsQueries = useSuspenseQueries({
    queries: steps.map((step) => allStepsFieldsQueryOptions(formId, step.id, namespaceId)),
  });

  const fields = useMemo(() => {
    const map: Record<string, FieldI[]> = {};
    steps.forEach((step, index) => {
      map[step.id] = fieldsQueries[index].data;
    });
    return map;
  }, [steps, fieldsQueries]);

  // 3. Fetch Select Configs for all select fields
  const allFields = useMemo(() => Object.values(fields).flat(), [fields]);
  const selectFields = useMemo(() => allFields.filter(f => f.type === 'select'), [allFields]);

  const selectConfigsQueries = useSuspenseQueries({
    queries: selectFields.map(f => allSelectConfigsQueryOptions(f.id, formId, f.step_id, namespaceId))
  });

  // Enrich fields with their select configs
  const enrichedFields = useMemo(() => {
    const enrichedMap: Record<string, FieldI[]> = { ...fields };
    selectFields.forEach((field, index) => {
      const config = selectConfigsQueries[index].data;
      const stepFields = enrichedMap[field.step_id];
      const fieldIndex = stepFields.findIndex(f => f.id === field.id);
      if (fieldIndex !== -1) stepFields[fieldIndex] = { ...stepFields[fieldIndex], config };
    });
    return enrichedMap;
  }, [fields, selectFields, selectConfigsQueries]);

  // Initial Form Data with Defaults - Use useEffect to run only once after fields are loaded
  useEffect(() => {
    if (steps.length === 0) return;

    const defaults: Record<string, unknown> = {};
    Object.values(enrichedFields).flat().forEach(field => {
      let defaultValue: unknown;
      if (typeof field.default_value === 'object' && field.default_value !== null && 'value' in field.default_value) {
        defaultValue = (field.default_value as { value: unknown }).value;
      } else defaultValue = field.default_value;

      if (defaultValue !== undefined && defaultValue !== null && defaultValue !== "") {
        // Standardize select field defaults to array if they aren't already
        if (field.type === 'select' && !Array.isArray(defaultValue)) {
          defaults[field.id] = [String(defaultValue)];
        } else defaults[field.id] = defaultValue;
      } else if (field.type === 'select') {
        // Ensure select fields always start as arrays
        defaults[field.id] = [];
      }
    });

    if (Object.keys(defaults).length > 0) {
      setState(prev => {
        // Only apply if formData is completely empty to avoid overwriting during session
        if (Object.keys(prev.formData).length === 0) return { ...prev, formData: defaults };
        return prev;
      });
    }
  }, [enrichedFields, steps.length]);

  const setFieldValue = useCallback((fieldId: string, value: unknown) => {
    setState((prev) => ({
      ...prev,
      formData: { ...prev.formData, [fieldId]: value },
      errors: { ...prev.errors, [fieldId]: "" },
    }));
  }, []);

  const lastStepHasFields = useMemo(() => {
    if (steps.length === 0) return false;
    const lastStepId = steps[steps.length - 1].id;
    return enrichedFields[lastStepId].length > 0;
  }, [steps, enrichedFields]);

  const totalStepsCount = useMemo(() => {
    return steps.length + (lastStepHasFields ? 1 : 0);
  }, [steps.length, lastStepHasFields]);

  const isReviewStep = useMemo(() => {
    if (state.currentStepIndex >= steps.length) return true;
    const currentStep = steps[state.currentStepIndex];
    return enrichedFields[currentStep.id].length === 0;
  }, [state.currentStepIndex, steps, enrichedFields]);

  const validateStep = useCallback((): boolean => {
    if (isReviewStep) return true;
    const step = steps[state.currentStepIndex];
    const stepFields = enrichedFields[step.id];
    const newErrors: Record<string, string> = {};
    let valid = true;

    for (const field of stepFields) {
      if (!field.required) continue;

      const value = state.formData[field.id];

      const isEmpty =
        value === undefined ||
        value === "" ||
        value === null ||
        (Array.isArray(value) && value.length === 0);

      if (isEmpty) {
        newErrors[field.id] = "This field is required";
        valid = false;
        continue;
      }

      if (field.type === "email" && typeof value === "string") {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(value)) {
          newErrors[field.id] = "Invalid email";
          valid = false;
        }
      }

      if (field.type === "url" && typeof value === "string") {
        try {
          new URL(value);
        } catch {
          newErrors[field.id] = "Invalid URL";
          valid = false;
        }
      }
    }

    setState((prev) => ({ ...prev, errors: newErrors }));
    return valid;
  }, [state.currentStepIndex, steps, enrichedFields, isReviewStep, state.formData]);

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
      const allFieldsList = Object.values(enrichedFields).flat();
      const emailField = allFieldsList.find((f) => f.type === "email");
      const email = emailField ? (state.formData[emailField.id] as string) : undefined;

      const answers: AnswerCreateI[] = allFieldsList
        .filter((f) => {
          const val = state.formData[f.id];
          return val !== undefined && val !== "" && val !== null;
        })
        .map((f) => ({
          field_id: f.id,
          answer: JSON.stringify(state.formData[f.id]),
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
  }, [formId, state.formData, enrichedFields, validateStep]);

  return {
    ...state,
    steps,
    fields: enrichedFields,
    isReviewStep,
    totalStepsCount,
    loading: false, // Handled by Suspense
    setFieldValue,
    goNext,
    goBack,
    submit,
  };
}