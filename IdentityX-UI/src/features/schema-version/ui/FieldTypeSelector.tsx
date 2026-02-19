import { cn } from '@/shared/lib/utils';
import { getFieldTypeIcon } from '../model/field-type-to-icon';
import type { VersionFieldType } from '../model/types';

interface FieldTypeSelectorProps {
  selectedType: VersionFieldType;
  onSelectType: (type: VersionFieldType) => void;
}

const fieldTypes: VersionFieldType[] = ["string", "email", "int", "select", "radio", "checkbox", "bool"];

export const FieldTypeSelector: React.FC<FieldTypeSelectorProps> = ({ selectedType, onSelectType }) => {
  return (
    <div className="grid grid-cols-2 gap-1 max-h-50 overflow-y-auto p-1">
      {fieldTypes.map((type) => {
        const Icon = getFieldTypeIcon(type);
        const isSelected = selectedType === type;
        return (
          <button
            type='button'
            key={type}
            className={cn(
              "flex flex-col items-center justify-center p-2 border rounded-md cursor-pointer h-17.5",
              "transition-colors duration-200",
              isSelected ? "border-primary bg-primary/10 text-primary" : "border-border bg-card hover:bg-muted/50"
            )}
            onClick={() => onSelectType(type)}
          >
            <Icon className="w-6 h-6 mb-1" />
            <span className="text-xs font-medium capitalize">{type === "int" ? "number" : type}</span>
          </button>
        );
      })}
    </div>
  );
};
