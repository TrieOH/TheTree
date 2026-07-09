import { useEffect, useMemo, useRef, useState } from 'react'
import { Eye, FileImage, FileText, X } from 'lucide-react'
import { jsPDF } from 'jspdf'
import { useQuery } from '@tanstack/react-query'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/shared/ui/shadcn/dialog'
import { Button } from '@/shared/ui/shadcn/button'
import { Spinner } from '@/shared/ui/loader/spinner'
import { useCertificateCanvas } from '@/features/editor/use-certificate-canvas'
import { type CanvasElement, type ImageCanvasElement, type SignatureCanvasElement, type TextCanvasElement } from '@/features/editor/types'
import { parseRichTextMarkup, type RichTextStyle, type RichTextStyleMap } from '@/features/editor/rich-text'
import type { CertificationTemplateI } from '@/features/certifications/model'
import { allSignaturesQueryOptions } from '@/features/signatures/api'

const CERTIFICATE_WIDTH = 1000
const CERTIFICATE_HEIGHT = 707
const DEFAULT_TEXT_COLOR = '#111827'
const DEFAULT_TEXT_FONT_SIZE = 16
const DEFAULT_TEXT_FONT_FAMILY = 'Inter, system-ui, sans-serif'
const DEFAULT_LINE_HEIGHT = 1.16

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
        fontSize: DEFAULT_TEXT_FONT_SIZE,
        fontWeight: 400,
        fontFamily: 'Inter, system-ui, sans-serif',
        color: DEFAULT_TEXT_COLOR,
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

function splitTextLines(text: string) {
  return text.split('\n')
}

function styleSignature(style: RichTextStyle | null | undefined) {
  if (!style) return 'base'
  return [
    style.fontWeight ?? '',
    style.fontStyle ?? '',
    style.underline ? 'u' : '',
    style.fill ?? '',
    style.fontSize ?? '',
    style.fontFamily ?? '',
  ].join('|')
}

function styleToCss(style: RichTextStyle | null | undefined) {
  return {
    fontWeight: style?.fontWeight ?? undefined,
    fontStyle: style?.fontStyle ?? undefined,
    textDecorationLine: style?.underline ? 'underline' : undefined,
    color: style?.fill ?? DEFAULT_TEXT_COLOR,
    fontSize: `${DEFAULT_TEXT_FONT_SIZE}px`,
    fontFamily: DEFAULT_TEXT_FONT_FAMILY,
  } as const
}

export function substituteTemplateVariables(text: string, values: Record<string, string>) {
  return text.replace(/\{\{(\w+)\}\}/g, (_, key: string) => values[key] ?? `{{${key}}}`)
}

function buildRuns(lineText: string, lineStyles: Record<string, RichTextStyle> | undefined) {
  const runs: Array<{ text: string; style: RichTextStyle | null }> = []
  let currentStyle: RichTextStyle | null = null
  let currentText = ''

  for (let i = 0; i < lineText.length; i += 1) {
    const char = lineText[i]
    const nextStyle = lineStyles?.[String(i)] ?? null
    if (currentText && styleSignature(nextStyle) !== styleSignature(currentStyle)) {
      runs.push({ text: currentText, style: currentStyle })
      currentText = ''
    }

    currentStyle = nextStyle
    currentText += char
  }

  if (currentText) runs.push({ text: currentText, style: currentStyle })
  return runs
}

function splitRunByUrl(text: string) {
  const match = text.match(/https?:\/\/[^\s]+/i)
  if (!match?.[0]) return [{ text, isUrl: false }]

  const url = match[0]
  const [before, after] = text.split(url)
  const parts: Array<{ text: string; isUrl: boolean }> = []
  if (before) parts.push({ text: before, isUrl: false })
  parts.push({ text: url, isUrl: true })
  if (after) parts.push({ text: after, isUrl: false })
  return parts
}

function renderRichTextLine(
  lineText: string,
  styles: RichTextStyleMap,
  lineIndex: number
) {
  const runs = buildRuns(lineText, styles[String(lineIndex)])

  return runs.flatMap((run, runIndex) => {
    const parts = splitRunByUrl(run.text)
    const css = styleToCss(run.style)

    if (parts.length === 1) {
      return [(
        <span key={`${lineIndex}-${runIndex}`} style={css}>
          {run.text}
        </span>
      )]
    }

    return parts.map((part, partIndex) => {
      if (!part.isUrl) {
        return (
          <span key={`${lineIndex}-${runIndex}-${partIndex}`} style={css}>
            {part.text}
          </span>
        )
      }

      return (
        <a
          key={`${lineIndex}-${runIndex}-${partIndex}`}
          href={part.text}
          target="_blank"
          rel="noreferrer"
          style={css}
          className="underline decoration-current"
        >
          {part.text}
        </a>
      )
    })
  })
}

function renderTemplateTextOverlay(template: Pick<CertificationTemplateI, 'data' | 'id'>) {
  return template.data.elements
    .filter((element) => element.type === 'text')
    .map((element, index) => {
      const parsed = parseRichTextMarkup(element.content, {
        fontFamily: 'Inter, system-ui, sans-serif',
        fontSize: DEFAULT_TEXT_FONT_SIZE,
        color: DEFAULT_TEXT_COLOR,
      })
      const lines = splitTextLines(parsed.plainText)

      return (
        <div
          key={`${template.id}-${element.type}-${index}`}
          className="absolute select-text wrap-break-word"
          style={{
            left: `${element.xPct}%`,
            top: `${element.yPct}%`,
            width: `${element.widthPct}%`,
            color: DEFAULT_TEXT_COLOR,
            fontSize: `${DEFAULT_TEXT_FONT_SIZE}px`,
            lineHeight: `${DEFAULT_LINE_HEIGHT}`,
            textAlign: element.textAlign ?? 'left',
            fontFamily: DEFAULT_TEXT_FONT_FAMILY,
            fontWeight: 400,
            userSelect: 'text',
            whiteSpace: 'pre-wrap',
            transform: 'translate(-50%, -50%)',
          }}
        >
          {lines.map((line, lineIndex) => (
            <span key={`${template.id}-${index}-${lineIndex}`}>
              {renderRichTextLine(line, parsed.styles, lineIndex)}
              {lineIndex < lines.length - 1 ? <br /> : null}
            </span>
          ))}
        </div>
      )
    })
}

function triggerDownloadHref(href: string, filename: string) {
  const link = document.createElement('a')
  link.href = href
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
}

async function snapshotCanvasWithoutText(fabricCanvas: any) {
  const textObjects = (fabricCanvas.getObjects?.() ?? []).filter((obj: any) => obj?.type === 'textbox' || obj?.type === 'text')
  const previousState: Array<{ obj: any; visible: boolean; opacity: number }> = textObjects.map((obj: any) => ({
    obj,
    visible: obj.visible,
    opacity: obj.opacity,
  }))

  previousState.forEach((item: { obj: any; visible: boolean; opacity: number }) => {
    item.obj.set({ visible: false, opacity: 0 })
  })
  fabricCanvas.requestRenderAll?.()

  await new Promise<void>((resolve) => {
    requestAnimationFrame(() => resolve())
  })

  const dataUrl = fabricCanvas.toDataURL({
    format: 'png',
    multiplier: 2,
  })

  previousState.forEach((item: { obj: any; visible: boolean; opacity: number }) => {
    item.obj.set({ visible: item.visible, opacity: item.opacity })
  })
  fabricCanvas.requestRenderAll?.()

  return dataUrl
}

async function snapshotCanvasWithText(fabricCanvas: any) {
  const textObjects = (fabricCanvas.getObjects?.() ?? []).filter((obj: any) => obj?.type === 'textbox' || obj?.type === 'text')
  const previousState: Array<{ obj: any; visible: boolean; opacity: number }> = textObjects.map((obj: any) => ({
    obj,
    visible: obj.visible,
    opacity: obj.opacity,
  }))

  previousState.forEach((item: { obj: any; visible: boolean; opacity: number }) => {
    item.obj.set({ visible: true, opacity: 1 })
  })
  fabricCanvas.requestRenderAll?.()

  await new Promise<void>((resolve) => {
    requestAnimationFrame(() => resolve())
  })

  const dataUrl = fabricCanvas.toDataURL({
    format: 'png',
    multiplier: 2,
  })

  previousState.forEach((item: { obj: any; visible: boolean; opacity: number }) => {
    item.obj.set({ visible: item.visible, opacity: item.opacity })
  })
  fabricCanvas.requestRenderAll?.()

  return dataUrl
}

function drawTemplateTextOnPdf(pdf: jsPDF, template: Pick<CertificationTemplateI, 'data' | 'id'>) {
  template.data.elements.forEach((element) => {
    if (element.type !== 'text') return

    const parsed = parseRichTextMarkup(element.content, {
      fontFamily: 'Inter, system-ui, sans-serif',
      fontSize: DEFAULT_TEXT_FONT_SIZE,
      color: DEFAULT_TEXT_COLOR,
    })
    const lines = splitTextLines(parsed.plainText)
    const centerX = (element.xPct / 100) * CERTIFICATE_WIDTH
    const centerY = (element.yPct / 100) * CERTIFICATE_HEIGHT
    const width = (element.widthPct / 100) * CERTIFICATE_WIDTH
    const align = element.textAlign ?? 'left'
    const fontSize = DEFAULT_TEXT_FONT_SIZE
    const lineHeight = fontSize * DEFAULT_LINE_HEIGHT
    const totalHeight = Math.max(lineHeight, lines.length * lineHeight)
    const startY = centerY - (totalHeight / 2) + fontSize
    const baseX =
      align === 'center'
        ? centerX
        : align === 'right'
          ? centerX + (width / 2) - 2
          : centerX - (width / 2) + 2

    pdf.setFont('helvetica', 'normal')
    pdf.setFontSize(fontSize)
    pdf.setTextColor(17, 24, 39)

    lines.forEach((line, index) => {
      const y = startY + (index * lineHeight)
      pdf.text(line, baseX, y, { align })

      const linkMatch = line.match(/https?:\/\/[^\s]+/i)
      if (linkMatch?.[0]) {
        const textWidth = Math.max(1, pdf.getTextWidth(line))
        const linkX =
          align === 'center'
            ? baseX - (textWidth / 2)
            : align === 'right'
              ? baseX - textWidth
              : baseX
        pdf.link(linkX, y - fontSize, textWidth, lineHeight, {
          url: linkMatch[0],
        })
      }
    })
  })
}

export function CertificationTemplateStaticView({
  template,
  eventId,
  editionId,
  className,
}: {
  template: Pick<CertificationTemplateI, 'id' | 'title' | 'url' | 'data'>
  eventId: string
  editionId: string
  className?: string
}) {
  const previewHostRef = useRef<HTMLDivElement>(null)
  const elements = useMemo(() => toCanvasElements(template), [template])
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

  useEffect(() => {
    const fabricCanvas = canvasRef.current as any
    if (!isReady || !fabricCanvas) return

    const textObjects = (fabricCanvas.getObjects?.() ?? []).filter((obj: any) => obj?.type === 'textbox' || obj?.type === 'text')
    textObjects.forEach((obj: any) => {
      obj.set({ visible: false, opacity: 0 })
    })
    fabricCanvas.requestRenderAll?.()
  }, [canvasRef, isReady])

  return (
    <div className={className ?? 'space-y-3'}>
      <div className="flex flex-wrap gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => {
            const fabricCanvas = canvasRef.current as any
            if (!fabricCanvas) return
            void snapshotCanvasWithText(fabricCanvas).then((dataUrl) => {
              triggerDownloadHref(dataUrl, `${template.title}.png`)
            })
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
          onClick={async () => {
            const fabricCanvas = canvasRef.current as any
            if (!fabricCanvas) return
            const pngDataUrl = await snapshotCanvasWithoutText(fabricCanvas)
            const pdf = new jsPDF({
              orientation: 'landscape',
              unit: 'px',
              format: [CERTIFICATE_WIDTH, CERTIFICATE_HEIGHT],
              compress: true,
            })
            pdf.addImage(pngDataUrl, 'PNG', 0, 0, CERTIFICATE_WIDTH, CERTIFICATE_HEIGHT)
            drawTemplateTextOnPdf(pdf, template)
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
        <div className="absolute inset-0 z-20">
          {renderTemplateTextOverlay(template)}
        </div>
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
            {open && <CertificationTemplateStaticView eventId={eventId} editionId={editionId} template={template} />}
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}
