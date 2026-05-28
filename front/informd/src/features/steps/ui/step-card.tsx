export interface Step {
  id: string;
  form_id: string;
  title: string;
  description?: string;
  position_hint: number /* int */;
}

interface StepCardProps {
  step: Step;
  onClick?: (step: Step) => void;
}

export function StepCard({ step, onClick }: StepCardProps) {
  return (
    <button
      type="button"
      onClick={() => onClick?.(step)}
      className="
        w-full text-left
        flex items-start gap-3 sm:gap-4
        px-4 py-3 sm:px-5 sm:py-4
        bg-card text-card-foreground
        border border-border rounded-sm
        hover:bg-muted
        transition-colors duration-150
        focus-visible:outline-none
      "
    >
      <span
        aria-label={`Step ${step.position_hint}`}
        className="
          shrink-0
          w-7 h-7 sm:w-8 sm:h-8
          flex items-center justify-center
          rounded-full
          border border-border
          text-xs sm:text-sm font-semibold
          text-muted-foreground
          mt-0.5
        "
      >
        {step.position_hint}
      </span>

      <div className="flex-1 min-w-0">
        <p className="text-sm sm:text-base font-semibold text-foreground leading-snug">
          {step.title}
        </p>

        {step.description && (
          <p className="mt-1 text-xs sm:text-sm text-muted-foreground leading-relaxed">
            {step.description}
          </p>
        )}
      </div>
    </button>
  );
}