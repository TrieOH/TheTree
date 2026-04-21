import type { Meta, StoryObj } from '@storybook/react-vite';
import SignUpWithProvider from './components/SignUpWithProvider';

const meta = {
  title: "Authentication/SignUp",
  component: SignUpWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {
    isProjectMode: {
      control: 'boolean',
      description: 'Whether the project is in Project Mode (requires projectId) or Auth Mode',
    }
  }
} satisfies Meta<typeof SignUpWithProvider>;

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

export const WithLoginLink: Story = {
  args: {
    loginRedirect: () => alert("Redirect to Login"),
    isProjectMode: true,
  },
};
