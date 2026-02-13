import {
  Type,
  Mail,
  Hash,
  List,
  CheckSquare,
  ToggleLeft,
  HelpCircle,
  CircleDot,
} from "lucide-react";

import type { LucideIcon } from "lucide-react";

type FieldType =
  | "string"
  | "email"
  | "int"
  | "select"
  | "radio"
  | "checkbox"
  | "bool";

const fieldTypeIconMap = {
  string: Type,
  email: Mail,
  int: Hash,
  select: List,
  radio: CircleDot,
  checkbox: CheckSquare,
  bool: ToggleLeft,
} satisfies Record<FieldType, LucideIcon>;


export function getFieldTypeIcon(type: FieldType): LucideIcon {
  return fieldTypeIconMap[type] ?? HelpCircle;
}
