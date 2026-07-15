import { useEffect, useRef, useState } from 'react'
import { ImageUp, Upload, X } from 'lucide-react'
import { Button } from '@/shared/ui/shadcn/button'
import { cn } from '@/shared/lib/utils'

export interface SignatureImageSelectorProps {
  file: File | null
  onChange: (file: File | null) => void
}

export function SignatureImageSelector({ file, onChange }: SignatureImageSelectorProps) {
  const [isDragging, setIsDragging] = useState(false)
  const [previewUrl, setPreviewUrl] = useState<string | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const dragDepthRef = useRef(0)

  useEffect(() => {
    if (!file) {
      setPreviewUrl(null)
      return
    }

    const objectUrl = URL.createObjectURL(file)
    setPreviewUrl(objectUrl)

    return () => {
      URL.revokeObjectURL(objectUrl)
    }
  }, [file])

  const openFilePicker = () => {
    fileInputRef.current?.click()
  }

  const selectFile = (nextFile: File | null) => {
    onChange(nextFile)
  }

  const previewAvailable = Boolean(previewUrl)

  return (
    <div className="space-y-2">
      <input
        ref={fileInputRef}
        type="file"
        accept="image/png,image/jpeg,image/webp"
        className="hidden"
        onChange={(e) => {
          selectFile(e.target.files?.[0] ?? null)
        }}
      />

      <div
        onDragEnter={(e) => {
          e.preventDefault()
          e.stopPropagation()
          dragDepthRef.current += 1
          setIsDragging(true)
        }}
        onDragOver={(e) => {
          e.preventDefault()
          e.stopPropagation()
          setIsDragging(true)
        }}
        onDragLeave={(e) => {
          e.preventDefault()
          e.stopPropagation()
          dragDepthRef.current = Math.max(dragDepthRef.current - 1, 0)
          if (dragDepthRef.current === 0) setIsDragging(false)
        }}
        onDrop={(e) => {
          e.preventDefault()
          e.stopPropagation()
          dragDepthRef.current = 0
          setIsDragging(false)
          selectFile(e.dataTransfer.files[0])
        }}
        className={cn(
          'group relative overflow-hidden rounded-2xl border border-dashed transition-all',
          previewAvailable ? 'min-h-48 border-border bg-muted/10' : 'min-h-56 border-border bg-muted/10',
          isDragging ? 'border-primary bg-primary/5' : 'hover:border-primary/40',
        )}
        onClick={openFilePicker}
        role="button"
        tabIndex={0}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault()
            openFilePicker()
          }
        }}
      >
        {previewAvailable ? (
          <div className="absolute inset-0">
            <img src={previewUrl ?? undefined} alt="Prévia da assinatura" className="h-full w-full object-contain bg-background p-4" />
            <div className="absolute inset-0 bg-linear-to-t from-background/90 via-background/30 to-transparent" />
            <div className="absolute inset-x-0 bottom-0 flex items-end justify-between gap-3 p-4">
              <div className="min-w-0 space-y-1">
                <div className="inline-flex items-center gap-2 rounded-full border bg-background/80 px-2.5 py-1 text-[11px] font-medium text-foreground backdrop-blur-sm">
                  <ImageUp className="size-3.5" />
                  Imagem selecionada
                </div>
                <p className="truncate text-sm font-medium text-foreground">
                  {file?.name ?? 'Arquivo selecionado'}
                </p>
                <p className="text-xs text-muted-foreground">
                  Clique para trocar ou use o botão remover.
                </p>
              </div>
              <div className="flex shrink-0 items-center gap-2">
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  className="h-8 gap-2"
                  onClick={(event) => {
                    event.stopPropagation()
                    openFilePicker()
                  }}
                >
                  <Upload className="size-4" />
                  Trocar
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  size="icon-sm"
                  onClick={(event) => {
                    event.stopPropagation()
                    selectFile(null)
                  }}
                >
                  <X className="size-4" />
                </Button>
              </div>
            </div>
          </div>
        ) : (
          <div className="flex min-h-56 flex-col items-center justify-center gap-3 px-4 py-6 text-center">
            <div className="flex size-12 items-center justify-center rounded-full border bg-background shadow-sm">
              <Upload className="size-5 text-muted-foreground" />
            </div>
            <div className="space-y-1">
              <p className="text-sm font-medium text-foreground">
                Arraste e solte a assinatura aqui
              </p>
              <p className="text-xs text-muted-foreground">
                ou clique para selecionar um arquivo
              </p>
            </div>
            <p className="text-[11px] text-muted-foreground">
              PNG, JPG ou WEBP
            </p>
          </div>
        )}
      </div>
    </div>
  )
}
