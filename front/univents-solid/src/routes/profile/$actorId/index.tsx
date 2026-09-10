import { createFileRoute } from '@tanstack/solid-router'
import { useAuth } from '@trieoh/identityx-sdk-ts-solid';
import { createMemo } from 'solid-js';
import z from 'zod'
import { ProfileView } from "@/features/profile/ui/ProfileView";

const PROFILE_TAB_TITLES: Record<string, string> = {
  about: "Perfil",
  badges: "Crachás",
  certificates: "Certificados",
  purchases: "Compras",
};

export const Route = createFileRoute('/profile/$actorId/')({
  validateSearch: z.object({
    tab: z
      .enum(["about", "badges", "certificates", "purchases"])
      .catch("about"),
  }),
  // `head` receives the match, so changing the tab changes the title too.
  head: ({ match, params }) => ({
    meta: [
      {
        title: `${PROFILE_TAB_TITLES[match.search.tab] ?? "Perfil"} de ${params.actorId} - Univents`,
      },
    ],
  }),
  component: RouteComponent,
})

function RouteComponent() {
  const params = Route.useParams();
  const search = Route.useSearch();

  const { auth } = useAuth();

  const navigate = Route.useNavigate();

  const viewerActorId = createMemo(() => auth.profile()?.id);

  const loadProfile = async (identifier: string) => {
    const isActorId =
      /^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(identifier);

    const response = isActorId
      ? await auth.getActorProfile(identifier)
      : await auth.getProfileByHandle(identifier);

    if (!response.success) return response;

    const handle = response.data?.handle;

    if (isActorId && handle && handle !== identifier) {
      await navigate({
        to: "/profile/$actorId",
        params: {
          actorId: handle,
        },
        search: {
          tab: search().tab,
        },
        replace: true,
      });
    }

    return response;
  };

  return (
    <ProfileView
      actorId={params().actorId}
      loadProfile={loadProfile}
      ownProfile={params().actorId === viewerActorId()}
      viewerActorId={viewerActorId()}
      activeTab={search().tab}
      onTabChange={(nextTab) => {
        void navigate({
          to: "/profile/$actorId",
          params: {
            actorId: params().actorId,
          },
          search: {
            tab: nextTab,
          },
        });
      }}
    />
  );
}
