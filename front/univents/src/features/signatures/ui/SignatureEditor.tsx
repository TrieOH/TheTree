import { useLayoutEffect, useRef, useState } from 'react'
import { PenLine, Eraser, Upload } from 'lucide-react'
import { Button } from '@/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/shadcn/card'
import { Input } from '@/shared/ui/shadcn/input'
import { Label } from '@/shared/ui/shadcn/label'
import { uploadFile } from '@/features/storage/api'
import { toast } from 'sonner'
import { useNavigate } from '@tanstack/react-router'
import { useCreateSignatureMutation } from '@/features/signatures/api/mutations'
import { SignatureImageSelector } from '@/features/signatures/ui/SignatureImageSelector'

type Mode = 'draw' | 'upload'

const SIGNATURE_CANVAS_WIDTH = 1200
const SIGNATURE_CANVAS_HEIGHT = 420

export interface SignatureEditorProps {
  eventId: string
  editionId: string
}

export function SignatureEditor({ eventId, editionId }: SignatureEditorProps) {
  const navigate = useNavigate()
  const [title, setTitle] = useState('Assinatura')
  const [mode, setMode] = useState<Mode>('draw')
  const [importedFile, setImportedFile] = useState<File | null>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const drawingRef = useRef(false)
  const lastPointRef = useRef<{ x: number; y: number } | null>(null)
  const canvasReadyRef = useRef(false)

  const syncCanvasSize = () => {
    const canvas = canvasRef.current
    if (!canvas || mode !== 'draw') return
    const ctx = canvas.getContext('2d')
    if (!ctx) return
    const nextWidth = SIGNATURE_CANVAS_WIDTH
    const nextHeight = SIGNATURE_CANVAS_HEIGHT
    if (canvas.width === nextWidth && canvas.height === nextHeight) return
    canvas.width = nextWidth
    canvas.height = nextHeight
    ctx.setTransform(1, 0, 0, 1, 0, 0)
    ctx.lineWidth = 3
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'
    ctx.strokeStyle = '#111827'
    ctx.clearRect(0, 0, nextWidth, nextHeight)
    canvasReadyRef.current = true
  }

  useLayoutEffect(() => {
    syncCanvasSize()
    const canvas = canvasRef.current
    if (!canvas || mode !== 'draw' || typeof ResizeObserver === 'undefined') return
    const observer = new ResizeObserver(() => {
      syncCanvasSize()
    })
    observer.observe(canvas)
    return () => observer.disconnect()
  }, [mode])

  const getPoint = (e: React.PointerEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current
    if (!canvas) return null
    const rect = canvas.getBoundingClientRect()
    return {
      x: (e.clientX - rect.left) * (canvas.width / rect.width),
      y: (e.clientY - rect.top) * (canvas.height / rect.height),
    }
  }

  const handlePointerDown = (e: React.PointerEvent<HTMLCanvasElement>) => {
    const point = getPoint(e)
    const canvas = canvasRef.current
    const ctx = canvas?.getContext('2d')
    if (!point || !ctx || !canvasReadyRef.current) return
    drawingRef.current = true
    lastPointRef.current = point
    ctx.beginPath()
    ctx.moveTo(point.x, point.y)
  }

  const handlePointerMove = (e: React.PointerEvent<HTMLCanvasElement>) => {
    if (!drawingRef.current) return
    const point = getPoint(e)
    const canvas = canvasRef.current
    const ctx = canvas?.getContext('2d')
    const last = lastPointRef.current
    if (!point || !ctx || !last) return
    ctx.lineTo(point.x, point.y)
    ctx.stroke()
    lastPointRef.current = point
  }

  const stopDrawing = () => {
    drawingRef.current = false
    lastPointRef.current = null
  }

  const clearCanvas = () => {
    const canvas = canvasRef.current
    const ctx = canvas?.getContext('2d')
    if (!canvas || !ctx) return
    ctx.clearRect(0, 0, canvas.width, canvas.height)
  }

  const saveMutation = useCreateSignatureMutation()

  return (
    <div className="mx-auto max-w-6xl px-4 py-6 pb-28!">
      <div className="mb-6 space-y-1">
        <p className="text-xs font-medium uppercase tracking-[0.24em] text-muted-foreground">Admin</p>
        <h1 className="text-2xl font-semibold">Nova assinatura</h1>
        <p className="max-w-2xl text-sm text-muted-foreground">
          Crie uma assinatura desenhando no canvas ou importando uma imagem pronta.
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,360px)_minmax(0,1fr)]">
        <Card className="h-fit">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-semibold">Configuração</CardTitle>
            <CardDescription className="text-xs">Nome e origem da assinatura.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-1.5">
              <Label className="text-xs text-muted-foreground">Título</Label>
              <Input
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Assinatura do responsável"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-xs text-muted-foreground">Modo</Label>
              <div className="grid grid-cols-2 gap-2">
                <Button
                  type="button"
                  variant={mode === 'draw' ? 'default' : 'outline'}
                  className="h-9 gap-2"
                  onClick={() => setMode('draw')}
                >
                  <PenLine className="size-4" />
                  Desenhar
                </Button>
                <Button
                  type="button"
                  variant={mode === 'upload' ? 'default' : 'outline'}
                  className="h-9 gap-2"
                  onClick={() => setMode('upload')}
                >
                  <Upload className="size-4" />
                  Importar
                </Button>
              </div>
            </div>

            <div className="rounded-2xl border bg-muted/20 p-3 text-xs text-muted-foreground">
              <p className="font-medium text-foreground">Dica</p>
              <p className="mt-1">
                Use uma imagem com fundo limpo para melhor legibilidade ou desenhe diretamente aqui.
              </p>
            </div>

            <Button
              type="button"
              className="h-9 w-full gap-2"
              onClick={async () => {
                try {
                  const trimmedTitle = title.trim()
                  if (!trimmedTitle) {
                    toast.error('Título é obrigatório')
                    return
                  }

                  let url: string | null = null
                  if (mode === 'draw') {
                    const canvas = canvasRef.current
                    if (!canvas) throw new Error('Canvas indisponível')
                    const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, 'image/png'))
                    if (!blob) throw new Error('Falha ao gerar imagem da assinatura')
                    const file = new File([blob], `${Date.now()}-signature.png`, { type: 'image/png' })
                    url = await uploadFile(file, `events/${eventId}/editions/${editionId}/signatures`)
                  } else if (importedFile) {
                    url = await uploadFile(importedFile, `events/${eventId}/editions/${editionId}/signatures`)
                  }

                  if (!url) {
                    toast.error('Selecione ou desenhe uma assinatura')
                    return
                  }

                  const res = await saveMutation.mutateAsync({
                    eventId,
                    editionId,
                    data: {
                      title: trimmedTitle,
                      url,
                    },
                  })

                  if (res.success) {
                    toast.success('Assinatura criada com sucesso')
                    void navigate({
                      to: '/admin/events/$eventId/editions/$editionId/signatures',
                      params: { eventId, editionId },
                    })
                    return
                  }

                  toast.error(res.message || 'Erro ao criar assinatura')
                } catch (error) {
                  toast.error(error instanceof Error ? error.message : 'Erro ao criar assinatura')
                }
              }}
              disabled={saveMutation.isPending}
            >
              Salvar assinatura
            </Button>
          </CardContent>
        </Card>

        <Card className="min-w-0">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-semibold">Prévia</CardTitle>
            <CardDescription className="text-xs">
              Veja o que vai ser salvo antes de concluir.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {mode === 'upload' ? (
              <SignatureImageSelector file={importedFile} onChange={setImportedFile} />
            ) : (
              <div className="space-y-3">
                <div className="flex items-center justify-between gap-3">
                  <div className="space-y-0.5">
                    <p className="text-sm font-medium text-foreground">Desenho</p>
                    <p className="text-xs text-muted-foreground">
                      Use o mouse ou toque para assinar no quadro abaixo.
                    </p>
                  </div>
                  <Button type="button" variant="outline" size="sm" onClick={clearCanvas}>
                    <Eraser className="size-4" />
                    Limpar
                  </Button>
                </div>
                <div className="rounded-2xl border bg-muted/10 p-2">
                  <canvas
                    ref={canvasRef}
                    className="h-44 w-full min-w-full touch-none rounded-xl bg-white"
                    style={{
                      aspectRatio: `${SIGNATURE_CANVAS_WIDTH} / ${SIGNATURE_CANVAS_HEIGHT}`,
                    }}
                    onPointerDown={handlePointerDown}
                    onPointerMove={handlePointerMove}
                    onPointerUp={stopDrawing}
                    onPointerLeave={stopDrawing}
                    onPointerCancel={stopDrawing}
                  />
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
