import type { Meta, StoryObj } from '@storybook/react-vite';
import BasicLogoutWithProvider from './components/BasicLogoutWithProvider';

const meta = {
  title: "Authentication/BasicLogoutButton",
  component: BasicLogoutWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {
    isProjectMode: {
      control: 'boolean',
      description: 'Whether the project is in Project Mode (requires projectId) or Auth Mode',
    }
  }
} satisfies Meta<typeof BasicLogoutWithProvider>;

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
