import {
  AlignCenter,
  AlignJustify,
  AlignLeft,
  AlignRight,
  Bold,
  Braces,
  Italic,
  Minus,
  Plus,
  Rows3,
  Underline,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import { CERTIFICATE_FONT_FAMILIES, CERTIFICATE_VARIABLES } from '../constants'
import { useCertificateEditorState } from '../store'
import { ToolbarCombobox } from './toolbar-combobox'
import { Button } from '@/shared/ui/shadcn/button'
import { Input } from '@/shared/ui/shadcn/input'

const FONT_OPTIONS = CERTIFICATE_FONT_FAMILIES.map((font) => ({
  value: font.value,
  label: font.label,
}))

const VARIABLE_OPTIONS = CERTIFICATE_VARIABLES.map((variable) => ({
  value: variable.token,
  label: variable.label,
  description: variable.description,
}))

const LINE_HEIGHT_OPTIONS = [
  { value: '1', label: 'Simples' },
  { value: '1.15', label: '1,15' },
  { value: '1.25', label: '1,25' },
  { value: '1.5', label: '1,5' },
  { value: '2', label: 'Duplo' },
]

export function CertificateTextToolbar() {
  const controller = useCertificateEditorState(
    (state) => state.richTextController,
  )
  const selectionStyles = useCertificateEditorState(
    (state) => state.textSelectionStyles,
  )
  const [fontSize, setFontSize] = useState('24')
  const disabled = !controller

  useEffect(() => {
    if (selectionStyles) {
      setFontSize(
        selectionStyles.fontSize === null
          ? '—'
          : String(selectionStyles.fontSize),
      )
    }
  }, [selectionStyles])

  const alignments = [
    ['left', AlignLeft, 'Alinhar à esquerda'],
    ['center', AlignCenter, 'Centralizar'],
    ['right', AlignRight, 'Alinhar à direita'],
    ['justify', AlignJustify, 'Justificar'],
  ] as const

  function applyFontSize(value: number) {
    const nextValue = Math.min(400, Math.max(6, Math.round(value)))
    setFontSize(String(nextValue))
    controller?.setFontSize(nextValue)
  }

  function stepFontSize(step: number) {
    const currentValue = Number(fontSize)
    applyFontSize(
      Number.isFinite(currentValue)
        ? currentValue + step
        : (selectionStyles?.fontSize ?? 24) + step,
    )
  }

  return (
    <div
      className="relative z-30 flex h-10 shrink-0 items-center justify-center gap-0.5 border-b bg-muted/70 px-2 shadow-sm"
      onPointerDown={(event) => event.stopPropagation()}
    >
      <Button
        type="button"
        size="icon-xs"
        variant={selectionStyles?.bold ? 'secondary' : 'ghost'}
        className="relative"
        title="Negrito"
        disabled={disabled}
        onClick={() => controller?.toggleBold()}
      >
        <Bold className="size-3.5" />
        {selectionStyles?.bold === null ? (
          <span className="absolute right-0.5 bottom-0 text-[9px] leading-none">
            −
          </span>
        ) : null}
      </Button>
      <Button
        type="button"
        size="icon-xs"
        variant={selectionStyles?.italic ? 'secondary' : 'ghost'}
        className="relative"
        title="Itálico"
        disabled={disabled}
        onClick={() => controller?.toggleItalic()}
      >
        <Italic className="size-3.5" />
        {selectionStyles?.italic === null ? (
          <span className="absolute right-0.5 bottom-0 text-[9px] leading-none">
            −
          </span>
        ) : null}
      </Button>
      <Button
        type="button"
        size="icon-xs"
        variant={selectionStyles?.underline ? 'secondary' : 'ghost'}
        className="relative"
        title="Sublinhado"
        disabled={disabled}
        onClick={() => controller?.toggleUnderline()}
      >
        <Underline className="size-3.5" />
        {selectionStyles?.underline === null ? (
          <span className="absolute right-0.5 bottom-0 text-[9px] leading-none">
            −
          </span>
        ) : null}
      </Button>
      <span className="mx-1 h-5 w-px shrink-0 bg-border" />
      {alignments.map(([align, Icon, label]) => (
        <Button
          key={align}
          type="button"
          size="icon-xs"
          variant={selectionStyles?.align === align ? 'secondary' : 'ghost'}
          title={label}
          disabled={disabled}
          onClick={() => controller?.setAlign(align)}
        >
          <Icon className="size-3.5" />
        </Button>
      ))}
      <span className="mx-1 h-5 w-px shrink-0 bg-border" />
      <ToolbarCombobox
        value={selectionStyles?.fontFamily ?? undefined}
        options={FONT_OPTIONS}
        placeholder={
          controller && selectionStyles?.fontFamily === null ? '—' : 'Fonte'
        }
        searchPlaceholder="Buscar fonte…"
        disabled={disabled}
        className="w-32"
        onChange={(fontFamily) => controller?.setFontFamily(fontFamily)}
      />
      <div className="flex h-7 shrink-0 items-center rounded-md border border-input bg-background">
        <Button
          type="button"
          size="icon-xs"
          variant="ghost"
          className="rounded-r-none"
          title="Diminuir tamanho da fonte"
          disabled={disabled}
          onClick={() => stepFontSize(-1)}
        >
          <Minus className="size-3" />
        </Button>
        <Input
          type="text"
          inputMode="numeric"
          value={fontSize}
          disabled={disabled}
          aria-label="Tamanho da fonte"
          className="h-6 w-10 rounded-none border-0 bg-transparent px-1 text-center text-xs shadow-none focus-visible:ring-0"
          onFocus={(event) => event.currentTarget.select()}
          onChange={(event) => setFontSize(event.target.value)}
          onBlur={() => {
            const value = Number(fontSize)
            if (Number.isFinite(value) && value >= 6 && value <= 400) {
              applyFontSize(value)
            } else {
              setFontSize(
                selectionStyles?.fontSize === null
                  ? '—'
                  : String(selectionStyles?.fontSize ?? 24),
              )
            }
          }}
          onKeyDown={(event) => {
            if (event.key === 'Enter') {
              event.preventDefault()
              event.stopPropagation()
              event.currentTarget.blur()
            }
          }}
        />
        <Button
          type="button"
          size="icon-xs"
          variant="ghost"
          className="rounded-l-none"
          title="Aumentar tamanho da fonte"
          disabled={disabled}
          onClick={() => stepFontSize(1)}
        >
          <Plus className="size-3" />
        </Button>
      </div>
      <div className="relative h-7 w-8 shrink-0">
        <Input
          type="color"
          value={selectionStyles?.color ?? '#111827'}
          disabled={disabled}
          aria-label="Cor do texto"
          className="h-7 w-8 rounded-md p-1"
          onChange={(event) => controller?.setColor(event.target.value)}
        />
        {controller && selectionStyles?.color === null ? (
          <span className="pointer-events-none absolute inset-0 flex items-center justify-center rounded-md bg-background/75 text-xs text-muted-foreground">
            —
          </span>
        ) : null}
      </div>
      <ToolbarCombobox
        options={VARIABLE_OPTIONS}
        placeholder="Inserir informação dinâmica"
        searchPlaceholder="Buscar informação…"
        disabled={disabled}
        iconOnly
        icon={<Braces className="size-3.5" />}
        onChange={(token) => controller?.insertText(token)}
      />
      <ToolbarCombobox
        value={
          selectionStyles?.lineHeight === null
            ? undefined
            : String(selectionStyles?.lineHeight ?? 1.25)
        }
        options={LINE_HEIGHT_OPTIONS}
        placeholder="Espaçamento entre linhas"
        searchPlaceholder="Buscar espaçamento…"
        disabled={disabled}
        iconOnly
        icon={<Rows3 className="size-3.5" />}
        onChange={(lineHeight) => controller?.setLineHeight(Number(lineHeight))}
      />
    </div>
  )
}
