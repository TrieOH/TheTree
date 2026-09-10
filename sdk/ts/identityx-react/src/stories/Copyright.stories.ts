import type { Meta, StoryObj } from '@storybook/react-vite';
import { Copyright } from '../react';

const meta = {
  title: "Common/Copyright",
  component: Copyright,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
} satisfies Meta<typeof Copyright>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const ExtraSmall: Story = {
  args: {
    size: 'xs',
  },
};

export const Small: Story = {
  args: {
    size: 'sm',
  },
};

export const Medium: Story = {
  args: {
    size: 'md',
  },
};

export const Large: Story = {
  args: {
    size: 'lg',
  },
};

export const ExtraLarge: Story = {
  args: {
    size: 'xl',
  },
};

export const DoubleExtraLarge: Story = {
  args: {
    size: '2xl',
  },
};