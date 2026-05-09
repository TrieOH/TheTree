import { useState, useRef, useEffect } from "react";
import { PencilLine, ImageIcon } from "lucide-react";
import { cn } from '@/shared/lib/utils';

export interface InlineImageEditProps {
  value: string | null;
  onChange: (url: string) => void;
  isEditEnabled: boolean;
  isEditing: boolean;
  onStartEdit: () => void;
  onFinishEdit: () => void;
  onUpload?: (file: File) => Promise<string>;
  accept?: string;
  maxSize?: number;
  className?: string;
  placeholder?: string;
  renderDisplay?: (url: string | null) => React.ReactNode;
}

const InlineImageEdit = ({
  value,
  onChange,
  isEditEnabled,
  isEditing,
  onStartEdit,
  onFinishEdit,
  onUpload,
  accept = "image/png,image/jpeg,image/webp",
  maxSize = 5 * 1024 * 1024,
  className = '',
  placeholder = 'Clique para adicionar imagem...',
  renderDisplay,
}: InlineImageEditProps) => {
  const [preview, setPreview] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!value && !isEditing) setPreview(null);
  }, [value, isEditing]);

  useEffect(() => {
    if (!isEditing && preview) {
      URL.revokeObjectURL(preview);
      setPreview(null);
    }
  }, [isEditing]);

  const validateFile = (file: File): string | null => {
    if (maxSize && file.size > maxSize)
      return `Arquivo muito grande. Máximo: ${(maxSize / 1024 / 1024).toFixed(1)}MB`;
    const acceptedTypes = accept.split(',').map(t => t.trim());
    if (!acceptedTypes.some(type => file.type.match(type.replace('*', '.*'))))
      return `Tipo não suportado. Use: ${accept.replace(/image\//g, '').replace(/,/g, ', ')}`;
    return null;
  };

  const handleFile = async (file: File) => {
    const err = validateFile(file);
    if (err) { setError(err); return; }
    setError(null);

    if (preview) URL.revokeObjectURL(preview);
    const objectUrl = URL.createObjectURL(file);
    setPreview(objectUrl);
    onStartEdit();

    setIsUploading(true);
    try {
      if (onUpload) {
        const uploadedUrl = await onUpload(file);
        onChange(uploadedUrl);
        URL.revokeObjectURL(objectUrl);
        setPreview(null);
      } else {
        onChange(objectUrl);
      }
      onFinishEdit();
    } catch {
      setError('Erro ao fazer upload');
      setIsUploading(false);
    } finally {
      setIsUploading(false);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) void handleFile(file);
    e.target.value = '';
  };

  const openPicker = (e: React.MouseEvent) => {
    e.stopPropagation();
    inputRef.current?.click();
  };

  // Read only
  if (!isEditEnabled) {
    if (!value && !renderDisplay) return null;
    return (
      <div className={cn("relative", className)}>
        {renderDisplay
          ? renderDisplay(value)
          : value && <img src={value} alt="" className="w-full h-32 object-cover rounded-lg" />
        }
      </div>
    );
  }

  const displayUrl = preview ?? value;

  // Editable
  return (
    <button
      type="button"
      onClick={openPicker}
      className={cn(
        "group relative w-full text-left @container",
        "rounded-md transition-all cursor-pointer",
        "border border-dashed border-muted-foreground/25",
        "hover:border-muted-foreground/40",
        className
      )}
    >
      <input
        ref={inputRef}
        type="file"
        accept={accept}
        onChange={handleInputChange}
        className="hidden"
      />

      {renderDisplay
        ? renderDisplay(displayUrl)
        : displayUrl
          ? <img src={displayUrl} alt="" className="w-full h-32 object-cover rounded-[inherit]" />
          : <div className="flex flex-col items-center justify-center gap-2 py-4 px-2 min-h-24">
            <ImageIcon className="w-6 h-6 text-muted-foreground shrink-0" />
            <span className="text-xs text-muted-foreground text-center hidden @[120px]:inline-block line-clamp-2">
              {placeholder}
            </span>
          </div>
      }

      <div className={cn(
        "absolute inset-0 rounded-[inherit] flex flex-col items-center justify-center gap-1.5 transition-all duration-200 pointer-events-none",
        isUploading
          ? "bg-black/40"
          : "bg-black/0 group-hover:bg-black/30"
      )}>
        {isUploading ? (
          <div className="w-5 h-5 border-2 border-border border-t-transparent rounded-full animate-spin" />
        ) : (
          <span className="opacity-0 group-hover:opacity-100 transition-opacity text-white font-medium text-sm drop-shadow flex items-center gap-1.5 px-2">
            <PencilLine className="w-4 h-4 shrink-0" />
            <span className="hidden @[100px]:inline truncate">
              {displayUrl ? 'Trocar imagem' : 'Adicionar imagem'}
            </span>
          </span>
        )}
        {error && (
          <span className="text-[10px] leading-tight text-destructive px-2 text-center drop-shadow mt-1">
            {error}
          </span>
        )}
      </div>

      <span className={cn(
        "absolute -top-2 -right-2",
        "flex items-center justify-center",
        "h-6 w-6 rounded-md",
        "bg-background border border-muted-foreground/20",
        "shadow-sm cursor-pointer",
        "transition-all duration-200",
        isEditing || isUploading
          ? "bg-primary text-primary-foreground border-primary scale-110 -translate-y-0.5 translate-x-0.5"
          : "group-hover:bg-primary group-hover:text-primary-foreground group-hover:border-primary group-hover:scale-110 group-hover:-translate-y-0.5 group-hover:translate-x-0.5"
      )}>
        {isUploading
          ? <div className="w-2.5 h-2.5 border-[1.5px] border-current border-t-transparent rounded-full animate-spin" />
          : <PencilLine className="h-3 w-3" />
        }
      </span>
    </button>
  );
};

export default InlineImageEdit;