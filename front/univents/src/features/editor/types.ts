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
  { key: 'nome', label: 'Nome do participante', description: 'Nome completo do participante', defaultValue: 'Nome do Participante' },
  { key: 'evento', label: 'Nome do evento', description: 'Nome do evento/curso', defaultValue: 'Nome do Evento' },
  { key: 'data', label: 'Data de conclusão', description: 'Data de conclusão formatada', defaultValue: '01 de Janeiro de 2026' },
  { key: 'carga_horaria', label: 'Carga horária', description: 'Carga horária total', defaultValue: '40 horas' },
  { key: 'edicao', label: 'Edição do evento', description: 'Número/ano da edição', defaultValue: '1ª Edição' },
  { key: 'orgao', label: 'Órgão realizador', description: 'Nome da instituição/órgão', defaultValue: 'Univents' },
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
