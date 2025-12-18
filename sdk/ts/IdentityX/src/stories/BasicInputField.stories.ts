import type { Meta, StoryObj } from '@storybook/react-vite';
import BasicInputField from '../react/components/Form/BasicInputField';

const meta = {
  title: "Example/BasicInputField",
  component: BasicInputField,
  parameters: {
    // Optional parameter to center the component in the Canvas. More info: https://storybook.js.org/docs/configure/story-layout
    layout: 'centered',
  },
  argTypes: {
    type: { 
      control: "select",
      options: ['text', 'number', 'email', 'password']
    },
  },
  args: {
    type: "text"
  },
  tags: ['autodocs'],
} satisfies Meta<typeof BasicInputField>;

export default meta;
type Story = StoryObj<typeof meta>;

export const NameInput: Story = {
  args: {
    name: "name",
    label: "Nome",
    placeholder: "Your Name...",
    autoComplete: "Nome"
  }
};

export const PasswordInput: Story = {
  args: {
    name: "password",
    label: "Senha",
    placeholder: "**********",
    autoComplete: "password",
    type: "password",
  }
};