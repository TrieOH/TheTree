import { createFileRoute } from "@tanstack/solid-router";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { createMemo } from "solid-js";
import { setOwnProfilePath } from "@/features/profile/lib/own-profile-path";
import { ProfileView } from "@/features/profile/ui/ProfileView";
import {
  cacheProfile,
  loadProfile as fetchProfile,
  prefetchProfileTab,
} from "@/features/profile/api";
import {
  PROFILE_TABS,
  PROFILE_TAB_LABELS,
} from "@/features/profile/model/profile-tabs";
import { Schema, createSchema } from "@/shared/lib/effect-schema";
import { appQueryClient } from "@/shared/lib/query-client";
import { boundLoader } from "@/shared/lib/route-data";

const profileActorSearchSchema = createSchema(
  Schema.Struct({
    tab: Schema.optional(Schema.Literals(PROFILE_TABS)),
  }),
);

export const Route = createFileRoute("/profile/$actorId/")({
  validateSearch: (search) => {
    const parsed = profileActorSearchSchema.parse(search);
    return { tab: parsed.tab ?? "about" };
  },
  loaderDeps: ({ search }) => ({ tab: search.tab }),
  loader: ({ deps, params, context }) => {
    const auth = context.session?.service;
    if (!auth) return;
    return boundLoader(
      prefetchProfileTab(appQueryClient, {
        auth,
        identifier: params.actorId,
        viewerId: auth.profile()?.id,
        tab: deps.tab,
      }),
    );
  },
  head: ({ match, params }) => ({
    meta: [
      {
        title: `${
          match.search.tab === "about"
            ? "Perfil"
            : PROFILE_TAB_LABELS[match.search.tab]
        } de ${params.actorId} - Univents`,
      },
    ],
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const params = Route.useParams();
  const search = Route.useSearch();
  const { auth } = useAuth();
  const navigate = Route.useNavigate();
  const viewerActorId = createMemo(() => auth.profile()?.id);

  const loadProfile = async (identifier: string) => {
    const response = await fetchProfile(identifier, auth);

    if (!response.success) return response;

    // Seed the actor-id and handle keys too: the redirect below changes the
    // identifier, and without this the new key starts empty and the page goes
    // content -> skeleton -> content on every canonicalization.
    cacheProfile(appQueryClient, response, identifier);

    const handle = response.data?.handle;

    if (
      identifier === response.data?.actor_id &&
      handle &&
      handle !== identifier
    ) {
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

    setOwnProfilePath(
      response.data?.actor_id && response.data.actor_id === viewerActorId()
        ? `/profile/${handle ?? identifier}`
        : null,
    );

    return response;
  };

  return (
    <ProfileView
      actorId={params().actorId}
      loadProfile={loadProfile}
      viewerActorId={viewerActorId}
      activeTab={search().tab}
      onTabChange={(nextTab) => {
        if (nextTab === search().tab) return;

        void navigate({
          to: "/profile/$actorId",
          params: {
            actorId: params().actorId,
          },
          search: {
            tab: nextTab,
          },
          resetScroll: false,
        });
      }}
    />
  );
}
