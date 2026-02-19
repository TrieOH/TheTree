import type { CustomDiff } from "../hooks/useEditableList";
import type { FieldDefinitionResultI, RuleFieldCreateRequestI, RuleResultI } from "../model/types";
import { mapRuleResultToRuleFieldCreateRequest } from "./convert-field-utils";

function areRuleListsEqual(a: RuleResultI[], b: RuleResultI[]): boolean {
  if (a.length !== b.length) return false;

  const simplifiedA = a;
  const simplifiedB = b;

  return JSON.stringify(simplifiedA) === JSON.stringify(simplifiedB);
}

function diffFieldRules(oldField: FieldDefinitionResultI, newField: FieldDefinitionResultI) {
  const requiredRulesChanged = !areRuleListsEqual(oldField.required_rules, newField.required_rules);
  const visibilityRulesChanged = !areRuleListsEqual(oldField.visibility_rules, newField.visibility_rules);

  if (requiredRulesChanged || visibilityRulesChanged) {
    return {
      type: "put" as const,
      requiredRules: newField.required_rules,
      visibilityRules: newField.visibility_rules,
    };
  }

  return { type: "none" as const };
}

export function rulesDiff(api: {
  putRequiredRules: (fieldId: string, rules: RuleFieldCreateRequestI[]) => Promise<void>;
  putVisibilityRules: (fieldId: string, rules: RuleFieldCreateRequestI[]) => Promise<void>;
}): CustomDiff<FieldDefinitionResultI> {
  return async ({ getOriginalById, diff }) => {
    // 1. Identify deleted fields
    const deletedFieldsInfo: { id: string; key: string }[] = [];
    for (const [originalId, originalField] of diff.originalMap.entries()) {
      if (!diff.currentMap.has(originalId))
        deletedFieldsInfo.push({ id: originalId, key: originalField.key });
    }

    // 2. If there are deleted fields, clean up rules in other existing fields that depend on them
    const fields = [...diff.currentMap.values()];
    if (deletedFieldsInfo.length > 0) {
      for (const [currentFieldId, currentField] of diff.currentMap.entries()) {
        let requiredRulesModified = false;
        const newRequiredRules = currentField.required_rules.filter(rule => {
          const dependsOnDeleted = deletedFieldsInfo.some(deleted =>
            deleted.id === rule.depends_on_field_id
          );
          if (dependsOnDeleted) requiredRulesModified = true;
          
          return !dependsOnDeleted;
        });

        let visibilityRulesModified = false;
        const newVisibilityRules = currentField.visibility_rules.filter(rule => {
          const dependsOnDeleted = deletedFieldsInfo.some(deleted =>
            deleted.id === rule.depends_on_field_id
          );
          if (dependsOnDeleted) visibilityRulesModified = true;
          return !dependsOnDeleted;
        });

        // Only call API if rules actually changed after cleanup
        if (requiredRulesModified && !areRuleListsEqual(currentField.required_rules, newRequiredRules)) {
          await api.putRequiredRules(
            currentFieldId,  mapRuleResultToRuleFieldCreateRequest(newRequiredRules, fields)
          );
        }
        
        if (visibilityRulesModified && !areRuleListsEqual(currentField.visibility_rules, newVisibilityRules)) {
          await api.putVisibilityRules(
            currentFieldId, mapRuleResultToRuleFieldCreateRequest(newVisibilityRules, fields)
          );
        }
      }
    }

    // 3. Handle updated fields: existing logic for changes within a field
    for (const [id, newField] of diff.currentMap.entries()) {
      const oldField = getOriginalById(id);
      if (!oldField) continue;

      const change = diffFieldRules(oldField, newField);

      if (change.type === "put") {
        if (!areRuleListsEqual(oldField.required_rules, change.requiredRules))
          await api.putRequiredRules(id, mapRuleResultToRuleFieldCreateRequest(change.requiredRules, fields));
        if (!areRuleListsEqual(oldField.visibility_rules, change.visibilityRules))
          await api.putVisibilityRules(id, mapRuleResultToRuleFieldCreateRequest(change.visibilityRules, fields));
      }
    }
  };
}
