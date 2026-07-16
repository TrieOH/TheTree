import { toBlob } from 'html-to-image'
import { jsPDF } from 'jspdf'
import { DEFAULT_CERTIFICATE_CANVAS } from '../editor/constants'

export type CertificateExportFormat = 'png' | 'pdf'

function safeFilename(filename: string) {
  const normalized = filename
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-zA-Z0-9._-]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .toLowerCase()
  return normalized || 'certificado'
}

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  window.setTimeout(() => URL.revokeObjectURL(url), 0)
}

export async function createCertificatePng(canvas: HTMLElement): Promise<Blob> {
  const { width, height } = DEFAULT_CERTIFICATE_CANVAS
  const blob = await toBlob(canvas, {
    width,
    height,
    pixelRatio: 2,
    cacheBust: true,
    backgroundColor: '#ffffff',
    style: {
      transform: 'none',
      transformOrigin: 'top left',
    },
  })
  if (!blob) throw new Error('Não foi possível gerar a imagem do certificado.')
  return blob
}

export async function createCertificatePdf(canvas: HTMLElement): Promise<Blob> {
  const { width, height } = DEFAULT_CERTIFICATE_CANVAS
  const links = readCanvasLinks(canvas, width, height)
  const png = await createCertificatePng(canvas)
  const imageUrl = await blobToDataUrl(png)
  const pdf = new jsPDF({
    orientation: width >= height ? 'landscape' : 'portrait',
    unit: 'px',
    format: [width, height],
    hotfixes: ['px_scaling'],
  })
  pdf.addImage(imageUrl, 'PNG', 0, 0, width, height, undefined, 'FAST')
  links.forEach((link) => {
    pdf.link(link.x, link.y, link.width, link.height, { url: link.url })
  })
  return pdf.output('blob')
}

export async function downloadCertificate(
  canvas: HTMLElement,
  filename: string,
  format: CertificateExportFormat,
) {
  const blob =
    format === 'png'
      ? await createCertificatePng(canvas)
      : await createCertificatePdf(canvas)
  downloadBlob(blob, `${safeFilename(filename)}.${format}`)
}

function blobToDataUrl(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result))
    reader.onerror = () => reject(reader.error)
    reader.readAsDataURL(blob)
  })
}

function readCanvasLinks(
  canvas: HTMLElement,
  logicalWidth: number,
  logicalHeight: number,
) {
  const canvasRect = canvas.getBoundingClientRect()
  if (canvasRect.width <= 0 || canvasRect.height <= 0) return []

  return Array.from(canvas.querySelectorAll<HTMLAnchorElement>('a[href]'))
    .map((anchor) => {
      const rect = anchor.getBoundingClientRect()
      return {
        url: anchor.href,
        x: ((rect.left - canvasRect.left) / canvasRect.width) * logicalWidth,
        y: ((rect.top - canvasRect.top) / canvasRect.height) * logicalHeight,
        width: (rect.width / canvasRect.width) * logicalWidth,
        height: (rect.height / canvasRect.height) * logicalHeight,
      }
    })
    .filter(
      (link) =>
        /^https?:\/\//i.test(link.url) && link.width > 0 && link.height > 0,
    )
}
