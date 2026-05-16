import { Controller } from "react-hook-form"
import { AlertCircle } from "lucide-react"
import ImageUploadField from "./image-upload-field"
import GalleryUploadField from "./gallery-upload-field"
import { DateTimePicker } from "./date-time-picker"
import { FormFieldNumber } from "./form-field-number"
import type { Control, FieldValues, Path, PathValue, UseFormRegister, UseFormSetValue } from "react-hook-form"
import type { FormFieldI } from "@/shared/model/field"
import { Input } from '@/shared/ui/shadcn/input'
import { Label } from '@/shared/ui/shadcn/label'
import { cn } from '@/shared/lib/utils'

interface FormFieldProps<T extends FieldValues> {
  idPrefix?: string
  field: FormFieldI<T>
  register: UseFormRegister<T>
  control: Control<T>
  setValue: UseFormSetValue<T>
  error?: { message?: string }
  loading?: boolean
  pendingFile?: File | File[] | null
  onFileSelect: (file: File | File[] | null) => void
}

export function FormField<T extends FieldValues>({
  idPrefix = '',
  field,
  register,
  control,
  setValue,
  error,
  loading,
  pendingFile,
  onFileSelect,
}: FormFieldProps<T>) {
  const fieldName = field.name
  const uniqueId = `${idPrefix}${fieldName.replace(/\./g, '-')}`

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
            id={uniqueId}
            {...register(fieldName)}
            className="w-4 h-4 rounded border-border accent-primary"
          />
          <Label htmlFor={uniqueId} className="text-sm text-muted-foreground font-normal">
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
              id={uniqueId}
              name={fieldName}
              value={typeof value === 'string' ? value : ''}
              onChange={(url) => {
                onChange(url)
                if (url) onFileSelect(null)
              }}
              onBlur={onBlur}
              onFileSelect={(file) => { onFileSelect(file); }}
              accept={field.accept}
              maxSize={field.maxSize}
              disabled={loading}
              placeholder={field.placeholder}
            />
          )}
        />
      )
    }

    // Gallery Upload
    if (field.type === 'gallery-upload') {
      return (
        <Controller
          name={fieldName}
          control={control}
          render={({ field: { onChange, onBlur, value } }) => (
            <GalleryUploadField
              id={uniqueId}
              name={fieldName}
              value={Array.isArray(value) ? value : []}
              pendingFiles={Array.isArray(pendingFile) ? pendingFile : []}
              onChange={onChange}
              onBlur={onBlur}
              onFileSelect={(files) => { onFileSelect(files); }}
              accept={field.accept}
              maxSize={field.maxSize}
              disabled={loading}
              itemActions={field.itemActions?.map(action => ({
                label: action.label,
                icon: action.icon,
                onClick: (url: string) => {
                  const typedSetValue = (name: Path<T>, val: PathValue<T, Path<T>>) => {
                    setValue(name, val, { shouldDirty: true, shouldValidate: true })
                  }
                  action.onClick?.(url, typedSetValue)
                }
              }))}
            />
          )}
        />
      )
    }

    // DateTime Picker
    if (field.type === 'datetime') {
      return (
        <Controller
          name={fieldName}
          control={control}
          render={({ field: { onChange, value } }) => (
            <DateTimePicker
              id={uniqueId}
              value={typeof value === 'string' ? value : undefined}
              onChange={onChange}
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
          id={uniqueId}
          {...register(fieldName)}
          className={baseInputClass}
          autoComplete={field.autocomplete}
          autoFocus={field.autoFocus}
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
          id={uniqueId}
          placeholder={field.placeholder}
          rows={field.rows ?? 3}
          {...register(fieldName)}
          className={cn(baseInputClass, 'resize-none')}
          autoComplete={field.autocomplete}
          autoFocus={field.autoFocus}
        />
      )
    }

    // Number and Percentage fields
    if (field.type === 'number' || field.type === 'percentage') {
      return (
        <FormFieldNumber
          idPrefix={idPrefix}
          field={field}
          register={register}
          error={error}
          loading={loading}
        />
      );
    }

    // Default input types
    return (
      <Input
        id={uniqueId}
        type={field.type}
        placeholder={field.placeholder}
        className={baseInputClass}
        autoComplete={field.autocomplete}
        autoFocus={field.autoFocus}
        {...register(fieldName)}
        disabled={loading}
      />
    )
  }

  return (
    <div className={cn("space-y-1.5", field.span === 'full' ? 'sm:col-span-2' : '')}>
      <Label htmlFor={uniqueId} className="text-sm font-medium flex items-center gap-1">
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
