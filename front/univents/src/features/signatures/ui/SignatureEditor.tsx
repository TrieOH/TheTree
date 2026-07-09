import { useLayoutEffect, useRef, useState } from 'react'
import { PenLine, Upload, Eraser } from 'lucide-react'
import { Button } from '@/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/shadcn/card'
import { Input } from '@/shared/ui/shadcn/input'
import { Label } from '@/shared/ui/shadcn/label'
import ImageUploadField from '@/widgets/form/ui/image-upload-field'
import { uploadAndModerateFile } from '@/features/storage/api'
import { createSignatureFn } from '@/features/signatures/api'
import { useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { useNavigate } from '@tanstack/react-router'

type Mode = 'draw' | 'upload'

export interface SignatureEditorProps {
  eventId: string
  editionId: string
}

export function SignatureEditor({ eventId, editionId }: SignatureEditorProps) {
  const navigate = useNavigate()
  const [title, setTitle] = useState('Assinatura')
  const [mode, setMode] = useState<Mode>('draw')
  const [imagePreview, setImagePreview] = useState<string | null>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const drawingRef = useRef(false)
  const lastPointRef = useRef<{ x: number; y: number } | null>(null)

  const syncCanvasSize = () => {
    const canvas = canvasRef.current
    if (!canvas || mode !== 'draw') return
    const ctx = canvas.getContext('2d')
    if (!ctx) return
    const dpr = window.devicePixelRatio || 1
    const width = canvas.clientWidth
    const height = canvas.clientHeight
    const nextWidth = Math.floor(width * dpr)
    const nextHeight = Math.floor(height * dpr)
    if (canvas.width === nextWidth && canvas.height === nextHeight) return
    const snapshot = ctx.getImageData(0, 0, canvas.width, canvas.height)
    canvas.width = nextWidth
    canvas.height = nextHeight
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
    ctx.lineWidth = 3
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'
    ctx.strokeStyle = '#111827'
    ctx.clearRect(0, 0, width, height)
    if (snapshot.width > 0 && snapshot.height > 0) {
      ctx.putImageData(snapshot, 0, 0)
    }
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
    return { x: e.clientX - rect.left, y: e.clientY - rect.top }
  }

  const handlePointerDown = (e: React.PointerEvent<HTMLCanvasElement>) => {
    const point = getPoint(e)
    const canvas = canvasRef.current
    const ctx = canvas?.getContext('2d')
    if (!point || !ctx) return
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

  const saveMutation = useMutation({
    mutationFn: async () => {
      let url = imagePreview
      if (mode === 'draw') {
        const canvas = canvasRef.current
        if (!canvas) throw new Error('Canvas indisponível')
        const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, 'image/png'))
        if (!blob) throw new Error('Falha ao gerar imagem da assinatura')
        const file = new File([blob], `${Date.now()}-signature.png`, { type: 'image/png' })
        url = await uploadAndModerateFile(file, `events/${eventId}/editions/${editionId}/signatures`)
      }
      if (!url) throw new Error('Selecione ou desenhe uma assinatura')
      return createSignatureFn(eventId, editionId, {
        title,
        url,
        pos_x: 0,
        pos_y: 0,
      })
    },
    onSuccess: (res) => {
      if (res.success) {
        toast.success('Assinatura criada com sucesso')
        void navigate({
          to: '/admin/events/$eventId/editions/$editionId/signatures',
          params: { eventId, editionId },
        })
      } else {
        toast.error(res.message || 'Erro ao criar assinatura')
      }
    },
    onError: () => toast.error('Erro ao conectar com o servidor'),
  })

  return (
    <div className="mx-auto max-w-5xl px-4 py-6 md:px-6 md:py-10">
      <div className="mb-6 space-y-1">
        <p className="text-xs font-medium uppercase tracking-[0.24em] text-muted-foreground">Admin</p>
        <h1 className="text-2xl font-semibold">Nova assinatura</h1>
        <p className="text-sm text-muted-foreground">
          Desenhe no canvas ou importe uma imagem. O resultado final vira a assinatura salva.
        </p>
      </div>

      <div className="w-full gap-6">
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-semibold">Dados</CardTitle>
            <CardDescription className="text-xs">Título e origem da assinatura.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-1">
              <Label className="text-xs text-muted-foreground">Título</Label>
              <Input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Assinatura do responsável" />
            </div>

            <div className="space-y-2">
              <Label className="text-xs text-muted-foreground">Modo</Label>
              <div className="grid grid-cols-2 gap-2">
                <Button type="button" variant={mode === 'draw' ? 'default' : 'outline'} onClick={() => setMode('draw')}>
                  <PenLine className="size-4" />
                  Desenhar
                </Button>
                <Button type="button" variant={mode === 'upload' ? 'default' : 'outline'} onClick={() => setMode('upload')}>
                  <Upload className="size-4" />
                  Importar
                </Button>
              </div>
            </div>

            {mode === 'upload' ? (
              <ImageUploadField
                value={imagePreview ?? undefined}
                onChange={(url) => setImagePreview(url || null)}
                onFileSelect={async (file) => {
                  if (!file) {
                    setImagePreview(null)
                    return
                  }
                  const preview = URL.createObjectURL(file)
                  setImagePreview(preview)
                }}
                accept="image/png,image/jpeg,image/webp"
                placeholder="Selecionar assinatura"
              />
            ) : (
              <div className="space-y-2">
                <div className="flex gap-2">
                  <Button type="button" variant="outline" size="sm" onClick={clearCanvas}>
                    <Eraser className="size-4" />
                    Limpar
                  </Button>
                </div>
                <div className="w-full rounded-2xl border bg-muted/10 p-2">
                  <canvas
                    ref={canvasRef}
                    className="h-56 w-full min-w-full touch-none rounded-xl bg-white"
                    onPointerDown={handlePointerDown}
                    onPointerMove={handlePointerMove}
                    onPointerUp={stopDrawing}
                    onPointerLeave={stopDrawing}
                    onPointerCancel={stopDrawing}
                  />
                </div>
              </div>
            )}

            <Button
              type="button"
              className="w-full"
              onClick={() => { void saveMutation.mutateAsync() }}
              disabled={saveMutation.isPending}
            >
              Salvar assinatura
            </Button>

            {mode === 'upload' && imagePreview && (
              <div className="rounded-2xl border bg-muted/10 p-3">
                <p className="mb-2 text-xs font-medium uppercase tracking-[0.2em] text-muted-foreground">Preview</p>
                <img src={imagePreview} alt={title} className="max-h-56 w-full object-contain" />
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
