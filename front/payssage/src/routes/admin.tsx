import { createFileRoute, Outlet } from "@tanstack/react-router";
import { AdminLayout } from "#/features/admin/ui/admin-layout";
import { requireAuth } from "#/features/auths/lib/route-guard";

export const Route = createFileRoute("/admin")({
  beforeLoad: (ctx) => {
    requireAuth(ctx);
  },
  component: AdminLayoutWrapper,
});

function AdminLayoutWrapper() {
  return (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  );
}
