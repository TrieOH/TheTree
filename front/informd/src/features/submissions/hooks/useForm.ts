import { useState, useCallback, useEffect } from "react";
import { getForm, getSteps, getFields, submitForm as apiSubmit } from "../model/mock";
import type { FormI } from "#/features/forms/model";
import type { StepI } from "#/features/steps/model";
import type { FieldI } from "#/features/fields/model";
import { toast } from "sonner";
import type { AnswerI, SubmitRequestI } from "../model";

interface FormState {
  form: FormI | null;
  steps: StepI[];
  fields: Record<string, FieldI[]>;
  currentStepIndex: number;
  formData: Record<string, unknown>;
  errors: Record<string, string>;
  loading: boolean;
  submitting: boolean;
  submitted: boolean;
}

export function useForm() {
  const [state, setState] = useState<FormState>({
    form: null,
    steps: [],
    fields: {},
    currentStepIndex: 0,
    formData: {},
    errors: {},
    loading: true,
    submitting: false,
    submitted: false,
  });

  const load = useCallback(async () => {
    try {
      const [formData, stepsData] = await Promise.all([getForm(), getSteps()]);
      const fieldsMap: Record<string, FieldI[]> = {};

      await Promise.all(
        stepsData.map(async (step) => {
          fieldsMap[step.id] = await getFields(step.id);
        })
      );

      setState((prev) => ({
        ...prev,
        form: formData,
        steps: stepsData,
        fields: fieldsMap,
        loading: false,
      }));
    } catch (err) {
      toast.error("Error loading form");
      setState((prev) => ({ ...prev, loading: false }));
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const setFieldValue = useCallback((fieldId: string, value: unknown) => {
    setState((prev) => ({
      ...prev,
      formData: { ...prev.formData, [fieldId]: value },
      errors: { ...prev.errors, [fieldId]: "" },
    }));
  }, []);

  const validateStep = useCallback((): boolean => {
    const { steps, fields, currentStepIndex, formData } = state;
    const step = steps[currentStepIndex];

    const stepFields = fields[step.id];
    const newErrors: Record<string, string> = {};
    let valid = true;

    for (const field of stepFields) {
      if (!field.required) continue;

      const value = formData[field.id];
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
  }, [state]);

  const goNext = useCallback(() => {
    if (!validateStep()) {
      toast.error("Please fill in all required fields.");
      return;
    }

    setState((prev) => ({
      ...prev,
      currentStepIndex: Math.min(prev.currentStepIndex + 1, prev.steps.length - 1),
    }));
  }, [validateStep]);

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
      const allFields = Object.values(state.fields).flat();
      const emailField = allFields.find((f) => f.type === "email");
      const email = emailField ? (state.formData[emailField.id] as string) : undefined;

      const answers: AnswerI[] = allFields
        .filter((f) => {
          const val = state.formData[f.id];
          return val !== undefined && val !== "" && val !== null;
        })
        .map((f) => ({
          id: "temp",
          response_id: "temp",
          field_id: f.id,
          answer: JSON.stringify(state.formData[f.id]),
          answered_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }));

      const request: SubmitRequestI = {
        email,
        answers,
      };

      await apiSubmit(request);
      setState((prev) => ({ ...prev, submitting: false, submitted: true }));
      toast.success("Form submitted successfully!");
    } catch (err) {
      setState((prev) => ({ ...prev, submitting: false }));
      toast.error("Error submitting form");
    }
  }, [state, validateStep]);

  return {
    ...state,
    setFieldValue,
    goNext,
    goBack,
    submit,
  };
}