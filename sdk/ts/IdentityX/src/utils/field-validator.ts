export type Rule = {
  id?: string;
  message: string;
  test: (value: string) => boolean;
};

export type RuleStatus = {
  id?: string;
  message: string;
  passed: boolean;
};

export function evaluateRules(rules: Rule[], value: string): RuleStatus[] {
  return rules.map(r => ({ 
    id: r.id,
    message: r.message,
    passed: !!r.test(value) 
  }));
}