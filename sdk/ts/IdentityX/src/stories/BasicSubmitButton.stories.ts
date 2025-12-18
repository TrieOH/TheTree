import type { Meta, StoryObj } from '@storybook/react-vite';
import BasicSubmitButton from '../react/components/Form/BasicSubmitButton';
import { fn } from 'storybook/test';

const meta = {
  title: "Example/BasicSubmitButton",
  component: BasicSubmitButton,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  args: { onSubmit: fn() },
} satisfies Meta<typeof BasicSubmitButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const DefaultSubmitButton: Story = {
  args: {
    label: "Enviar...",
    loading: false,
  }
};

export const Loading: Story = {
  args: {
    label: "Enviando...",
    loading: true,
  }
};