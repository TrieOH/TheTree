export type FieldType = 'text' | 'email' | 'textarea' | 'checkbox' | 'select' |
  'number' | 'url' | 'percentage'

export interface FormFieldI<T> {
  name: keyof T
  label: string
  type: FieldType
  placeholder?: string
  required?: boolean
  options?: { value: string; label: string }[]
  rows?: number
  span?: 'full' | 'half'
}
