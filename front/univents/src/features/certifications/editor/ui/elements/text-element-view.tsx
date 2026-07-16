import type { TextCertificateElement } from '../../types'

interface TextElementViewProps {
  element: TextCertificateElement
}

export function TextElementView({ element }: TextElementViewProps) {
  return (
    <div
      className="h-full w-full overflow-hidden whitespace-pre-wrap wrap-break-word"
      style={{ lineHeight: 1.25, overflowWrap: 'anywhere' }}
    >
      {element.paragraphs.map((paragraph, paragraphIndex) => (
        <div key={paragraphIndex} style={{ textAlign: paragraph.align }}>
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
