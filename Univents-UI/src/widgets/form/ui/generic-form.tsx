import { useForm, get } from "react-hook-form"
import { useState, useCallback, useEffect } from "react"
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema"
import { Loader2 } from "lucide-react"
import { toast } from "sonner"
import { FormField } from "./form-field"
import type { ZodType } from "zod"
import type { DefaultValues, FieldValues, Path, PathValue } from "react-hook-form"
import type { FormFieldI } from "@/shared/model/field"
import { Button } from '@/shared/ui/shadcn/button'

interface GenericFormProps<T extends FieldValues> {
  idPrefix?: string
  fields: readonly FormFieldI<T>[]
  schema: ZodType<T>
  onSubmit: (data: T) => void | Promise<void>
  onCancel: () => void
  defaultValues?: DefaultValues<T>
  submitLabel?: string
  loading?: boolean
}

type PendingFileType = File | File[] | null

export function GenericForm<T extends FieldValues>({
  idPrefix,
  fields,
  schema,
  onSubmit,
  onCancel,
  defaultValues,
  submitLabel = 'Salvar',
  loading: externalLoading = false,
}: GenericFormProps<T>) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [pendingFiles, setPendingFiles] = useState<Record<string, PendingFileType>>({})

  const loading = externalLoading || isSubmitting

  const { register, handleSubmit, formState: { errors }, reset, control, setValue, getValues } = useForm<T>({
    resolver: standardSchemaResolver(schema),
    defaultValues,
  })

  useEffect(() => {
    reset(defaultValues)
  }, [defaultValues, reset])

  const handleFileSelect = useCallback((fieldName: string, file: PendingFileType) => {
    setPendingFiles(prev => ({ ...prev, [fieldName]: file }))
  }, [])

  const processUploads = async (): Promise<Record<string, string | string[]>> => {
    const urls: Record<string, string | string[]> = {}

    for (const [fieldName, files] of Object.entries(pendingFiles)) {
      if (!files) continue

      const field = fields.find(f => f.name === fieldName)
      const uploadFn = field?.uploadFn
      if (!uploadFn) continue

      try {
        if (Array.isArray(files)) {
          const uploadedUrls = await Promise.all(files.map(file => uploadFn(file)))
          urls[fieldName] = uploadedUrls
        } else {
          const url = await uploadFn(files)
          urls[fieldName] = url
        }
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

      Object.entries(uploadedUrls).forEach(([key, value]) => {
        const fieldName = key as Path<T>
        if (Array.isArray(value)) {
          const current = data[key] as unknown
          const existing = Array.isArray(current) ? (current as unknown[]) : []
          setValue(fieldName, [...existing, ...value] as PathValue<T, Path<T>>)
        } else {
          setValue(fieldName, value as PathValue<T, Path<T>>)
        }
      })

      await onSubmit(getValues())
      setPendingFiles({})
      reset()
    } catch (error) {
      console.error('Erro no submit:', error)
      toast.error(error instanceof Error ? error.message : 'Erro ao processar formulário')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    setPendingFiles({})
    reset()
    onCancel()
  }

  const hasPendingUploads = Object.values(pendingFiles).some(f => f !== null && (Array.isArray(f) ? f.length > 0 : true))

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
            key={field.name}
            idPrefix={idPrefix}
            field={field}
            register={register}
            control={control}
            setValue={setValue}
            error={get(errors, field.name) as { message?: string } | undefined}
            loading={loading}
            pendingFile={pendingFiles[field.name] ?? null}
            onFileSelect={(file) => { handleFileSelect(field.name, file); }}
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
