import type { JSX } from "@solidjs/web";
import { Show, createMemo, type Accessor } from "solid-js";
import { useQuery } from "@trieoh/front-core-solid";
import { EventQueryError } from "@/features/events/ui/EventQueryError";
import {
  type LoadProfile,
  profileKeys,
  profileQueryOptions,
} from "@/features/profile/api";
import {
  asUniventsProfile,
  profileDisplayName,
  type ActorProfile,
} from "../model/profile-data";
import { ProfileAboutSection } from "./ProfileAboutSection";
import { ProfileEmptyState } from "./ProfileEmptyState";
import { ProfileHeader } from "./ProfileHeader";
import { ProfileSkeleton } from "./ProfileSkeleton";
import { ProfileTabContent } from "./ProfileTabContent";
import type { ProfileTab } from "./profile-tabs";

export interface ProfileViewProps {
  actorId: string;
  loadProfile?: LoadProfile;
  ownProfile?: boolean | Accessor<boolean>;
  viewerActorId?: string | Accessor<string | undefined>;
  activeTab: ProfileTab;
  onTabChange: (tab: ProfileTab) => void;
}

export function ProfileView(props: ProfileViewProps): JSX.Element {
  const query = useQuery(() =>
    props.loadProfile
      ? profileQueryOptions(props.actorId, props.loadProfile)
      : {
          queryKey: profileKeys.detail(props.actorId),
          queryFn: () => null,
        },
  );
  const data = createMemo<ActorProfile | undefined>(
    () => query().data ?? undefined,
  );
  const profile = createMemo(() =>
    asUniventsProfile({
      ...data()?.profile,
      ...(data()?.pfp_url !== undefined ? { pfpUrl: data()?.pfp_url } : {}),
    }),
  );
  const own = createMemo(() => {
    const isExplicitlyOwn =
      typeof props.ownProfile === "function"
        ? props.ownProfile()
        : props.ownProfile;
    const viewerActorId =
      typeof props.viewerActorId === "function"
        ? props.viewerActorId()
        : props.viewerActorId;

    return (
      Boolean(isExplicitlyOwn) ||
      Boolean(
        viewerActorId &&
        (data()?.actor_id === viewerActorId || props.actorId === viewerActorId),
      )
    );
  });
  const publicIdentifier = createMemo(
    () => data()?.handle ?? data()?.actor_id ?? props.actorId,
  );
  const participantName = createMemo(
    () => profile().legalName || profile().preferredName || "",
  );
  const queryError = createMemo(
    () =>
      query().error as
        | (Error & {
            code?: number;
            envelope?: { code?: number; message?: string };
          })
        | null,
  );
  const errorCode = createMemo(
    () => queryError()?.envelope?.code ?? queryError()?.code,
  );
  const isProfileNotFound = createMemo(
    () => query().isError && errorCode() === 404,
  );
  const errorMessage = createMemo(() => {
    if (isProfileNotFound() && own()) {
      return "Não encontramos os dados do seu usuário. Você ainda pode editar ou configurar o perfil.";
    }
    if (isProfileNotFound()) {
      return "Não encontramos este perfil.";
    }
    return "Não foi possível carregar este perfil. Verifique sua conexão e tente novamente.";
  });

  return (
    <Show
      when={!query().isPending}
      fallback={<ProfileSkeleton ownProfile={own()} />}
    >
      <Show
        when={!query().isError}
        fallback={
          <main class="min-h-dvh bg-background">
            <div class="mx-auto max-w-7xl px-4 py-8">
              <EventQueryError
                message={errorMessage()}
                onRetry={() => void query().refetch()}
              />
            </div>
          </main>
        }
      >
        <Show
          when={data()}
          fallback={
            <div class="mx-auto max-w-7xl px-4 py-8">
              <ProfileEmptyState message="Não encontramos os dados do seu usuário. Você ainda pode editar ou configurar o perfil." />
            </div>
          }
        >
          <main class="min-h-dvh bg-background pb-28">
            <ProfileHeader
              profile={profile()}
              name={profileDisplayName(profile())}
              profileUrl={`${typeof window !== "undefined" ? window.location.origin : ""}/profile/${data()?.handle ?? data()?.actor_id ?? ""}`}
              handle={data()?.handle ?? undefined}
              ownProfile={own()}
              activeTab={props.activeTab}
              onTabChange={props.onTabChange}
            />
            <Show
              when={props.activeTab === "about"}
              fallback={
                <ProfileTabContent
                  tab={props.activeTab as Exclude<ProfileTab, "about">}
                  actorId={data()?.actor_id || props.actorId}
                  publicIdentifier={publicIdentifier()}
                  participantName={participantName()}
                  ownProfile={own()}
                />
              }
            >
              <ProfileAboutSection profile={profile()} ownProfile={own()} />
            </Show>
          </main>
        </Show>
      </Show>
    </Show>
  );
}
