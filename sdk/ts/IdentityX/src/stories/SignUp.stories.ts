import type { Meta, StoryObj } from '@storybook/react-vite';
import SignUpWithProvider from './components/SignUpWithProvider';
import { MOCK_FIELDS } from './mocks/field-mocks';

const meta = {
  title: "Authentication/SignUp",
  component: SignUpWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {
    flow_id: { control: 'text' },
  },
  args: {
    flow_id: 'default',
  },
} satisfies Meta<typeof SignUpWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    flow_id: "test"
  }
};

export const WithLoginLink: Story = {
  args: {
    flow_id: 'test',
    loginRedirect: () => alert("Redirect to Login")
  },
};

export const FullProfile: Story = {
  args: {
    flow_id: 'full-profile',
    fields: [
      MOCK_FIELDS.USER_TYPE,
      MOCK_FIELDS.COMPANY_NAME("f_utype"),
      MOCK_FIELDS.COUNTRY,
      MOCK_FIELDS.ZIPCODE("f_country"),
      MOCK_FIELDS.CPF("f_country"),
      MOCK_FIELDS.AGE,
      MOCK_FIELDS.BIO("f_age"),
      MOCK_FIELDS.GENDER,
      MOCK_FIELDS.INTERESTS,
      MOCK_FIELDS.NEWSLETTER,
    ]
  }
};

export const MultiStepSimulation: Story = {
  parameters: {
    docs: {
      description: {
        story: 'Simulação de regras complexas: O campo Bio só aparece para maiores de 18, e o campo Nome da Empresa só aparece para Business.',
      },
    },
  },
  args: {
    flow_id: 'complex-rules',
    fields: [
      MOCK_FIELDS.AGE,
      MOCK_FIELDS.BIO("f_age"),
      MOCK_FIELDS.USER_TYPE,
      MOCK_FIELDS.COMPANY_NAME("f_utype"),
    ]
  }
};

export const ConditionalRequirements: Story = {
  parameters: {
    docs: {
      description: {
        story: 'Regras de obrigatoriedade: CEP é obrigatório para Brasil e EUA. CPF é obrigatório apenas para Brasil.',
      },
    },
  },
  args: {
    flow_id: 'cond-req',
    fields: [
      MOCK_FIELDS.COUNTRY,
      MOCK_FIELDS.ZIPCODE("f_country"),
      MOCK_FIELDS.CPF("f_country"),
    ]
  }
};
