import { AlertCircle } from "lucide-react"
import type { FieldValues, UseFormRegister } from "react-hook-form"
import type { FormFieldI } from "@/shared/model/field"
import { Input } from '@/shared/ui/shadcn/input'
import { Label } from '@/shared/ui/shadcn/label'
import { cn } from '@/shared/lib/utils'

export interface NumberFormFieldProps<T extends FieldValues> extends FormFieldI<T> {
  min?: number;
  max?: number;
  step?: number;
}

interface FormFieldNumberComponentProps<T extends FieldValues> {
  idPrefix?: string;
  field: NumberFormFieldProps<T>;
  register: UseFormRegister<T>;
  error?: { message?: string };
  loading?: boolean;
}

export function FormFieldNumber<T extends FieldValues>({
  idPrefix = '',
  field,
  register,
  error,
  loading,
}: FormFieldNumberComponentProps<T>) {
  const fieldName = field.name;
  const uniqueId = `${idPrefix}${fieldName.replace(/\./g, '-')}`;

  const baseInputClass = cn(
    "w-full px-3 py-2.5 rounded-xl border bg-background text-sm transition-colors",
    "focus:outline-none focus:ring-2 focus:ring-primary/20",
    error ? "border-destructive focus:border-destructive" : "border-input focus:border-primary"
  );

  const inputType = field.type === 'percentage' ? 'number' : field.type;
  const stepValue = field.step ? field.step.toString() : (field.type === 'percentage' ? '0.01' : undefined);

  const validationOptions = {
    valueAsNumber: true,
    min: field.min,
    max: field.max,
  };

  return (
    <div className={cn("space-y-1.5", field.span === 'full' ? 'sm:col-span-2' : '')}>
      <Label htmlFor={uniqueId} className="text-sm font-medium flex items-center gap-1">
        {field.label}
        {field.required && <span className="text-destructive">*</span>}
      </Label>

      <Input
        id={uniqueId}
        type={inputType}
        placeholder={field.placeholder}
        min={field.min !== undefined ? field.min.toString() : undefined}
        max={field.max !== undefined ? field.max.toString() : undefined}
        step={stepValue}
        className={baseInputClass}
        autoComplete={field.autocomplete}
        autoFocus={field.autoFocus}
        {...register(fieldName, validationOptions)}
        disabled={loading}
      />

      {error && (
        <span className="text-xs text-destructive flex items-center gap-1">
          <AlertCircle className="w-3 h-3" />
          {typeof error.message === "string" ? error.message : ""}
        </span>
      )}
    </div>
  );
}
