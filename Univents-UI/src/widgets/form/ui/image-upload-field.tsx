import { useState, useCallback, useRef, useEffect } from "react"
import { Upload, X } from "lucide-react"
import { cn } from '@/shared/lib/utils'
import { Button } from '@/shared/ui/shadcn/button'

interface ImageUploadFieldProps {
  value?: string
  onChange: (url: string) => void
  onBlur?: () => void
  onFileSelect?: (file: File | null) => void
  accept?: string
  maxSize?: number // bytes
  disabled?: boolean
  placeholder?: string
}

export default function ImageUploadField({
  value,
  onChange,
  onBlur,
  onFileSelect,
  accept = "image/png,image/jpeg,image/webp",
  maxSize = 5 * 1024 * 1024,
  disabled,
  placeholder = "Arraste uma imagem ou clique para selecionar",
}: ImageUploadFieldProps) {
  const [isDragging, setIsDragging] = useState(false)
  const [preview, setPreview] = useState<string | null>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [error, setError] = useState<string | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (!value && !selectedFile) setPreview(null)
  }, [value, selectedFile])

  const validateFile = (file: File): string | null => {
    if (maxSize && file.size > maxSize) {
      return `Arquivo muito grande. Máximo: ${(maxSize / 1024 / 1024).toFixed(1)}MB`
    }
    const acceptedTypes = accept.split(',').map(t => t.trim())
    if (!acceptedTypes.some(type => file.type.match(type.replace('*', '.*')))) {
      return `Tipo não suportado. Use: ${accept.replace(/image\//g, '').replace(/,/g, ', ')}`
    }
    return null
  }

  const handleFileSelection = (file: File) => {
    const validationError = validateFile(file)
    if (validationError) {
      setError(validationError)
      return
    }

    setError(null)
    const previewUrl = URL.createObjectURL(file)

    if (preview) URL.revokeObjectURL(preview)

    setPreview(previewUrl)
    setSelectedFile(file)

    onFileSelect?.(file)
    onBlur?.()
  }

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) handleFileSelection(file)
  }

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (!e.currentTarget.contains(e.relatedTarget as Node)) {
      setIsDragging(false)
    }
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)
    const file = e.dataTransfer.files[0]
    handleFileSelection(file)
  }, [])

  const handleRemove = () => {
    if (preview) {
      URL.revokeObjectURL(preview)
      setPreview(null)
    }
    setSelectedFile(null)
    setError(null)
    onFileSelect?.(null)
    onChange('')
  }

  const handleClick = () => {
    if (!disabled) inputRef.current?.click()
  }

  const displayImage = preview ?? value

  if (displayImage) {
    return (
      <div className="relative group">
        <div className="relative rounded-xl overflow-hidden border border-border bg-muted">
          <img
            src={displayImage}
            alt="Preview"
            className="w-full h-40 object-cover"
          />

          {selectedFile && (
            <div className="absolute top-2 left-2 bg-primary text-primary-foreground text-xs px-2 py-1 rounded-md">
              Pronto para upload
            </div>
          )}

          <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              onClick={handleRemove}
              disabled={disabled}
              className="rounded-lg bg-white/90 hover:bg-white text-destructive"
            >
              <X className="w-4 h-4 mr-1" />
              Remover
            </Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-2">
      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={handleClick}
        className={cn(
          "relative border-2 border-dashed rounded-xl p-4 transition-all duration-200",
          "flex flex-col items-center justify-center gap-2",
          "min-h-30 cursor-pointer",
          isDragging
            ? "border-primary bg-primary/5 scale-[1.02]"
            : "border-muted-foreground/25 hover:border-muted-foreground/50 hover:bg-muted/30",
          disabled && "opacity-50 cursor-not-allowed",
          error && "border-destructive bg-destructive/5"
        )}
      >
        <input
          ref={inputRef}
          type="file"
          accept={accept}
          onChange={handleFileSelect}
          disabled={disabled}
          className="hidden"
        />

        <div className={cn(
          "w-10 h-10 rounded-xl flex items-center justify-center transition-colors",
          isDragging ? "bg-primary/20 text-primary" : "bg-muted text-muted-foreground"
        )}>
          <Upload className="w-5 h-5" />
        </div>

        <div className="text-center space-y-0.5">
          <p className="text-sm font-medium text-foreground">
            {isDragging ? 'Solte aqui' : placeholder}
          </p>
          <p className="text-xs text-muted-foreground">
            {accept.replace(/image\//g, '').replace(/,/g, ', ')} até {(maxSize / 1024 / 1024).toFixed(0)}MB
          </p>
        </div>
      </div>

      {error && (
        <p className="text-xs text-destructive">
          {error}
        </p>
      )}
    </div>
  )
}