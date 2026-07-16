import { useEffect, useLayoutEffect, useRef } from 'react'
import type { TextCertificateElement } from '../../types'
import {
  computedColorToHex,
  domToParagraphs,
  paragraphsToHtml,
} from '../../rich-text/dom-serializer'
import { certificateEditorActions } from '../../store'
import type { CertificateTextSelectionStyles } from '../../store'
import {
  CERTIFICATE_FONT_FAMILIES,
  DEFAULT_CERTIFICATE_FONT,
} from '../../constants'

interface TextElementViewProps {
  element: TextCertificateElement
  editing?: boolean
}

export function TextElementView({
  element,
  editing = false,
}: TextElementViewProps) {
  return editing ? (
    <EditableTextElement element={element} />
  ) : (
    <StaticTextElement element={element} />
  )
}

function EditableTextElement({ element }: { element: TextCertificateElement }) {
  const editorRef = useRef<HTMLDivElement>(null)
  const savedRange = useRef<Range | null>(null)

  function commit() {
    const editor = editorRef.current
    if (!editor) return
    const paragraphs = domToParagraphs(editor)
    certificateEditorActions.updateElement(element.id, (current) =>
      current.type === 'text' ? { ...current, paragraphs } : current,
    )
  }

  function restoreSelection() {
    const editor = editorRef.current
    if (!editor) return
    editor.focus()
    if (!savedRange.current) return
    const selection = window.getSelection()
    selection?.removeAllRanges()
    selection?.addRange(savedRange.current)
  }

  function updateSelectionStyles() {
    const editor = editorRef.current
    const selection = window.getSelection()
    if (!editor || !selection?.rangeCount) return
    const range = selection.getRangeAt(0)
    if (!editor.contains(range.commonAncestorContainer)) return

    const textNodes: Text[] = []
    const walker = document.createTreeWalker(editor, NodeFilter.SHOW_TEXT)
    let current = walker.nextNode()
    while (current) {
      if (range.collapsed) {
        if (current === selection.anchorNode) textNodes.push(current as Text)
      } else if (range.intersectsNode(current)) {
        textNodes.push(current as Text)
      }
      current = walker.nextNode()
    }

    if (textNodes.length === 0 && selection.anchorNode) {
      if (selection.anchorNode.nodeType === Node.TEXT_NODE) {
        textNodes.push(selection.anchorNode as Text)
      } else {
        const children = selection.anchorNode.childNodes
        if (children.length > 0) {
          let fallback = children[Math.max(0, selection.anchorOffset - 1)]
          while (fallback.lastChild) fallback = fallback.lastChild
          if (fallback.nodeType === Node.TEXT_NODE)
            textNodes.push(fallback as Text)
        }
      }
    }

    const styles = textNodes.map((textNode) => {
      const selectedElement = textNode.parentElement
      const computed = selectedElement
        ? window.getComputedStyle(selectedElement)
        : null
      const paragraph = selectedElement?.closest('p, div')
      const computedAlign = paragraph
        ? window.getComputedStyle(paragraph).textAlign
        : 'left'
      const align: CertificateTextSelectionStyles['align'] =
        computedAlign === 'center' ||
        computedAlign === 'right' ||
        computedAlign === 'justify'
          ? computedAlign
          : ('left' as const)
      const computedFamily = computed?.fontFamily ?? DEFAULT_CERTIFICATE_FONT
      const normalizedFamily = computedFamily.replace(/["']/g, '').toLowerCase()
      const fontFamily =
        CERTIFICATE_FONT_FAMILIES.find((font) => {
          const primaryFamily = font.value
            .split(',')[0]
            ?.replace(/["']/g, '')
            .trim()
            .toLowerCase()
          return primaryFamily
            ? normalizedFamily.startsWith(primaryFamily)
            : false
        })?.value ?? DEFAULT_CERTIFICATE_FONT
      const numericWeight = Number(computed?.fontWeight)

      return {
        bold:
          computed?.fontWeight === 'bold' ||
          (!Number.isNaN(numericWeight) && numericWeight >= 600),
        italic: computed?.fontStyle === 'italic',
        underline: computed?.textDecorationLine.includes('underline') ?? false,
        align,
        lineHeight:
          paragraph &&
          Number.parseFloat(window.getComputedStyle(paragraph).fontSize) > 0
            ? Number.parseFloat(window.getComputedStyle(paragraph).lineHeight) /
              Number.parseFloat(window.getComputedStyle(paragraph).fontSize)
            : 1.25,
        color: computedColorToHex(computed?.color ?? '#111827'),
        fontSize: Math.round(Number.parseFloat(computed?.fontSize ?? '24')),
        fontFamily,
      }
    })

    function uniform<T>(values: T[]): T | null {
      const first = values[0]
      if (first === undefined) return null
      return values.every((value) => value === first) ? first : null
    }

    certificateEditorActions.setTextSelectionStyles({
      bold: uniform(styles.map((style) => style.bold)),
      italic: uniform(styles.map((style) => style.italic)),
      underline: uniform(styles.map((style) => style.underline)),
      align: uniform(styles.map((style) => style.align)),
      lineHeight: uniform(styles.map((style) => style.lineHeight)),
      color: uniform(styles.map((style) => style.color)),
      fontSize: uniform(styles.map((style) => style.fontSize)),
      fontFamily: uniform(styles.map((style) => style.fontFamily)),
    })
  }

  function runCommand(command: string, value?: string) {
    restoreSelection()
    document.execCommand(command, false, value)
    const selection = window.getSelection()
    if (selection?.rangeCount)
      savedRange.current = selection.getRangeAt(0).cloneRange()
    commit()
    updateSelectionStyles()
  }

  function applyInlineStyle(style: string, applyAtCaret: () => void) {
    restoreSelection()
    const selection = window.getSelection()
    if (!selection?.rangeCount) return
    const range = selection.getRangeAt(0)

    if (range.collapsed) {
      applyAtCaret()
    } else {
      const safeStyle = style.replace(/"/g, '&quot;')
      const marker = `format-${crypto.randomUUID()}`
      const content = document.createElement('div')
      const fragment = range.cloneContents()
      const walker = document.createTreeWalker(fragment, NodeFilter.SHOW_TEXT)
      const textNodes: Text[] = []
      let current = walker.nextNode()
      while (current) {
        if (current.textContent) textNodes.push(current as Text)
        current = walker.nextNode()
      }
      textNodes.forEach((textNode) => {
        const wrapper = document.createElement('span')
        wrapper.dataset.certificateFormat = marker
        wrapper.setAttribute('style', safeStyle)
        textNode.parentNode?.insertBefore(wrapper, textNode)
        wrapper.appendChild(textNode)
      })
      content.appendChild(fragment)
      document.execCommand('insertHTML', false, content.innerHTML)
      const formatted = editorRef.current?.querySelectorAll(
        `[data-certificate-format="${marker}"]`,
      )
      if (formatted && formatted.length > 0) {
        const formattedRange = document.createRange()
        formattedRange.setStartBefore(formatted[0])
        formattedRange.setEndAfter(formatted[formatted.length - 1])
        selection.removeAllRanges()
        selection.addRange(formattedRange)
        formatted.forEach((node) =>
          node.removeAttribute('data-certificate-format'),
        )
      }
    }

    if (selection.rangeCount) {
      savedRange.current = selection.getRangeAt(0).cloneRange()
    }
    commit()
    updateSelectionStyles()
  }

  useLayoutEffect(() => {
    const editor = editorRef.current
    if (!editor) return
    editor.innerHTML = paragraphsToHtml(element.paragraphs)
    editor.focus()
    const range = document.createRange()
    range.selectNodeContents(editor)
    range.collapse(false)
    const selection = window.getSelection()
    selection?.removeAllRanges()
    selection?.addRange(range)
    savedRange.current = range.cloneRange()
  }, [element.id])

  useEffect(() => {
    function rememberSelection() {
      const editor = editorRef.current
      const selection = window.getSelection()
      if (!editor || !selection?.rangeCount) return
      const range = selection.getRangeAt(0)
      if (editor.contains(range.commonAncestorContainer)) {
        savedRange.current = range.cloneRange()
        updateSelectionStyles()
      }
    }

    document.addEventListener('selectionchange', rememberSelection)
    certificateEditorActions.setRichTextController({
      elementId: element.id,
      commit,
      toggleBold: () => runCommand('bold'),
      toggleItalic: () => runCommand('italic'),
      toggleUnderline: () => runCommand('underline'),
      setAlign: (align) =>
        runCommand(
          align === 'left'
            ? 'justifyLeft'
            : align === 'center'
              ? 'justifyCenter'
              : align === 'right'
                ? 'justifyRight'
                : 'justifyFull',
        ),
      setLineHeight: (lineHeight) => {
        restoreSelection()
        const selection = window.getSelection()
        const editor = editorRef.current
        if (!selection?.rangeCount || !editor) return
        const range = selection.getRangeAt(0)
        editor.querySelectorAll('p, div').forEach((paragraph) => {
          if (range.intersectsNode(paragraph)) {
            const paragraphElement = paragraph as HTMLElement
            paragraphElement.style.lineHeight = String(lineHeight)
          }
        })
        commit()
        updateSelectionStyles()
      },
      setColor: (color) => runCommand('foreColor', color),
      setFontFamily: (fontFamily) =>
        applyInlineStyle(`font-family:${fontFamily}`, () => {
          const primaryFamily = fontFamily
            .split(',')[0]
            .trim()
            .replace(/^['"]|['"]$/g, '')
          document.execCommand('fontName', false, primaryFamily)
        }),
      setFontSize: (fontSize) => {
        applyInlineStyle(`font-size:${fontSize}px`, () => {
          document.execCommand('styleWithCSS', false, 'false')
          document.execCommand('fontSize', false, '7')
          editorRef.current
            ?.querySelectorAll('font[size="7"]')
            .forEach((node) => {
              const font = node as HTMLElement
              font.removeAttribute('size')
              font.style.fontSize = `${fontSize}px`
            })
        })
      },
      insertText: (text) => runCommand('insertText', text),
    })
    updateSelectionStyles()

    return () => {
      document.removeEventListener('selectionchange', rememberSelection)
      certificateEditorActions.setRichTextController(null)
      certificateEditorActions.setTextSelectionStyles(null)
    }
  }, [element.id])

  return (
    <div
      ref={editorRef}
      contentEditable
      suppressContentEditableWarning
      role="textbox"
      aria-multiline="true"
      aria-label="Conteúdo do texto"
      className="h-full w-full cursor-text overflow-auto whitespace-pre-wrap wrap-break-word outline-none"
      style={{ lineHeight: 1.25, overflowWrap: 'anywhere' }}
      onPointerDown={(event) => event.stopPropagation()}
      onInput={commit}
      onKeyDown={(event) => {
        if (event.key === 'Escape') {
          event.preventDefault()
          certificateEditorActions.stopEditing()
        }
      }}
    />
  )
}

function StaticTextElement({ element }: { element: TextCertificateElement }) {
  return (
    <div
      className="h-full w-full overflow-hidden whitespace-pre-wrap wrap-break-word"
      style={{ lineHeight: 1.25, overflowWrap: 'anywhere' }}
    >
      {element.paragraphs.map((paragraph, paragraphIndex) => (
        <div
          key={paragraphIndex}
          style={{
            textAlign: paragraph.align,
            lineHeight: paragraph.lineHeight,
          }}
        >
          {paragraph.runs.length === 0 ? (
            <br />
          ) : (
            paragraph.runs.map((run, runIndex) => (
              <span
                key={runIndex}
                style={{
                  color: run.color,
                  fontSize: run.fontSize,
                  fontFamily: run.fontFamily,
                  fontWeight: run.bold ? 700 : 400,
                  fontStyle: run.italic ? 'italic' : 'normal',
                  textDecoration: run.underline ? 'underline' : 'none',
                }}
              >
                {run.text}
              </span>
            ))
          )}
        </div>
      ))}
    </div>
  )
}
