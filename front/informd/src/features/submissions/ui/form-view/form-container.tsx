import { useForm } from "../../hooks/useForm";
import { FieldRenderer } from "./field-renderer";
import { FormFooter } from "./form-footer";
import ProgressBar from "./progress-bar";
import { ReviewStep } from "./review-step";
import { motion, AnimatePresence } from "motion/react";
import { Check, AlertCircle } from "lucide-react";
import { CatchBoundary } from "@tanstack/react-router";

interface FormContainerProps {
  formId: string;
}

export default function FormContainer({ formId }: FormContainerProps) {
  return (
    <CatchBoundary
      getResetKey={() => formId}
      errorComponent={({ error }) => {
        const err = error as any;
        const message = err?.error?.fields?.[0]?.message || err?.message || "Something went wrong";

        return (
          <div className="flex min-h-[40vh] w-full max-w-xl flex-col items-center justify-center gap-4 rounded-lg border border-destructive/20 bg-destructive/5 p-8 text-center shadow-sm">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10 text-destructive">
              <AlertCircle className="size-6" />
            </div>
            <div>
              <h2 className="text-lg font-bold text-foreground">Form not found</h2>
              <p className="mt-1 text-sm text-muted-foreground">
                {message}
              </p>
            </div>
            <button
              onClick={() => window.location.reload()}
              className="mt-2 text-sm font-semibold text-primary hover:underline"
            >
              Try again
            </button>
          </div>
        );
      }}
    >
      <FormContent formId={formId} />
    </CatchBoundary>
  );
}

function FormContent({ formId }: { formId: string }) {
  const {
    steps,
    fields,
    currentStepIndex,
    formData,
    errors,
    isReviewStep,
    totalStepsCount,
    loading,
    submitting,
    submitted,
    setFieldValue,
    goNext,
    goBack,
    submit,
  } = useForm(formId);

  if (loading) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4">
        <div className="h-10 w-10 animate-spin rounded-full border-3 border-muted border-t-primary" />
        <p className="text-sm text-muted-foreground">Loading form...</p>
      </div>
    );
  }

  if (submitted) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4">
        <div className="flex h-16 w-16 items-center justify-center rounded-full bg-green-500 text-white">
          <Check className="size-8" />
        </div>
        <h2 className="text-xl font-bold text-foreground">Form submitted!</h2>
        <p className="text-sm text-muted-foreground">Your responses have been recorded successfully.</p>
      </div>
    );
  }

  if (steps.length === 0) {
    return (
      <div className="flex min-h-[40vh] w-full max-w-xl flex-col items-center justify-center gap-4 rounded-lg border border-border bg-card p-8 text-center shadow-sm">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted text-muted-foreground">
          <AlertCircle className="size-6" />
        </div>
        <div>
          <h2 className="text-lg font-bold text-foreground">Empty Form</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            This form doesn't have any steps or fields to display yet.
          </p>
        </div>
      </div>
    );
  }

  const currentStep = steps[currentStepIndex];
  const stepFields = isReviewStep ? [] : fields[currentStep.step.id];
  const isLastStep = currentStepIndex === totalStepsCount - 1;

  return (
    <div className="w-full max-w-xl overflow-hidden rounded-lg border border-border bg-card shadow-sm">
      <div className="h-2 w-full bg-primary" />
      <ProgressBar currentStep={currentStepIndex} totalSteps={totalStepsCount} />

      <div className="px-6 pb-6 pt-6">
        <AnimatePresence mode="wait">
          <motion.div
            key={isReviewStep ? "review" : currentStep.step.id}
            initial={{ opacity: 0, x: 10 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -10 }}
            transition={{ duration: 0.2 }}
          >
            <h1 className="mb-1 text-xl font-bold text-foreground sm:text-lg">
              {isReviewStep ? "Review & Submit" : currentStep.step.title}
            </h1>
            {(isReviewStep || currentStep.step.description) && (
              <p className="mb-6 text-sm leading-relaxed text-muted-foreground">
                {isReviewStep
                  ? "Please review your information before final submission."
                  : currentStep.step.description}
              </p>
            )}

            {isReviewStep ? (
              <ReviewStep steps={steps} fields={fields} formData={formData} />
            ) : (
              <div className="space-y-1">
                {stepFields
                  .sort((a, b) => a.field.position_hint - b.field.position_hint)
                  .map((field) => (
                    <FieldRenderer
                      key={field.field.id}
                      field={field}
                      value={formData[field.field.id]}
                      error={errors[field.field.id]}
                      onChange={setFieldValue}
                    />
                  ))}
              </div>
            )}
          </motion.div>
        </AnimatePresence>
      </div>

      <FormFooter
        showBack={currentStepIndex > 0}
        isLastStep={isLastStep}
        submitting={submitting}
        onBack={goBack}
        onNext={isLastStep ? submit : goNext}
      />
    </div>
  );
}