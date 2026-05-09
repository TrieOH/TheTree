import { useState, useRef, useCallback, useEffect } from "react"
import { X, Plus, Image as ImageIcon, Layout, Star } from "lucide-react"
import { cn } from '@/shared/lib/utils'

interface GalleryItemAction {
  label: string
  icon?: 'image' | 'layout' | 'star'
  onClick: (url: string) => void
}

interface GalleryUploadFieldProps {
  id?: string
  name?: string
  value?: string[]
  pendingFiles?: File[]
  onChange: (urls: string[]) => void
  onBlur?: () => void
  onFileSelect?: (files: File[]) => void
  accept?: string
  maxSize?: number
  disabled?: boolean
  itemActions?: GalleryItemAction[]
}

const IconMap = {
  image: ImageIcon,
  layout: Layout,
  star: Star,
}

export default function GalleryUploadField({
  id,
  name,
  value = [],
  pendingFiles = [],
  onChange,
  onBlur,
  onFileSelect,
  accept = "image/png,image/jpeg,image/webp",
  maxSize = 5 * 1024 * 1024,
  disabled,
  itemActions = [],
}: GalleryUploadFieldProps) {
  const [previews, setPreviews] = useState<string[]>([])
  const [error, setError] = useState<string | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    const newPreviews = pendingFiles.map(file => URL.createObjectURL(file))
    setPreviews(newPreviews)
    return () => { newPreviews.forEach(url => { URL.revokeObjectURL(url); }); }
  }, [pendingFiles])

  const validateFile = (file: File): string | null => {
    if (maxSize && file.size > maxSize) return `Arquivo muito grande (${file.name}). Máximo: ${(maxSize / 1024 / 1024).toFixed(1)}MB`
    const acceptedTypes = accept.split(',').map(t => t.trim())
    if (!acceptedTypes.some(type => file.type.match(type.replace('*', '.*')))) return `Tipo não suportado (${file.name}).`
    return null
  }

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files ?? [])
    const validFiles: File[] = []
    let lastError: string | null = null

    for (const file of files) {
      const err = validateFile(file)
      if (err) lastError = err
      else validFiles.push(file)
    }

    if (lastError) setError(lastError)
    else setError(null)

    if (validFiles.length > 0) {
      onFileSelect?.([...pendingFiles, ...validFiles])
      onBlur?.()
    }
    e.target.value = ''
  }

  const removeExisting = (index: number) => {
    const next = [...value]
    next.splice(index, 1)
    onChange(next)
    onBlur?.()
  }

  const removePending = (index: number) => {
    const next = [...pendingFiles]
    next.splice(index, 1)
    onFileSelect?.(next)
    onBlur?.()
  }

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    const files = Array.from(e.dataTransfer.files)
    const validFiles: File[] = []
    files.forEach(file => { if (!validateFile(file)) validFiles.push(file) })
    if (validFiles.length > 0) onFileSelect?.([...pendingFiles, ...validFiles])
  }, [pendingFiles, onFileSelect])

  return (
    <div className="space-y-4">
      <div
        className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3"
        onDragOver={handleDragOver}
        onDrop={handleDrop}
      >
        {value.map((url, i) => (
          <div key={`existing-${i}`} className="relative aspect-square group rounded-xl overflow-hidden border border-border bg-muted">
            <img src={url} alt={`Gallery ${i}`} className="w-full h-full object-cover" />

            {/* Overlay */}
            <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex flex-wrap items-center justify-center p-2 gap-1.5">
              {itemActions.map((action, idx) => {
                const Icon = action.icon ? IconMap[action.icon] : ImageIcon
                return (
                  <button
                    key={idx}
                    type="button"
                    onClick={() => { action.onClick(url); }}
                    title={action.label}
                    className="p-1.5 rounded-full bg-white/90 text-foreground hover:bg-white hover:scale-110 transition-all active:scale-95 shadow-sm"
                  >
                    <Icon className="w-3.5 h-3.5" />
                  </button>
                )
              })}

              <button
                type="button"
                onClick={() => { removeExisting(i); }}
                disabled={disabled}
                title="Remover"
                className="p-1.5 rounded-full bg-destructive text-destructive-foreground hover:scale-110 transition-all active:scale-95 shadow-sm"
              >
                <X className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        ))}

        {previews.map((url, i) => (
          <div key={`pending-${i}`} className="relative aspect-square group rounded-xl overflow-hidden border-2 border-primary/30 bg-primary/5">
            <img src={url} alt={`Pending ${i}`} className="w-full h-full object-cover opacity-70" />
            <div className="absolute top-1 right-1 bg-primary text-primary-foreground text-[8px] px-1 rounded uppercase font-bold">Pendente</div>
            <div className="absolute inset-0 flex items-center justify-center">
              <button
                type="button"
                onClick={() => { removePending(i); }}
                disabled={disabled}
                className="p-1.5 rounded-full bg-destructive text-destructive-foreground hover:scale-110 transition-all shadow-sm"
              >
                <X className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        ))}

        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          disabled={disabled}
          className={cn(
            "aspect-square rounded-xl border-2 border-dashed transition-all",
            "flex flex-col items-center justify-center gap-2",
            "border-muted-foreground/25 hover:border-primary hover:bg-primary/5 hover:text-primary",
            disabled && "opacity-50 cursor-not-allowed"
          )}
        >
          <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
            <Plus className="w-5 h-5 text-muted-foreground" />
          </div>
          <span className="text-xs font-medium">Adicionar</span>
        </button>
      </div>

      <input ref={inputRef} id={id} name={name} type="file" multiple accept={accept} onChange={handleFileSelect} disabled={disabled} className="hidden" />
      {error && <p className="text-xs text-destructive">{error}</p>}
    </div>
  )
}
