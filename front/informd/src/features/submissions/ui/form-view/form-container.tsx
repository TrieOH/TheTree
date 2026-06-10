import { useForm } from "../../hooks/useForm";
import { FieldRenderer } from "./field-renderer";
import { FormFooter } from "./form-footer";
import ProgressBar from "./progress-bar";
import { ReviewStep } from "./review-step";
import { motion, AnimatePresence } from "motion/react";
import { Check } from "lucide-react";

export function FormContainer() {
  const {
    steps,
    fields,
    currentStepIndex,
    formData,
    errors,
    loading,
    submitting,
    submitted,
    setFieldValue,
    goNext,
    goBack,
    submit,
  } = useForm();

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

  const currentStep = steps[currentStepIndex];
  const stepFields = fields[currentStep.id];
  const isLastStep = currentStepIndex === steps.length - 1;

  return (
    <div className="w-full max-w-xl overflow-hidden rounded-lg border border-border bg-card shadow-sm">
      <div className="h-2 w-full bg-primary" />
      <ProgressBar currentStep={currentStepIndex} totalSteps={steps.length} />

      <div className="px-6 pb-6 pt-6">
        <AnimatePresence mode="wait">
          <motion.div
            key={currentStep.id}
            initial={{ opacity: 0, x: 10 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -10 }}
            transition={{ duration: 0.2 }}
          >
            <h1 className="mb-1 text-xl font-bold text-foreground sm:text-lg">
              {currentStep.title}
            </h1>
            {currentStep.description && (
              <p className="mb-6 text-sm leading-relaxed text-muted-foreground">
                {currentStep.description}
              </p>
            )}

            {isLastStep || stepFields.length === 0 ? (
              <ReviewStep steps={steps} fields={fields} formData={formData} />
            ) : (
              <div className="space-y-1">
                {stepFields
                  .sort((a, b) => a.position_hint - b.position_hint)
                  .map((field) => (
                    <FieldRenderer
                      key={field.id}
                      field={field}
                      value={formData[field.id]}
                      error={errors[field.id]}
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