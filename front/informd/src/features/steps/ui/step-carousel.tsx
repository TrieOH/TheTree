import { useCallback, useEffect, useRef, useState } from "react";
import { motion, AnimatePresence } from "motion/react";
import { Plus } from "lucide-react";
import type { StepI } from "../model";
import { StepCard } from "./step-card";
import { cn } from "#/shared/lib/utils";
import { Button } from "#/shared/ui/shadcn/button";

interface AddButtonProps {
  positionHint: number;
  label: string;
  onAdd: (positionHint: number) => void;
}

function AddButton({ positionHint, label, onAdd }: AddButtonProps) {
  return (
    <Button
      variant="ghost"
      aria-label={label}
      onClick={() => onAdd(positionHint)}
      className={cn(
        "shrink-0 w-8 h-8 rounded-full",
        "flex items-center justify-center",
        "border border-border bg-background text-muted-foreground",
        "transition-all duration-150",
        "hover:border-primary/50 hover:text-primary hover:bg-primary/5",
        "active:scale-90",
        "outline-none cursor-pointer"
      )}
    >
      <Plus size={14} strokeWidth={2.5} />
    </Button>
  );
}

interface StepCarouselProps {
  steps: StepI[];
  onAddBefore: (positionHint: number) => void;
  onAddAfter: (positionHint: number) => void;
  onStepClick?: (step: StepI) => void;
}


export function StepCarousel({
  steps,
  onAddBefore,
  onAddAfter,
  onStepClick,
}: StepCarouselProps) {
  const [[page, direction], setPage] = useState([0, 0]);
  const count = steps.length;
  const lastMoveTime = useRef(0);

  // Wrap the infinite page to the actual steps index
  const activeIndex = count > 0 ? ((page % count) + count) % count : 0;

  useEffect(() => {
    // If steps are deleted and index is out of bounds, reset page
    if (activeIndex >= count && count > 0) setPage([0, 0]);
  }, [count, activeIndex]);

  const goTo = useCallback(
    (targetPage: number, forcedDirection?: number) => {
      const now = Date.now();
      if (now - lastMoveTime.current < 200) return;
      lastMoveTime.current = now;

      const diff = forcedDirection ?? (targetPage > page ? 1 : -1);
      setPage([targetPage, diff]);
    },
    [page]
  );

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (count <= 1) return;
      if (e.key === "ArrowLeft") {
        goTo(page - 1, -1);
      } else if (e.key === "ArrowRight") {
        goTo(page + 1, 1);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [page, count, goTo]);

  const onDragEnd = (_: any, info: { offset: { x: number } }) => {
    const swipe = info.offset.x;
    const threshold = 50;
    if (swipe < -threshold) goTo(page + 1, 1);
    else if (swipe > threshold) goTo(page - 1, -1);
  };

  if (count === 0) {
    return (
      <div className="w-full h-70 flex items-center justify-center border-2 border-dashed border-border rounded-sm bg-muted/30">
        <Button
          onClick={() => onAddAfter(1)}
          variant="ghost"
          className="group flex w-full h-full flex-col items-center justify-center gap-2 text-sm font-medium text-muted-foreground hover:text-primary transition-all"
        >
          <div className="w-10 h-10 rounded-full flex items-center justify-center border border-border bg-background group-hover:border-primary/30 group-hover:bg-primary/5 transition-all">
            <Plus size={20} />
          </div>
          Add your first step
        </Button>
      </div>
    );
  }

  const prevStep = steps[((activeIndex - 1) % count + count) % count];
  const nextStep = steps[((activeIndex + 1) % count + count) % count];
  const activeStep = steps[activeIndex];
  const showSiblings = count > 1;

  const variants = {
    enter: (dir: number) => ({
      x: dir > 0 ? "60%" : "-60%",
      opacity: 0,
      scale: 0.9,
    }),
    center: {
      x: 0,
      opacity: 1,
      scale: 1,
      zIndex: 1,
    },
    exit: (dir: number) => ({
      x: dir > 0 ? "-60%" : "60%",
      opacity: 0,
      scale: 0.9,
      zIndex: 0,
    }),
  };

  return (
    <div className="w-full flex flex-col gap-4 select-none">
      <div className="relative w-full h-70 overflow-hidden flex items-center justify-center">
        <AnimatePresence initial={false} custom={direction} mode="popLayout">
          <motion.div
            key={page}
            custom={direction}
            variants={variants}
            initial="enter"
            animate="center"
            exit="exit"
            drag="x"
            dragConstraints={{ left: 0, right: 0 }}
            dragElastic={0.2}
            onDragEnd={onDragEnd}
            transition={{
              x: { type: "spring", stiffness: 400, damping: 40 },
              opacity: { duration: 0.2 },
            }}
            className="absolute inset-0 flex items-center justify-center w-full h-full cursor-grab active:cursor-grabbing"
          >
            {showSiblings ? (
              <div className="flex items-center justify-center gap-4 w-[140%] sm:w-full">
                {/* Left Peek */}
                <div
                  className="shrink-0 w-[20%] opacity-60 hover:opacity-100 transition-opacity cursor-pointer scale-90"
                  onClick={() => goTo(page - 1, -1)}
                >
                  <StepCard step={prevStep} active={false} className="pointer-events-none" />
                </div>

                {/* Main Active Card */}
                <div className="shrink-0 w-[90%] sm:w-[50%] flex items-center gap-3">
                  <AddButton
                    positionHint={Math.max(1, activeStep.position_hint - 1)}
                    label={`Add step before "${activeStep.title}"`}
                    onAdd={onAddBefore}
                  />
                  <div className="flex-1 min-w-0">
                    <StepCard
                      step={activeStep}
                      active
                      onClick={() => onStepClick?.(activeStep)}
                      className="shadow-xl"
                    />
                  </div>
                  <AddButton
                    positionHint={activeStep.position_hint + 1}
                    label={`Add step after "${activeStep.title}"`}
                    onAdd={onAddAfter}
                  />
                </div>

                {/* Right Peek */}
                <div
                  className="shrink-0 w-[20%] opacity-60 hover:opacity-100 transition-opacity cursor-pointer scale-90"
                  onClick={() => goTo(page + 1, 1)}
                >
                  <StepCard step={nextStep} active={false} className="pointer-events-none" />
                </div>
              </div>
            ) : (
              <div className="w-full max-w-lg flex items-center gap-3 px-4">
                <AddButton
                  positionHint={Math.max(1, activeStep.position_hint - 1)}
                  label={`Add step before "${activeStep.title}"`}
                  onAdd={onAddBefore}
                />
                <div className="flex-1 min-w-0">
                  <StepCard
                    step={activeStep}
                    active
                    onClick={() => onStepClick?.(activeStep)}
                    className="shadow-lg"
                  />
                </div>
                <AddButton
                  positionHint={activeStep.position_hint + 1}
                  label={`Add step after "${activeStep.title}"`}
                  onAdd={onAddAfter}
                />
              </div>
            )}
          </motion.div>
        </AnimatePresence>
      </div>

      {showSiblings && (
        <div className="flex justify-center gap-2 py-2" aria-hidden="true">
          {steps.map((_, i) => (
            <button
              key={i}
              type="button"
              onClick={() => {
                // Calculate nearest page for this index
                const currentWrap = ((page % count) + count) % count;
                let diff = i - currentWrap;
                if (diff > count / 2) diff -= count;
                if (diff < -count / 2) diff += count;
                goTo(page + diff);
              }}
              className={cn(
                "h-1 rounded-full transition-all duration-300",
                i === activeIndex
                  ? "bg-primary w-8"
                  : "bg-muted-foreground/20 hover:bg-muted-foreground/40 w-2"
              )}
            />
          ))}
        </div>
      )}
    </div>
  );
}
