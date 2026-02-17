import type { RuleOperator, VersionFieldResult } from "../model/types";

export const getCompatibleOperators = (fieldType: VersionFieldResult['type']): RuleOperator[] => {
  const commonOperators: RuleOperator[] = ["equals", "not_equals", "exists", "not_exists"];
  const inNotInOperators: RuleOperator[] = ["in", "not_in"];
  const numericOperators: RuleOperator[] = ["gt", "gte", "lt", "lte"];
  const stringOperators: RuleOperator[] = ["contains"];

  switch (fieldType) {
    case 'string':
    case 'email':
      return [...commonOperators, ...inNotInOperators, ...stringOperators];
    case 'int':
      return [...commonOperators, ...inNotInOperators, ...numericOperators];
    case 'bool':
      return [...commonOperators];
    case 'select':
    case 'radio':
    case 'checkbox':
      return [...commonOperators, ...inNotInOperators];
    default:
      return commonOperators;
  }
};
