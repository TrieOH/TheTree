export type CanvasElementType = 'text' | 'signature' | 'image'

export interface CanvasElementBase {
  id: string
  type: CanvasElementType
  xPct: number
  yPct: number
  widthPct: number
  heightPct: number
  zIndex: number
  fixed?: boolean
}

export interface TextCanvasElement extends CanvasElementBase {
  type: 'text'
  content: string
  fontSize: number
  fontWeight: number | string
  fontFamily: string
  color: string
  textAlign?: 'left' | 'center' | 'right'
  fixed?: boolean
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
  { key: 'cert_hash', label: 'Hash do certificado', description: 'Identificador único para validação do certificado', defaultValue: 'HASH' },
  { key: 'verify_url', label: 'URL de verificação', description: 'Endereço público para validar o certificado', defaultValue: 'http://localhost:3002/verify/HASH' },
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
