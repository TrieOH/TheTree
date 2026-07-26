import type { FieldValues } from "react-hook-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/shadcn/dialog";
import type { MultiStepFormController } from "../hooks/use-multi-step-form";
import { MultiStepForm } from "./multi-step-form";

export interface MultiStepFormModalProps<
  TInput extends FieldValues,
  TOutput = TInput,
> {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  controller: MultiStepFormController<TInput, TOutput>;
  submitLabel?: string;
}

/**
 * Thin modal shell around <MultiStepForm />. Kept separate so the form
 * itself can be reused inline (a page, a drawer, etc.) without a dialog.
 */
export function MultiStepFormModal<
  TInput extends FieldValues,
  TOutput = TInput,
>({
  open,
  onOpenChange,
  title,
  controller,
  submitLabel,
}: MultiStepFormModalProps<TInput, TOutput>) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="flex h-[calc(100dvh-1rem)] max-h-168 flex-col overflow-hidden sm:max-w-lg sm:h-[min(90dvh,42rem)]">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>
        <MultiStepForm
          controller={controller}
          submitLabel={submitLabel}
          onCancel={() => onOpenChange(false)}
        />
      </DialogContent>
    </Dialog>
  );
}
