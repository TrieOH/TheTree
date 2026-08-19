import { Check, ListChecks, X } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";
import { useState } from "react";

/** A compact, dismissible checklist with an optional action for each step. */
export interface ChecklistItem {
  id: string | number;
  title: string;
  description?: string;
  completed: boolean;
  action?: { label: string; onClick: () => void };
}

interface StepChecklistProps {
  title?: string;
  items: ChecklistItem[];
  onDismiss?: () => void;
  defaultOpen?: boolean;
  className?: string;
}

export function StepChecklist({
  title,
  items,
  onDismiss,
  defaultOpen = false,
  className = "",
}: StepChecklistProps) {
  const [open, setOpen] = useState(defaultOpen);
  const pending = items.filter((item) => !item.completed).length;

  return (
    <div className={`relative inline-block ${className}`}>
      <AnimatePresence initial={false}>
        {!open && (
          <motion.button
            key="trigger"
            type="button"
            onClick={() => setOpen(true)}
            aria-label="Open checklist"
            initial={{ opacity: 0, scale: 0.85 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.85 }}
            transition={{ duration: 0.15 }}
            className="fixed right-4 top-24 z-50 flex h-11 w-11 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-lg shadow-primary/20 hover:bg-primary/90 sm:right-6"
          >
            <ListChecks size={18} />
            {pending > 0 && (
              <span className="absolute -right-0.5 -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-destructive px-1 text-[9px] font-bold text-destructive-foreground ring-2 ring-background">
                {pending}
              </span>
            )}
          </motion.button>
        )}
      </AnimatePresence>

      <AnimatePresence>
        {open && (
          <motion.div
            key="panel"
            initial={{ opacity: 0, scale: 0.92, x: 12 }}
            animate={{ opacity: 1, scale: 1, x: 0 }}
            exit={{ opacity: 0, scale: 0.92, x: 12 }}
            transition={{ duration: 0.18, ease: [0.16, 1, 0.3, 1] }}
            style={{ transformOrigin: "top right" }}
            className="fixed right-4 top-24 z-50 w-80 max-w-[calc(100vw-2rem)] rounded-2xl border border-primary/15 bg-card/95 p-5 text-card-foreground shadow-xl shadow-primary/10 backdrop-blur-md sm:right-6"
          >
            <div className="mb-4 flex items-start justify-between gap-3">
              {title ? (
                <h3 className="text-sm font-semibold text-foreground">
                  {title}
                </h3>
              ) : null}
              <button
                type="button"
                onClick={() => setOpen(false)}
                aria-label="Close checklist"
                className="-mr-1 -mt-1 shrink-0 rounded-md p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
              >
                <X size={16} />
              </button>
            </div>

            <ul className="flex flex-col">
              {items.map((item, index) => {
                const isLast = index === items.length - 1;
                const nextCompleted = !isLast && items[index + 1].completed;
                const lineFilled = item.completed && nextCompleted;

                return (
                  <li key={item.id} className="flex gap-3">
                    <div className="flex flex-col items-center">
                      <span
                        className={`flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-[11px] font-semibold transition-colors ${item.completed ? "bg-primary text-primary-foreground" : "border-2 border-border bg-background text-muted-foreground"}`}
                      >
                        {item.completed ? (
                          <Check size={13} strokeWidth={3} />
                        ) : (
                          index + 1
                        )}
                      </span>
                      {!isLast && (
                        <span
                          aria-hidden="true"
                          className={`my-1 w-px flex-1 transition-colors ${lineFilled ? "bg-primary/60" : "bg-border"}`}
                        />
                      )}
                    </div>

                    <div
                      className={`flex min-w-0 flex-1 flex-col gap-1 ${isLast ? "" : "pb-5"}`}
                    >
                      <div className="flex min-h-6 items-center justify-between gap-3">
                        <p
                          className={`text-sm font-medium leading-tight ${item.completed ? "text-foreground" : "text-muted-foreground"}`}
                        >
                          {item.title}
                        </p>
                        {item.action ? (
                          <button
                            type="button"
                            onClick={item.action.onClick}
                            className="shrink-0 rounded-md border border-border px-2.5 py-1 text-xs font-medium text-muted-foreground transition-colors hover:border-primary/40 hover:bg-muted hover:text-foreground"
                          >
                            {item.action.label}
                          </button>
                        ) : null}
                      </div>
                      {item.description ? (
                        <p className="text-xs leading-snug text-muted-foreground">
                          {item.description}
                        </p>
                      ) : null}
                    </div>
                  </li>
                );
              })}
            </ul>

            {onDismiss ? (
              <button
                type="button"
                onClick={onDismiss}
                className="mt-3 text-xs font-medium text-muted-foreground transition-colors hover:text-foreground"
              >
                Hide checklist
              </button>
            ) : null}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
