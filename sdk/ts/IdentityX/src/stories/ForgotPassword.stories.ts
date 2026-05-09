import type { Meta, StoryObj } from '@storybook/react';
import ForgotPasswordWithProvider from './components/ForgotPasswordWithProvider';

const meta = {
  title: 'Authentication/ForgotPassword',
  component: ForgotPasswordWithProvider,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    isProjectMode: {
      control: 'boolean',
      description: 'Whether the project is in Project Mode (requires projectId) or Auth Mode',
    }
  }
} satisfies Meta<typeof ForgotPasswordWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const ProjectMode: Story = {
  args: {
    isProjectMode: true,
  }
};

export const AuthMode: Story = {
  args: {
    isProjectMode: false,
  }
};

export const WithLogin: Story = {
  args: {
    loginRedirect: () => {},
    isProjectMode: true,
  }
};
