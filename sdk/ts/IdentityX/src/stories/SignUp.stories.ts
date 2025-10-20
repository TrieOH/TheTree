import type { Meta, StoryObj } from '@storybook/react-vite';
import SignUpWithProvider from './components/SignUpWithProvider';

const meta = {
  title: "Example/SignUp",
  component: SignUpWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
} satisfies Meta<typeof SignUpWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;
export const Default: Story = {};