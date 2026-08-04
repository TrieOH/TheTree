import { createFileRoute } from "@tanstack/react-router";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { ProfileConfigPage } from "@/features/profile/ui/profile-config-page";

export const Route = createFileRoute("/profile/config")({
  beforeLoad: requireAuth,
  component: ProfileConfigPage,
});
