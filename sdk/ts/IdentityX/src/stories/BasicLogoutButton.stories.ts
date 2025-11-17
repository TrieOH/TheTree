import type { Meta, StoryObj } from '@storybook/react-vite';
import { fn } from 'storybook/test';
import BasicLogoutWithProvider from './components/BasicLogoutWithProvider';

const meta = {
  title: "Example/BasicLogoutButton",
  component: BasicLogoutWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  args: { onSubmit: fn() },
} satisfies Meta<typeof BasicLogoutWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};