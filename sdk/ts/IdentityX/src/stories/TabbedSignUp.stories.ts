import type { Meta, StoryObj } from "@storybook/react";
import TabbedSignUpWithProvider from "./components/TabbedSignUpWithProvider";
import type { FieldDefinitionResultI } from "../types/fields-types";

const mockField = (id: string, key: string, title: string, type: FieldDefinitionResultI['type'] = 'string'): FieldDefinitionResultI => ({
  id,
  object_id: "obj_1",
  key,
  title,
  type,
  placeholder: `Enter your ${title.toLowerCase()}`,
  description: "",
  position: 0,
  options: [],
  default_value: "",
  mutable: true,
  required: true,
  owner: "user",
  visibility_rules: [],
  required_rules: [],
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
});

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
        label: "Personal", 
        value: "personal_flow",
        fields: [
          mockField("f1", "first_name", "First Name"),
          mockField("f2", "last_name", "Last Name"),
        ]
      },
      { 
        label: "Business", 
        value: "business_flow",
        fields: [
          mockField("f3", "company_name", "Company Name"),
          mockField("f4", "tax_id", "Tax ID"),
        ]
      },
    ],
  },
} satisfies Meta<typeof TabbedSignUpWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const ManyTabs: Story = {
  args: {
    flowIds: [
      { label: "Individual", value: "individual", fields: [mockField("f5", "nickname", "Nickname")] },
      { label: "Developer", value: "developer", fields: [mockField("f6", "github", "GitHub Username")] },
      { label: "Company", value: "company", fields: [mockField("f7", "org", "Organization")] },
      { label: "Non-Profit", value: "non_profit", fields: [mockField("f8", "cause", "Cause")] },
    ],
  },
};

export const WithLoginRedirect: Story = {
  args: {
    loginRedirect: () => alert("Redirect to login"),
  },
};