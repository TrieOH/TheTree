
export interface FieldDefinition<T> {
  name: keyof T;
  placeholder?: string;
  label: string;
  type: 'text' | 'number' | 'select' | 'checkbox' | 'radio';
  options?: { label: string; value: string | number }[];
}