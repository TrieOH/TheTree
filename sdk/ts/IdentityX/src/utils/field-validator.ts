export type Rule = {
  id?: string;
  message: string;
  test: (value: string) => boolean;
};

export type RuleStatus = Rule & { passed: boolean };

export function evaluateRules(rules: Rule[], value: string): RuleStatus[] {
  return rules.map(r => ({ ...r, passed: !!r.test(value) }));
}