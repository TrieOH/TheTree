import { useEffect, useState } from 'react';
import type { VersionField } from '../model/types';
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { ShadowInput } from '@/shared/ui/form/ShadowInput';
import { FieldTypeSelector } from './FieldTypeSelector';

interface FieldEditFormProps {
  field: VersionField;
  onSave: (updatedField: VersionField) => void;
  onCancel: () => void;
}

export const FieldEditForm: React.FC<FieldEditFormProps> = ({ field, onSave, onCancel }) => {
  const [editedField, setEditedField] = useState<VersionField>(field);

  useEffect(() => {
    setEditedField(field);
  }, [field]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(editedField);
  };

  return (
    <form onSubmit={handleSubmit} className="p-4 space-y-4">
      <h3 className="text-lg font-semibold">Edit Field</h3>
      <div>
        <label htmlFor="title" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Title</label>
        <ShadowInput
          id="title"
          value={editedField.title}
          onChange={(value) => setEditedField((prev) => ({ ...prev, title: value }))}
        />
      </div>
      <div>
        <label htmlFor="key" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Key</label>
        <ShadowInput
          id="key"
          value={editedField.key}
          onChange={(value) => setEditedField((prev) => ({ ...prev, key: value }))}
        />
      </div>
      <div>
        <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 mb-2">Type</label>
        <FieldTypeSelector selectedType={editedField.type} onSelectType={(type) => setEditedField((prev) => ({ ...prev, type }))} />
      </div>

      <div className="flex justify-end gap-2">
        <ShadowButton type="button" variant="ghost" onClick={onCancel} value="Cancel" />
        <ShadowButton type="submit" value="Save" />
      </div>
    </form>
  );
};
