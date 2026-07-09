import {
  Bold,
  ImagePlus,
  Italic,
  Move,
  Signature,
  Trash2,
  Type,
  Underline,
  Image as ImageIcon,
  GripVertical,
} from 'lucide-react'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'
import ImageUploadField from '@/widgets/form/ui/image-upload-field'
import { Badge } from '@/shared/ui/shadcn/badge'
import { Button } from '@/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/shadcn/card'
import { Input } from '@/shared/ui/shadcn/input'
import { Label } from '@/shared/ui/shadcn/label'
import { Separator } from '@/shared/ui/shadcn/separator'
import { cn } from '@/shared/lib/utils'
import { Spinner } from '@/shared/ui/loader/spinner'
import { useCertificateCanvas } from '@/features/editor/use-certificate-canvas'
import VariableInput from '@/features/editor/variable-input'
import {
  toggleBoldMarkup,
  toggleColorMarkup,
  toggleItalicMarkup,
  toggleSizeMarkup,
  toggleUnderlineMarkup,
  parseRichTextMarkup,
  plainIndexToMarkupIndex,
} from '@/features/editor/rich-text'
import {
  type CanvasElement,
  type TextCanvasElement,
  type SignatureCanvasElement,
  type ImageCanvasElement,
} from '@/features/editor/types'
import { useNavigate } from '@tanstack/react-router'
import {
  createCertificationTemplateFn,
} from '@/features/certifications/api'
import {
  certificationTemplateCreateSchema,
  type CertificationTemplateCreateI,
} from '@/features/certifications/model'
import { uploadAndModerateFile } from '@/features/storage/api'
import { allSignaturesQueryOptions } from '@/features/signatures/api'

let elementCounter = 0
const MIN_EDITOR_WIDTH = 1280
const MIN_EDITOR_HEIGHT = 720

function generateId(): string {
  elementCounter += 1
  return `el-${elementCounter}-${Date.now()}`
}

function getViewportSize() {
  if (typeof window === 'undefined') return { width: 0, height: 0 }
  return { width: window.innerWidth, height: window.innerHeight }
}

export interface CertificationTemplateEditorProps {
  eventId: string
  editionId: string
}

export function CertificationTemplateEditor({ eventId, editionId }: CertificationTemplateEditorProps) {
  const navigate = useNavigate()
  const [templateTitle, setTemplateTitle] = useState('Modelo de Certificado')
  const [backgroundDataUrl, setBackgroundDataUrl] = useState<string | null>(null)
  const [backgroundFile, setBackgroundFile] = useState<File | null>(null)
  const [imageFilesById, setImageFilesById] = useState<Record<string, File | null>>({})
  const [elements, setElements] = useState<CanvasElement[]>([])
  const [selectedElementId, setSelectedElementId] = useState<string | null>(null)
  const [showGrid, setShowGrid] = useState(true)
  const [viewportSize, setViewportSize] = useState(getViewportSize)
  const [showInitialLoading, setShowInitialLoading] = useState(true)

  const canvasHostRef = useRef<HTMLDivElement>(null)
  const contentTextareaRef = useRef<HTMLTextAreaElement>(null)
  const { data: signatures = [] } = useQuery(allSignaturesQueryOptions(eventId, editionId))
  const signatureUrlsById = useMemo(
    () => Object.fromEntries(signatures.map((signature) => [signature.id, signature.url])),
    [signatures]
  )

  const { canvasRef, isReady } = useCertificateCanvas({
    canvasHostRef,
    backgroundUrl: backgroundDataUrl,
    elements,
    selectedElementId,
    signatureUrlsById,
    onElementsChange: setElements,
    onElementSelect: setSelectedElementId,
  })

  useEffect(() => {
    if (isReady) setShowInitialLoading(false)
  }, [isReady])

  useEffect(() => {
    const updateViewportSize = () => setViewportSize(getViewportSize())
    updateViewportSize()
    window.addEventListener('resize', updateViewportSize)
    return () => window.removeEventListener('resize', updateViewportSize)
  }, [])

  useEffect(() => {
    return () => {
      if (backgroundDataUrl?.startsWith('blob:')) {
        URL.revokeObjectURL(backgroundDataUrl)
      }
    }
  }, [backgroundDataUrl])

  const selectedElement = useMemo(
    () => elements.find((e) => e.id === selectedElementId) ?? null,
    [elements, selectedElementId]
  )

  const selectedTextElement = selectedElement?.type === 'text' ? selectedElement : null
  const selectedTextMarkup = useMemo(() => {
    if (!selectedTextElement) return null
    return parseRichTextMarkup(selectedTextElement.content, {
      fontFamily: selectedTextElement.fontFamily,
      fontSize: selectedTextElement.fontSize,
      color: selectedTextElement.color,
    })
  }, [
    selectedTextElement?.content,
    selectedTextElement?.fontFamily,
    selectedTextElement?.fontSize,
    selectedTextElement?.color,
  ])

  const isViewportAllowed =
    viewportSize.width >= MIN_EDITOR_WIDTH && viewportSize.height >= MIN_EDITOR_HEIGHT
  const isEditorDisabled = !isViewportAllowed

  const buildTemplatePayload = useCallback(async (): Promise<CertificationTemplateCreateI> => {
    const backgroundUrl = backgroundFile
      ? await uploadAndModerateFile(backgroundFile, `events/${eventId}/editions/${editionId}/certifications`)
      : backgroundDataUrl

    const payloadElements = await Promise.all(elements.map(async (e) => {
      const base = {
        xPct: e.xPct,
        yPct: e.yPct,
        widthPct: e.widthPct,
        heightPct: e.heightPct,
      }

      if (e.type === 'text') {
        return {
          type: 'text' as const,
          ...base,
          content: e.content,
        }
      }

      if (e.type === 'signature') {
        return {
          type: 'signature' as const,
          ...base,
          title: e.title ?? null,
          signatureId: e.signatureId,
        }
      }

      const file = imageFilesById[e.id]
      const src = (e as ImageCanvasElement).src
      if (src?.startsWith('blob:') || src?.startsWith('data:')) {
        if (!file) {
          throw new Error('A imagem do componente precisa ser reanexada antes de salvar')
        }
        const uploadedUrl = await uploadAndModerateFile(file, `events/${eventId}/editions/${editionId}/certifications`)
        return {
          type: 'image' as const,
          ...base,
          src: uploadedUrl,
          fileName: file.name,
        }
      }

      return {
        type: 'image' as const,
        ...base,
        src,
        fileName: (e as ImageCanvasElement).fileName ?? null,
      }
    }))

    const payload = {
      title: templateTitle,
      url: backgroundUrl,
      data: {
        background: backgroundUrl,
        elements: payloadElements,
      },
    } satisfies CertificationTemplateCreateI

    return certificationTemplateCreateSchema.parse(payload)
  }, [backgroundDataUrl, backgroundFile, elements, editionId, eventId, imageFilesById, templateTitle])

  const getDisplayName = (el: CanvasElement) => {
    if (el.type === 'text') return 'Texto'
    if (el.type === 'image') return 'Imagem'
    return 'Assinatura'
  }

  const selectedSignatureUrl =
    selectedElement?.type === 'signature' && selectedElement.signatureId
      ? signatureUrlsById[selectedElement.signatureId] ?? null
      : null
  const hasInvalidSignatureSelection = elements.some(
    (element) =>
      element.type === 'signature' &&
      (!element.signatureId || !signatureUrlsById[element.signatureId])
  )

  const addTextElement = () => {
    const newElement: TextCanvasElement = {
      id: generateId(),
      type: 'text',
      xPct: 50,
      yPct: 50,
      widthPct: 60,
      heightPct: 10,
      zIndex: elements.length + 1,
      content: '{{nome}}',
      fontSize: 16,
      fontWeight: 400,
      fontFamily: 'Inter, system-ui, sans-serif',
      color: '#111827',
    }
    setElements((prev) => [...prev, newElement])
    setSelectedElementId(newElement.id)
  }

  const addSignatureElement = () => {
    const newElement: SignatureCanvasElement = {
      id: generateId(),
      type: 'signature',
      xPct: 69,
      yPct: 74,
      widthPct: 23,
      heightPct: 10,
      zIndex: elements.length + 1,
      title: 'Assinatura',
      signatureId: null,
    }
    setElements((prev) => [...prev, newElement])
    setSelectedElementId(newElement.id)
  }

  const addImageElement = () => {
    const newElement: ImageCanvasElement = {
      id: generateId(),
      type: 'image',
      xPct: 50,
      yPct: 50,
      widthPct: 30,
      heightPct: 20,
      zIndex: elements.length + 1,
      src: null,
    }
    setElements((prev) => [...prev, newElement])
    setSelectedElementId(newElement.id)
  }

  const deleteSelectedElement = () => {
    if (!selectedElementId) return
    setImageFilesById((prev) => {
      const next = { ...prev }
      delete next[selectedElementId]
      return next
    })
    setElements((prev) => prev.filter((e) => e.id !== selectedElementId))
    setSelectedElementId(null)
  }

  const updateElement = (id: string, updates: Partial<CanvasElement>) => {
    setElements((prev) => prev.map((e) => (e.id === id ? ({ ...e, ...updates } as CanvasElement) : e)))
  }

  const updateTextContent = (content: string) => {
    if (!selectedElement || selectedElement.type !== 'text') return
    updateElement(selectedElement.id, { content } as Partial<CanvasElement>)
  }

  const getActiveTextInput = () => {
    const canvas = canvasRef.current
    const activeObject = canvas?.getActiveObject() as any

    if (activeObject?.type === 'textbox' && activeObject.isEditing) {
      const textarea = activeObject.hiddenTextarea ?? null
      return {
        kind: 'preview' as const,
        textarea,
        start: textarea?.selectionStart ?? activeObject.selectionStart ?? 0,
        end: textarea?.selectionEnd ?? activeObject.selectionEnd ?? 0,
      }
    }

    const textarea = contentTextareaRef.current
    if (!textarea) return null
    const { selectionStart, selectionEnd } = textarea
    if (typeof selectionStart !== 'number' || typeof selectionEnd !== 'number') return null
    return { kind: 'sidebar' as const, textarea, start: selectionStart, end: selectionEnd }
  }

  const applyMarkupChange = (
    updater: (value: string, selectionStart: number, selectionEnd: number) => {
      value: string
      selectionStart: number
      selectionEnd: number
    }
  ) => {
    if (!selectedElement || selectedElement.type !== 'text') return

    const currentValue = selectedElement.content
    const input = getActiveTextInput()
    const start = input?.start ?? currentValue.length
    const end = input?.end ?? currentValue.length
    const effectiveStart = input?.kind === 'preview' && selectedTextMarkup
      ? plainIndexToMarkupIndex(selectedTextMarkup, start)
      : start
    const effectiveEnd = input?.kind === 'preview' && selectedTextMarkup
      ? plainIndexToMarkupIndex(selectedTextMarkup, end)
      : end
    const next = updater(currentValue, effectiveStart, effectiveEnd)

    updateTextContent(next.value)
    window.requestAnimationFrame(() => {
      const nextInput = getActiveTextInput()
      const textarea = nextInput?.textarea
      if (!textarea) return
      textarea.focus()
      const selectionStart = input?.kind === 'preview' && selectedTextMarkup ? start : next.selectionStart
      const selectionEnd = input?.kind === 'preview' && selectedTextMarkup ? end : next.selectionEnd
      textarea.setSelectionRange(selectionStart, selectionEnd)
    })
  }

  const saveTemplateMutation = useMutation({
    mutationFn: (payload: CertificationTemplateCreateI) =>
      createCertificationTemplateFn(eventId, editionId, payload),
    onSuccess: (res) => {
      if (res.success) {
        toast.success('Template salvo com sucesso!')
        void navigate({
          to: '/admin/events/$eventId/editions/$editionId/certifications',
          params: { eventId, editionId },
        })
      } else {
        toast.error(res.message || 'Erro ao salvar template')
      }
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  const isTextSelected = selectedElement?.type === 'text'
  const pageTitle = 'Novo template'

  return (
    <div className="min-h-screen bg-background text-foreground overflow-x-hidden">
      {!isViewportAllowed && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-background/95 px-6 text-foreground backdrop-blur-sm">
          <div className="max-w-md rounded-2xl border bg-card p-6 text-center shadow-sm">
            <p className="text-sm font-semibold uppercase tracking-widest text-muted-foreground">Tela insuficiente</p>
            <h1 className="mt-3 text-2xl font-semibold">Abra no desktop</h1>
            <p className="mt-3 text-sm text-muted-foreground">
              O builder precisa de no mínimo {MIN_EDITOR_WIDTH}px de largura e {MIN_EDITOR_HEIGHT}px de altura.
            </p>
          </div>
        </div>
      )}

      <div className="mx-auto max-w-7xl px-4 py-6 md:px-6 md:py-10">
        <div className={cn(isEditorDisabled && 'pointer-events-none select-none opacity-40')}>
          <div className="mb-6 flex flex-col gap-1">
            <p className="text-xs font-medium uppercase tracking-[0.24em] text-muted-foreground">Certificações</p>
            <h1 className="text-2xl font-semibold">{pageTitle}</h1>
            <p className="text-sm text-muted-foreground">
              Ajuste o layout do template e salve para voltar à listagem.
            </p>
          </div>
          <div className="grid gap-6 lg:grid-cols-[280px_minmax(0,1fr)_320px]">
            <aside className="space-y-4 lg:sticky lg:top-6 lg:self-start">
              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="flex items-center gap-2 text-sm font-semibold">
                    <Type className="size-4 text-muted-foreground" />
                    Template
                  </CardTitle>
                  <CardDescription className="text-xs">
                    Nome, fundo e conteúdo do certificado.
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="space-y-1">
                    <Label className="text-xs text-muted-foreground">Nome do template</Label>
                    <Input
                      value={templateTitle}
                      onChange={(e) => setTemplateTitle(e.target.value)}
                      className="h-7 text-xs"
                      placeholder="Meu template"
                    />
                  </div>

                  <div className="space-y-1">
                    <Label className="text-xs text-muted-foreground">Imagem de fundo</Label>
                    <ImageUploadField
                      value={backgroundDataUrl ?? undefined}
                      onChange={(url) => {
                        setBackgroundDataUrl(url || null)
                        setBackgroundFile(null)
                      }}
                      onFileSelect={async (file) => {
                        if (!file) {
                          if (backgroundDataUrl?.startsWith('blob:')) {
                            URL.revokeObjectURL(backgroundDataUrl)
                          }
                          setBackgroundDataUrl(null)
                          setBackgroundFile(null)
                          return
                        }
                        setBackgroundFile(file)
                        if (backgroundDataUrl?.startsWith('blob:')) {
                          URL.revokeObjectURL(backgroundDataUrl)
                        }
                        setBackgroundDataUrl(URL.createObjectURL(file))
                      }}
                      accept="image/*"
                      placeholder="Enviar imagem de fundo"
                    />
                  </div>

                  <Button
                    type="button"
                    className="w-full text-xs"
                    onClick={async () => {
                      if (hasInvalidSignatureSelection) {
                        toast.error('Escolha uma assinatura válida para cada componente de assinatura')
                        return
                      }
                      try {
                        const payload = await buildTemplatePayload()
                        void saveTemplateMutation.mutateAsync(payload)
                      } catch (error) {
                        toast.error(error instanceof Error ? error.message : 'Erro ao preparar o template')
                      }
                    }}
                    disabled={saveTemplateMutation.isPending || hasInvalidSignatureSelection}
                  >
                    Salvar template
                  </Button>
                  {hasInvalidSignatureSelection && (
                    <p className="text-[11px] text-destructive">
                      Escolha uma assinatura válida para cada componente antes de salvar.
                    </p>
                  )}
                </CardContent>
              </Card>
            </aside>

            <section className="space-y-4">
              <Card>
                <CardHeader className="border-b pb-3">
                  <div className="flex items-center justify-between gap-2">
                    <div className="space-y-0.5">
                      <CardTitle className="flex items-center gap-2 text-sm font-semibold">
                        <ImagePlus className="size-4 text-muted-foreground" />
                        Preview
                      </CardTitle>
                      <CardDescription className="text-xs">
                        Monte o layout do template direto no canvas.
                      </CardDescription>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant="outline" className="gap-1 text-xs font-normal">
                        <Move className="size-3" />
                        {elements.length} el{elements.length !== 1 ? 's' : ''}
                      </Badge>
                      <Button type="button" variant="outline" size="sm" className="h-7 text-xs" onClick={() => setShowGrid((v) => !v)}>
                        {showGrid ? 'Ocultar grade' : 'Mostrar grade'}
                      </Button>
                    </div>
                  </div>
                </CardHeader>

                <div
                  className={cn(
                    'flex items-center gap-1 border-b px-3 py-1.5 transition-opacity flex-wrap',
                    isTextSelected ? 'bg-muted/20 border-border/50' : 'bg-muted/10 border-transparent'
                  )}
                >
                  <button type="button" onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); applyMarkupChange(toggleBoldMarkup) }} disabled={!isTextSelected} className={cn('rounded p-1 transition-colors', isTextSelected ? 'text-foreground hover:bg-accent' : 'text-muted-foreground/40 cursor-not-allowed')} title="Negrito">
                    <Bold className="size-3.5" />
                  </button>
                  <button type="button" onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); applyMarkupChange(toggleItalicMarkup) }} disabled={!isTextSelected} className={cn('rounded p-1 transition-colors', isTextSelected ? 'text-foreground hover:bg-accent' : 'text-muted-foreground/40 cursor-not-allowed')} title="Itálico">
                    <Italic className="size-3.5" />
                  </button>
                  <button type="button" onMouseDown={(e) => { e.preventDefault(); e.stopPropagation(); applyMarkupChange(toggleUnderlineMarkup) }} disabled={!isTextSelected} className={cn('rounded p-1 transition-colors', isTextSelected ? 'text-foreground hover:bg-accent' : 'text-muted-foreground/40 cursor-not-allowed')} title="Sublinhado">
                    <Underline className="size-3.5" />
                  </button>
                  <Separator orientation="vertical" className="mx-1 h-5" />

                  <input
                    type="number"
                    min={10}
                    max={120}
                    disabled={!isTextSelected}
                    className={cn('h-7 w-16 rounded border border-input bg-transparent px-1.5 text-xs', isTextSelected ? '' : 'text-muted-foreground/40')}
                    value={selectedElement && 'fontSize' in selectedElement ? (selectedElement as TextCanvasElement).fontSize : 32}
                    onChange={(e) => {
                      if (!selectedElement || selectedElement.type !== 'text') return
                      const v = Number(e.target.value)
                      if (Number.isNaN(v)) return
                      const clamped = Math.min(120, Math.max(10, v))
                      updateElement(selectedElement.id, { fontSize: clamped } as Partial<CanvasElement>)
                      applyMarkupChange((value, start, end) => toggleSizeMarkup(value, start, end, clamped))
                    }}
                    title="Tamanho da fonte"
                  />
                  <Separator orientation="vertical" className="mx-1 h-5" />
                  <input
                    type="color"
                    disabled={!isTextSelected}
                    className={cn('h-7 w-8 rounded border border-input bg-transparent p-0.5', isTextSelected ? '' : 'opacity-40')}
                    value={selectedElement && 'color' in selectedElement ? (selectedElement as TextCanvasElement).color : '#111827'}
                    onChange={(e) => {
                      if (!selectedElement || selectedElement.type !== 'text') return
                      const newColor = e.target.value
                      updateElement(selectedElement.id, { color: newColor } as Partial<CanvasElement>)
                      applyMarkupChange((value, start, end) => toggleColorMarkup(value, start, end, newColor))
                    }}
                    title="Cor do texto"
                  />
                  <span className="ml-auto text-[10px] text-muted-foreground">
                    {isTextSelected ? 'Selecione um trecho no conteúdo para aplicar markup' : 'Selecione um elemento de texto'}
                  </span>
                </div>

                <CardContent className="p-4">
                  <div className="relative aspect-297/210 w-full overflow-hidden rounded-xl border bg-background">
                    <div ref={canvasHostRef} className="absolute inset-0" />
                    <div
                      className={cn('pointer-events-none absolute inset-0 transition-opacity duration-200', showGrid ? 'opacity-30' : 'opacity-0')}
                      style={{
                        backgroundImage:
                          'linear-gradient(to right, rgba(0,0,0,0.12) 1px, transparent 1px), linear-gradient(to bottom, rgba(0,0,0,0.12) 1px, transparent 1px)',
                        backgroundSize: '50px 50px',
                      }}
                    />
                    {showInitialLoading && (
                      <div className="pointer-events-none absolute inset-0 z-10 flex items-center justify-center bg-background/75 backdrop-blur-sm">
                        <div className="flex items-center gap-2 rounded-full border bg-background px-4 py-2 text-sm text-muted-foreground shadow-sm">
                          <Spinner className="size-4" />
                          Carregando elementos
                        </div>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            </section>

            <aside className="space-y-4 lg:sticky lg:top-6 lg:self-start">
              {selectedElement ? (
                <Card>
                  <CardHeader className="pb-3">
                    <CardTitle className="flex items-center justify-between gap-2 text-sm font-semibold">
                      <span>{getDisplayName(selectedElement)}</span>
                      <Button type="button" variant="ghost" size="icon" className="size-6 text-muted-foreground hover:text-destructive" onClick={deleteSelectedElement}>
                        <Trash2 className="size-3.5" />
                      </Button>
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    <div className="grid grid-cols-2 gap-2">
                      <div className="space-y-1">
                        <Label className="text-xs text-muted-foreground">X</Label>
                        <div className="flex items-center gap-1">
                          <Input type="number" min={-15} max={115} value={Math.round(selectedElement.xPct)} onChange={(e) => { const v = Number(e.target.value); if (!Number.isNaN(v)) updateElement(selectedElement.id, { xPct: Math.min(115, Math.max(-15, v)) }) }} className="h-7 text-xs" />
                          <span className="text-xs text-muted-foreground shrink-0">%</span>
                        </div>
                      </div>
                      <div className="space-y-1">
                        <Label className="text-xs text-muted-foreground">Y</Label>
                        <div className="flex items-center gap-1">
                          <Input type="number" min={-15} max={115} value={Math.round(selectedElement.yPct)} onChange={(e) => { const v = Number(e.target.value); if (!Number.isNaN(v)) updateElement(selectedElement.id, { yPct: Math.min(115, Math.max(-15, v)) }) }} className="h-7 text-xs" />
                          <span className="text-xs text-muted-foreground shrink-0">%</span>
                        </div>
                      </div>
                      <div className="space-y-1">
                        <Label className="text-xs text-muted-foreground">Larg.</Label>
                        <div className="flex items-center gap-1">
                          <Input type="number" min={10} max={150} value={Math.round(selectedElement.widthPct)} onChange={(e) => { const v = Number(e.target.value); if (!Number.isNaN(v)) updateElement(selectedElement.id, { widthPct: Math.min(150, Math.max(10, v)) }) }} className="h-7 text-xs" />
                          <span className="text-xs text-muted-foreground shrink-0">%</span>
                        </div>
                      </div>
                      {selectedElement.type !== 'text' && (
                        <div className="space-y-1">
                          <Label className="text-xs text-muted-foreground">Alt.</Label>
                          <div className="flex items-center gap-1">
                            <Input type="number" min={6} max={80} value={Math.round(selectedElement.heightPct)} onChange={(e) => { const v = Number(e.target.value); if (!Number.isNaN(v)) updateElement(selectedElement.id, { heightPct: Math.min(80, Math.max(6, v)) }) }} className="h-7 text-xs" />
                            <span className="text-xs text-muted-foreground shrink-0">%</span>
                          </div>
                        </div>
                      )}
                    </div>

                    {selectedElement.type === 'text' && (
                      <VariableInput
                        id={`text-${selectedElement.id}`}
                        label="Conteúdo"
                        value={(selectedElement as TextCanvasElement).content}
                        onChange={(v) => updateElement(selectedElement.id, { content: v } as Partial<CanvasElement>)}
                        placeholder="Use **negrito**, *itálico*, __sublinhado__, [color=#22c55e], [size=24]"
                        multiline
                        textareaRef={contentTextareaRef}
                      />
                    )}

                    {selectedElement.type === 'signature' && (
                      <div className="space-y-3">
                        <div className="space-y-1">
                          <Label className="text-xs text-muted-foreground">Título</Label>
                          <Input
                            value={(selectedElement as SignatureCanvasElement).title ?? ''}
                            onChange={(e) => updateElement(selectedElement.id, { title: e.target.value } as Partial<CanvasElement>)}
                            className="h-7 text-xs"
                            placeholder="Assinatura"
                          />
                        </div>

                        <div className="space-y-1">
                          <Label className="text-xs text-muted-foreground">Assinatura usada</Label>
                          <select
                            className="h-7 w-full rounded-md border border-input bg-background px-2 text-xs"
                            value={(selectedElement as SignatureCanvasElement).signatureId ?? ''}
                            onChange={(e) => updateElement(selectedElement.id, { signatureId: e.target.value || null } as Partial<CanvasElement>)}
                          >
                            <option value="">Selecione uma assinatura</option>
                            {signatures.map((signature) => (
                              <option key={signature.id} value={signature.id}>
                                {signature.title}
                              </option>
                            ))}
                          </select>
                        </div>

                        {selectedSignatureUrl ? (
                          <div className="rounded-xl border bg-muted/10 p-3">
                            <p className="mb-2 text-[10px] uppercase tracking-wider text-muted-foreground">Preview da assinatura</p>
                            <img
                              src={selectedSignatureUrl}
                              alt={(selectedElement as SignatureCanvasElement).title ?? 'Assinatura'}
                              className="max-h-24 w-full object-contain"
                            />
                          </div>
                        ) : (
                          <p className="text-[11px] text-muted-foreground">
                            Selecione uma assinatura salva para renderizar a imagem no template.
                          </p>
                        )}
                      </div>
                    )}

                    {selectedElement.type === 'image' && (
                      <div className="space-y-1.5">
                        <Label className="text-xs text-muted-foreground">Imagem</Label>
                        <ImageUploadField
                          value={(selectedElement as ImageCanvasElement).src ?? undefined}
                          onChange={(v) => {
                            if (!v) {
                              setImageFilesById((prev) => {
                                const next = { ...prev }
                                delete next[selectedElement.id]
                                return next
                              })
                              updateElement(selectedElement.id, { src: null, fileName: null } as Partial<CanvasElement>)
                            }
                          }}
                          onFileSelect={async (file) => {
                            if (!file) return
                            const nextUrl = URL.createObjectURL(file)
                            const current = selectedElement as ImageCanvasElement
                            if (current.src?.startsWith('blob:')) URL.revokeObjectURL(current.src)
                            setImageFilesById((prev) => ({ ...prev, [selectedElement.id]: file }))
                            updateElement(selectedElement.id, {
                              src: nextUrl,
                              fileName: file.name,
                            } as Partial<CanvasElement>)
                          }}
                          accept="image/*"
                          placeholder="Escolher imagem"
                        />
                      </div>
                    )}
                  </CardContent>
                </Card>
              ) : elements.length > 0 ? (
                <Card className="border-dashed h-fit">
                  <CardContent className="p-6 text-center text-sm text-muted-foreground">
                    <p>Selecione um elemento</p>
                    <p className="text-xs mt-1">Clique em um elemento na tela ou na lista ao lado.</p>
                  </CardContent>
                </Card>
              ) : null}

              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="flex items-center gap-2 text-sm font-semibold">
                    <GripVertical className="size-4 text-muted-foreground" />
                    Componentes
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="grid grid-cols-3 gap-2">
                    <Button type="button" variant="outline" size="sm" onClick={addTextElement}>
                      <Type className="size-3.15" />
                      Texto
                    </Button>
                    <Button type="button" variant="outline" size="sm" onClick={addSignatureElement}>
                      <Signature className="size-3.5" />
                      Assin.
                    </Button>
                    <Button type="button" variant="outline" size="sm" onClick={addImageElement}>
                      <ImageIcon className="size-3.5" />
                      Imagem
                    </Button>
                  </div>

                  {elements.length > 0 && (
                    <div className="space-y-1 rounded-lg border bg-muted/30 p-1.5">
                      <p className="px-1.5 text-[10px] font-medium uppercase tracking-wider text-muted-foreground">Na tela</p>
                      {elements.map((el) => (
                        <button
                          key={el.id}
                          type="button"
                          onClick={() => setSelectedElementId(el.id)}
                          className={cn(
                            'flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-xs transition-colors',
                            selectedElementId === el.id ? 'bg-accent text-accent-foreground' : 'hover:bg-muted'
                          )}
                        >
                          {el.type === 'text' ? (
                            <Type className="size-3 shrink-0 text-muted-foreground" />
                          ) : el.type === 'image' ? (
                            <ImageIcon className="size-3 shrink-0 text-muted-foreground" />
                          ) : (
                            <Signature className="size-3 shrink-0 text-muted-foreground" />
                          )}
                          <span className="truncate">{getDisplayName(el)}</span>
                          <span className="ml-auto text-[10px] text-muted-foreground">
                            {Math.round(el.xPct)}%, {Math.round(el.yPct)}%
                          </span>
                        </button>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </aside>
          </div>
        </div>
      </div>
    </div>
  )
}
