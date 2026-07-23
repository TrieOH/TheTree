import { useRef, useState } from 'react'
import type { FieldValues } from 'react-hook-form'
import { ImagePlus, UploadCloud } from 'lucide-react'
import type { FieldConfig, FieldFormApi } from '../../model/types'
import { useImageUploadField } from '../../hooks/use-image-upload-field'
import { ImageThumb } from './helper/image-thumb'
import { Button } from '@/shared/ui/shadcn/button'
import { cn } from '@/shared/lib/utils'
import { Label } from '@/shared/ui/shadcn/label'

export interface ImageFieldRendererProps<TFieldValues extends FieldValues> {
  field: FieldConfig<TFieldValues>
  form: FieldFormApi<TFieldValues>
}

export function ImageFieldRenderer<TFieldValues extends FieldValues>({
  field,
  form,
}: ImageFieldRendererProps<TFieldValues>) {
  if (field.kind !== 'image') return null

  const inputRef = useRef<HTMLInputElement>(null)
  const dragDepthRef = useRef(0)
  const [isDragging, setIsDragging] = useState(false)

  // Captured once (lazy initializer) so it doesn't change on every
  // render — see the note on `initialUrls` in useImageUploadField.
  const [initialUrls] = useState<string[]>(() => {
    const current = form.getValues(field.name)
    return typeof current === 'string' && current.length > 0 ? [current] : []
  })

  const { items, addFiles, removeItem } = useImageUploadField({
    fieldKey: String(field.name),
    initialUrls,
    maxItems: 1,
    accept: field.accept,
    maxSizeMB: field.maxSizeMB,
    onValueChange: (urls) => {
      form.setValue(field.name, (urls[0] ?? null) as never, {
        shouldDirty: true,
      })
    },
    onTrackingChange: field.onTrackingChange,
  })

  const item = items.at(0)
  const hasItem = item !== undefined

  return (
    <div className="space-y-1.5">
      <Label className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
        {field.label}
        {field.optional ? (
          <span className="ml-1 font-normal normal-case text-muted-foreground/70">
            (opcional)
          </span>
        ) : null}
      </Label>

      <div
        onDragEnter={(event) => {
          event.preventDefault()
          if (field.disabled) return
          dragDepthRef.current += 1
          setIsDragging(true)
        }}
        onDragOver={(event) => {
          event.preventDefault()
          if (field.disabled) return
          setIsDragging(true)
        }}
        onDragLeave={(event) => {
          event.preventDefault()
          dragDepthRef.current = Math.max(dragDepthRef.current - 1, 0)
          if (dragDepthRef.current === 0) setIsDragging(false)
        }}
        onDrop={(event) => {
          event.preventDefault()
          if (field.disabled) return
          dragDepthRef.current = 0
          setIsDragging(false)
          if (event.dataTransfer.files.length > 0)
            addFiles(event.dataTransfer.files)
        }}
        className={cn(
          'relative flex min-h-72 w-full overflow-hidden rounded-xl border border-dashed text-left transition-colors',
          isDragging
            ? 'border-primary bg-primary/5'
            : 'border-border bg-muted/20',
          field.disabled && 'cursor-not-allowed opacity-60',
        )}
      >
        {hasItem ? (
          <ImageThumb
            item={item}
            onRemove={() => removeItem(item.id)}
            className="absolute inset-0 h-full w-full rounded-none border-0"
          />
        ) : (
          <button
            type="button"
            onClick={() => {
              if (!field.disabled) inputRef.current?.click()
            }}
            disabled={field.disabled}
            className="flex h-full min-h-72 w-full flex-1 flex-col items-center justify-center gap-3 px-6 py-8 text-center"
          >
            <div className="flex size-14 items-center justify-center rounded-full border border-border bg-background shadow-sm">
              <UploadCloud className="h-6 w-6 text-muted-foreground" />
            </div>
            <div className="space-y-1">
              <p className="text-sm font-semibold text-foreground">
                Arraste e solte a imagem aqui
              </p>
              <p className="text-xs text-muted-foreground">
                {field.disabled
                  ? 'Upload disponível em breve'
                  : 'ou clique para selecionar um arquivo'}
              </p>
            </div>
            {field.hint ? (
              <p className="max-w-sm text-xs text-muted-foreground">
                {field.hint}
              </p>
            ) : null}
          </button>
        )}

        <div className="absolute bottom-4 left-1/2 z-10 -translate-x-1/2">
          <Button
            type="button"
            onClick={() => inputRef.current?.click()}
            disabled={field.disabled}
            variant="outline"
            size="sm"
            className="gap-2 bg-background/90 shadow-sm backdrop-blur-sm"
          >
            <ImagePlus className="h-4 w-4" />
            {hasItem ? 'Trocar imagem' : 'Selecionar imagem'}
          </Button>
        </div>

        <input
          ref={inputRef}
          type="file"
          disabled={field.disabled}
          accept={field.accept}
          className="hidden"
          onChange={(event) => {
            if (event.target.files && event.target.files.length > 0) {
              addFiles(event.target.files)
            }
            event.target.value = ''
          }}
        />
      </div>
    </div>
  )
}
