import type { VersionField } from '../model/types';
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { ShadowInput } from '@/shared/ui/form/ShadowInput';
import { FieldTypeSelector } from './FieldTypeSelector';
import { useForm } from '@tanstack/react-form';
import { FormField } from '@/shared/ui/form/FormField';
import { useEffect } from 'react';
import { versionFieldSchema } from '../model/types';

interface FieldEditFormProps {
  field: VersionField;
  onSave: (updatedField: VersionField) => void;
  onCancel: () => void;
}

export const FieldEditForm: React.FC<FieldEditFormProps> = ({ field, onSave, onCancel }) => {
  const form = useForm({
    defaultValues: field,
    onSubmit: async ({ value }) => {
      onSave(value);
    },
    validators: {
      onChange: versionFieldSchema,
    }
  });

  useEffect(() => {
    form.reset(field);
  }, [field, form.reset]);

  return (
    <form
      onSubmit={(e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
        className="p-4 space-y-4"
      >
        <FormField<VersionField, 'title'> name="title" label="Title" form={form}>
          {(field) => (
            <ShadowInput
              id={field.name}
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={field.handleChange}
            />
          )}
        </FormField>
        <FormField<VersionField, 'key'> name="key" label="Key" form={form}>
          {(field) => (
            <ShadowInput
              id={field.name}
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={field.handleChange}
            />
          )}
        </FormField>
        <FormField<VersionField, 'default_value'> name="default_value" label="Default Value" form={form}>
          {(field) => (
            <ShadowInput
              id={field.name}
              value={String(field.state.value ?? '')}
              onBlur={field.handleBlur}
              onChange={field.handleChange}
            />
          )}
        </FormField>
        <FormField<VersionField, 'type'> name="type" label="Type" form={form}>
          {(field) => (
            <FieldTypeSelector
              selectedType={field.state.value}
              onSelectType={field.handleChange}
            />
          )}
        </FormField>            
        <div className="flex justify-end gap-2">
          <ShadowButton type="button" variant="ghost" onClick={onCancel} value="Cancel" />
          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
            children={([canSubmit, isSubmitting]) => (
              <ShadowButton type="submit" value="Save" disabled={!canSubmit || isSubmitting} />
            )}
          />
        </div>
      </form>
  );
};
