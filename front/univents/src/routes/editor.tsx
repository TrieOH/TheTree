import { createFileRoute } from '@tanstack/react-router'
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
import { useEffect, useMemo, useRef, useState, useCallback } from 'react'
import ImageUploadField from '@/widgets/form/ui/image-upload-field'
import { Badge } from '@/shared/ui/shadcn/badge'
import { Button } from '@/shared/ui/shadcn/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/shadcn/card'
import { Input } from '@/shared/ui/shadcn/input'
import { Label } from '@/shared/ui/shadcn/label'
import { Separator } from '@/shared/ui/shadcn/separator'
import { cn } from '@/shared/lib/utils'
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

let elementCounter = 0
const MIN_EDITOR_WIDTH = 1280
const MIN_EDITOR_HEIGHT = 720

function generateId(): string {
  elementCounter += 1
  return `el-${elementCounter}-${Date.now()}`
}

export const Route = createFileRoute('/editor')({
  component: RouteComponent,
})

function readFileAsDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

function getViewportSize() {
  if (typeof window === 'undefined') {
    return { width: 0, height: 0 }
  }

  return {
    width: window.innerWidth,
    height: window.innerHeight,
  }
}

function RouteComponent() {
  const [templateTitle, setTemplateTitle] = useState('Modelo de Certificado')
  const [backgroundDataUrl, setBackgroundDataUrl] = useState<string | null>(null)
  const [elements, setElements] = useState<CanvasElement[]>([])
  const [selectedElementId, setSelectedElementId] = useState<string | null>(null)
  const [showGrid, setShowGrid] = useState(true)
  const [viewportSize, setViewportSize] = useState(getViewportSize)

  const canvasHostRef = useRef<HTMLDivElement>(null)
  const canvasElRef = useRef<HTMLCanvasElement>(null)
  const contentTextareaRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    const updateViewportSize = () => setViewportSize(getViewportSize())
    updateViewportSize()
    window.addEventListener('resize', updateViewportSize)
    return () => window.removeEventListener('resize', updateViewportSize)
  }, [])

  const handleBackgroundFile = useCallback((file: File | null) => {
    if (!file) return
    void readFileAsDataUrl(file).then(setBackgroundDataUrl)
  }, [])

  const removeBackground = useCallback(() => {
    setBackgroundDataUrl(null)
  }, [])

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

  const getDisplayName = (el: CanvasElement) => {
    if (el.type === 'text') return 'Texto'
    if (el.type === 'image') return 'Imagem'
    return 'Assinatura'
  }

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
    setElements((prev) => prev.filter((e) => e.id !== selectedElementId))
    setSelectedElementId(null)
  }

  const updateElement = (id: string, updates: Partial<CanvasElement>) => {
    setElements((prev) =>
      prev.map((e) => (e.id === id ? { ...e, ...updates } as CanvasElement : e))
    )
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

  const toggleTextMarkup = (styleType: 'bold' | 'italic' | 'underline') => {
    switch (styleType) {
      case 'bold':
        applyMarkupChange(toggleBoldMarkup)
        break
      case 'italic':
        applyMarkupChange(toggleItalicMarkup)
        break
      case 'underline':
        applyMarkupChange(toggleUnderlineMarkup)
        break
    }
  }

  const { canvasRef } = useCertificateCanvas({
    canvasHostRef,
    canvasElRef,
    backgroundUrl: backgroundDataUrl,
    elements,
    selectedElementId,
    signatureUrl: null,
    onElementsChange: setElements,
    onElementSelect: setSelectedElementId,
  })

  const isTextSelected = selectedElement?.type === 'text'

  if (!isViewportAllowed) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background px-6 text-foreground">
        <div className="max-w-md rounded-2xl border bg-card p-6 text-center shadow-sm">
          <p className="text-sm font-semibold uppercase tracking-widest text-muted-foreground">
            Tela insuficiente
          </p>
          <h1 className="mt-3 text-2xl font-semibold">Abra no desktop</h1>
          <p className="mt-3 text-sm text-muted-foreground">
            O editor precisa de no mínimo {MIN_EDITOR_WIDTH}px de largura e {MIN_EDITOR_HEIGHT}px de altura para funcionar corretamente.
          </p>
          <p className="mt-2 text-xs text-muted-foreground">
            Amplie a janela ou use um monitor maior para acessar esta página.
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background text-foreground overflow-x-hidden">
      <div className="mx-auto max-w-7xl px-4 py-6 md:px-6 md:py-10">
        <div className="grid gap-6 lg:grid-cols-[280px_minmax(0,1fr)_320px]">
          {/* Sidebar Esquerda */}
          <aside className="space-y-4 lg:sticky lg:top-6 lg:self-start">
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="flex items-center gap-2 text-sm font-semibold">
                  <Type className="size-4 text-muted-foreground" />
                  Template
                </CardTitle>
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
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="w-full text-xs"
                  onClick={() => {
                    const exportData = {
                      title: templateTitle,
                      background: backgroundDataUrl,
                      elements: elements.map((e) => {
                        const base = {
                          xPct: e.xPct,
                          yPct: e.yPct,
                          widthPct: e.widthPct,
                          heightPct: e.heightPct,
                        }

                        if (e.type === 'text') {
                          return {
                            type: 'text',
                            ...base,
                            content: e.content,
                          }
                        }

                        if (e.type === 'signature') {
                          return { type: 'signature', ...base, title: e.title ?? '' }
                        }

                        return {
                          type: 'image',
                          ...base,
                          src: (e as ImageCanvasElement).src,
                          fileName: (e as ImageCanvasElement).fileName ?? null,
                        }
                      }),
                    }
                    console.log('=== EXPORT TEMPLATE ===')
                    console.log(JSON.stringify(exportData, null, 2))
                    console.log('=== FIM EXPORT ===')
                  }}
                >
                  Exportar template (console)
                </Button>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-3">
                <CardTitle className="flex items-center gap-2 text-sm font-semibold">
                  <ImageIcon className="size-4 text-muted-foreground" />
                  Fundo
                </CardTitle>
              </CardHeader>
              <CardContent>
                {backgroundDataUrl ? (
                  <div className="relative group">
                    <div className="relative rounded-lg overflow-hidden border bg-muted">
                      <img src={backgroundDataUrl} alt="Fundo" className="w-full h-28 object-cover" />
                      <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                        <button
                          type="button"
                          onClick={removeBackground}
                          className="p-1.5 rounded-full bg-destructive text-destructive-foreground hover:scale-110 transition-all"
                        >
                          <Trash2 className="size-3.5" />
                        </button>
                      </div>
                    </div>
                  </div>
                ) : (
                  <label className="relative flex flex-col items-center justify-center gap-2 rounded-lg border-2 border-dashed border-muted-foreground/25 p-4 min-h-24 cursor-pointer hover:border-muted-foreground/50 hover:bg-muted/30 transition-colors">
                    <div className="w-9 h-9 rounded-xl bg-muted flex items-center justify-center text-muted-foreground">
                      <ImageIcon className="size-4" />
                    </div>
                    <div className="text-center space-y-0.5">
                      <p className="text-xs font-medium text-foreground">Imagem de fundo</p>
                      <p className="text-[10px] text-muted-foreground">PNG, JPG até 5MB</p>
                    </div>
                    <input
                      type="file"
                      accept="image/*"
                      className="absolute inset-0 opacity-0 cursor-pointer"
                      onChange={async (e) => {
                        const file = e.target.files?.[0]
                        if (file) handleBackgroundFile(file)
                        e.target.value = ''
                      }}
                    />
                  </label>
                )}
              </CardContent>
            </Card>

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
                    <p className="px-1.5 text-[10px] font-medium uppercase tracking-wider text-muted-foreground">
                      Na tela
                    </p>
                    {elements.map((el) => (
                      <button
                        key={el.id}
                        type="button"
                        onClick={() => setSelectedElementId(el.id)}
                        className={cn(
                          'flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-xs transition-colors',
                          selectedElementId === el.id
                            ? 'bg-accent text-accent-foreground'
                            : 'hover:bg-muted'
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

            {elements.length === 0 && (
              <Card className="border-dashed">
                <CardContent className="p-6 text-center text-sm text-muted-foreground">
                  <p className="font-medium">Nenhum elemento</p>
                  <p className="mt-1 text-xs">Clique em "Texto", "Assin." ou "Imagem" para começar.</p>
                </CardContent>
              </Card>
            )}
          </aside>

          {/* Preview central */}
          <section className="space-y-4">
            <Card>
              <CardHeader className="border-b pb-3">
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <CardTitle className="flex items-center gap-2 text-sm font-semibold">
                      <ImagePlus className="size-4 text-muted-foreground" />
                      Preview
                    </CardTitle>
                    <CardDescription className="text-xs">
                      Clique no texto do preview para editar direto no próprio elemento ou use a barra para aplicar marcação.
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

              {/* Rich text toolbar */}
              <div className={cn(
                'flex items-center gap-1 border-b px-3 py-1.5 transition-opacity flex-wrap',
                isTextSelected ? 'bg-muted/20 border-border/50' : 'bg-muted/10 border-transparent'
              )}>
                <button
                  type="button"
                  onMouseDown={(e) => {
                    e.preventDefault()
                    e.stopPropagation()
                    toggleTextMarkup('bold')
                  }}
                  disabled={!isTextSelected}
                  className={cn(
                    'rounded p-1 transition-colors',
                    isTextSelected
                      ? 'cursor-pointer'
                      : 'text-muted-foreground/40 cursor-not-allowed',
                    isTextSelected ? 'text-foreground hover:bg-accent' : ''
                  )}
                  title="Negrito"
                >
                  <Bold className="size-3.5" />
                </button>
                <button
                  type="button"
                  onMouseDown={(e) => {
                    e.preventDefault()
                    e.stopPropagation()
                    toggleTextMarkup('italic')
                  }}
                  disabled={!isTextSelected}
                  className={cn(
                    'rounded p-1 transition-colors',
                    isTextSelected
                      ? 'cursor-pointer'
                      : 'text-muted-foreground/40 cursor-not-allowed',
                    isTextSelected ? 'text-foreground hover:bg-accent' : ''
                  )}
                  title="Itálico"
                >
                  <Italic className="size-3.5" />
                </button>
                <button
                  type="button"
                  onMouseDown={(e) => {
                    e.preventDefault()
                    e.stopPropagation()
                    toggleTextMarkup('underline')
                  }}
                  disabled={!isTextSelected}
                  className={cn(
                    'rounded p-1 transition-colors',
                    isTextSelected
                      ? 'cursor-pointer'
                      : 'text-muted-foreground/40 cursor-not-allowed',
                    isTextSelected ? 'text-foreground hover:bg-accent' : ''
                  )}
                  title="Sublinhado"
                >
                  <Underline className="size-3.5" />
                </button>
                <Separator orientation="vertical" className="mx-1 h-5" />

                <input
                  type="number"
                  min={10}
                  max={120}
                  disabled={!isTextSelected}
                  className={cn(
                    'h-7 w-16 rounded border border-input bg-transparent px-1.5 text-xs',
                    isTextSelected ? '' : 'text-muted-foreground/40'
                  )}
                  value={selectedElement && 'fontSize' in selectedElement ? (selectedElement as TextCanvasElement).fontSize : 32}
                  onChange={(e) => {
                    if (!selectedElement || selectedElement.type !== 'text') return
                    const v = Number(e.target.value)
                    if (isNaN(v)) return
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
                  className={cn(
                    'h-7 w-8 rounded border border-input bg-transparent p-0.5',
                    isTextSelected ? '' : 'opacity-40'
                  )}
                  value={selectedElement && 'color' in selectedElement ? (selectedElement as TextCanvasElement).color : '#111827'}
                  onChange={(e) => {
                    if (!selectedElement || selectedElement.type !== 'text') return
                    const newColor = e.target.value
                    updateElement(selectedElement.id, { color: newColor } as Partial<CanvasElement>)
                    applyMarkupChange((value, start, end) => toggleColorMarkup(value, start, end, newColor))
                  }}
                  title="Cor do texto"
                />

                <span className="text-[10px] text-muted-foreground ml-auto">
                  {isTextSelected ? 'Selecione um trecho no conteúdo para aplicar markup' : 'Selecione um elemento de texto'}
                </span>
              </div>

              <CardContent className="p-4">
                <div
                  ref={canvasHostRef}
                  className={cn(
                    'relative overflow-hidden rounded-xl border',
                    'aspect-297/210 w-full bg-background'
                  )}
                >
                  <canvas
                    ref={canvasElRef}
                    className="h-full w-full"
                    style={{ touchAction: 'none' }}
                  />
                  <div
                    className={cn(
                      'pointer-events-none absolute inset-0 transition-opacity duration-200',
                      showGrid ? 'opacity-30' : 'opacity-0'
                    )}
                    style={{
                      backgroundImage:
                        'linear-gradient(to right, rgba(0,0,0,0.12) 1px, transparent 1px), linear-gradient(to bottom, rgba(0,0,0,0.12) 1px, transparent 1px)',
                      backgroundSize: '50px 50px',
                    }}
                  />
                </div>
              </CardContent>
            </Card>
          </section>

          {/* Sidebar Direita */}
          <aside className="space-y-4 lg:sticky lg:top-6 lg:self-start">
            {selectedElement ? (
              <Card>
                <CardHeader className="pb-3">
                  <CardTitle className="flex items-center justify-between gap-2 text-sm font-semibold">
                    <span>{getDisplayName(selectedElement)}</span>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="size-6 text-muted-foreground hover:text-destructive"
                      onClick={deleteSelectedElement}
                    >
                      <Trash2 className="size-3.5" />
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="grid grid-cols-2 gap-2">
                    <div className="space-y-1">
                      <Label className="text-xs text-muted-foreground">X</Label>
                      <div className="flex items-center gap-1">
                        <Input type="number" min={-15} max={115} value={Math.round(selectedElement.xPct)}
                          onChange={(e) => { const v = Number(e.target.value); if (!isNaN(v)) updateElement(selectedElement.id, { xPct: Math.min(115, Math.max(-15, v)) }) }}
                          className="h-7 text-xs" />
                        <span className="text-xs text-muted-foreground shrink-0">%</span>
                      </div>
                    </div>
                    <div className="space-y-1">
                      <Label className="text-xs text-muted-foreground">Y</Label>
                      <div className="flex items-center gap-1">
                        <Input type="number" min={-15} max={115} value={Math.round(selectedElement.yPct)}
                          onChange={(e) => { const v = Number(e.target.value); if (!isNaN(v)) updateElement(selectedElement.id, { yPct: Math.min(115, Math.max(-15, v)) }) }}
                          className="h-7 text-xs" />
                        <span className="text-xs text-muted-foreground shrink-0">%</span>
                      </div>
                    </div>
                    <div className="space-y-1">
                      <Label className="text-xs text-muted-foreground">Larg.</Label>
                      <div className="flex items-center gap-1">
                        <Input type="number" min={10} max={150} value={Math.round(selectedElement.widthPct)}
                          onChange={(e) => { const v = Number(e.target.value); if (!isNaN(v)) updateElement(selectedElement.id, { widthPct: Math.min(150, Math.max(10, v)) }) }}
                          className="h-7 text-xs" />
                        <span className="text-xs text-muted-foreground shrink-0">%</span>
                      </div>
                    </div>
                    {selectedElement.type !== 'text' && (
                      <div className="space-y-1">
                        <Label className="text-xs text-muted-foreground">Alt.</Label>
                        <div className="flex items-center gap-1">
                          <Input type="number" min={6} max={80} value={Math.round(selectedElement.heightPct)}
                            onChange={(e) => { const v = Number(e.target.value); if (!isNaN(v)) updateElement(selectedElement.id, { heightPct: Math.min(80, Math.max(6, v)) }) }}
                            className="h-7 text-xs" />
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
                    <div className="space-y-1">
                      <Label className="text-xs text-muted-foreground">Título</Label>
                      <Input
                        value={(selectedElement as SignatureCanvasElement).title ?? ''}
                        onChange={(e) => updateElement(selectedElement.id, { title: e.target.value } as Partial<CanvasElement>)}
                        className="h-7 text-xs"
                        placeholder="Assinatura"
                      />
                    </div>
                  )}

                  {selectedElement.type === 'image' && (
                    <div className="space-y-1.5">
                      <Label className="text-xs text-muted-foreground">Imagem</Label>
                      <ImageUploadField
                        value={(selectedElement as ImageCanvasElement).src ?? undefined}
                        onChange={(v) => { if (!v) updateElement(selectedElement.id, { src: null } as Partial<CanvasElement>) }}
                        onFileSelect={async (file) => {
                          if (file) {
                            const dataUrl = await readFileAsDataUrl(file)
                            updateElement(selectedElement.id, { src: dataUrl } as Partial<CanvasElement>)
                          }
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
          </aside>
        </div>
      </div>
    </div>
  )
}
