import type { Meta, StoryObj } from "@storybook/react";
import TabbedSignUpWithProvider from "./components/TabbedSignUpWithProvider";
import { MOCK_FIELDS, createMockField } from "./mocks/field-mocks";

const meta = {
  title: "Authentication/TabbedSignUp",
  component: TabbedSignUpWithProvider,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  argTypes: {
    flowIds: {
      control: "object",
      description: "Array of flow IDs with labels and fields to display as tabs",
    },
  },
  args: {
    flowIds: [
      { 
        label: "Básico", 
        value: "basic_flow",
        fields: [
          createMockField("f1", "first_name", "Nome"),
          createMockField("f2", "last_name", "Sobrenome"),
          MOCK_FIELDS.COUNTRY,
        ]
      },
      { 
        label: "Perfil", 
        value: "profile_flow",
        fields: [
          MOCK_FIELDS.AGE,
          MOCK_FIELDS.GENDER,
          MOCK_FIELDS.BIO("f_age"),
        ]
      },
      { 
        label: "Configurações", 
        value: "settings_flow",
        fields: [
          MOCK_FIELDS.USER_TYPE,
          MOCK_FIELDS.COMPANY_NAME("f_utype"),
          MOCK_FIELDS.INTERESTS,
          MOCK_FIELDS.NEWSLETTER,
        ]
      },
    ],
  },
} satisfies Meta<typeof TabbedSignUpWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const BusinessOnboarding: Story = {
  args: {
    flowIds: [
      {
        label: "Dados Pessoais",
        value: "personal",
        fields: [
          createMockField("p1", "full_name", "Nome Completo"),
          MOCK_FIELDS.COUNTRY,
          MOCK_FIELDS.CPF("f_country"),
        ]
      },
      {
        label: "Dados Profissionais",
        value: "professional",
        fields: [
          MOCK_FIELDS.USER_TYPE,
          MOCK_FIELDS.COMPANY_NAME("f_utype"),
          MOCK_FIELDS.ZIPCODE("f_country"),
        ]
      }
    ]
  }
};

export const ManyTabs: Story = {
  args: {
    flowIds: [
      { label: "Individual", value: "individual", fields: [createMockField("f5", "nickname", "Nickname")] },
      { label: "Developer", value: "developer", fields: [createMockField("f6", "github", "GitHub Username")] },
      { label: "Company", value: "company", fields: [createMockField("f7", "org", "Organization")] },
      { label: "Non-Profit", value: "non_profit", fields: [createMockField("f8", "cause", "Cause")] },
    ],
  },
};
