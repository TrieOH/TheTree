import type React from 'react';
import { ShadowInput } from '@/shared/ui/form/ShadowInput';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/shadcn/select';
import { RadioGroup, RadioGroupItem } from '@/shared/ui/shadcn/radio-group';
import type { Option, RuleOperator, VersionFieldResult } from '../model/types';

interface RuleValueInputProps {
  value: unknown;
  onChange: (value: unknown) => void;
  fieldType: VersionFieldResult['type'];
  options?: Option[];
  operator: RuleOperator;
  id: string; // for input ID
}

const NONE_VALUE = "__NONE__"; // Unique string to represent "None"

export const RuleValueInput: React.FC<RuleValueInputProps> = ({
  value,
  onChange,
  fieldType,
  options,
  operator,
  id,
}) => {
  // For 'exists' and 'not_exists' operators, the value field is irrelevant.
  // We can disable the input or show a static message.
  const isValueIrrelevant = ['exists', 'not_exists'].includes(operator);

  if (isValueIrrelevant) {
    return (
      <ShadowInput
        id={id}
        value="N/A"
        disabled
        onChange={() => {}} // Added no-op onChange
        className="flex-1 h-8 text-xs min-w-20"
      />
    );
  }

  let inputComponent: React.ReactNode | undefined;

  switch (fieldType) {
    case 'string':
    case 'email':
      inputComponent = (
        <ShadowInput
          id={id}
          value={String(value ?? '')}
          onChange={onChange}
          className="flex-1 h-8 text-xs min-w-20"
        />
      );
      break;
    case 'int':
      inputComponent = (
        <ShadowInput
          id={id}
          type="number"
          value={String(value ?? '')}
          onChange={(val) => onChange(Number(val))} // Convert to number
          className="flex-1 h-8 text-xs min-w-20"
        />
      );
      break;
    case 'bool':
      inputComponent = (
        <RadioGroup
          value={value === true ? "true" : value === false ? "false" : NONE_VALUE}
          onValueChange={(val) => onChange(val === "true" ? true : val === "false" ? false : undefined)}
          className="flex space-x-4" // Simple flex layout for radio buttons
        >
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="true" id={`${id}-true`} />
            <label htmlFor={`${id}-true`}>True</label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="false" id={`${id}-false`} />
            <label htmlFor={`${id}-false`}>False</label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value={NONE_VALUE} id={`${id}-none`} />
            <label htmlFor={`${id}-none`}>None</label>
          </div>
        </RadioGroup>
      );
      break;
    case 'select':
    case 'radio':
    case 'checkbox':
      inputComponent = (
        <Select
          value={value === undefined || value === null ? NONE_VALUE : String(value)}
          onValueChange={(val) => onChange(val === NONE_VALUE ? undefined : val)}
        >
          <SelectTrigger className="flex-1 h-8 text-xs min-w-20 border border-input bg-background shadow-[1px_1px_0_0_var(--color-input)] focus-within:shadow-[2px_2px_0_0_var(--color-input)] transition-all duration-300 ease-out">
            <SelectValue placeholder="Select a value" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value={NONE_VALUE}>None</SelectItem>
            {options?.map(option => (
              <SelectItem key={option.value} value={option.value}>
                {option.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      );
      break;
    default:
      inputComponent = (
        <ShadowInput
          id={id}
          value={String(value ?? '')}
          onChange={onChange}
          className="flex-1 h-8 text-xs min-w-20"
        />
      );
  }

  return inputComponent;
};