import { DEFAULT_CERTIFICATE_FONT } from '../constants'
import type { TextCertificateElement } from '../types'

type RichParagraph = TextCertificateElement['paragraphs'][number]
type RichRun = RichParagraph['runs'][number]

function escapeHtml(text: string): string {
  return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

function runStyle(run: RichRun): string {
  return [
    `color:${run.color}`,
    `font-size:${run.fontSize}px`,
    `font-family:${run.fontFamily}`,
    `font-weight:${run.bold ? 700 : 400}`,
    `font-style:${run.italic ? 'italic' : 'normal'}`,
    `text-decoration:${run.underline ? 'underline' : 'none'}`,
  ].join(';')
}

export function paragraphsToHtml(
  paragraphs: TextCertificateElement['paragraphs'],
): string {
  return paragraphs
    .map((paragraph, paragraphIndex) => {
      const nearestRun =
        paragraph.runs[0] ??
        paragraphs
          .slice(0, paragraphIndex)
          .reverse()
          .find((candidate) => candidate.runs.length > 0)
          ?.runs.at(-1) ??
        paragraphs
          .slice(paragraphIndex + 1)
          .find((candidate) => candidate.runs.length > 0)?.runs[0]
      const content =
        paragraph.runs
          .map(
            (run) =>
              `<span style="${runStyle(run)}">${escapeHtml(run.text) || '\u200b'}</span>`,
          )
          .join('') ||
        (nearestRun
          ? `<span style="${runStyle(nearestRun)}">\u200b</span>`
          : '\u200b')
      return `<p style="text-align:${paragraph.align};line-height:${paragraph.lineHeight};margin:0;min-height:1.25em">${content}</p>`
    })
    .join('')
}

const convertedColorCache = new Map<string, string>()

function channelsToHex(channels: number[]): string {
  return `#${channels
    .map((channel) =>
      Math.max(0, Math.min(255, Math.round(channel)))
        .toString(16)
        .padStart(2, '0'),
    )
    .join('')}`
}

export function computedColorToHex(color: string): string {
  if (color.startsWith('#')) return color
  const cached = convertedColorCache.get(color)
  if (cached) return cached

  const rgb = color.match(
    /^rgba?\(\s*([\d.]+)(?:\s*,\s*|\s+)([\d.]+)(?:\s*,\s*|\s+)([\d.]+)/i,
  )
  if (rgb) {
    const converted = channelsToHex([
      Number(rgb[1]),
      Number(rgb[2]),
      Number(rgb[3]),
    ])
    convertedColorCache.set(color, converted)
    return converted
  }

  const canvas = document.createElement('canvas')
  canvas.width = 1
  canvas.height = 1
  const context = canvas.getContext('2d')
  if (!context) return '#111827'
  context.clearRect(0, 0, 1, 1)
  context.fillStyle = color
  context.fillRect(0, 0, 1, 1)
  const pixel = context.getImageData(0, 0, 1, 1).data
  const converted = channelsToHex([pixel[0], pixel[1], pixel[2]])
  convertedColorCache.set(color, converted)
  return converted
}

function readRunStyle(element: Element): Omit<RichRun, 'text'> {
  const style = window.getComputedStyle(element)
  const numericWeight = Number(style.fontWeight)
  return {
    bold:
      style.fontWeight === 'bold' ||
      (!Number.isNaN(numericWeight) && numericWeight >= 600),
    italic: style.fontStyle === 'italic',
    underline: style.textDecorationLine.includes('underline'),
    color: computedColorToHex(style.color),
    fontSize: Math.round(Number.parseFloat(style.fontSize) || 16),
    fontFamily: style.fontFamily || DEFAULT_CERTIFICATE_FONT,
  }
}

function collectRuns(root: Node): RichRun[] {
  const runs: RichRun[] = []

  function visit(node: Node) {
    if (node.nodeType === Node.TEXT_NODE) {
      const text = node.textContent?.replace(/\u200b/g, '') ?? ''
      if (!text) return
      const style = node.parentElement
        ? readRunStyle(node.parentElement)
        : {
            bold: false,
            italic: false,
            underline: false,
            color: '#111827',
            fontSize: 16,
            fontFamily: DEFAULT_CERTIFICATE_FONT,
          }
      const previous = runs.at(-1)
      if (
        previous &&
        previous.bold === style.bold &&
        previous.italic === style.italic &&
        previous.underline === style.underline &&
        previous.color === style.color &&
        previous.fontSize === style.fontSize &&
        previous.fontFamily === style.fontFamily
      ) {
        previous.text += text
      } else {
        runs.push({ ...style, text })
      }
      return
    }

    node.childNodes.forEach(visit)
  }

  visit(root)
  return runs
}

function paragraphAlign(element: Element): RichParagraph['align'] {
  const align = window.getComputedStyle(element).textAlign
  return align === 'center' || align === 'right' || align === 'justify'
    ? align
    : 'left'
}

export function domToParagraphs(
  container: HTMLElement,
): TextCertificateElement['paragraphs'] {
  const blocks = Array.from(container.children).filter(
    (child) => child.tagName === 'DIV' || child.tagName === 'P',
  )
  if (blocks.length === 0) {
    return [{ align: 'left', lineHeight: 1.25, runs: collectRuns(container) }]
  }
  return blocks.map((block) => ({
    align: paragraphAlign(block),
    lineHeight:
      Number.parseFloat(window.getComputedStyle(block).lineHeight) /
        Number.parseFloat(window.getComputedStyle(block).fontSize) || 1.25,
    runs: collectRuns(block),
  }))
}
