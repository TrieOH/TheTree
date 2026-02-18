import type { Meta, StoryObj } from "@storybook/react";
import TabbedSignUpWithProvider from "./components/TabbedSignUpWithProvider";
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
      description: "Array of flow IDs with labels to display as tabs",
    },
  },
  args: {
    flowIds: [
      { label: "Email Flow", value: "email_flow_id_1" },
      { label: "Social Flow", value: "social_flow_id_2" },
      { label: "Partner Flow", value: "partner_flow_id_3" },
      { label: "Partner Flow", value: "partner_flow_id_4" },
      { label: "Partner Flow", value: "partner_flow_id5" },
    ],
  },
} satisfies Meta<typeof TabbedSignUpWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const CustomFlowIds: Story = {
  args: {
    flowIds: [
      { label: "Type A", value: "type_a" },
      { label: "Type B", value: "type_b" },
      { label: "Type C", value: "type_c" },
    ],
  },
};


export const WithLogin: Story = {
  args: {
    loginRedirect: () => {}
  },
};