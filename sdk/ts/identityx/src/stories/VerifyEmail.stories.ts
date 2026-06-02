import type { Meta, StoryObj } from '@storybook/react-vite';
import VerifyEmailWithProvider from './components/VerifyEmailWithProvider';

const meta = {
  title: "Authentication/VerifyEmail",
  component: VerifyEmailWithProvider,
  parameters: {
    layout: 'fullscreen',
  },
  tags: ['autodocs'],
  argTypes: {
    isProjectMode: {
      control: 'boolean',
      description: 'Whether the project is in Project Mode (requires projectId) or Auth Mode',
    },
    mockState: {
      control: 'select',
      options: ['verifying', 'success', 'error', 'already_verified'],
      description: 'Force a specific state for visual testing',
    }
  }
} satisfies Meta<typeof VerifyEmailWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Verifying: Story = {
  args: {
    token: 'fake-token',
    mockState: 'verifying',
  }
};

export const Success: Story = {
  args: {
    token: 'fake-token',
    mockState: 'success',
  }
};

export const AlreadyVerified: Story = {
  args: {
    token: 'fake-token',
    mockState: 'already_verified',
  }
};

export const Error: Story = {
  args: {
    token: 'fake-token',
    mockState: 'error',
  }
};

export const RealBehavior: Story = {
  args: {
    token: 'real-token-simulation',
  }
};
