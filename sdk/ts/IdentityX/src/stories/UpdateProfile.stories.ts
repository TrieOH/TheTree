import type { Meta, StoryObj } from '@storybook/react-vite';
import UpdateProfileWithProvider from './components/UpdateProfileWithProvider';
import { MOCK_FIELDS } from './mocks/field-mocks';

const meta = {
  title: "Authentication/UpdateProfile",
  component: UpdateProfileWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {},
  args: {},
} satisfies Meta<typeof UpdateProfileWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    fields: [
      MOCK_FIELDS.AGE,
      MOCK_FIELDS.GENDER,
      MOCK_FIELDS.NEWSLETTER,
    ],
    initialValues: {
      age: 25,
      gender: 'm',
      newsletter: true,
    }
  }
};

export const CompletingData: Story = {
  parameters: {
    docs: {
      description: {
        story: 'Usado quando o usuário já tem conta mas precisa preencher dados adicionais obrigatórios.',
      },
    },
  },
  args: {
    fields: [
      MOCK_FIELDS.COUNTRY,
      MOCK_FIELDS.ZIPCODE("f_country"),
      MOCK_FIELDS.CPF("f_country"),
    ],
    onSuccess: async () => {
      alert("Perfil atualizado com sucesso!");
    }
  }
};

export const ComplexRules: Story = {
  parameters: {
    docs: {
      description: {
        story: 'Demonstração de regras de visibilidade e obrigatoriedade dinâmica.',
      },
    },
  },
  args: {
    fields: [
      MOCK_FIELDS.USER_TYPE,
      MOCK_FIELDS.COMPANY_NAME("f_utype"),
      MOCK_FIELDS.AGE,
      MOCK_FIELDS.BIO("f_age"),
      MOCK_FIELDS.INTERESTS,
    ]
  }
};
