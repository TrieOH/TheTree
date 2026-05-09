import type { Meta, StoryObj } from '@storybook/react-vite';
import ResetPasswordWithProvider from './components/ResetPasswordWithProvider';

const meta = {
  title: "Authentication/ResetPassword",
  component: ResetPasswordWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {
    isProjectMode: {
      control: 'boolean',
      description: 'Whether the project is in Project Mode (requires projectId) or Auth Mode',
    }
  }
} satisfies Meta<typeof ResetPasswordWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    token: 'fake-token',
    isProjectMode: true,
  }
};

export const WithLoginRedirect: Story = {
  args: {
    token: 'fake-token',
    loginRedirect: () => alert('Redirect to login'),
    isProjectMode: true,
  }
};
