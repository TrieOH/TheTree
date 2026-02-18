import type { Meta, StoryObj } from '@storybook/react-vite';
import SignUpWithProvider from './components/SignUpWithProvider';

const meta = {
  title: "Authentication/SignUp",
  component: SignUpWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {
    flow_id: { control: 'text' },
  },
  args: {
    flow_id: 'default',
  },
} satisfies Meta<typeof SignUpWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    flow_id: "test"
  }
};

export const WithLogin: Story = {
  args: {
    flow_id: 'test',
    loginRedirect: () => {}
  },
};