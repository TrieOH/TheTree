export interface FieldDefinition<T> {
  name: keyof T;
  placeholder?: string;
  label: string;
  type: 'text' | 'number' | 'select' | 'checkbox' | 'radio' | 'textarea' | 'boolean';
  options?: { label: string; value: string | number }[];
  min?: number;
  max?: number;
  rows?: number;
  disabled?: boolean;
  required?: boolean;
  dependsOn?: {
    field: keyof T;
    /** Single value or array of accepted values. When the watched field matches, the field is shown. */
    value: unknown | unknown[];
  };
}