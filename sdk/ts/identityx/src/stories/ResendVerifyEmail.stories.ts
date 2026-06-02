import type { Meta, StoryObj } from '@storybook/react-vite';
import ResendVerifyEmailWithProvider from './components/ResendVerifyEmailWithProvider';

const meta = {
  title: "Authentication/ResendVerifyEmail",
  component: ResendVerifyEmailWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {
    isProjectMode: {
      control: 'boolean',
      description: 'Whether the project is in Project Mode (requires projectId) or Auth Mode',
    }
  }
} satisfies Meta<typeof ResendVerifyEmailWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    isProjectMode: true,
  }
};
