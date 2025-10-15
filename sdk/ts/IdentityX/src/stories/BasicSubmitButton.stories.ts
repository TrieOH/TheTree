import type { Meta, StoryObj } from '@storybook/react-vite';
import BasicSubmitButton from '../next/components/Form/BasicSubmitButton';

const meta = {
  title: "Example/BasicSubmitButton",
  component: BasicSubmitButton,
  parameters: {
    // Optional parameter to center the component in the Canvas. More info: https://storybook.js.org/docs/configure/story-layout
    layout: 'centered',
  },
  tags: ['autodocs'],
} satisfies Meta<typeof BasicSubmitButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const DefaultSubmitButton: Story = {
  args: {
    label: "Enviar...",
  }
};