import { Controller } from "react-hook-form"
import { AlertCircle } from "lucide-react"
import ImageUploadField from "./image-upload-field"
import type { Control, FieldValues, Path, UseFormRegister } from "react-hook-form"
import type { FormFieldI } from "@/shared/model/field"
import { Input } from '@/shared/ui/shadcn/input'
import { Label } from '@/shared/ui/shadcn/label'
import { cn } from '@/shared/lib/utils'

interface FormFieldProps<T extends FieldValues> {
  field: FormFieldI<T>
  register: UseFormRegister<T>
  control: Control<T>
  error?: { message?: string }
  loading?: boolean
  pendingFile?: File | null
  onFileSelect: (file: File | null) => void
}

export function FormField<T extends FieldValues>({
  field,
  register,
  control,
  error,
  loading,
  pendingFile,
  onFileSelect,
}: FormFieldProps<T>) {
  const fieldName = field.name as Path<T>

  const baseInputClass = cn(
    "w-full px-3 py-2.5 rounded-xl border bg-background text-sm transition-colors",
    "focus:outline-none focus:ring-2 focus:ring-primary/20",
    error ? "border-destructive focus:border-destructive" : "border-input focus:border-primary"
  )

  const renderInput = () => {
    // Checkbox
    if (field.type === 'checkbox') {
      return (
        <div className="flex items-center gap-2 pt-1">
          <input
            type="checkbox"
            id={fieldName}
            {...register(fieldName)}
            className="w-4 h-4 rounded border-border accent-primary"
          />
          <Label htmlFor={fieldName} className="text-sm text-muted-foreground font-normal">
            {field.placeholder ?? 'Ativar'}
          </Label>
        </div>
      )
    }

    // Image Upload
    if (field.type === 'image-upload') {
      return (
        <Controller
          name={fieldName}
          control={control}
          render={({ field: { onChange, onBlur, value } }) => (
            <ImageUploadField
              value={value || ''}
              onChange={(url) => {
                onChange(url)
                if (url) onFileSelect(null)
              }}
              onBlur={onBlur}
              onFileSelect={onFileSelect}
              accept={field.accept}
              maxSize={field.maxSize}
              disabled={loading}
              placeholder={field.placeholder}
            />
          )}
        />
      )
    }

    // Select
    if (field.type === 'select') {
      return (
        <select
          id={fieldName}
          {...register(fieldName)}
          className={baseInputClass}
        >
          <option value="">Selecione...</option>
          {field.options?.map(opt => (
            <option key={opt.value} value={opt.value}>{opt.label}</option>
          ))}
        </select>
      )
    }

    // Textarea
    if (field.type === 'textarea') {
      return (
        <textarea
          id={fieldName}
          placeholder={field.placeholder}
          rows={field.rows ?? 3}
          {...register(fieldName)}
          className={cn(baseInputClass, 'resize-none')}
        />
      )
    }

    // Default input types
    return (
      <Input
        id={fieldName}
        type={field.type === 'percentage' ? 'number' : field.type}
        placeholder={field.placeholder}
        step={field.type === 'percentage' ? '0.01' : undefined}
        className={baseInputClass}
        {...register(fieldName, field.type === 'percentage' ? { valueAsNumber: true } : undefined)}
      />
    )
  }

  return (
    <div className={cn("space-y-1.5", field.span === 'full' ? 'sm:col-span-2' : '')}>
      <Label htmlFor={fieldName} className="text-sm font-medium flex items-center gap-1">
        {field.label}
        {field.required && <span className="text-destructive">*</span>}
        {pendingFile && (
          <span className="text-xs bg-primary/10 text-primary px-1.5 py-0.5 rounded">
            ↑
          </span>
        )}
      </Label>

      {renderInput()}

      {error && (
        <span className="text-xs text-destructive flex items-center gap-1">
          <AlertCircle className="w-3 h-3" />
          {typeof error.message === "string" ? error.message : ""}
        </span>
      )}
    </div>
  )
}
