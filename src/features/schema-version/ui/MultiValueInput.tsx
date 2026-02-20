import type React from 'react';
import { useState, useRef } from 'react';
import { ShadowInput } from '@/shared/ui/form/ShadowInput';
import { Badge } from '@/shared/ui/shadcn/badge';
import { XIcon } from 'lucide-react';
import { cn } from '@/shared/lib/utils';

interface MultiValueInputProps {
  id: string;
  value?: (string | number)[] | null;
  onChange: (value: (string | number)[] | undefined) => void;
  inputType: 'text' | 'number';
  placeholder?: string;
  className?: string;
}

export const MultiValueInput: React.FC<MultiValueInputProps> = ({
  id,
  value,
  onChange,
  inputType,
  placeholder = "Add values (comma-separated)",
  className,
}) => {
  const [inputValue, setInputValue] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  const currentValues = Array.isArray(value) ? value : [];

  const handleAddValue = (e: React.KeyboardEvent<HTMLInputElement> | React.FocusEvent<HTMLInputElement>) => {
    if (e.type === 'blur' || (e.type === 'keydown' && (e as React.KeyboardEvent).key === 'Enter')) {
      e.preventDefault();
      const trimmedInput = inputValue.trim();
      if (trimmedInput) {
        const newValues = trimmedInput.split(',').map(s => s.trim()).filter(Boolean);

        let processedValues: (string | number)[] = [];
        if (inputType === 'number') {
          processedValues = newValues
            .map(val => Number(val))
            .filter(num => !Number.isNaN(num));
        } else {
          processedValues = newValues;
        }

        const uniqueValues = Array.from(new Set([...currentValues, ...processedValues]));
        onChange(uniqueValues.length > 0 ? uniqueValues : undefined);
        setInputValue('');
      }
    }
  };

  const handleRemoveValue = (valueToRemove: string | number) => {
    let updatedValues: (string | number)[];
    if (inputType === 'number' && typeof valueToRemove === 'number' && Number.isNaN(valueToRemove)) {
      updatedValues = currentValues.filter(val => !(typeof val === 'number' && Number.isNaN(val)));
    } else {
      updatedValues = currentValues.filter(val => val !== valueToRemove);
    }
    onChange(updatedValues.length > 0 ? updatedValues : undefined);
  };



  return (
    <div className={cn("flex flex-col gap-1", className)}>
      <ShadowInput
        inputRef={inputRef}
        id={id}
        type={inputType === 'number' ? 'text' : inputType}
        placeholder={placeholder}
        value={inputValue}
        onChange={(value) => setInputValue(value)}
        onKeyDown={handleAddValue}
        onBlur={handleAddValue}
        className="h-8 text-xs min-w-40"
      />
      {(currentValues.length > 0) && (
        <div className="flex flex-wrap items-center gap-1 p-1 border rounded-md min-w-40 max-h-28 overflow-y-auto"> {/* Badges container */}
          {currentValues.map((val, _index) => (
            <Badge key={String(val)} variant="secondary" className="pr-1">
              {String(val)}
              <button
                type="button"
                onClick={() => handleRemoveValue(val)}
                className="ml-1 rounded-full p-0.5 hover:bg-secondary-foreground/20"
              >
                <XIcon className="h-3 w-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}
    </div>
  );
};