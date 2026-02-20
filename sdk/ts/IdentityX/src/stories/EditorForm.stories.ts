import type { Meta, StoryObj } from '@storybook/react-vite';
import { MOCK_FIELDS } from './mocks/field-mocks';
import { EditorForm } from '../react';

const meta = {
  title: "Common/EditorForm",
  component: EditorForm,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
} satisfies Meta<typeof EditorForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    title: "Editor Dinâmico",
    description: "Este formulário demonstra regras de visibilidade baseadas em outros campos.",
    fields: [
      MOCK_FIELDS.USER_TYPE,
      MOCK_FIELDS.COMPANY_NAME("f_utype"),
      MOCK_FIELDS.AGE,
      MOCK_FIELDS.BIO("f_age"),
      MOCK_FIELDS.INTERESTS,
    ]
  }
};

export const WithoutHeader: Story = {
  args: {
    fields: [
      MOCK_FIELDS.AGE,
      MOCK_FIELDS.GENDER,
    ],
  }
};

export const NoFields: Story = {
  args: {
    title: "Empty Form",
    description: "This form has no fields currently.",
    fields: [],
    noFieldsMessage: "Custom empty message: No fields found."
  }
};
