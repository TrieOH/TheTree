import { ImagePlus, PenLine, Type } from 'lucide-react'
import { useRef, useState } from 'react'
import type { ChangeEvent } from 'react'
import {
  createImageElement,
  createSignatureElement,
  createTextElement,
} from '../factories'
import { certificateEditorActions, useCertificateEditorState } from '../store'
import { loadCertificateImageDimensions, readCertificateFile } from '../utils'
import { Button } from '@/shared/ui/shadcn/button'

const ACCEPTED_IMAGE_TYPES = new Set(['image/png', 'image/jpeg', 'image/webp'])

export function CertificateToolsSidebar() {
  const canvas = useCertificateEditorState((state) => state.canvas)
  const signatures = useCertificateEditorState(
    (state) => state.availableSignatures,
  )
  const imageInputRef = useRef<HTMLInputElement>(null)
  const [imageError, setImageError] = useState<string | null>(null)
  const [readingImage, setReadingImage] = useState(false)

  async function addImage(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0]
    event.target.value = ''
    if (!file) return

    if (!ACCEPTED_IMAGE_TYPES.has(file.type)) {
      setImageError('Use uma imagem PNG, JPEG ou WebP.')
      return
    }

    setReadingImage(true)
    setImageError(null)
    try {
      const src = await readCertificateFile(file)
      const naturalSize = await loadCertificateImageDimensions(src).catch(
        () => undefined,
      )
      certificateEditorActions.addElement(
        createImageElement(src, canvas, naturalSize),
      )
    } catch {
      setImageError('Não foi possível carregar a imagem.')
    } finally {
      setReadingImage(false)
    }
  }

  return (
    <aside className="flex w-64 shrink-0 flex-col gap-6 overflow-y-auto border-r border-border bg-card p-4 text-card-foreground">
      <section className="space-y-2.5">
        <h2 className="text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          Adicionar
        </h2>
        <div className="grid grid-cols-2 gap-2">
          <Button
            type="button"
            variant="outline"
            className="h-16 flex-col gap-1"
            onClick={() =>
              certificateEditorActions.addElement(createTextElement(canvas))
            }
          >
            <Type className="size-4" />
            Texto
          </Button>
          <Button
            type="button"
            variant="outline"
            className="h-16 flex-col gap-1"
            disabled={readingImage}
            onClick={() => imageInputRef.current?.click()}
          >
            <ImagePlus className="size-4" />
            {readingImage ? 'Carregando…' : 'Imagem'}
          </Button>
          <input
            ref={imageInputRef}
            type="file"
            accept="image/png,image/jpeg,image/webp"
            className="hidden"
            onChange={(event) => void addImage(event)}
          />
        </div>
        {imageError ? (
          <p className="text-xs text-destructive" role="alert">
            {imageError}
          </p>
        ) : null}
      </section>

      <section className="space-y-2.5">
        <div className="flex items-center gap-1.5 text-xs font-semibold tracking-wide text-muted-foreground uppercase">
          <PenLine className="size-3.5" />
          Assinaturas
        </div>
        {signatures.length === 0 ? (
          <p className="rounded-md border border-dashed p-2.5 text-center text-xs text-muted-foreground">
            Nenhuma assinatura disponível
          </p>
        ) : (
          <div className="grid grid-cols-2 gap-2">
            {signatures.map((signature) => (
              <button
                key={signature.id}
                type="button"
                onClick={() =>
                  certificateEditorActions.addElement(
                    createSignatureElement(signature, canvas),
                  )
                }
                className="overflow-hidden rounded-md border bg-popover text-left transition-colors hover:border-ring hover:bg-muted focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/50 focus-visible:outline-none"
                title={`Adicionar assinatura de ${signature.name}`}
              >
                <span className="flex h-16 items-center justify-center bg-white p-2">
                  <img
                    src={signature.url}
                    alt=""
                    className="max-h-full max-w-full object-contain"
                  />
                </span>
                <span className="block truncate border-t px-2 py-1.5 text-xs">
                  {signature.name}
                </span>
              </button>
            ))}
          </div>
        )}
      </section>
    </aside>
  )
}
