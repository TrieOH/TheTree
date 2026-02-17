import type React from 'react';
import type { AnyFieldApi } from '@tanstack/react-form';
import { ShadowInput } from '@/shared/ui/form/ShadowInput';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/shadcn/select';
import { RadioGroup, RadioGroupItem } from '@/shared/ui/shadcn/radio-group';
import type { Option, VersionFieldResult } from '../model/types';

interface DefaultValueInputProps {
  field: AnyFieldApi;
  fieldType: VersionFieldResult['type'];
  options?: Option[];
}

export const DefaultValueInput: React.FC<DefaultValueInputProps> = ({ field, fieldType, options }) => {
  let defaultValueInput: React.ReactNode | undefined;

  switch (fieldType) {
    case 'string':
    case 'email':
      defaultValueInput = (
        <ShadowInput
          id={field.name}
          value={String(field.state.value ?? '')}
          onBlur={field.handleBlur}
          onChange={field.handleChange}
        />
      );
      break;
    case 'int':
      defaultValueInput = (
        <ShadowInput
          id={field.name}
          type="number"
          value={String(field.state.value ?? '')}
          onBlur={field.handleBlur}
          onChange={field.handleChange}
        />
      );
      break;
    case 'bool':
      defaultValueInput = (
        <RadioGroup
          value={field.state.value === true ? "true" : "false"}
          onValueChange={(value) => field.handleChange(value === "true")}
          className="flex justify-center gap-4 mt-2"
        >
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="true" id={`${field.name}-true`} />
            <label htmlFor={`${field.name}-true`}>True</label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="false" id={`${field.name}-false`} />
            <label htmlFor={`${field.name}-false`}>False</label>
          </div>
        </RadioGroup>
      );
      break;
    case 'select':
    case 'radio':
    case 'checkbox':
      defaultValueInput = (
        <Select
          value={String(field.state.value)}
          onValueChange={(value) => field.handleChange(value)}
        >
          <SelectTrigger className="flex items-center w-full rounded-sm border border-input bg-background shadow-[1px_1px_0_0_var(--color-input)] focus-within:shadow-[2px_2px_0_0_var(--color-input)] transition-all duration-300 ease-out h-9 px-3 text-sm">
            <SelectValue placeholder="Select a default option" />
          </SelectTrigger>
          <SelectContent>
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
      defaultValueInput = (
        <ShadowInput
          id={field.name}
          value={String(field.state.value ?? '')}
          onBlur={field.handleBlur}
          onChange={field.handleChange}
        />
      );
  }

  return defaultValueInput;
};
