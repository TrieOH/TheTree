import type { CustomDiff } from "../hooks/useEditableList";
import type { Rule, VersionFieldResult } from "../model/types";

function simplifyRule(rule: Rule) {
  return {
    depends_on_field_key: rule.depends_on_field_key,
    operator: rule.operator,
    value: rule.value,
  };
}

function areRuleListsEqual(a: Rule[], b: Rule[]): boolean {
  if (a.length !== b.length) return false;

  const simplifiedA = a.map(simplifyRule);
  const simplifiedB = b.map(simplifyRule);

  return JSON.stringify(simplifiedA) === JSON.stringify(simplifiedB);
}

function diffFieldRules(oldField: VersionFieldResult, newField: VersionFieldResult) {
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
  putRequiredRules: (fieldId: string, rules: Rule[]) => Promise<void>;
  putVisibilityRules: (fieldId: string, rules: Rule[]) => Promise<void>;
}): CustomDiff<VersionFieldResult> {
  return async ({ getOriginalById, diff }) => {
    // 1. Identify deleted fields
    const deletedFieldsInfo: { id: string; key: string }[] = [];
    for (const [originalId, originalField] of diff.originalMap.entries()) {
      if (!diff.currentMap.has(originalId))
        deletedFieldsInfo.push({ id: originalId, key: originalField.key });
    }

    // 2. If there are deleted fields, clean up rules in other existing fields that depend on them
    if (deletedFieldsInfo.length > 0) {
      for (const [currentFieldId, currentField] of diff.currentMap.entries()) {
        let requiredRulesModified = false;
        const newRequiredRules = currentField.required_rules.filter(rule => {
          const dependsOnDeleted = deletedFieldsInfo.some(deleted =>
            deleted.id === rule.depends_on_field_id || deleted.key === rule.depends_on_field_key
          );
          if (dependsOnDeleted) requiredRulesModified = true;
          
          return !dependsOnDeleted;
        });

        let visibilityRulesModified = false;
        const newVisibilityRules = currentField.visibility_rules.filter(rule => {
          const dependsOnDeleted = deletedFieldsInfo.some(deleted =>
            deleted.id === rule.depends_on_field_id || deleted.key === rule.depends_on_field_key
          );
          if (dependsOnDeleted) {
            visibilityRulesModified = true;
          }
          return !dependsOnDeleted;
        });

        // Only call API if rules actually changed after cleanup
        if (requiredRulesModified && !areRuleListsEqual(currentField.required_rules, newRequiredRules)) {
          await api.putRequiredRules(currentFieldId, newRequiredRules);
        }
        if (visibilityRulesModified && !areRuleListsEqual(currentField.visibility_rules, newVisibilityRules)) {
          await api.putVisibilityRules(currentFieldId, newVisibilityRules);
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
          await api.putRequiredRules(id, change.requiredRules);
        if (!areRuleListsEqual(oldField.visibility_rules, change.visibilityRules))
          await api.putVisibilityRules(id, change.visibilityRules);
      }
    }
  };
}
