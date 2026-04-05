import type { Path, PathValue, FieldValues } from "react-hook-form"

export type FieldType = 'text' | 'email' | 'textarea' | 'checkbox' | 'select' |
  'number' | 'url' | 'percentage' | 'image-upload' | 'gallery-upload' | 'datetime'

export interface FormFieldI<T extends FieldValues> {
  name: Path<T>
  label: string
  type: FieldType
  placeholder?: string
  required?: boolean
  options?: { value: string; label: string }[]
  rows?: number
  span?: 'full' | 'half'
  autocomplete?: string
  autoFocus?: boolean
  // Image and Gallery
  accept?: string
  maxSize?: number
  uploadFn?: (file: File) => Promise<string>
  itemActions?: {
    label: string
    icon?: 'image' | 'layout' | 'star'
    onClick?: (url: string, setValue: (name: Path<T>, value: PathValue<T, Path<T>>) => void) => void
  }[]
}
