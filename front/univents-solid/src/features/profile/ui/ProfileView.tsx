import {
  Loading,
  Show,
  createMemo,
  createSignal,
  onSettled,
  untrack,
} from "solid-js";
import type { JSX } from "@solidjs/web";
import { useQueryClient } from "@trieoh/front-core/solid";
import { userBadgesQueryOptions } from "@/features/badges/api";
import { myCertificationsQueryOptions } from "@/features/certifications/api";
import { myPurchasesQueryOptions } from "@/features/purchases/api";
import GlobeIcon from "~icons/lucide/globe";
import MailIcon from "~icons/lucide/mail";
import GithubIcon from "~icons/lucide/github";
import InstagramIcon from "~icons/lucide/instagram";
import LinkedinIcon from "~icons/lucide/linkedin";
import YoutubeIcon from "~icons/lucide/youtube";
import MessageCircleIcon from "~icons/lucide/message-circle";
import {
  asUniventsProfile,
  profileCompleteness,
  profileDisplayName,
  socialHref,
} from "../model/profile-data";
import { ProfileHeader } from "./ProfileHeader";
import { ProfileCollectionItem } from "./ProfileCollectionItem";

const Globe = GlobeIcon as unknown as () => JSX.Element;
const Mail = MailIcon as unknown as () => JSX.Element;
const Github = GithubIcon as unknown as () => JSX.Element;
const Instagram = InstagramIcon as unknown as () => JSX.Element;
const Linkedin = LinkedinIcon as unknown as () => JSX.Element;
const Youtube = YoutubeIcon as unknown as () => JSX.Element;
const MessageCircle = MessageCircleIcon as unknown as () => JSX.Element;

type Tab = "about" | "badges" | "certificates" | "purchases";
type ProfileResponse = {
  success: boolean;
  data?: {
    actor_id?: string;
    handle?: string | null;
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
  const [result, setResult] = createSignal<ProfileResponse>();
  onSettled(() => {
    void (async () =>
      setResult(
        await (props.actorId
          ? props.loadProfile(props.actorId)
          : { success: false }),
      ))();
  });
  return (
    <Loading fallback={<ProfileSkeleton />}>
      <Show
        when={result()}
        fallback={<main class="p-12 text-center">Perfil não encontrado.</main>}
      >
        {(response) => {
          const profile = asUniventsProfile(response().data?.profile ?? {});
          const name = profileDisplayName(profile);
          const own =
            props.ownProfile ||
            response().data?.actor_id === props.viewerActorId;
          const socials = Object.entries(profile.socials ?? {}).filter(
            ([, value]) => Boolean(value),
          ) as [string, string][];
          return (
            <main class="min-h-dvh bg-background pb-28">
              <ProfileHeader
                profile={profile}
                name={name}
                profileUrl={`${window.location.origin}/profile/${response().data?.handle ?? response().data?.actor_id ?? ""}`}
                handle={response().data?.handle ?? undefined}
                ownProfile={own}
                activeTab={props.activeTab}
                onTabChange={props.onTabChange}
              />
              <Show
                when={props.activeTab === "about"}
                fallback={
                  <ProfileTabContent
                    tab={props.activeTab as Exclude<Tab, "about">}
                    actorId={response().data?.actor_id ?? ""}
                    ownProfile={own}
                    queryClient={queryClient}
                  />
                }
              >
                <div class="mx-auto mt-4 grid max-w-7xl gap-4 px-4 md:grid-cols-[minmax(0,1fr)_280px] md:gap-5">
                  <div class="space-y-5">
                    <Show when={own && profileCompleteness(profile) < 100}>
                      <Card title="Integridade do Perfil">
                        <div class="flex items-center justify-between text-sm">
                          <span>Complete seu perfil</span>
                          <span class="rounded-full bg-primary/10 px-3 py-1 text-xs text-primary">
                            {profileCompleteness(profile)}% completo
                          </span>
                        </div>
                        <div class="mt-3 h-2 overflow-hidden rounded-full bg-muted">
                          <div
                            class="h-full bg-primary"
                            style={{
                              width: `${profileCompleteness(profile)}%`,
                            }}
                          />
                        </div>
                        <p class="mt-3 text-sm italic text-muted-foreground">
                          {completenessHint(profile)}
                        </p>
                      </Card>
                    </Show>
                    <Card title="Sobre mim">
                      <Show
                        when={profile.aboutMe}
                        fallback={
                          <EmptyState message="Nada aqui ainda. Conte um pouco sobre você!" />
                        }
                      >
                        <p class="whitespace-pre-wrap text-[15px] leading-[1.7] text-muted-foreground">
                          {profile.aboutMe}
                        </p>
                      </Show>
                    </Card>
                  </div>
                  <div class="space-y-5">
                    <Card title="Idiomas">
                      <Show
                        when={profile.languages?.length}
                        fallback={
                          <EmptyState message="Adicione idiomas para destacar seu perfil." />
                        }
                      >
                        <div class="flex flex-wrap gap-2">
                          {profile.languages?.map((item) => (
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
                          profile.website ||
                          profile.contactEmail ||
                          socials.length
                        }
                        fallback={
                          <EmptyState message="Adicione formas de contato para que outros possam se conectar." />
                        }
                      >
                        <div class="grid grid-cols-2 gap-1">
                          {profile.website && (
                            <Social href={profile.website} label="Website" />
                          )}{" "}
                          {profile.contactEmail && (
                            <Social
                              href={`mailto:${profile.contactEmail}`}
                              label="E-mail"
                            />
                          )}{" "}
                          {socials.map(([network, value]) => (
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
    return props.queryClient.fetchQuery(myPurchasesQueryOptions());
  });

  return (
    <Loading fallback={<ProfileSkeleton />}>
      <Show when={data()}>
        {(loaded) => <CollectionCard tab={props.tab} data={loaded()} />}
      </Show>
    </Loading>
  );
}

function CollectionCard(props: { tab: Exclude<Tab, "about">; data: unknown }) {
  const title = () =>
    props.tab === "badges"
      ? "Crachás"
      : props.tab === "certificates"
        ? "Certificados"
        : "Compras";
  const items = () => {
    const data = props.data;
    return Array.isArray(data)
      ? data
      : data && typeof data === "object"
        ? Object.values(data as Record<string, unknown>).flatMap((value) =>
          Array.isArray(value) ? value : [],
        )
        : [];
  };
  return (
    <div class="mx-auto mt-5 max-w-7xl px-4">
      <Card title={title()}>
        <Show
          when={items().length > 0}
          fallback={
            <EmptyState
              message={`Nenhum item em ${title().toLowerCase()} ainda.`}
            />
          }
        >
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
function ProfileSkeleton() {
  return <main class="min-h-dvh animate-pulse bg-muted" />;
}
