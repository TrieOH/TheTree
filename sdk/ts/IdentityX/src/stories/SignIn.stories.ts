import type { Meta, StoryObj } from '@storybook/react-vite';
import { fn } from 'storybook/internal/test';
import { SignIn } from '../next';

const meta = {
  title: "Example/SignIn",
  component: SignIn,
  parameters: {
    // Optional parameter to center the component in the Canvas. More info: https://storybook.js.org/docs/configure/story-layout
    layout: 'centered',
  },
  tags: ['autodocs'],
  args: { onSubmit: fn() },
} satisfies Meta<typeof SignIn>;

export default meta;
type Story = StoryObj<typeof meta>;
export const Primary: Story = {
  args: {
    
  }
};