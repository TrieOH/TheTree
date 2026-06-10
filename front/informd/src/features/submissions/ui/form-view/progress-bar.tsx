interface ProgressBarProps {
  currentStep: number;
  totalSteps: number;
}

export default function ProgressBar({ currentStep, totalSteps }: ProgressBarProps) {
  const progress = ((currentStep + 1) / totalSteps) * 100;

  return (
    <div className="px-6 pt-4">
      <div className="mb-2 flex items-center justify-between">
        <span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
          Step {currentStep + 1} of {totalSteps}
        </span>
        <span className="text-[10px] font-medium text-muted-foreground/60">
          {Math.round(progress)}%
        </span>
      </div>
      <div className="h-1 w-full overflow-hidden rounded-full bg-muted/30">
        <div
          className="h-full rounded-full bg-primary transition-all duration-500 ease-out"
          style={{ width: `${progress}%` }}
        />
      </div>
    </div>
  );
}