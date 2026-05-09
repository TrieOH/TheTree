import { useForm, get, useWatch } from "react-hook-form"
import { useState, useCallback } from "react"
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema"
import { Loader2, Save, X } from "lucide-react"
import { toast } from "sonner"
import { FormField } from "./form-field"
import type { ZodType } from "zod"
import type {
  DeepPartialSkipArrayKey,
  DefaultValues,
  FieldValues,
  Path,
  PathValue,
} from "react-hook-form"
import type { FormFieldI } from "@/shared/model/field"
import { Button } from "@/shared/ui/shadcn/button"

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
  submitLabel = "Salvar",
  loading: externalLoading = false,
}: GenericFormProps<T>) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [pendingFiles, setPendingFiles] = useState<Record<string, PendingFileType>>({})

  const loading = externalLoading || isSubmitting

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    control,
    setValue,
    getValues,
  } = useForm<T>({
    resolver: standardSchemaResolver(schema),
    defaultValues,
  })

  const currentValues = useWatch<T>({
    control,
    defaultValue: defaultValues as DeepPartialSkipArrayKey<T>,
  })


  const handleFileSelect = useCallback((fieldName: string, file: PendingFileType) => {
    setPendingFiles((prev) => ({ ...prev, [fieldName]: file }))
  }, [])

  const processUploads = async (): Promise<Record<string, string | string[]>> => {
    const urls: Record<string, string | string[]> = {}

    for (const [fieldName, files] of Object.entries(pendingFiles)) {
      if (!files) continue

      const field = fields.find((f) => f.name === fieldName)
      const uploadFn = field?.uploadFn
      if (!uploadFn) continue

      try {
        if (Array.isArray(files)) {
          const uploadedUrls = await Promise.all(files.map((file) => uploadFn(file)))
          urls[fieldName] = uploadedUrls
        } else {
          urls[fieldName] = await uploadFn(files)
        }
      } catch (error) {
        throw new Error(
          `Falha no upload de ${field.label}: ${error instanceof Error ? error.message : "Erro desconhecido"
          }`
        )
      }
    }

    return urls
  }

  const handleFormSubmit = async (data: T) => {
    setIsSubmitting(true)

    try {
      const uploadedUrls = await processUploads()

      for (const [key, value] of Object.entries(uploadedUrls)) {
        const fieldName = key as Path<T>
        if (Array.isArray(value)) {
          const current = data[key] as unknown
          const existing = Array.isArray(current) ? (current as unknown[]) : []
          setValue(fieldName, [...existing, ...value] as PathValue<T, Path<T>>)
        } else {
          setValue(fieldName, value as PathValue<T, Path<T>>)
        }
      }

      await onSubmit(getValues())
      setPendingFiles({})
      reset()
    } catch (error) {
      console.error("Erro no submit:", error)
      toast.error(
        error instanceof Error ? error.message : "Erro ao processar formulário"
      )
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    setPendingFiles({})
    reset()
    onCancel()
  }

  const hasPendingUploads = Object.values(pendingFiles).some(
    (f) => f !== null && (Array.isArray(f) ? f.length > 0 : true)
  )

  return (
    <form
      onSubmit={(e) => void handleSubmit(handleFormSubmit)(e)}
      className="flex flex-col gap-6 p-6 overflow-y-auto"
    >
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4">
        {fields.map((field) => {
          if (
            field.rules?.isVisible &&
            !field.rules.isVisible(currentValues as T)
          ) {
            return null
          }

          return (
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
          )
        })}
      </div>

      <div className="flex items-center justify-end gap-3 pt-2 border-t border-border">
        <Button
          type="button"
          variant="ghost"
          onClick={handleCancel}
          disabled={loading}
          className="gap-2 rounded-lg"
        >
          <X className="w-4 h-4" />
          Cancelar
        </Button>

        <Button
          type="submit"
          disabled={loading}
          className="gap-2 rounded-lg min-w-30"
        >
          {isSubmitting ? (
            <>
              <Loader2 className="w-4 h-4 animate-spin" />
              {hasPendingUploads ? "Enviando..." : "Salvando..."}
            </>
          ) : (
            <>
              <Save className="w-4 h-4" />
              {submitLabel}
            </>
          )}
        </Button>
      </div>
    </form>
  )
}