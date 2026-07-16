import { useEffect } from 'react'
import { EditorContent, useEditor, useEditorState } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import TextAlign from '@tiptap/extension-text-align'
import { TextStyleKit } from '@tiptap/extension-text-style'
import { Selection } from '@tiptap/extensions'
import type { TextCertificateElement } from '../../types'
import {
  domToParagraphs,
  paragraphsToHtml,
} from '../../rich-text/dom-serializer'
import { CertificateLineHeight } from '../../rich-text/line-height-extension'
import { certificateEditorActions } from '../../store'
import { DEFAULT_CERTIFICATE_FONT } from '../../constants'

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
  const editor = useEditor({
    immediatelyRender: false,
    extensions: [
      StarterKit.configure({
        heading: false,
        blockquote: false,
        bulletList: false,
        orderedList: false,
        codeBlock: false,
        horizontalRule: false,
      }),
      TextStyleKit,
      TextAlign.configure({ types: ['paragraph'] }),
      CertificateLineHeight,
      Selection.configure({ className: 'certificate-preserved-selection' }),
    ],
    content: paragraphsToHtml(element.paragraphs),
    autofocus: 'end',
    editorProps: {
      attributes: {
        class:
          'certificate-rich-text h-full w-full cursor-text overflow-auto whitespace-pre-wrap wrap-break-word outline-none',
        'aria-label': 'Conteúdo do texto',
      },
      handleKeyDown: (_view, event) => {
        if (event.key !== 'Escape') return false
        certificateEditorActions.stopEditing()
        return true
      },
    },
    onUpdate: ({ editor: currentEditor }) => {
      updateParagraphs(element.id, domToParagraphs(currentEditor.view.dom))
    },
  })

  const selectionStyles = useEditorState({
    editor,
    selector: ({ editor: currentEditor }) =>
      currentEditor ? readSelectionStyles(currentEditor) : null,
  })

  useEffect(() => {
    if (!editor) return
    const commit = () =>
      updateParagraphs(element.id, domToParagraphs(editor.view.dom))
    certificateEditorActions.setRichTextController({
      elementId: element.id,
      commit,
      toggleBold: () => void editor.chain().focus().toggleBold().run(),
      toggleItalic: () => void editor.chain().focus().toggleItalic().run(),
      toggleUnderline: () =>
        void editor.chain().focus().toggleUnderline().run(),
      setAlign: (align) =>
        void editor.chain().focus().setTextAlign(align).run(),
      setLineHeight: (lineHeight) =>
        void editor.chain().focus().setCertificateLineHeight(lineHeight).run(),
      setColor: (color) => void editor.chain().focus().setColor(color).run(),
      setFontFamily: (fontFamily) =>
        void editor.chain().focus().setFontFamily(fontFamily).run(),
      setFontSize: (fontSize) =>
        void editor.chain().focus().setFontSize(`${fontSize}px`).run(),
      insertText: (text) =>
        void editor.chain().focus().insertContent(text).run(),
    })

    return () => {
      certificateEditorActions.setRichTextController(null)
      certificateEditorActions.setTextSelectionStyles(null)
    }
  }, [editor, element.id])

  useEffect(() => {
    certificateEditorActions.setTextSelectionStyles(selectionStyles)
  }, [selectionStyles])

  if (!editor) return null
  return (
    <EditorContent
      editor={editor}
      className="h-full w-full"
      onPointerDown={(event) => event.stopPropagation()}
    />
  )
}

function updateParagraphs(
  elementId: string,
  paragraphs: TextCertificateElement['paragraphs'],
) {
  certificateEditorActions.updateElement(elementId, (current) =>
    current.type === 'text' ? { ...current, paragraphs } : current,
  )
}

function uniform<T>(values: T[]): T | null {
  const first = values[0]
  return first !== undefined && values.every((value) => value === first)
    ? first
    : null
}

type RichEditor = NonNullable<ReturnType<typeof useEditor>>

function readSelectionStyles(editor: RichEditor) {
  const { from, to, empty, $from } = editor.state.selection
  const markSets: Array<
    readonly { type: { name: string }; attrs: Record<string, unknown> }[]
  > = []
  const aligns: Array<'left' | 'center' | 'right' | 'justify'> = []
  const lineHeights: number[] = []

  const readParagraph = (attrs: Record<string, unknown>) => {
    const align = String(attrs.textAlign ?? 'left')
    aligns.push(
      align === 'center' || align === 'right' || align === 'justify'
        ? align
        : 'left',
    )
    lineHeights.push(Number(attrs.lineHeight) || 1.25)
  }

  if (empty) {
    markSets.push(editor.state.storedMarks ?? $from.marks())
    readParagraph($from.parent.attrs)
  } else {
    editor.state.doc.nodesBetween(from, to, (node) => {
      if (node.isText) markSets.push(node.marks)
      if (node.type.name === 'paragraph') readParagraph(node.attrs)
    })
  }

  const hasMark = (
    marks: readonly { type: { name: string } }[],
    name: string,
  ) => marks.some((mark) => mark.type.name === name)
  const styleValue = (
    marks: readonly {
      type: { name: string }
      attrs: Record<string, unknown>
    }[],
    key: string,
    fallback: string,
  ) => {
    const value = marks.find((mark) => mark.type.name === 'textStyle')?.attrs[
      key
    ]
    return typeof value === 'string' && value ? value : fallback
  }

  return {
    bold: uniform(markSets.map((marks) => hasMark(marks, 'bold'))),
    italic: uniform(markSets.map((marks) => hasMark(marks, 'italic'))),
    underline: uniform(markSets.map((marks) => hasMark(marks, 'underline'))),
    align: uniform(aligns),
    lineHeight: uniform(lineHeights),
    color: uniform(
      markSets.map((marks) => styleValue(marks, 'color', '#111827')),
    ),
    fontSize: uniform(
      markSets.map((marks) =>
        Number.parseFloat(styleValue(marks, 'fontSize', '24px')),
      ),
    ),
    fontFamily: uniform(
      markSets.map((marks) =>
        styleValue(marks, 'fontFamily', DEFAULT_CERTIFICATE_FONT),
      ),
    ),
  }
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
            fontSize:
              paragraph.runs[0]?.fontSize ??
              findNearestRun(element.paragraphs, paragraphIndex)?.fontSize,
            fontFamily:
              paragraph.runs[0]?.fontFamily ??
              findNearestRun(element.paragraphs, paragraphIndex)?.fontFamily,
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

function findNearestRun(
  paragraphs: TextCertificateElement['paragraphs'],
  paragraphIndex: number,
) {
  return (
    paragraphs
      .slice(0, paragraphIndex)
      .reverse()
      .find((paragraph) => paragraph.runs.length > 0)
      ?.runs.at(-1) ??
    paragraphs
      .slice(paragraphIndex + 1)
      .find((paragraph) => paragraph.runs.length > 0)?.runs[0]
  )
}
