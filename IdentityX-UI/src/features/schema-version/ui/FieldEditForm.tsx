import type { VersionField } from '../model/types';
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { ShadowInput } from '@/shared/ui/form/ShadowInput';
import { FieldTypeSelector } from './FieldTypeSelector';
import { useForm } from '@tanstack/react-form';
import { FormField } from '@/shared/ui/form/FormField';
import { OptionsEditor } from './OptionsEditor';
import { versionFieldSchema } from '../model/types';
import { useEffect } from 'react';
import { Checkbox } from '@/shared/ui/shadcn/checkbox';

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
        <FormField<VersionField, 'placeholder'> name="placeholder" label="Placeholder" form={form}>
          {(field) => (
            <ShadowInput
              id={field.name}
              value={field.state.value ?? ''}
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
        
        <div className='flex items-center justify-between px-4'>
          <FormField<VersionField, 'required'> name="required" label="" form={form}>
            {(field) => (
              <div className="flex items-center space-x-2">
                <Checkbox
                  id={field.name}
                  checked={field.state.value}
                  onCheckedChange={field.handleChange}
                />
                <label
                  htmlFor={field.name}
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Required
                </label>
              </div>
            )}
          </FormField>
          <FormField<VersionField, 'mutable'> name="mutable" label="" form={form}>
            {(field) => (
              <div className="flex items-center space-x-2">
                <Checkbox
                  id={field.name}
                  checked={field.state.value}
                  onCheckedChange={field.handleChange}
                />
                <label
                  htmlFor={field.name}
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Mutable
                </label>
              </div>
            )}
          </FormField>
        </div>
        <FormField<VersionField, 'type'> name="type" label="Type" form={form}>
          {(field) => (
            <FieldTypeSelector
              selectedType={field.state.value}
              onSelectType={field.handleChange}
            />
          )}
        </FormField>
        <form.Subscribe
          selector={(state) => state.values.type}
          children={(type) => {
            const showOptions = ['select', 'radio', 'checkbox'].includes(type);
            return showOptions ? (
              <FormField<VersionField, 'options'> name="options" label="Options" form={form}>
                {(field) => (
                  <OptionsEditor
                    options={field.state.value || []}
                    onChange={field.handleChange}
                  />
                )}
              </FormField> 
            ) : null;
          }}
        />
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
