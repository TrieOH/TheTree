import type { Meta, StoryObj } from '@storybook/react-vite';
import SessionsWithProvider from './components/SessionsWithProvider';

const meta = {
  title: "Authentication/Sessions",
  component: SessionsWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {
    isProjectMode: {
      control: 'boolean',
      description: 'Whether the project is in Project Mode (requires projectId) or Auth Mode',
    }
  }
} satisfies Meta<typeof SessionsWithProvider>;

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
