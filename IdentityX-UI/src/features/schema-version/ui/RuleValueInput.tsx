import type React from 'react';
import { ShadowInput } from '@/shared/ui/form/ShadowInput';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/shadcn/select';
import { RadioGroup, RadioGroupItem } from '@/shared/ui/shadcn/radio-group';
import { MultiValueInput } from './MultiValueInput';
import { MultiSelectOptions } from './MultiSelectOptions';
import type { FieldDefinitionResultI, Operator, OptionResultI } from '../model/types';

interface RuleValueInputProps {
  value: unknown;
  onChange: (value: unknown) => void;
  fieldType: FieldDefinitionResultI['type'];
  options?: OptionResultI[];
  operator: Operator;
  id: string;
}

export const RuleValueInput: React.FC<RuleValueInputProps> = ({
  value,
  onChange,
  fieldType,
  options,
  operator,
  id,
}) => {
  const isValueIrrelevant = ['exists', 'not_exists'].includes(operator);

  if (isValueIrrelevant) {
    return (
      <ShadowInput
        id={id}
        value="N/A"
        disabled
        onChange={() => {}}
        className="w-full h-8 text-xs"
      />
    );
  }

  const renderShadowInput = (type: React.HTMLInputTypeAttribute = "text", placeholder: string = "Value...") => (
    <ShadowInput
      id={id}
      type={type}
      value={String(value ?? '')}
      onChange={
        type === "number"
          ? (value) => onChange(value === '' ? null : Number(value))
          : (value) => onChange(value)
      }
      placeholder={placeholder}
      className="w-full h-8 text-xs"
    />
  );

  switch (fieldType) {
    case 'string':
    case 'email':
      if (operator === 'contains') {
        return renderShadowInput("text", "Substring...");
      } else if (['in', 'not_in'].includes(operator)) {
        return (
          <MultiValueInput
            id={id}
            value={Array.isArray(value) && value.every(v => typeof v === 'string') ? value as string[] : undefined}
            onChange={(val) => onChange(val)}
            inputType="text"
            placeholder="value1, value2..."
            className="w-full"
          />
        );
      } else {
        return renderShadowInput("text", "Value...");
      }

    case 'int':
      if (['gt', 'gte', 'lt', 'lte'].includes(operator)) {
        return renderShadowInput("number", "Number...");
      } else if (['in', 'not_in'].includes(operator)) {
        return (
          <MultiValueInput
            id={id}
            value={Array.isArray(value) && value.every(v => typeof v === 'number') ? value as number[] : undefined}
            onChange={(val) => onChange(val)}
            inputType="number"
            placeholder="1, 2, 3..."
            className="w-full"
          />
        );
      } else {
        return renderShadowInput("number", "Number...");
      }

    case 'bool':
      return (
        <RadioGroup
          value={value === true ? "true" : value === false ? "false" : undefined}
          onValueChange={(val) => onChange(val === "true" ? true : val === "false" ? false : undefined)}
          className="flex flex-wrap gap-3 sm:gap-4 h-8 items-center"
        >
          <div className="flex items-center space-x-1.5">
            <RadioGroupItem value="true" id={`${id}-true`} className="h-4 w-4" />
            <label htmlFor={`${id}-true`} className="text-xs sm:text-sm cursor-pointer">True</label>
          </div>
          <div className="flex items-center space-x-1.5">
            <RadioGroupItem value="false" id={`${id}-false`} className="h-4 w-4" />
            <label htmlFor={`${id}-false`} className="text-xs sm:text-sm cursor-pointer">False</label>
          </div>
        </RadioGroup>
      );

    case 'select':
    case 'radio':
    case 'checkbox': {
      if (['in', 'not_in'].includes(operator)) {
        return (
          <MultiSelectOptions
            id={id}
            options={options || []}
            value={Array.isArray(value) && value.every(v => typeof v === 'string') ? value as string[] : undefined}
            onChange={(val) => onChange(val)}
            placeholder="Select options..."
            className="w-full"
          />
        );
      } else {
        return (
          <Select
            value={String(value ?? '')}
            onValueChange={(val) => onChange(val)}
          >
            <SelectTrigger className="w-full h-8 text-xs border border-input bg-background shadow-[1px_1px_0_0_var(--color-input)] focus-within:shadow-[2px_2px_0_0_var(--color-input)] transition-all duration-300 ease-out">
              <SelectValue placeholder="Select a value" />
            </SelectTrigger>
            <SelectContent>
              {options?.map(option => (
                <SelectItem key={option.value} value={option.value} className="text-xs">
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        );
      }
    }

    default:
      return renderShadowInput("text", "Value...");
  }
};