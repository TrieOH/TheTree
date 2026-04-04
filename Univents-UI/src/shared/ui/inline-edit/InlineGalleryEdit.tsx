import { useState, useRef, useCallback, useEffect } from 'react'
import { X, Plus, Image as ImageIcon, Layout, Star, Upload, Loader2 } from 'lucide-react'
import { cn } from '@/shared/lib/utils'

interface GalleryItemAction {
  label: string
  icon?: 'image' | 'layout' | 'star'
  onClick: (url: string) => void
}

interface InlineGalleryEditProps {
  value: string[]
  onChange: (urls: string[]) => void
  isEditEnabled: boolean
  onUpload: (file: File) => Promise<string>
  itemActions?: GalleryItemAction[]
  className?: string
  accept?: string
  maxSize?: number
}

const IconMap = {
  image: ImageIcon,
  layout: Layout,
  star: Star,
}

interface PendingItem { file: File; previewUrl: string }

export default function InlineGalleryEdit({
  value,
  onChange,
  isEditEnabled,
  onUpload,
  itemActions = [],
  className,
  accept = 'image/png,image/jpeg,image/webp',
  maxSize = 5 * 1024 * 1024,
}: InlineGalleryEditProps) {
  const [pending, setPending] = useState<PendingItem[]>([])
  const [error, setError] = useState<string | null>(null)
  const [isDragOver, setIsDragOver] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    return () => {
      pending.forEach(({ previewUrl }) => { URL.revokeObjectURL(previewUrl); })
    }
  }, [])

  const validateFile = (file: File): string | null => {
    if (file.size > maxSize)
      return `Arquivo muito grande (${file.name}). Máximo: ${(maxSize / 1024 / 1024).toFixed(1)}MB`
    const acceptedTypes = accept.split(',').map((t) => t.trim())
    if (!acceptedTypes.some((type) => file.type.match(type.replace('*', '.*'))))
      return `Tipo não suportado (${file.name}).`
    return null
  }

  const processFiles = useCallback(
    async (files: File[]) => {
      setError(null)

      const validFiles: File[] = []
      for (const file of files) {
        const err = validateFile(file)
        if (err) { setError(err); continue }
        validFiles.push(file)
      }
      if (validFiles.length === 0) return

      const newItems: PendingItem[] = validFiles.map((file) => ({
        file,
        previewUrl: URL.createObjectURL(file),
      }))

      setPending((prev) => [...prev, ...newItems])

      for (const item of newItems) {
        try {
          const url = await onUpload(item.file)
          onChange([...value, url])
        } catch {
          setError(`Erro ao enviar ${item.file.name}`)
        } finally {
          URL.revokeObjectURL(item.previewUrl)
          setPending((prev) => prev.filter((p) => p.previewUrl !== item.previewUrl))
        }
      }
    },
    [value, onChange, onUpload],
  )

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files ?? [])
    void processFiles(files)
    e.target.value = ''
  }

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragOver(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
  }, [])

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault()
      e.stopPropagation()
      setIsDragOver(false)
      if (!isEditEnabled) return
      void processFiles(Array.from(e.dataTransfer.files))
    },
    [isEditEnabled, processFiles],
  )

  const removeImage = (index: number) => {
    const next = [...value]
    next.splice(index, 1)
    onChange(next)
  }

  // View-only
  if (!isEditEnabled) {
    if (value.length === 0) return null
    return (
      <div className={cn('grid grid-cols-3 gap-1.5', className)}>
        {value.map((url, i) => (
          <a
            key={i}
            href={url}
            target="_blank"
            rel="noopener noreferrer"
            className="aspect-square rounded-lg overflow-hidden bg-muted block"
          >
            <img
              src={url}
              alt={`Galeria ${i + 1}`}
              className="h-full w-full object-cover hover:scale-105 transition-transform duration-200"
            />
          </a>
        ))}
      </div>
    )
  }

  // Edit mode
  return (
    <div className={cn('space-y-2', className)}>
      <div
        className={cn(
          'grid grid-cols-3 gap-2 rounded-xl p-2 transition-colors',
          isDragOver && 'bg-primary/5 ring-2 ring-primary/20 ring-dashed',
        )}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        {/* Committed images */}
        {value.map((url, i) => (
          <div
            key={`existing-${i}`}
            className="relative aspect-square group rounded-xl overflow-hidden border border-border bg-muted"
          >
            <img src={url} alt={`Galeria ${i + 1}`} className="w-full h-full object-cover" />
            <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex flex-wrap items-center justify-center p-1.5 gap-1">
              {itemActions.map((action, idx) => {
                const Icon = action.icon ? IconMap[action.icon] : ImageIcon
                return (
                  <button
                    key={idx}
                    type="button"
                    onClick={() => { action.onClick(url); }}
                    title={action.label}
                    className="p-1.5 rounded-lg bg-background/90 hover:bg-background text-foreground hover:scale-110 transition-all active:scale-95 shadow-sm"
                  >
                    <Icon className="w-3.5 h-3.5" />
                  </button>
                )
              })}
              <button
                type="button"
                onClick={() => { removeImage(i); }}
                title="Remover"
                className="p-1.5 rounded-lg bg-destructive/90 hover:bg-destructive text-destructive-foreground hover:scale-110 transition-all active:scale-95 shadow-sm"
              >
                <X className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        ))}

        {/* Uploading previews */}
        {pending.map(({ previewUrl }, i) => (
          <div
            key={`pending-${i}`}
            className="relative aspect-square rounded-xl overflow-hidden border-2 border-primary/30 bg-primary/5"
          >
            <img src={previewUrl} alt="" className="w-full h-full object-cover opacity-50" />
            <div className="absolute inset-0 flex flex-col items-center justify-center gap-1.5">
              <Loader2 className="h-5 w-5 text-primary animate-spin" />
              <span className="text-[9px] font-bold uppercase tracking-wider text-primary bg-background/70 px-1.5 py-0.5 rounded">
                Enviando
              </span>
            </div>
          </div>
        ))}

        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          className={cn(
            'aspect-square rounded-xl border-2 border-dashed transition-all',
            'flex flex-col items-center justify-center gap-1.5',
            'border-muted-foreground/20 text-muted-foreground',
            'hover:border-primary/50 hover:bg-primary/5 hover:text-primary',
            'active:scale-[.97]',
          )}
        >
          <div className="w-7 h-7 rounded-lg bg-muted flex items-center justify-center">
            <Plus className="w-4 h-4" />
          </div>
          <span className="text-[10px] font-medium">Adicionar</span>
        </button>
      </div>

      {isDragOver && (
        <div className="flex items-center gap-2 text-xs text-primary/70">
          <Upload className="h-3.5 w-3.5" />
          <span>Solte para adicionar à galeria</span>
        </div>
      )}

      {error && <p className="text-xs text-destructive">{error}</p>}

      <input ref={inputRef} type="file" multiple accept={accept} onChange={handleFileSelect} className="hidden" />
    </div>
  )
}