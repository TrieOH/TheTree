import { useCallback, useEffect, useRef, useState } from "react";
import { motion, AnimatePresence } from "motion/react";
import { Plus } from "lucide-react";
import type { StepI } from "../model";
import type { FieldI } from "#/features/fields/model";
import { StepCard } from "./step-card";
import { cn } from "#/shared/lib/utils";
import { Button } from "#/shared/ui/shadcn/button";

interface StepCarouselProps {
  steps: StepI[];
  onAddAfter: (positionHint: number) => void;
  onStepClick?: (step: StepI) => void;
  onEditStep?: (step: StepI) => void;
  onMoveStep?: (stepId: string, direction: 'left' | 'right') => void;
  onAddField?: (step: StepI) => void;
  onEditField?: (field: FieldI) => void;
  onDeleteField?: (field: FieldI) => void;
  onReorderFields?: (step: StepI, fieldIds: string[]) => void;
  fieldsByStepId?: Record<string, FieldI[]>;
  focusedStepId?: string | null;
  focusKey?: number;
}


export function StepCarousel({
  steps,
  onAddAfter,
  onStepClick,
  onEditStep,
  onMoveStep,
  onAddField,
  onEditField,
  onDeleteField,
  onReorderFields,
  fieldsByStepId,
  focusedStepId,
  focusKey,
}: StepCarouselProps) {
  const [[page, direction], setPage] = useState([0, 0]);
  const [isDraggingField, setIsDraggingField] = useState(false);
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

  // Navigate to focusedStepId whenever focusKey increments
  useEffect(() => {
    if (!focusedStepId || count <= 1) return;

    setPage(([currentPage]) => {
      const targetIndex = steps.findIndex(s => s.id === focusedStepId);
      if (targetIndex === -1) return [currentPage, 0];

      const currentWrap = ((currentPage % count) + count) % count;
      let diff = targetIndex - currentWrap;
      if (diff > count / 2) diff -= count;
      if (diff < -count / 2) diff += count;

      return [currentPage + diff, diff > 0 ? 1 : -1];
    });
  }, [focusKey]);

  const onDragEnd = (_: any, info: { offset: { x: number } }) => {
    const swipe = info.offset.x;
    const threshold = 50;
    if (swipe < -threshold) goTo(page + 1, 1);
    else if (swipe > threshold) goTo(page - 1, -1);
  };

  if (count === 0) {
    return (
      <div className="w-full min-h-64 flex items-center justify-center border-2 border-dashed border-border rounded-sm bg-muted/30">
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

  const activeFields = fieldsByStepId?.[activeStep.id] ?? [];
  const canMoveLeft = count > 1 && steps.some(s => s.position_hint === activeStep.position_hint - 1);
  const canMoveRight = count > 1 && steps.some(s => s.position_hint === activeStep.position_hint + 1);
  const maxPositionHint = count > 0 ? Math.max(...steps.map(s => s.position_hint)) : 0;

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
      <div className="relative w-full min-h-70 overflow-hidden">
        <AnimatePresence initial={false} custom={direction} mode="popLayout">
          <motion.div
            key={page}
            custom={direction}
            variants={variants}
            initial="enter"
            animate="center"
            exit="exit"
            drag={isDraggingField ? false : "x"}
            dragConstraints={{ left: 0, right: 0 }}
            dragElastic={0.2}
            onDragEnd={onDragEnd}
            transition={{
              x: { type: "spring", stiffness: 400, damping: 40 },
              opacity: { duration: 0.2 },
            }}
            className="flex items-start justify-center w-full cursor-grab active:cursor-grabbing"
          >
            {showSiblings ? (
              <>
                {/* Mobile: show only active step */}
                <div className="sm:hidden w-full px-4">
                  <StepCard
                    step={activeStep}
                    fields={activeFields}
                    active
                    onClick={() => onStepClick?.(activeStep)}
                    onEdit={onEditStep}
                    onMoveLeft={onMoveStep ? (s) => onMoveStep(s.id, 'left') : undefined}
                    onMoveRight={onMoveStep ? (s) => onMoveStep(s.id, 'right') : undefined}
                    onAddField={onAddField}
                    onEditField={onEditField}
                    onDeleteField={onDeleteField}
                    onReorderFields={onReorderFields}
                    onFieldDragChange={setIsDraggingField}
                    canMoveLeft={canMoveLeft}
                    canMoveRight={canMoveRight}
                    className="shadow-xl"
                  />
                </div>

                {/* Desktop */}
                <div className="hidden sm:flex items-center justify-center gap-1 w-full">
                  <div
                    className="shrink-0 w-[18%] opacity-40 hover:opacity-70 transition-opacity cursor-pointer scale-[0.70]"
                    onClick={() => goTo(page - 1, -1)}
                  >
                    <StepCard step={prevStep} active={false} className="pointer-events-none" />
                  </div>

                  <div className="shrink-0 w-[56%] min-w-0">
                    <StepCard
                      step={activeStep}
                      fields={activeFields}
                      active
                      onClick={() => onStepClick?.(activeStep)}
                      onEdit={onEditStep}
                      onMoveLeft={onMoveStep ? (s) => onMoveStep(s.id, 'left') : undefined}
                      onMoveRight={onMoveStep ? (s) => onMoveStep(s.id, 'right') : undefined}
                      onAddField={onAddField}
                      onEditField={onEditField}
                      onDeleteField={onDeleteField}
                      onReorderFields={onReorderFields}
                      onFieldDragChange={setIsDraggingField}
                      canMoveLeft={canMoveLeft}
                      canMoveRight={canMoveRight}
                      className="shadow-xl"
                    />
                  </div>

                  <div
                    className="shrink-0 w-[18%] opacity-40 hover:opacity-70 transition-opacity cursor-pointer scale-[0.70]"
                    onClick={() => goTo(page + 1, 1)}
                  >
                    <StepCard step={nextStep} active={false} className="pointer-events-none" />
                  </div>
                </div>
              </>
            ) : (
              <div className="w-full max-w-lg px-4">
                <StepCard
                  step={activeStep}
                  fields={activeFields}
                  active
                  onClick={() => onStepClick?.(activeStep)}
                  onEdit={onEditStep}
                  onMoveLeft={onMoveStep ? (s) => onMoveStep(s.id, 'left') : undefined}
                  onMoveRight={onMoveStep ? (s) => onMoveStep(s.id, 'right') : undefined}
                  onAddField={onAddField}
                  onEditField={onEditField}
                  onDeleteField={onDeleteField}
                  onReorderFields={onReorderFields}
                  onFieldDragChange={setIsDraggingField}
                  canMoveLeft={canMoveLeft}
                  canMoveRight={canMoveRight}
                  className="shadow-lg"
                />
              </div>
            )}
          </motion.div>
        </AnimatePresence>
      </div>

      {/* Pagination dots */}
      {showSiblings && (
        <div className="flex justify-center gap-2" aria-hidden="true">
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

      {/* Add step button at the end */}
      {count > 0 && (
        <div className="flex justify-center">
          <Button
            variant="ghost"
            onClick={() => onAddAfter(maxPositionHint + 1)}
            className={cn(
              "group flex items-center gap-2 px-4 py-2 rounded-xs",
              "border border-dashed border-border bg-transparent text-muted-foreground",
              "transition-all duration-150",
              "hover:border-primary/40 hover:text-primary hover:bg-primary/5",
              "cursor-pointer"
            )}
          >
            <Plus size={14} strokeWidth={2.5} />
            <span className="text-xs font-medium">Add Step</span>
          </Button>
        </div>
      )}
    </div>
  );
}
