import type { Meta, StoryObj } from '@storybook/react-vite';
import SignInWithProvider from './components/SignInWithProvider';

const meta = {
  title: "Authentication/SignIn",
  component: SignInWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
} satisfies Meta<typeof SignInWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithForgotPassword: Story = {
  args: {
    forgotPasswordRedirect: () => {},
  }
};

export const WithSignUp: Story = {
  args: {
    signUpRedirect: () => {},
  }
};
export const WithAll: Story = {
  args: {
    forgotPasswordRedirect: () => {},
    signUpRedirect: () => {},
  }
};