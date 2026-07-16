import type { CertificationTemplateElement } from '../../../model'
import { HashElementView } from './hash-element-view'
import { ImageElementView } from './image-element-view'
import { SignatureElementView } from './signature-element-view'
import { TextElementView } from './text-element-view'

interface CertificateElementViewProps {
  element: CertificationTemplateElement
}

export function CertificateElementView({
  element,
}: CertificateElementViewProps) {
  switch (element.type) {
    case 'hash':
      return <HashElementView element={element} />
    case 'text':
      return <TextElementView element={element} />
    case 'image':
      return <ImageElementView element={element} />
    case 'signature':
      return <SignatureElementView element={element} />
  }
}
