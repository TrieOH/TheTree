export interface RichTextStyle {
  fontWeight?: number | string
  fontStyle?: 'normal' | 'italic'
  underline?: boolean
  fill?: string
  fontSize?: number
  fontFamily?: string
}

export type RichTextStyleMap = Record<string, Record<string, RichTextStyle>>

export interface ParsedRichText {
  plainText: string
  styles: RichTextStyleMap
  plainToMarkupIndex: number[]
  sourceLength: number
}

export interface MarkupChangeResult {
  value: string
  selectionStart: number
  selectionEnd: number
}

function buildStyle(active: {
  bold: boolean
  italic: boolean
  underline: boolean
  color: string | null
  size: number | null
  fontFamily: string | null
}): RichTextStyle | null {
  const style: RichTextStyle = {}
  if (active.bold) style.fontWeight = 700
  if (active.italic) style.fontStyle = 'italic'
  if (active.underline) style.underline = true
  if (active.color) style.fill = active.color
  if (active.size) style.fontSize = active.size
  if (active.fontFamily) style.fontFamily = active.fontFamily
  return Object.keys(style).length > 0 ? style : null
}

function applyCharStyle(
  styles: RichTextStyleMap,
  lineIndex: number,
  charIndex: number,
  style: RichTextStyle | null
) {
  if (!style) return
  const lineKey = String(lineIndex)
  const charKey = String(charIndex)
  styles[lineKey] ??= {}
  styles[lineKey][charKey] = style
}

export function parseRichTextMarkup(
  content: string,
  baseStyle: {
    fontFamily?: string
    fontSize?: number
    color?: string
    fontWeight?: number | string
  } = {}
): ParsedRichText {
  const styles: RichTextStyleMap = {}
  let plainText = ''
  const plainToMarkupIndex: number[] = [0]
  let lineIndex = 0
  let charIndex = 0

  const colorStack: Array<string | null> = [baseStyle.color ?? null]
  const sizeStack: Array<number | null> = [baseStyle.fontSize ?? null]
  const active = {
    bold: false,
    italic: false,
    underline: false,
    color: baseStyle.color ?? null,
    size: baseStyle.fontSize ?? null,
    fontFamily: baseStyle.fontFamily ?? null,
  }

  const pushChar = (char: string, sourceIndex: number) => {
    plainText += char
    plainToMarkupIndex[plainText.length] = sourceIndex + 1
    if (char === '\n') {
      lineIndex += 1
      charIndex = 0
      return
    }

    applyCharStyle(styles, lineIndex, charIndex, buildStyle(active))
    charIndex += 1
  }

  for (let i = 0; i < content.length; i += 1) {
    const char = content[i]

    if (char === '\\' && i + 1 < content.length) {
      pushChar(content[i + 1], i + 1)
      i += 1
      continue
    }

    if (content.startsWith('**', i)) {
      active.bold = !active.bold
      i += 1
      continue
    }

    if (content.startsWith('__', i)) {
      active.underline = !active.underline
      i += 1
      continue
    }

    if (char === '*') {
      active.italic = !active.italic
      continue
    }

    if (content.startsWith('[color=', i)) {
      const end = content.indexOf(']', i)
      if (end !== -1) {
        const value = content.slice(i + 7, end).trim()
        colorStack.push(value || null)
        active.color = colorStack[colorStack.length - 1] ?? null
        i = end
        continue
      }
    }

    if (content.startsWith('[/color]', i)) {
      if (colorStack.length > 1) colorStack.pop()
      active.color = colorStack[colorStack.length - 1] ?? null
      i += '[/color]'.length - 1
      continue
    }

    if (content.startsWith('[size=', i)) {
      const end = content.indexOf(']', i)
      if (end !== -1) {
        const rawSize = content.slice(i + 6, end).trim()
        const parsedSize = Number(rawSize)
        sizeStack.push(Number.isFinite(parsedSize) ? parsedSize : null)
        active.size = sizeStack[sizeStack.length - 1] ?? null
        i = end
        continue
      }
    }

    if (content.startsWith('[/size]', i)) {
      if (sizeStack.length > 1) sizeStack.pop()
      active.size = sizeStack[sizeStack.length - 1] ?? null
      i += '[/size]'.length - 1
      continue
    }

    pushChar(char, i)
  }

  return {
    plainText,
    styles,
    plainToMarkupIndex,
    sourceLength: content.length,
  }
}

export function plainIndexToMarkupIndex(parsed: ParsedRichText, plainIndex: number): number {
  if (plainIndex <= 0) return 0
  if (plainIndex >= parsed.plainToMarkupIndex.length) return parsed.sourceLength
  return parsed.plainToMarkupIndex[plainIndex]
}

function toggleDelimitedMarkup(
  value: string,
  selectionStart: number,
  selectionEnd: number,
  open: string,
  close = open
): MarkupChangeResult {
  const start = Math.min(selectionStart, selectionEnd)
  const end = Math.max(selectionStart, selectionEnd)
  const hasSelection = start !== end
  const before = value.slice(0, start)
  const selected = value.slice(start, end)
  const after = value.slice(end)

  if (!hasSelection) {
    const nextValue = `${before}${open}${close}${after}`
    const cursor = start + open.length
    return { value: nextValue, selectionStart: cursor, selectionEnd: cursor }
  }

  const startPrefix = value.slice(Math.max(0, start - open.length), start)
  const endSuffix = value.slice(end, end + close.length)

  if (startPrefix === open && endSuffix === close) {
    const nextValue =
      value.slice(0, start - open.length) +
      selected +
      value.slice(end + close.length)
    return {
      value: nextValue,
      selectionStart: start - open.length,
      selectionEnd: end - open.length,
    }
  }

  const nextValue = `${before}${open}${selected}${close}${after}`
  return {
    value: nextValue,
    selectionStart: start + open.length,
    selectionEnd: end + open.length,
  }
}

export function toggleBoldMarkup(value: string, selectionStart: number, selectionEnd: number): MarkupChangeResult {
  return toggleDelimitedMarkup(value, selectionStart, selectionEnd, '**')
}

export function toggleItalicMarkup(value: string, selectionStart: number, selectionEnd: number): MarkupChangeResult {
  return toggleDelimitedMarkup(value, selectionStart, selectionEnd, '*')
}

export function toggleUnderlineMarkup(value: string, selectionStart: number, selectionEnd: number): MarkupChangeResult {
  return toggleDelimitedMarkup(value, selectionStart, selectionEnd, '__')
}

export function toggleColorMarkup(
  value: string,
  selectionStart: number,
  selectionEnd: number,
  color: string
): MarkupChangeResult {
  return toggleDelimitedMarkup(value, selectionStart, selectionEnd, `[color=${color}]`, '[/color]')
}

export function toggleSizeMarkup(
  value: string,
  selectionStart: number,
  selectionEnd: number,
  size: number
): MarkupChangeResult {
  return toggleDelimitedMarkup(value, selectionStart, selectionEnd, `[size=${size}]`, '[/size]')
}
