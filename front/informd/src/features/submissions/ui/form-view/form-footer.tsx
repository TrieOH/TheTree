import { cn } from "#/shared/lib/utils";
import { ChevronLeft, ChevronRight, Check } from "lucide-react";

interface FormFooterProps {
  showBack: boolean;
  isLastStep: boolean;
  submitting: boolean;
  onBack: () => void;
  onNext: () => void;
}

export function FormFooter({
  showBack,
  isLastStep,
  submitting,
  onBack,
  onNext,
}: FormFooterProps) {
  return (
    <div className="flex items-center justify-between border-t border-border bg-card px-6 py-4">
      <div>
        {showBack && (
          <button
            type="button"
            onClick={onBack}
            className="inline-flex items-center gap-1 rounded-md px-4 py-2 text-sm font-semibold text-muted-foreground transition hover:bg-muted hover:text-foreground"
          >
            <ChevronLeft className="size-4" />
            Back
          </button>
        )}
      </div>

      <button
        type="button"
        onClick={onNext}
        disabled={submitting}
        className={cn(
          "inline-flex items-center gap-2 rounded-md bg-primary px-6 py-2 text-sm font-semibold text-primary-foreground transition",
          submitting ? "opacity-50 cursor-not-allowed" : "hover:bg-primary/90"
        )}
      >
        {submitting ? (
          <>
            <span className="h-4 w-4 animate-spin rounded-full border-2 border-primary-foreground border-t-transparent" />
            Submitting...
          </>
        ) : isLastStep ? (
          <>
            Submit
            <Check className="size-4" />
          </>
        ) : (
          <>
            Next
            <ChevronRight className="size-4" />
          </>
        )}
      </button>
    </div>
  );
}