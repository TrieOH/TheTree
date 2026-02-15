import type React from 'react';
import { ShadowInput } from '@/shared/ui/form/ShadowInput';
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { PlusIcon, TrashIcon } from 'lucide-react';

interface Option {
  id: string;
  label: string;
  value: string;
  position: number;
}

interface OptionsEditorProps {
  options: Option[];
  onChange: (options: Option[]) => void;
}

export const OptionsEditor: React.FC<OptionsEditorProps> = ({ options, onChange }) => {
  const handleOptionChange = (index: number, key: keyof Option, value: string | number) => {
    const newOptions = [...options];
    if (key === 'position' && typeof value === 'string') {
      newOptions[index] = { ...newOptions[index], [key]: parseInt(value, 10) };
    } else newOptions[index] = { ...newOptions[index], [key]: value };
    
    onChange(newOptions);
  };

  const handleAddOption = () => {
    const newOptions = [...options, { id: crypto.randomUUID(), label: '', value: '', position: options.length }];
    onChange(newOptions);
  };

  const handleRemoveOption = (index: number) => {
    const newOptions = options.filter((_, i) => i !== index).map((opt, i) => ({ ...opt, position: i }));
    onChange(newOptions);
  };

  return (
    <div className="space-y-2">
      {options.map((option, index) => (
        <div key={option.id} className="flex items-center gap-2">
          <ShadowInput
            placeholder="Label"
            value={option.label}
            onChange={(value) => handleOptionChange(index, 'label', value)}
            className="flex-1"
          />
          <ShadowInput
            placeholder="Value"
            value={option.value}
            onChange={(value) => handleOptionChange(index, 'value', value)}
            className="flex-1"
          />
          <ShadowButton
            type="button"
            variant="destructive"
            leftIcon={<TrashIcon className="h-4 w-4" />}
            onClick={() => handleRemoveOption(index)}
          />
        </div>
      ))}
      <ShadowButton 
        type="button"
        variant="ghost"
        onClick={handleAddOption} 
        className="w-full h-8 text-xs text-muted-foreground hover:text-foreground border border-dashed border-muted-foreground/30 hover:border-muted-foreground/60"
        leftIcon={<PlusIcon className="h-4 w-4" />}
        value='Add Option'
      />
    </div>
  );
};
