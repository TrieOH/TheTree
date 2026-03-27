import { useForm } from "react-hook-form"
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema"
import { AlertCircle } from "lucide-react"
import type { ZodType } from "zod"
import type { DefaultValues, FieldValues, Path } from "react-hook-form"
import type { FormFieldI } from "@/shared/model/field"
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
} from '@/shared/ui/shadcn/drawer'
import { Input } from '@/shared/ui/shadcn/input'
import { Label } from '@/shared/ui/shadcn/label'
import { Button } from '@/shared/ui/shadcn/button'
import { cn } from '@/shared/lib/utils'

interface FormDrawerProps<T extends FieldValues> {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  fields: readonly FormFieldI<T>[]
  schema: ZodType<T>
  onSubmit: (data: T) => void | Promise<void>
  defaultValues?: DefaultValues<T>
  submitLabel?: string
  loading?: boolean
}

export function FormDrawer<T extends FieldValues>({
  open,
  onOpenChange,
  title,
  fields,
  schema,
  onSubmit,
  defaultValues,
  submitLabel = 'Salvar',
  loading = false,
}: FormDrawerProps<T>) {
  const { register, handleSubmit, formState: { errors }, reset } = useForm<T>({
    resolver: standardSchemaResolver(schema),
    defaultValues,
  })

  const handleFormSubmit = async (data: T) => {
    await onSubmit(data)
    reset()
    onOpenChange(false)
  }

  const handleCancel = () => {
    reset()
    onOpenChange(false)
  }

  const renderField = (field: FormFieldI<T>) => {
    const fieldName = field.name as Path<T>
    const error = errors[fieldName]

    const baseInputClass = cn(
      "w-full px-3 py-2.5 rounded-xl border bg-background text-sm transition-colors",
      "focus:outline-none focus:ring-2 focus:ring-primary/20",
      error ? "border-destructive focus:border-destructive" : "border-input focus:border-primary"
    )

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
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent className="z-60 rounded-t-2xl border-t border-border bg-card max-h-[90vh]">
        <DrawerHeader className="pb-4 border-b border-border">
          <DrawerTitle className="text-base font-semibold text-left">
            {title}
          </DrawerTitle>
        </DrawerHeader>

        <form
          onSubmit={(e) => {
            void handleSubmit(handleFormSubmit)(e)
          }}
          className="p-4 overflow-y-auto"
        >
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {fields.map((field) => {
              const fieldName = field.name as Path<T>
              const error = errors[fieldName]

              return (
                <div
                  key={String(field.name)}
                  className={cn("space-y-1.5", field.span === 'full' ? 'sm:col-span-2' : '')}
                >
                  <Label htmlFor={fieldName} className="text-sm font-medium flex items-center gap-1">
                    {field.label}
                    {field.required && <span className="text-destructive">*</span>}
                  </Label>

                  {renderField(field)}

                  {error && (
                    <span className="text-xs text-destructive flex items-center gap-1">
                      <AlertCircle className="w-3 h-3" />
                      {typeof error.message === "string" ? error.message : ""}
                    </span>
                  )}
                </div>
              )
            })}
          </div>

          <div className="flex gap-2 pt-6">
            <Button
              type="submit"
              disabled={loading}
              className="flex-1 rounded-xl font-medium"
            >
              {loading ? 'Salvando...' : submitLabel}
            </Button>
            <Button
              type="button"
              variant="secondary"
              onClick={handleCancel}
              className="rounded-xl font-medium"
            >
              Cancelar
            </Button>
          </div>
        </form>
      </DrawerContent>
    </Drawer>
  )
}