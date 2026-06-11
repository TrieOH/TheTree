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
  dependsOn?: { field: keyof T; value: unknown };
}