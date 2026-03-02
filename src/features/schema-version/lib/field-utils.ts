import type { FieldDefinitionResultI } from "../model/types";

export const areFieldsEqual = (a: FieldDefinitionResultI, b: FieldDefinitionResultI): boolean => {
  if (a.id !== b.id) return false;
  if (a.object_id !== b.object_id) return false;
  if (a.key !== b.key) return false;
  if (a.title !== b.title) return false;
  if (a.type !== b.type) return false;
  if (a.required !== b.required) return false;
  if (a.placeholder !== b.placeholder) return false;
  if (a.position !== b.position) return false;
  if (String(a.default_value) !== String(b.default_value)) return false;
  if (a.description !== b.description) return false;
  if (a.owner !== b.owner) return false;
  if (a.mutable !== b.mutable) return false;

  return true;
};
