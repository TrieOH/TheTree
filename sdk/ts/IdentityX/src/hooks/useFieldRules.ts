import { useMemo } from "react";
import type { FieldValue, Operator, RuleResultI } from "../types/fields-types";
import type { RuleStatus } from "../utils/field-validator";

export function useFieldRules(
  rules: RuleResultI[] | undefined,
  values: Record<string, FieldValue>,
  fieldsMap: Record<string, { key: string; title: string }>
): { 
  satisfied: boolean; 
  statuses: RuleStatus[] 
} {
  return useMemo(() => {
    if (!rules || rules.length === 0) return { satisfied: true, statuses: [] };
    
    const results = rules.map((rule) => evaluateRule(rule, values, fieldsMap));
    const allPassed = results.every((r) => r.passed);

    return {
      satisfied: allPassed,
      statuses: results,
    };
  }, [rules, values, fieldsMap]);
}

function evaluateRule(
  rule: RuleResultI,
  values: Record<string, FieldValue>,
  fieldsMap: Record<string, { key: string; title: string }>
): RuleStatus {
  const dependentField = fieldsMap[rule.depends_on_field_id];
  
  if (!dependentField) {
    const message = `Campo dependente não encontrado.`;
    return { id: rule.id, message, passed: false };
  }

  const fieldValue = values[dependentField.key];
  const operator = rule.operator as Operator;
  const compareValue = rule.value;
  const fieldTitle = dependentField.title;

  let passed: boolean;
  let message: string;

  const normalize = (val: FieldValue) => (val === undefined || val === null) ? "" : String(val);

  switch (operator) {
    case "exists": {
      passed = fieldValue !== undefined && fieldValue !== null && fieldValue !== "";
      message = `O campo "${fieldTitle}" deve estar preenchido.`;
      break;
    }
    
    case "not_exists": {
      passed = fieldValue === undefined || fieldValue === null || fieldValue === "";
      message = `O campo "${fieldTitle}" deve estar vazio.`;
      break;
    }
    
    case "equals": {
      passed = normalize(fieldValue) === normalize(compareValue);
      message = `O campo "${fieldTitle}" deve ser igual a "${normalize(compareValue)}".`;
      break;
    }
    
    case "not_equals": {
      passed = normalize(fieldValue) !== normalize(compareValue);
      message = `O campo "${fieldTitle}" deve ser diferente de "${normalize(compareValue)}".`;
      break;
    }
    
    case "contains": {
      passed = normalize(fieldValue).includes(normalize(compareValue));
      message = `O campo "${fieldTitle}" deve conter "${normalize(compareValue)}".`;
      break;
    }
    
    case "greater_than": {
      passed = Number(fieldValue) > Number(compareValue);
      message = `O campo "${fieldTitle}" deve ser maior que ${normalize(compareValue)}.`;
      break;
    }
    
    case "greater_than_equal": {
      passed = Number(fieldValue) >= Number(compareValue);
      message = `O campo "${fieldTitle}" deve ser maior ou igual a ${normalize(compareValue)}.`;
      break;
    }
    
    case "lower_than": {
      passed = Number(fieldValue) < Number(compareValue);
      message = `O campo "${fieldTitle}" deve ser menor que ${normalize(compareValue)}.`;
      break;
    }
    
    case "lower_than_equal": {
      passed = Number(fieldValue) <= Number(compareValue);
      message = `O campo "${fieldTitle}" deve ser menor ou igual a ${normalize(compareValue)}.`;
      break;
    }
    
    case "in": {
      const options = normalize(compareValue).split(",").map(v => v.trim());
      passed = options.includes(normalize(fieldValue));
      const optionsStr = options.join(", ");
      message = `O campo "${fieldTitle}" deve ser um dos seguintes: ${optionsStr}.`;
      break;
    }
    
    case "not_in": {
      const options = normalize(compareValue).split(",").map(v => v.trim());
      passed = !options.includes(normalize(fieldValue));
      const optionsStr = options.join(", ");
      message = `O campo "${fieldTitle}" não pode ser um dos seguintes: ${optionsStr}.`;
      break;
    }
    
    default: {
      passed = false;
      message = `Operador "${operator}" desconhecido.`;
    }
  }

  return {
    passed,
    id: rule.id,
    message,
  };
}