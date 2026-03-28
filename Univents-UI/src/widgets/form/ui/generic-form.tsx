import { useForm } from "react-hook-form"
import { useState, useCallback, useEffect } from "react"
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema"
import { Loader2 } from "lucide-react"
import { FormField } from "./form-field"
import type { ZodType } from "zod"
import type { DefaultValues, FieldValues } from "react-hook-form"
import type { FormFieldI } from "@/shared/model/field"
import { Button } from '@/shared/ui/shadcn/button'

interface GenericFormProps<T extends FieldValues> {
  fields: readonly FormFieldI<T>[]
  schema: ZodType<T>
  onSubmit: (data: T) => void | Promise<void>
  onCancel: () => void
  defaultValues?: DefaultValues<T>
  submitLabel?: string
  loading?: boolean
}

export function GenericForm<T extends FieldValues>({
  fields,
  schema,
  onSubmit,
  onCancel,
  defaultValues,
  submitLabel = 'Salvar',
  loading: externalLoading = false,
}: GenericFormProps<T>) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [pendingFiles, setPendingFiles] = useState<Record<string, File | null>>({})

  const loading = externalLoading || isSubmitting

  const { register, handleSubmit, formState: { errors }, reset, control } = useForm<T>({
    resolver: standardSchemaResolver(schema),
    defaultValues,
  })

  useEffect(() => {
    reset(defaultValues)
  }, [defaultValues, reset])

  const handleFileSelect = useCallback((fieldName: string, file: File | null) => {
    setPendingFiles(prev => ({ ...prev, [fieldName]: file }))
  }, [])

  const processUploads = async (): Promise<Record<string, string>> => {
    const urls: Record<string, string> = {}

    for (const [fieldName, file] of Object.entries(pendingFiles)) {
      if (!file) continue

      const field = fields.find(f => String(f.name) === fieldName)
      if (!field?.uploadFn) continue

      try {
        const url = await field.uploadFn(file)
        urls[fieldName] = url
      } catch (error) {
        throw new Error(`Falha no upload de ${field.label}: ${error instanceof Error ? error.message : 'Erro desconhecido'}`)
      }
    }

    return urls
  }

  const handleFormSubmit = async (data: T) => {
    setIsSubmitting(true)

    try {
      const uploadedUrls = await processUploads()
      const finalData = { ...data, ...uploadedUrls } as T

      await onSubmit(finalData)
      setPendingFiles({})
      reset()
    } catch (error) {
      console.error('Erro no submit:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    setPendingFiles({})
    reset()
    onCancel()
  }

  const hasPendingUploads = Object.values(pendingFiles).some(f => f !== null)

  return (
    <form
      onSubmit={(e) => {
        void handleSubmit(handleFormSubmit)(e)
      }}
      className="p-4 overflow-y-auto"
    >
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {fields.map((field) => (
          <FormField
            key={String(field.name)}
            field={field}
            register={register}
            control={control}
            error={errors[field.name] as { message?: string }}
            loading={loading}
            pendingFile={pendingFiles[String(field.name)]}
            onFileSelect={(file) => { handleFileSelect(String(field.name), file); }}
          />
        ))}
      </div>

      <div className="flex gap-2 pt-6">
        <Button
          type="submit"
          disabled={loading}
          className="flex-1 rounded-xl font-medium"
        >
          {isSubmitting ? (
            <span className="flex items-center gap-2">
              <Loader2 className="w-4 h-4 animate-spin" />
              {hasPendingUploads ? 'Enviando imagens...' : 'Salvando...'}
            </span>
          ) : (
            submitLabel
          )}
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
  )
}
