import { useMemo, useRef, useState } from 'react'
import { Eye, FileImage, FileText, X } from 'lucide-react'
import { jsPDF } from 'jspdf'
import { useQuery } from '@tanstack/react-query'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/shared/ui/shadcn/dialog'
import { Button } from '@/shared/ui/shadcn/button'
import { Spinner } from '@/shared/ui/loader/spinner'
import { useCertificateCanvas } from '@/features/editor/use-certificate-canvas'
import { type CanvasElement, type ImageCanvasElement, type SignatureCanvasElement, type TextCanvasElement } from '@/features/editor/types'
import type { CertificationTemplateI } from '@/features/certifications/model'
import { allSignaturesQueryOptions } from '@/features/signatures/api'

const CERTIFICATE_WIDTH = 1000
const CERTIFICATE_HEIGHT = 707

function toCanvasElements(template: Pick<CertificationTemplateI, 'data' | 'id'>): CanvasElement[] {
  return template.data.elements.map((element, index) => {
    const base = {
      id: `${template.id}-${element.type}-${index}`,
      xPct: element.xPct,
      yPct: element.yPct,
      widthPct: element.widthPct,
      heightPct: element.heightPct,
      zIndex: index + 1,
    }

    if (element.type === 'text') {
      return {
        ...base,
        type: 'text' as const,
        content: element.content,
        fontSize: 16,
        fontWeight: 400,
        fontFamily: 'Inter, system-ui, sans-serif',
        color: '#111827',
        textAlign: element.textAlign ?? 'left',
      } as TextCanvasElement
    }

    if (element.type === 'signature') {
      return {
        ...base,
        type: 'signature' as const,
        title: element.title ?? 'Assinatura',
        signatureId: element.signatureId ?? null,
      } as SignatureCanvasElement
    }

    return {
      ...base,
      type: 'image' as const,
      src: element.src,
      fileName: element.fileName ?? undefined,
    } as ImageCanvasElement
  })
}

function extractVerifyUrl(template: Pick<CertificationTemplateI, 'data' | 'id'>) {
  for (const element of template.data.elements) {
    if (element.type !== 'text') continue
    const match = element.content.match(/https?:\/\/[^\s]+/i)
    if (match?.[0]) return {
      url: match[0],
      xPct: element.xPct,
      yPct: element.yPct,
      widthPct: element.widthPct,
      heightPct: element.heightPct,
    }
  }
  return null
}

function triggerDownloadHref(href: string, filename: string) {
  const link = document.createElement('a')
  link.href = href
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
}

function CertificationTemplatePreviewCanvas({
  template,
  eventId,
  editionId,
}: {
  template: Pick<CertificationTemplateI, 'id' | 'title' | 'url' | 'data'>
  eventId: string
  editionId: string
}) {
  const previewHostRef = useRef<HTMLDivElement>(null)
  const elements = useMemo(() => toCanvasElements(template), [template])
  const verifyLink = useMemo(() => extractVerifyUrl(template), [template])
  const { data: signatures = [] } = useQuery(allSignaturesQueryOptions(eventId, editionId))
  const signatureUrlsById = useMemo(
    () => Object.fromEntries(signatures.map((signature) => [signature.id, signature.url])),
    [signatures]
  )

  const { canvasRef, isReady } = useCertificateCanvas({
    canvasHostRef: previewHostRef,
    backgroundUrl: template.data.background ?? template.url ?? null,
    elements,
    selectedElementId: null,
    signatureUrlsById,
    onElementsChange: () => { },
    onElementSelect: () => { },
    previewMode: true,
  })

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => {
            const fabricCanvas = canvasRef.current as any
            if (!fabricCanvas) return
            const dataUrl = fabricCanvas.toDataURL({
              format: 'png',
              multiplier: 2,
            })
            triggerDownloadHref(dataUrl, `${template.title}.png`)
          }}
          disabled={!isReady}
          className="gap-2"
        >
          <FileImage className="size-4" />
          PNG
        </Button>
        <Button
          type="button"
          variant="outline"
          size="sm"
        onClick={() => {
          const fabricCanvas = canvasRef.current as any
          if (!fabricCanvas) return
          const pngDataUrl = fabricCanvas.toDataURL({
            format: 'png',
            multiplier: 2,
            })
            const pdf = new jsPDF({
              orientation: 'landscape',
              unit: 'px',
              format: [CERTIFICATE_WIDTH, CERTIFICATE_HEIGHT],
              compress: true,
            })
            pdf.addImage(pngDataUrl, 'PNG', 0, 0, CERTIFICATE_WIDTH, CERTIFICATE_HEIGHT)
            if (verifyLink?.url) {
              const textX = (verifyLink.xPct / 100) * CERTIFICATE_WIDTH
              const textY = (verifyLink.yPct / 100) * CERTIFICATE_HEIGHT
              const textW = (verifyLink.widthPct / 100) * CERTIFICATE_WIDTH
              const textH = (verifyLink.heightPct / 100) * CERTIFICATE_HEIGHT
              pdf.link(textX, textY, textW, textH, { url: verifyLink.url })
            }
            pdf.save(`${template.title}.pdf`)
          }}
          disabled={!isReady}
          className="gap-2"
        >
          <FileText className="size-4" />
          PDF
        </Button>
      </div>

      <div className="relative aspect-1000/707 w-full overflow-hidden rounded-2xl border bg-muted/10">
        <div ref={previewHostRef} className="absolute inset-0" />
        {verifyLink && (
          <a
            href={verifyLink.url}
            target="_blank"
            rel="noreferrer"
            className="absolute z-20 block rounded-sm outline-none ring-0"
            style={{
              left: `${verifyLink.xPct}%`,
              top: `${verifyLink.yPct}%`,
              width: `${verifyLink.widthPct}%`,
              height: `${verifyLink.heightPct}%`,
            }}
            aria-label="Abrir link de verificação"
          />
        )}
        {!isReady && (
          <div className="pointer-events-none absolute inset-0 z-10 flex items-center justify-center bg-background/80 backdrop-blur-sm">
            <div className="flex items-center gap-2 rounded-full border bg-background px-4 py-2 text-sm text-muted-foreground shadow-sm">
              <Spinner className="size-4" />
              Carregando certificado
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export interface CertificationTemplatePreviewProps {
  eventId: string
  editionId: string
  template: Pick<CertificationTemplateI, 'id' | 'title' | 'url' | 'data'>
  triggerLabel?: string
}

export function CertificationTemplatePreview({ eventId, editionId, template, triggerLabel = 'Ver certificado' }: CertificationTemplatePreviewProps) {
  const [open, setOpen] = useState(false)

  return (
    <>
      <Button
        type="button"
        variant="outline"
        size="sm"
        onClick={(e) => {
          e.stopPropagation()
          setOpen(true)
        }}
        className="gap-2"
      >
        <Eye className="size-4" />
        {triggerLabel}
      </Button>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-6xl p-0 sm:max-w-6xl" showCloseButton={false}>
          <div className="border-b px-6 py-4">
            <div className="flex items-start justify-between gap-4">
              <DialogHeader className="space-y-1">
                <DialogTitle className="text-base">{template.title}</DialogTitle>
                <DialogDescription>
                  Visualização do certificado montado com o fundo e os elementos salvos.
                </DialogDescription>
              </DialogHeader>
              <Button type="button" variant="ghost" size="sm" onClick={() => setOpen(false)} className="gap-2">
                <X className="size-4" />
                Fechar
              </Button>
            </div>
          </div>

          <div className="p-4 sm:p-6">
            {open && <CertificationTemplatePreviewCanvas eventId={eventId} editionId={editionId} template={template} />}
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}
