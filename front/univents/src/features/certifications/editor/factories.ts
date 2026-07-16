import { DEFAULT_CERTIFICATE_FONT } from './constants'
import type {
  CertificateCanvasSize,
  HashCertificateElement,
  ImageCertificateElement,
  SignatureCertificateElement,
  TextCertificateElement,
} from './types'
import { createCertificateElementId } from './utils'

export function createHashElement(canvas: CertificateCanvasSize): HashCertificateElement {
  const width = Math.min(420, canvas.width * 0.42)

  return {
    id: createCertificateElementId('hash'),
    type: 'hash',
    x: (canvas.width - width) / 2,
    y: canvas.height - 92,
    width,
    height: 60,
    hashLabel: 'Código de verificação',
    hash: '{{cert_hash}}',
    linkLabel: 'Verificar autenticidade',
    url: '{{verify_url}}',
    fontSize: 13,
    color: '#4B5563',
    align: 'center',
  }
}

export function createTextElement(
  canvas: CertificateCanvasSize,
  bounds: Partial<Pick<TextCertificateElement, 'x' | 'y' | 'width' | 'height'>> = {},
): TextCertificateElement {
  const width = bounds.width ?? Math.min(560, canvas.width * 0.6)
  const height = bounds.height ?? 120

  return {
    id: createCertificateElementId('text'),
    type: 'text',
    x: bounds.x ?? (canvas.width - width) / 2,
    y: bounds.y ?? (canvas.height - height) / 2,
    width,
    height,
    paragraphs: [
      {
        align: 'center',
        runs: [
          {
            text: 'Clique duas vezes para editar este texto',
            bold: false,
            italic: false,
            underline: false,
            color: '#111827',
            fontSize: 24,
            fontFamily: DEFAULT_CERTIFICATE_FONT,
          },
        ],
      },
    ],
  }
}

export function createImageElement(
  src: string,
  canvas: CertificateCanvasSize,
  naturalSize?: CertificateCanvasSize,
): ImageCertificateElement {
  const maxWidth = canvas.width * 0.4
  const maxHeight = canvas.height * 0.4
  let width = maxWidth
  let height = maxHeight

  if (naturalSize && naturalSize.width > 0 && naturalSize.height > 0) {
    const ratio = naturalSize.width / naturalSize.height
    if (maxWidth / ratio <= maxHeight) height = maxWidth / ratio
    else width = maxHeight * ratio
  }

  return {
    id: createCertificateElementId('image'),
    type: 'image',
    x: (canvas.width - width) / 2,
    y: (canvas.height - height) / 2,
    width,
    height,
    src,
    fit: 'contain',
    radius: 0,
    opacity: 1,
  }
}

export function createSignatureElement(
  signature: { id: string; url: string; name: string },
  canvas: CertificateCanvasSize,
): SignatureCertificateElement {
  const width = Math.min(260, canvas.width * 0.28)
  const height = Math.min(110, canvas.height * 0.16)

  return {
    id: createCertificateElementId('signature'),
    type: 'signature',
    signatureId: signature.id,
    src: signature.url,
    name: signature.name,
    fit: 'contain',
    radius: 0,
    opacity: 1,
    x: (canvas.width - width) / 2,
    y: (canvas.height - height) / 2,
    width,
    height,
  }
}
