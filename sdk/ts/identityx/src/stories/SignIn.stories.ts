import type { Meta, StoryObj } from '@storybook/react-vite';
import SignInWithProvider from './components/SignInWithProvider';

const meta = {
  title: "Authentication/SignIn",
  component: SignInWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {
    isProjectMode: {
      control: 'boolean',
      description: 'Whether the project is in Project Mode (requires projectId) or Auth Mode',
    }
  }
} satisfies Meta<typeof SignInWithProvider>;

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

export const WithForgotPassword: Story = {
  args: {
    forgotPasswordRedirect: () => {},
    isProjectMode: true,
  }
};

export const WithSignUp: Story = {
  args: {
    signUpRedirect: () => {},
    isProjectMode: true,
  }
};

export const WithAll: Story = {
  args: {
    forgotPasswordRedirect: () => {},
    signUpRedirect: () => {},
    isProjectMode: true,
  }
};
