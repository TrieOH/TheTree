import { Eye, FileText } from 'lucide-react'
import type { CertificationTemplateI } from '../model'
import { DEFAULT_CERTIFICATE_CANVAS } from '../editor/constants'
import { useElementSize } from '../editor/hooks/use-element-size'
import { CertificateElementView } from '../editor/ui/elements/certificate-element-view'
import { Button } from '@/shared/ui/shadcn/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/shared/ui/shadcn/dialog'

interface CertViewerProps {
  template: CertificationTemplateI
  triggerLabel?: string
}

export function CertViewer({
  template,
  triggerLabel = 'Visualizar',
}: CertViewerProps) {
  return (
    <Dialog>
      <DialogTrigger
        render={<Button type="button" variant="outline" size="sm" />}
      >
        <Eye className="size-4" />
        {triggerLabel}
      </DialogTrigger>
      <DialogContent
        className="!z-[100] grid-rows-[auto_minmax(0,1fr)] gap-0 overflow-hidden p-0 shadow-2xl"
        overlayClassName="!z-[99] bg-black/40 backdrop-blur-md"
        style={{
          width: 'min(72rem, calc(100vw - 2rem))',
          maxWidth: 'calc(100vw - 2rem)',
          height: 'min(48rem, calc(100dvh - 2rem))',
        }}
      >
        <DialogHeader className="border-b bg-background px-5 py-3 pr-14">
          <DialogTitle className="flex items-center gap-2">
            <FileText className="size-4 text-muted-foreground" />
            {template.title}
          </DialogTitle>
          <DialogDescription>
            Pré-visualização do certificado no tamanho e proporção de emissão.
          </DialogDescription>
        </DialogHeader>
        <div className="min-h-0 overflow-hidden bg-muted/60 p-3 sm:p-5">
          <CertificateTemplateStaticView template={template} />
        </div>
      </DialogContent>
    </Dialog>
  )
}

export function CertificateTemplateStaticView({
  template,
}: {
  template: CertificationTemplateI
}) {
  const { ref, size } = useElementSize<HTMLDivElement>()
  const canvas = DEFAULT_CERTIFICATE_CANVAS
  const scale = Math.max(
    0,
    Math.min(size.width / canvas.width, size.height / canvas.height),
  )
  const backgroundUrl = template.url ?? template.data.background

  return (
    <div
      ref={ref}
      className="flex h-full w-full items-center justify-center overflow-hidden"
    >
      {scale > 0 ? (
        <div
          className="relative shrink-0 overflow-hidden bg-white shadow-2xl ring-1 ring-black/10"
          style={{
            width: canvas.width * scale,
            height: canvas.height * scale,
          }}
        >
          <div
            className="relative origin-top-left overflow-hidden bg-white"
            style={{
              width: canvas.width,
              height: canvas.height,
              transform: `scale(${scale})`,
              backgroundImage: backgroundUrl
                ? `url(${JSON.stringify(backgroundUrl)})`
                : undefined,
              backgroundPosition: 'center',
              backgroundRepeat: 'no-repeat',
              backgroundSize: 'cover',
            }}
          >
            {template.data.elements.map((element) => (
              <div
                key={element.id}
                className="absolute overflow-hidden"
                style={{
                  left: element.x,
                  top: element.y,
                  width: element.width,
                  height: element.height,
                }}
              >
                <CertificateElementView element={element} />
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  )
}
