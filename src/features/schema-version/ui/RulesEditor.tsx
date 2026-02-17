import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { PlusIcon, TrashIcon } from 'lucide-react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/shadcn/select";
import { type RuleOperator, ruleOperatorSchema, type Rule, type VersionFieldResult } from '../model/types';
import { RuleValueInput } from './RuleValueInput';
import { getCompatibleOperators } from '../lib/rule-utils';

interface RulesEditorProps {
  rules: Rule[];
  allFieldKeys: string[];
  allFields: VersionFieldResult[];
  onChange: (rules: Rule[]) => void;
}

const operatorConfig: Record<RuleOperator, { label: string; symbol: string }> = {
  equals: { label: 'Equals', symbol: '=' },
  not_equals: { label: 'Not equals', symbol: '≠' },
  in: { label: 'Is in', symbol: '∈' },
  not_in: { label: 'Is not in', symbol: '∉' },
  exists: { label: 'Exists', symbol: '∃' },
  not_exists: { label: 'Not exists', symbol: '∄' },
  contains: { label: 'Contains', symbol: '⊇' },
  gt: { label: 'Greater than', symbol: '>' },
  gte: { label: 'Greater than or equals', symbol: '≥' },
  lt: { label: 'Less than', symbol: '<' },
  lte: { label: 'Less than or equals', symbol: '≤' },
};

export const RulesEditor: React.FC<RulesEditorProps> = ({ rules, allFieldKeys, allFields, onChange }) => {
  const handleRuleChange = (index: number, key: keyof Rule, value: unknown) => {
    const newRules = [...rules];
    newRules[index] = { ...newRules[index], [key]: value };
    onChange(newRules);
  };

  const handleAddRule = () => {
    const newRule: Rule = {
      id: crypto.randomUUID(),
      depends_on_field_key: allFieldKeys.length > 0 ? allFieldKeys[0] : '',
      depends_on_field_id: '',
      operator: ruleOperatorSchema.enum.equals,
      value: null,
    };
    onChange([...rules, newRule]);
  };

  const handleRemoveRule = (index: number) => {
    const newRules = rules.filter((_, i) => i !== index);
    onChange(newRules);
  };

  return (
    <div className="space-y-2 w-full min-w-0">
      {rules.map((rule, index) => {
        const dependentField = allFields.find(f => f.key === rule.depends_on_field_key);
        const fieldType = dependentField?.type || 'string';
        const fieldOptions = dependentField?.options;
        const compatibleOperatorsForField = getCompatibleOperators(fieldType);

        return (
          <div
            key={rule.id}
            className="flex flex-col gap-2 p-2 bg-muted/20 border rounded-md group hover:bg-muted/40 transition-colors"
          >
            <div className="flex flex-wrap items-center gap-1.5 sm:gap-2">
              <Select
                value={rule.depends_on_field_key ?? ''}
                onValueChange={(val) => {
                  const newDependentField = allFields.find(f => f.key === val);
                  const newFieldType = newDependentField?.type || 'string';
                  const newCompatibleOperators = getCompatibleOperators(newFieldType);

                  let newOperator = rule.operator;
                  let newValue = rule.value;

                  if (!newCompatibleOperators.includes(newOperator)) {
                    newOperator = newCompatibleOperators[0] || ruleOperatorSchema.enum.equals;
                    newValue = null;
                  }

                  const updatedRule: Rule = {
                    ...rule,
                    depends_on_field_key: val,
                    operator: newOperator,
                    value: newValue
                  };
                  const newRules = [...rules];
                  newRules[index] = updatedRule;
                  onChange(newRules);
                }}
              >
                <SelectTrigger className="flex-1 min-w-20 h-8 text-xs border border-input bg-background shadow-[1px_1px_0_0_var(--color-input)] focus:ring-offset-background focus:ring-2 focus:ring-ring focus:outline-none">
                  <SelectValue placeholder="Field..." />
                </SelectTrigger>
                <SelectContent>
                  {allFieldKeys.map((key) => (
                    <SelectItem key={key} value={key} className="text-xs">
                      {key}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>

              <Select
                value={rule.operator}
                onValueChange={(val) => handleRuleChange(index, 'operator', val as RuleOperator)}
              >
                <SelectTrigger
                  className="w-10 sm:w-12 h-8 px-0 justify-center text-base sm:text-lg font-medium hover:bg-muted transition-colors shrink-0"
                  title={operatorConfig[rule.operator].label}
                >
                  <SelectValue placeholder="=">
                    <span className="leading-none">{operatorConfig[rule.operator].symbol}</span>
                  </SelectValue>
                </SelectTrigger>
                <SelectContent align="center" className="min-w-35">
                  {Object.values(ruleOperatorSchema.enum)
                    .filter(op => compatibleOperatorsForField.includes(op))
                    .map((op) => (
                      <SelectItem key={op} value={op} className="text-xs gap-2">
                        <span className="w-6 text-center text-base font-medium">{operatorConfig[op].symbol}</span>
                        <span className="text-muted-foreground">{operatorConfig[op].label}</span>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

              <ShadowButton
                type="button"
                variant="ghost"
                leftIcon={<TrashIcon className="h-3.5 w-3.5" />}
                className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive justify-center shrink-0"
                onClick={() => handleRemoveRule(index)}
              />
            </div>

            <div className="w-full">
              <RuleValueInput
                id={`rule-value-${rule.id}`}
                value={rule.value}
                onChange={(val) => handleRuleChange(index, 'value', val)}
                fieldType={fieldType}
                options={fieldOptions}
                operator={rule.operator}
              />
            </div>
          </div>
        );
      })}

      <ShadowButton
        type="button"
        variant="ghost"
        leftIcon={<PlusIcon className="h-4 w-4" />}
        onClick={handleAddRule}
        value='Add Rule'
        className="w-full h-8 text-xs text-muted-foreground hover:text-foreground border border-dashed border-muted-foreground/30 hover:border-muted-foreground/60"
      />
    </div>
  );
};