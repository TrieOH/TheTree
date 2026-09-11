import {
  Loading,
  Show,
  createMemo,
  untrack,
} from "solid-js";
import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core/solid";
import { userBadgesQueryOptions } from "@/features/badges/api";
import { myCertificationsQueryOptions } from "@/features/certifications/api";
import GlobeIcon from "~icons/lucide/globe";
import MailIcon from "~icons/lucide/mail";
import GithubIcon from "~icons/lucide/github";
import InstagramIcon from "~icons/lucide/instagram";
import LinkedinIcon from "~icons/lucide/linkedin";
import YoutubeIcon from "~icons/lucide/youtube";
import MessageCircleIcon from "~icons/lucide/message-circle";
import CircleAlertIcon from "~icons/lucide/circle-alert";
import {
  asUniventsProfile,
  profileCompleteness,
  profileDisplayName,
  socialHref,
} from "../model/profile-data";
import { ProfileHeader } from "./ProfileHeader";
import { ProfileCollectionItem } from "./ProfileCollectionItem";
import { PurchasesContent } from "@/features/purchases/ui/PurchasesContent";

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
type ProfileResponse = {
  success: boolean;
  data?: {
    actor_id?: string;
    handle?: string | null;
    pfp_url?: string | null;
    profile?: Record<string, unknown>;
  };
  message?: string;
};
export interface ProfileViewProps {
  actorId?: string;
  loadProfile: (actorId: string) => Promise<ProfileResponse>;
  ownProfile?: boolean;
  viewerActorId?: string;
  activeTab: Tab;
  onTabChange: (tab: Tab) => void;
}

export function ProfileView(props: ProfileViewProps) {
  const queryClient = useQueryClient();
  const result = createMemo(async () => {
    const response = props.actorId
      ? await props.loadProfile(props.actorId)
      : undefined;

    return response?.success ? response : null;
  });
  return (
    <Loading fallback={<ProfileSkeleton />}>
      <Show when={result()} fallback={<MissingPublicProfile />}>
        {(response) => {
          const data = () => response().data;
          const profile = createMemo(() =>
            asUniventsProfile({
              ...data()?.profile,
              ...(data()?.pfp_url !== undefined && { pfpUrl: data()?.pfp_url }),
            }),
          );
          const name = () => profileDisplayName(profile());
          const own = () =>
            Boolean(
              props.ownProfile || data()?.actor_id === props.viewerActorId,
            );
          const socials = () =>
            Object.entries(profile().socials ?? {}).filter(([, value]) =>
              Boolean(value),
            ) as [string, string][];
          return (
            <main class="min-h-dvh bg-background pb-28">
              <ProfileHeader
                profile={profile()}
                name={name()}
                profileUrl={`${window.location.origin}/profile/${data()?.handle ?? data()?.actor_id ?? ""}`}
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
                    actorId={data()?.actor_id ?? ""}
                    ownProfile={own()}
                    queryClient={queryClient}
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
                          <EmptyState message="Adicione idiomas para destacar seu perfil." />
                        }
                      >
                        <div class="flex flex-wrap gap-2">
                          {profile().languages?.map((item) => (
                            <span class="rounded-md border border-border px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-muted">
                              {item}
                            </span>
                          ))}
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
                          <EmptyState message="Adicione formas de contato para que outros possam se conectar." />
                        }
                      >
                        <div class="grid grid-cols-2 gap-1">
                          <Show when={profile().website}>
                            {(website) => (
                              <Social href={website()} label="Website" />
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
                          {socials().map(([network, value]) => (
                            <Social
                              href={socialHref(network, value)}
                              label={
                                network[0].toUpperCase() + network.slice(1)
                              }
                            />
                          ))}
                        </div>
                      </Show>
                    </Card>
                  </div>
                </div>
              </Show>
            </main>
          );
        }}
      </Show>
    </Loading>
  );
}

function ProfileTabContent(props: {
  tab: Exclude<Tab, "about">;
  actorId: string;
  ownProfile: boolean;
  queryClient: ReturnType<typeof useQueryClient>;
}) {
  const data = createMemo(async () => {
    const tab = props.tab;
    const actorId = props.actorId;
    const ownProfile = props.ownProfile;
    if (tab === "badges") {
      return props.queryClient.fetchQuery(userBadgesQueryOptions(actorId));
    }
    if (!ownProfile) return [];
    if (tab === "certificates") {
      return props.queryClient.fetchQuery(myCertificationsQueryOptions());
    }
    return [];
  });

  // Purchases fetch inside <PurchasesContent>, which owns its own skeleton.
  // Wrapping it here would flash a second, differently shaped one first.
  return (
    <Show
      when={props.tab === "purchases"}
      fallback={
        <Loading fallback={<CollectionSkeleton tab={props.tab} />}>
          <Show when={data()}>
            {(loaded) => <CollectionCard tab={props.tab} data={loaded()} />}
          </Show>
        </Loading>
      }
    >
      <CollectionCard tab="purchases" data={[]} />
    </Show>
  );
}

const tabTitle = (tab: Exclude<Tab, "about">) =>
  tab === "badges"
    ? "Crachás"
    : tab === "certificates"
      ? "Certificados"
      : "Compras";

function CollectionCard(props: { tab: Exclude<Tab, "about">; data: unknown }) {
  const title = () => tabTitle(props.tab);
  const items = () => {
    const data = props.data;
    const collect = (value: unknown): unknown[] =>
      Array.isArray(value)
        ? value
        : value && typeof value === "object"
          ? Object.values(value as Record<string, unknown>).flatMap(collect)
          : [];
    return collect(data);
  };
  return (
    <div class="mx-auto mt-4 max-w-7xl px-4">
      <Card title={title()}>
        <Show when={props.tab === "purchases"}>
          <PurchasesContent />
        </Show>
        {/* TODO: I need to change this(one component for each) */}
        <Show when={props.tab !== "purchases"}>
          <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {items().map((item) => (
              <ProfileCollectionItem tab={props.tab} item={item} />
            ))}
          </div>
        </Show>
      </Card>
    </div>
  );
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
  return <p class="text-sm italic text-muted-foreground/70">{props.message}</p>;
}

function completenessHint(profile: ReturnType<typeof asUniventsProfile>) {
  if (!(profile.preferredName || profile.legalName))
    return "Adicione seu nome para completar o perfil.";
  if (!profile.pfpUrl) return "Adicione uma foto para completar o perfil.";
  if (!profile.bannerUrl)
    return "Adicione uma imagem de capa para completar o perfil.";
  if (!profile.aboutMe)
    return 'Preencha a seção "Sobre mim" para completar o perfil.';
  if (!(profile.role || profile.organization))
    return "Adicione sua função ou organização para completar o perfil.";
  if (!profile.languages?.length)
    return "Adicione ao menos um idioma para completar o perfil.";
  if (
    !(
      profile.website ||
      profile.contactEmail ||
      Object.values(profile.socials ?? {}).some(Boolean)
    )
  )
    return "Adicione uma forma de contato para completar o perfil.";
  return "Perfil completo!";
}
function CollectionSkeleton(props: { tab: Exclude<Tab, "about"> }) {
  return (
    <div class="mx-auto mt-4 max-w-7xl px-4">
      <Card title={tabTitle(props.tab)}>
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <CollectionItemSkeleton />
          <CollectionItemSkeleton />
          <CollectionItemSkeleton />
        </div>
      </Card>
    </div>
  );
}

function CollectionItemSkeleton() {
  return (
    <div class="overflow-hidden rounded-md border border-border bg-card shadow-sm">
      <div class="flex gap-3 p-4">
        <div class="size-9 shrink-0 animate-pulse rounded-md bg-muted" />

        <div class="min-w-0 flex-1 space-y-2">
          <div class="h-4 w-2/3 animate-pulse rounded bg-muted" />
          <div class="h-3 w-full animate-pulse rounded bg-muted" />
          <div class="h-3 w-4/5 animate-pulse rounded bg-muted" />
        </div>
      </div>
    </div>
  );
}

function ProfileSkeleton() {
  return (
    <main class="min-h-dvh bg-background pb-28">
      <div class="w-full">
        <div class="h-44 w-full animate-pulse bg-muted md:h-56" />

        <div class="mx-auto max-w-7xl px-4">
          <div class="-mt-14 size-28 animate-pulse rounded-full border-4 border-border bg-muted md:-mt-16 md:size-32" />

          <div class="mt-4 space-y-2">
            <div class="h-8 w-48 animate-pulse rounded bg-muted" />
            <div class="h-4 w-64 animate-pulse rounded bg-muted" />
            <div class="h-4 w-40 animate-pulse rounded bg-muted" />
          </div>
        </div>
      </div>

      <div class="mx-auto mt-5 hidden max-w-7xl gap-5 px-4 md:grid md:grid-cols-[1fr_280px]">
        <div class="space-y-5">
          <div class="h-36 animate-pulse rounded-md bg-muted" />
          <div class="h-32 animate-pulse rounded-md bg-muted" />
        </div>

        <div class="space-y-5">
          <div class="h-64 animate-pulse rounded-md bg-muted" />
          <div class="h-24 animate-pulse rounded-md bg-muted" />
        </div>
      </div>
    </main>
  );
}

function MissingPublicProfile() {
  return (
    <main class="relative min-h-dvh overflow-hidden bg-background">
      <div aria-hidden="true" class="pointer-events-none blur-md opacity-50">
        <ProfileSkeleton />
      </div>

      <div class="fixed inset-0 z-50 flex items-center justify-center bg-background/35 px-4 backdrop-blur-sm">
        <div
          role="alert"
          class="w-full max-w-md rounded-lg border border-border bg-card p-6 text-center shadow-xl"
        >
          <CircleAlert class="mx-auto size-8 text-muted-foreground" />

          <h1 class="mt-4 text-xl font-semibold">Perfil não encontrado</h1>

          <p class="mt-2 text-sm text-muted-foreground">
            Este perfil não existe ou não está mais disponível.
          </p>

          <Link
            to="/profile"
            search={{ tab: "about" }}
            class="mt-5 inline-flex h-10 w-full items-center justify-center rounded-md bg-primary px-4 text-sm text-primary-foreground shadow-sm"
          >
            Ir para o meu perfil
          </Link>
        </div>
      </div>
    </main>
  );
}
