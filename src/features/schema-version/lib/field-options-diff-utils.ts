import { CustomDiff } from "../hooks/useEditableList";
import type { Option, VersionFieldResult } from "../model/types";

function diffFieldOptions(oldField: VersionFieldResult, newField: VersionFieldResult) {
  const oldOptions = oldField.options ?? [];
  const newOptions = newField.options ?? [];

  const oldMap = new Map(oldOptions.map(o => [o.id, o]));
  const newMap = new Map(newOptions.map(o => [o.id, o]));

  for (const [id, newOpt] of newMap) {
    const oldOpt = oldMap.get(id);
    if (!oldOpt) return { type: "put" as const, values: newOptions };

    if (oldOpt.label !== newOpt.label || oldOpt.value !== newOpt.value)
      return { type: "put" as const, values: newOptions };
  }

  const removedIds: string[] = [];
  for (const [id] of oldMap)
    if (!newMap.has(id)) removedIds.push(id);

  if (removedIds.length > 0) return { type: "delete" as const, ids: removedIds };

  return { type: "none" as const };
}

export function optionsDiff(api: {
  deleteOptions: (fieldId: string, optionIds: string[]) => Promise<void>;
  putOptions: (fieldId: string, options: Option[]) => Promise<void>;
}): CustomDiff<VersionFieldResult> {

  return async ({ getOriginalById, diff }) => {
    for (const [id, newField] of diff.currentMap.entries()) {
      const oldField = getOriginalById(id);
      if (!oldField) continue;
      
      const change = diffFieldOptions(oldField, newField);
      if (change.type === "delete") await api.deleteOptions(id, change.ids);
      if (change.type === "put") await api.putOptions(id, change.values);
    }
  };
}
