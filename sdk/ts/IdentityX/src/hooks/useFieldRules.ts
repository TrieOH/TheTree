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
      passed = String(fieldValue) === String(compareValue);
      message = `O campo "${fieldTitle}" deve ser igual a "${String(compareValue)}".`;
      break;
    }
    
    case "not_equals": {
      passed = String(fieldValue) !== String(compareValue) || fieldValue === undefined || fieldValue === null;
      message = `O campo "${fieldTitle}" deve ser diferente de "${String(compareValue)}".`;
      break;
    }
    
    case "contains": {
      passed = String(fieldValue).includes(String(compareValue));
      message = `O campo "${fieldTitle}" deve conter "${String(compareValue)}".`;
      break;
    }
    
    case "gt": {
      passed = Number(fieldValue) > Number(compareValue);
      message = `O campo "${fieldTitle}" deve ser maior que ${String(compareValue)}.`;
      break;
    }
    
    case "gte": {
      passed = Number(fieldValue) >= Number(compareValue);
      message = `O campo "${fieldTitle}" deve ser maior ou igual a ${String(compareValue)}.`;
      break;
    }
    
    case "lt": {
      passed = Number(fieldValue) < Number(compareValue);
      message = `O campo "${fieldTitle}" deve ser menor que ${String(compareValue)}.`;
      break;
    }
    
    case "lte": {
      passed = Number(fieldValue) <= Number(compareValue);
      message = `O campo "${fieldTitle}" deve ser menor ou igual a ${String(compareValue)}.`;
      break;
    }
    
    case "in": {
      const options = String(compareValue).split(",").map(v => v.trim());
      passed = options.includes(String(fieldValue));
      const optionsStr = options.join(", ");
      message = `O campo "${fieldTitle}" deve ser um dos seguintes: ${optionsStr}.`;
      break;
    }
    
    case "not_in": {
      const options = String(compareValue).split(",").map(v => v.trim());
      passed = !options.includes(String(fieldValue));
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