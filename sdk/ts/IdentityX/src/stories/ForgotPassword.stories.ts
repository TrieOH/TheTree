import type { Meta, StoryObj } from '@storybook/react';
import ForgotPasswordWithProvider from './components/ForgotPasswordWithProvider';

const meta = {
  title: 'Authentication/ForgotPassword',
  component: ForgotPasswordWithProvider,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
} satisfies Meta<typeof ForgotPasswordWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithLogin: Story = {
  args: {
    loginRedirect: () => {},
  }
};
