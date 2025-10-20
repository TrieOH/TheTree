import type { Meta, StoryObj } from '@storybook/react-vite';
import { Copyright } from '../next';

const meta = {
  title: "Example/Copyright",
  component: Copyright,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
} satisfies Meta<typeof Copyright>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};