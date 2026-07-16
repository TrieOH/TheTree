import { ExternalLink } from 'lucide-react'
import type { HashCertificateElement } from '../../types'

interface HashElementViewProps {
  element: HashCertificateElement
}

export function HashElementView({ element }: HashElementViewProps) {
  const hashLabel = element.hashLabel.trim() || 'Código de verificação'
  const hash = element.hash.trim() || '{{cert_hash}}'
  const url = element.url.trim() || '{{verify_url}}'
  const linkLabel = element.linkLabel.trim() || url

  return (
    <div
      className="flex h-full w-full flex-col justify-center gap-1 overflow-hidden px-1"
      style={{
        textAlign: element.align,
        color: element.color,
        fontSize: element.fontSize,
      }}
    >
      <div className="truncate">
        <span className="opacity-80">{hashLabel}: </span>
        <span className="font-semibold tracking-wide">{hash}</span>
      </div>
      <a
        href={url}
        onClick={(event) => event.preventDefault()}
        title={`No certificado exportado, este link abre: ${url}`}
        className="inline-flex items-center gap-1 truncate font-medium underline decoration-current underline-offset-2"
        style={{
          justifyContent:
            element.align === 'center'
              ? 'center'
              : element.align === 'right'
                ? 'flex-end'
                : 'flex-start',
        }}
      >
        {linkLabel}
        <ExternalLink className="size-3 shrink-0" />
      </a>
    </div>
  )
}
