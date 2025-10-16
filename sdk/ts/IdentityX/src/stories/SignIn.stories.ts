import type { Meta, StoryObj } from '@storybook/react-vite';
import SignInWithProvider from './SignInWithProvider';

const meta = {
  title: "Example/SignIn",
  component: SignInWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  // args: { onSubmit: fn() },
} satisfies Meta<typeof SignInWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;
export const Primary: Story = {};