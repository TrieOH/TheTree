import type React from "react";

export interface FieldDefinition<T> {
  name: keyof T;
  placeholder?: string;
  label: string;
  type:
    | "text"
    | "number"
    | "select"
    | "checkbox"
    | "radio"
    | "percentage"
    | "option-picker";
  options?: {
    label: string;
    value: string | number;
    icon?: React.ComponentType<{ className?: string }>;
  }[];
}
