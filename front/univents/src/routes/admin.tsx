import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { AdminLayout } from "@/features/admin/ui/admin-layout";
import { requireAuth } from "@/features/auths/lib/route-guard";

export const Route = createFileRoute("/admin")({
  beforeLoad: (args) =>
    requireAuth(args, {
      onRedirect: () => {
        throw redirect({ to: "/" });
      },
    }),
  component: AdminLayoutWrapper,
});

function AdminLayoutWrapper() {
  return (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  );
}
