import { Navigate, createFileRoute } from '@tanstack/solid-router';

export const Route = createFileRoute('/')({
  component: () => <Navigate to="/auth" />,
});
