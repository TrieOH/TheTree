import React from 'react';
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { ShadowInput } from '@/shared/ui/form/ShadowInput';
import { PlusIcon, TrashIcon } from 'lucide-react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/shadcn/select";
import { RuleOperator, ruleOperatorSchema } from '../model/types';

interface Rule {
  depends_on_field_key: string;
  operator: RuleOperator;
  value: string;
}

interface RulesEditorProps {
  rules: Rule[];
  allFieldKeys: string[];
  onChange: (rules: Rule[]) => void;
}

const operatorConfig: Record<RuleOperator, {  label: string; symbol: string }> = {
  equals: { 
    label: 'Equals',
    symbol: '=' 
  },
  not_equals: { 
    label: 'Not equals',
    symbol: '≠' 
  },
  in: { 
    label: 'Contains',
    symbol: '∋' 
  },
  not_in: { 
    label: 'Not contains',
    symbol: '∌' 
  },
  exists: { 
    label: 'Exists',
    symbol: '∃' 
  },
  not_exists: { 
    label: 'Not exists',
    symbol: '∄' 
  },
};

export const RulesEditor: React.FC<RulesEditorProps> = ({ rules, allFieldKeys, onChange }) => {
  const handleRuleChange = (index: number, key: keyof Rule, value: string) => {
    const newRules = [...rules];
    newRules[index] = { ...newRules[index], [key]: value };
    onChange(newRules);
  };

  const handleAddRule = () => {
    const newRule: Rule = {
      depends_on_field_key: allFieldKeys.length > 0 ? allFieldKeys[0] : '',
      operator: ruleOperatorSchema.enum.equals,
      value: '',
    };
    onChange([...rules, newRule]);
  };

  const handleRemoveRule = (index: number) => {
    const newRules = rules.filter((_, i) => i !== index);
    onChange(newRules);
  };

  return (
    <div className="space-y-2">
      {rules.map((rule, index) => (
                <div
                  key={index}
                  className="flex flex-wrap items-center gap-2 p-2 bg-muted/20 border rounded-md group hover:bg-muted/40 transition-colors"
                >          <Select
            value={rule.depends_on_field_key}
            onValueChange={(val) => handleRuleChange(index, 'depends_on_field_key', val)}
          >
            <SelectTrigger className="flex-1 min-w-30 h-8 text-xs border border-input bg-background shadow-[1px_1px_0_0_var(--color-input)] focus:ring-offset-background focus:ring-2 focus:ring-ring focus:outline-none">
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
              className="w-12 h-8 px-0 justify-center text-lg font-medium hover:bg-muted transition-colors" 
              title={operatorConfig[rule.operator].label}
            >
              <SelectValue placeholder="=">
                <span className="leading-none">{operatorConfig[rule.operator].symbol}</span>
              </SelectValue>
            </SelectTrigger>
            <SelectContent align="center" className="min-w-30">
              {Object.values(ruleOperatorSchema.enum).map((op) => (
                <SelectItem key={op} value={op} className="text-xs gap-2">
                  <span className="w-6 text-center text-base font-medium">{operatorConfig[op].symbol}</span>
                  <span className="text-muted-foreground">{operatorConfig[op].label}</span>
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <ShadowInput
            placeholder="value..."
            value={rule.value}
            onChange={(val) => handleRuleChange(index, 'value', val)}
            className="flex-1 h-8 text-xs min-w-20"
          />

          <ShadowButton
            type="button"
            variant="ghost"
            leftIcon={<TrashIcon className="h-3.5 w-3.5" />}
            className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive justify-center"
            onClick={() => handleRemoveRule(index)}
          />
        </div>
      ))}

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