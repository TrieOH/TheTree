import type { Meta, StoryObj } from '@storybook/react-vite';
import { Session } from '../react/components/Session/Session';

const meta = {
  title: "Example/Session",
  component: Session,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
} satisfies Meta<typeof Session>;

export default meta;
type Story = StoryObj<typeof meta>;
export const Default: Story = {};