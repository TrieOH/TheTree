import type { FieldDefinitionResultI, Operator } from "../model/types";

export const getCompatibleOperators = (fieldType: FieldDefinitionResultI['type']): Operator[] => {
  const commonOperators: Operator[] = ["equals", "not_equals", "exists", "not_exists"];
  const inNotInOperators: Operator[] = ["in", "not_in"];
  const numericOperators: Operator[] = ["greater_than", "greater_than_equal", "lower_than", "lower_than_equal"];
  const stringOperators: Operator[] = ["contains"];

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
