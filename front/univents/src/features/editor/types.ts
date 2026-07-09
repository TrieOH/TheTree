export type CanvasElementType = 'text' | 'signature' | 'image'

export interface CanvasElementBase {
  id: string
  type: CanvasElementType
  xPct: number
  yPct: number
  widthPct: number
  heightPct: number
  zIndex: number
}

export interface TextCanvasElement extends CanvasElementBase {
  type: 'text'
  content: string
  fontSize: number
  fontWeight: number | string
  fontFamily: string
  color: string
}

export interface SignatureCanvasElement extends CanvasElementBase {
  type: 'signature'
  title?: string
  signatureId: string | null
}

export interface ImageCanvasElement extends CanvasElementBase {
  type: 'image'
  src: string | null
  fileName?: string
}

export type CanvasElement = TextCanvasElement | SignatureCanvasElement | ImageCanvasElement

export interface VariableDefinition {
  key: string
  label: string
  description: string
  defaultValue: string
}

export interface ValidationError {
  elementId: string
  field: string
  message: string
}

export const DEFAULT_VARIABLES: VariableDefinition[] = [
  { key: 'activity_name', label: 'Nome da atividade/edição', description: 'Nome da atividade ou edição', defaultValue: 'XXXXX' },
  { key: 'certified_at', label: 'Data de certificação', description: 'Data em que o certificado foi emitido', defaultValue: '' },
]

export function extractVariables(text: string): string[] {
  const matches = text.match(/\{\{(\w+)\}\}/g)
  if (!matches) return []
  return matches.map(m => m.slice(2, -2))
}

export function validateVariables(
  text: string,
  available: VariableDefinition[]
): { valid: string[]; invalid: string[] } {
  const used = extractVariables(text)
  const validKeys = new Set(available.map(v => v.key))
  return {
    valid: used.filter(k => validKeys.has(k)),
    invalid: used.filter(k => !validKeys.has(k)),
  }
}
