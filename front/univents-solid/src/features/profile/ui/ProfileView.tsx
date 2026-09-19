import type { JSX } from "@solidjs/web";
import {
  For,
  Loading,
  Show,
  createEffect,
  createMemo,
  createSignal,
  untrack,
} from "solid-js";
import CircleAlertIcon from "~icons/lucide/circle-alert";
import GithubIcon from "~icons/lucide/github";
import GlobeIcon from "~icons/lucide/globe";
import InstagramIcon from "~icons/lucide/instagram";
import LinkedinIcon from "~icons/lucide/linkedin";
import MailIcon from "~icons/lucide/mail";
import MessageCircleIcon from "~icons/lucide/message-circle";
import YoutubeIcon from "~icons/lucide/youtube";
import { useQuery } from "@trieoh/front-core-solid";
import { userBadgesQueryOptions } from "@/features/badges/api";
import type { BadgeProfileGroups } from "@/features/badges/model";
import { allProfileBadges } from "@/features/badges/model/profile-badges";
import { ProfileBadges } from "@/features/badges/ui/ProfileBadges";
import { myCertificationsQueryOptions } from "@/features/certifications/api";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import { getPublicEditionsFn } from "@/features/editions/api";
import {
  type LoadProfile,
  profileKeys,
  profileQueryOptions,
} from "@/features/profile/api";
import { PurchasesContent } from "@/features/purchases/ui/PurchasesContent";
import {
  asUniventsProfile,
  profileCompleteness,
  profileDisplayName,
  socialHref,
  type ActorProfile,
} from "../model/profile-data";
import { ProfileCollectionItem } from "./ProfileCollectionItem";
import { ProfileHeader } from "./ProfileHeader";

const Globe = GlobeIcon as unknown as () => JSX.Element;
const Mail = MailIcon as unknown as () => JSX.Element;
const Github = GithubIcon as unknown as () => JSX.Element;
const Instagram = InstagramIcon as unknown as () => JSX.Element;
const Linkedin = LinkedinIcon as unknown as () => JSX.Element;
const Youtube = YoutubeIcon as unknown as () => JSX.Element;
const MessageCircle = MessageCircleIcon as unknown as () => JSX.Element;
const CircleAlert = CircleAlertIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;

type Tab = "about" | "badges" | "certificates" | "purchases";

export interface ProfileViewProps {
  actorId: string;
  loadProfile?: LoadProfile;
  ownProfile?: boolean;
  viewerActorId?: string;
  activeTab: Tab;
  onTabChange: (tab: Tab) => void;
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
  const raw = createMemo(() => query().data as any);
  const data = createMemo<ActorProfile | undefined>(() => {
    const d = raw();
    if (!d) return undefined;
    return d.actor_id ? d : d.data;
  });
  const profile = createMemo(() =>
    asUniventsProfile({
      ...(data()?.profile),
      ...(data()?.pfp_url !== undefined ? { pfpUrl: data()?.pfp_url } : {}),
    }),
  );
  const own = createMemo(
    () =>
      Boolean(props.ownProfile) ||
      Boolean(
        props.viewerActorId &&
        (data()?.actor_id === props.viewerActorId ||
          props.actorId === props.viewerActorId),
      ),
  );
  const name = createMemo(() => profileDisplayName(profile()));
  const publicIdentifier = createMemo(
    () => data()?.handle ?? data()?.actor_id ?? props.actorId,
  );
  const participantName = createMemo(
    () => profile().preferredName || profile().legalName || "",
  );
  const socials = createMemo(() => Object.entries(profile().socials ?? {}));

  const isError = createMemo(() => {
    if (own()) return false;
    if (query().isError) return true;
    const code = raw()?.code;
    if (code !== undefined && code >= 400) return true;
    if (!query().isLoading && !data()) return true;
    return false;
  });

  return (
    <Loading fallback={<ProfileSkeleton />}>
      <Show
        when={!isError()}
        fallback={
          <div class="mx-auto max-w-7xl px-4 py-8">
            <div class="flex items-center gap-3 rounded-md border border-destructive/40 bg-destructive/10 p-4 text-destructive">
              <CircleAlert class="size-5 shrink-0" />
              <div>
                <p class="font-medium">
                  {raw()?.message || "Não foi possível carregar este perfil."}
                </p>
                <p class="text-sm text-destructive/80">
                  {raw()?.code === 404
                    ? "O perfil solicitado não foi encontrado."
                    : "Tente novamente mais tarde."}
                </p>
              </div>
            </div>
          </div>
        }
      >
        <Show
          when={data()}
          fallback={
            <div class="mx-auto max-w-7xl px-4 py-8">
              <div class="flex items-center gap-3 rounded-md border border-muted bg-card p-4 text-muted-foreground">
                <p class="text-sm">
                  Não encontramos os dados do seu usuário. Você ainda pode editar
                  ou configurar o perfil.
                </p>
              </div>
            </div>
          }
        >
          <main class="min-h-dvh bg-background pb-28">
            <ProfileHeader
              profile={profile()}
              name={name()}
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
                  tab={props.activeTab as Exclude<Tab, "about">}
                  actorId={data()?.actor_id || props.actorId}
                  publicIdentifier={publicIdentifier()}
                  participantName={participantName()}
                  ownProfile={own()}
                />
              }
            >
              <div class="mx-auto mt-4 grid max-w-7xl gap-4 px-4 md:grid-cols-[minmax(0,1fr)_280px] md:gap-5">
                <div class="space-y-5">
                  <Show when={own() && profileCompleteness(profile()) < 100}>
                    <Card title="Integridade do Perfil">
                      <div class="flex items-center justify-between text-sm">
                        <span>Complete seu perfil</span>
                        <span class="rounded-full bg-primary/10 px-3 py-1 text-xs text-primary">
                          {profileCompleteness(profile())}% completo
                        </span>
                      </div>
                      <div class="mt-3 h-2 overflow-hidden rounded-full bg-muted">
                        <div
                          class="h-full bg-primary"
                          style={{
                            width: `${profileCompleteness(profile())}%`,
                          }}
                        />
                      </div>
                      <p class="mt-3 text-sm italic text-muted-foreground">
                        {completenessHint(profile())}
                      </p>
                    </Card>
                  </Show>
                  <Card title="Sobre mim">
                    <Show
                      when={profile().aboutMe}
                      fallback={
                        <EmptyState message="Nada aqui ainda. Conte um pouco sobre você!" />
                      }
                    >
                      <p class="whitespace-pre-wrap text-[15px] leading-[1.7] text-muted-foreground">
                        {profile().aboutMe}
                      </p>
                    </Show>
                  </Card>
                </div>
                <div class="space-y-5">
                  <Card title="Idiomas">
                    <Show
                      when={profile().languages?.length}
                      fallback={
                        <EmptyState message="Nenhum idioma informado ainda." />
                      }
                    >
                      <div class="flex flex-wrap gap-2">
                        <For each={profile().languages}>
                          {(item) => (
                            <span class="rounded-md border border-border bg-muted/50 px-2.5 py-1 text-xs text-muted-foreground">
                              {item}
                            </span>
                          )}
                        </For>
                      </div>
                    </Show>
                  </Card>
                  <Card title="Contato">
                    <Show
                      when={
                        profile().website ||
                        profile().contactEmail ||
                        socials().length
                      }
                      fallback={
                        <EmptyState message="Nenhuma informação de contato compartilhada." />
                      }
                    >
                      <div class="space-y-1">
                        <Show when={profile().website}>
                          {(site) => (
                            <Social href={site()} label="Website" />
                          )}
                        </Show>{" "}
                        <Show when={profile().contactEmail}>
                          {(email) => (
                            <Social
                              href={`mailto:${email()}`}
                              label="E-mail"
                            />
                          )}
                        </Show>{" "}
                        <For each={socials()}>
                          {([network, value]) => (
                            <Social
                              href={socialHref(network, String(value ?? ""))}
                              label={
                                network[0].toUpperCase() + network.slice(1)
                              }
                            />
                          )}
                        </For>
                      </div>
                    </Show>
                  </Card>
                </div>
              </div>
            </Show>
          </main>
        </Show>
      </Show>
    </Loading>
  );
}

function ProfileTabContent(props: {
  tab: Exclude<Tab, "about">;
  actorId: string;
  publicIdentifier: string;
  participantName: string;
  ownProfile: boolean;
}) {
  const badgesQuery = useQuery(() => ({
    ...userBadgesQueryOptions(props.actorId),
    enabled: props.tab === "badges" && Boolean(props.actorId),
  }));

  const certsQuery = useQuery(() => ({
    ...myCertificationsQueryOptions(),
    enabled: props.tab === "certificates" && props.ownProfile,
  }));

  const data = createMemo(() => {
    if (props.tab === "badges") return badgesQuery().data;
    if (props.tab === "certificates") return certsQuery().data;
    return undefined;
  });

  const isLoading = createMemo(() => {
    if (props.tab === "badges") return badgesQuery().isLoading;
    if (props.tab === "certificates") return certsQuery().isLoading;
    return false;
  });

  const eventsQuery = useQuery(() => ({
    ...allPublicEventsQueryOptions(),
    enabled: props.tab === "badges",
  }));

  const [editionLocations, setEditionLocations] = createSignal<Map<string, string>>(new Map());

  createEffect(
    () => [props.tab, eventsQuery().data] as const,
    ([tab, evData]) => {
      if (tab !== "badges") return;
      const events = (evData ?? []) as EventI[];
      if (events.length === 0) return;
      void (async () => {
        const map = new Map<string, string>();
        await Promise.all(
          events.map(async (ev) => {
            try {
              const editions = await getPublicEditionsFn(ev.id);
              for (const ed of editions) {
                map.set(ed.id, ed.location_name ?? "");
              }
            } catch {
              // ignore
            }
          }),
        );
        setEditionLocations(new Map(map));
      })();
    },
  );

  return (
    <Show
      when={props.tab === "purchases"}
      fallback={
        <Show
          when={!isLoading()}
          fallback={<CollectionSkeleton tab={props.tab} />}
        >
          <Show when={data()}>
            <Show
              when={props.tab === "badges"}
              fallback={
                <div class="mx-auto mt-4 max-w-7xl px-4">
                  <Card title={tabTitle(props.tab)}>
                    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                      <For each={collectItems(data())}>
                        {(item) => (
                          <ProfileCollectionItem
                            tab={props.tab}
                            item={item}
                          />
                        )}
                      </For>
                    </div>
                  </Card>
                </div>
              }
            >
              <div class="mx-auto mt-5 max-w-7xl px-4">
                <ProfileBadges
                  badges={allProfileBadges(data() as BadgeProfileGroups)}
                  profileIdentifier={props.publicIdentifier}
                  participantName={props.participantName}
                  editionLocations={editionLocations()}
                />
              </div>
            </Show>
          </Show>
        </Show>
      }
    >
      <div class="mx-auto mt-4 max-w-7xl px-4">
        <PurchasesContent />
      </div>
    </Show>
  );
}

const tabTitle = (tab: Exclude<Tab, "about">) =>
  tab === "badges"
    ? "Crachás"
    : tab === "certificates"
      ? "Certificados"
      : "Compras";

function collectItems(data: unknown): unknown[] {
  if (Array.isArray(data)) return data;
  if (data && typeof data === "object") {
    return Object.values(data as Record<string, unknown>).flatMap(collectItems);
  }
  return [];
}

function Card(props: { title: string; children: unknown }) {
  return (
    <section class="rounded-md border border-border bg-card p-5 shadow-md shadow-foreground/5">
      <h2 class="mb-4 text-base font-semibold">{props.title}</h2>
      {props.children}
    </section>
  );
}

function Social(props: { href: string; label: string }) {
  return (
    <a
      href={props.href}
      target={props.href.startsWith("mailto:") ? undefined : "_blank"}
      rel="noreferrer"
      class="flex items-center gap-2.5 rounded-md p-2 text-sm hover:bg-muted"
    >
      <span class="flex size-7 items-center justify-center rounded-md bg-foreground text-background">
        <SocialIcon label={props.label} />
      </span>
      <span class="truncate">{props.label}</span>
    </a>
  );
}

function SocialIcon(props: { label: string }) {
  const label = untrack(() => props.label.toLowerCase());
  if (label === "website") return <Globe />;
  if (label === "e-mail") return <Mail />;
  if (label === "github") return <Github />;
  if (label === "instagram") return <Instagram />;
  if (label === "linkedin") return <Linkedin />;
  if (label === "youtube") return <Youtube />;
  if (label === "discord") return <MessageCircle />;
  return <Globe />;
}

function EmptyState(props: { message: string }) {
  return <p class="text-sm italic text-muted-foreground">{props.message}</p>;
}

function ProfileSkeleton() {
  return (
    <div class="mx-auto max-w-7xl space-y-6 px-4 py-8">
      <div class="h-32 rounded-md bg-muted animate-pulse" />
      <div class="grid gap-4 md:grid-cols-[minmax(0,1fr)_280px]">
        <div class="h-64 rounded-md bg-muted animate-pulse" />
        <div class="h-64 rounded-md bg-muted animate-pulse" />
      </div>
    </div>
  );
}

function CollectionSkeleton(props: { tab: Tab }) {
  return (
    <div class="mx-auto mt-4 max-w-7xl px-4">
      <Card title={tabTitle(props.tab as Exclude<Tab, "about">)}>
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <div class="h-44 rounded-md bg-muted animate-pulse" />
          <div class="h-44 rounded-md bg-muted animate-pulse" />
          <div class="h-44 rounded-md bg-muted animate-pulse" />
        </div>
      </Card>
    </div>
  );
}

function completenessHint(profile: ReturnType<typeof asUniventsProfile>) {
  if (!profile.preferredName && !profile.legalName)
    return "Adicione seu nome completo ou como prefere ser chamado.";
  if (!profile.aboutMe) return "Conte um pouco sobre você na bio.";
  if (!profile.languages?.length) return "Adicione seus idiomas de domínio.";
  return "Seu perfil está quase completo!";
}
