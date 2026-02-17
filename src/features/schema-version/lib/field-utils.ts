import type { VersionFieldResult } from "../model/types";

export const areFieldsEqual = (a: VersionFieldResult, b: VersionFieldResult): boolean => {
  if (a.title !== b.title) return false;
  if (a.type !== b.type) return false;
  if (a.required !== b.required) return false;
  if (a.placeholder !== b.placeholder) return false;
  if (a.position !== b.position) return false;
  if (String(a.default_value) !== String(b.default_value)) return false;
  if (a.description !== b.description) return false;
  if (a.owner !== b.owner) return false;
  if (a.mutable !== b.mutable) return false;
  if (a.key !== b.key) return false;

  const simplifyRule = (rule: any) => ({
    operator: rule.operator,
    value: rule.value,
    depends_on_field_key: rule.depends_on_field_key
  });

  const simplifiedARequiredRules = a.required_rules.map(simplifyRule);
  const simplifiedBRequiredRules = b.required_rules.map(simplifyRule);
  if (JSON.stringify(simplifiedARequiredRules) !== JSON.stringify(simplifiedBRequiredRules)) return false;

  const simplifiedAVisibilityRules = a.visibility_rules.map(simplifyRule);
  const simplifiedBVisibilityRules = b.visibility_rules.map(simplifyRule);
  if (JSON.stringify(simplifiedAVisibilityRules) !== JSON.stringify(simplifiedBVisibilityRules)) return false;
  
  return true;
};
