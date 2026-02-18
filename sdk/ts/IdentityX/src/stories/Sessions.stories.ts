import type { Meta, StoryObj } from '@storybook/react-vite';
import SessionsWithProvider from './components/SessionsWithProvider';

const meta = {
  title: "Authentication/Sessions",
  component: SessionsWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
} satisfies Meta<typeof SessionsWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;
export const Default: Story = {};