import type { JSX } from "@solidjs/web";
import { For, Show, untrack } from "solid-js";
import GlobeIcon from "~icons/lucide/globe";
import MailIcon from "~icons/lucide/mail";
import MessageCircleIcon from "~icons/lucide/message-circle";
import InstagramIcon from "~icons/simple-icons/instagram";
import LinkedinIcon from "~icons/simple-icons/linkedin";
import GithubIcon from "~icons/simple-icons/github";
import YoutubeIcon from "~icons/simple-icons/youtube";
import XIcon from "~icons/simple-icons/x";
import TwitterIcon from "~icons/simple-icons/twitter";
import TwitchIcon from "~icons/simple-icons/twitch";
import BlueskyIcon from "~icons/simple-icons/bluesky";
import DiscordIcon from "~icons/simple-icons/discord";
import { cn } from "@trieoh/ui-solid";
import {
  profileCompleteness,
  socialHref,
  type UniventsProfile,
} from "../model/profile-data";
import { ProfileCard } from "./ProfileCard";
import { ProfileEmptyState } from "./ProfileEmptyState";

const Globe = GlobeIcon as unknown as (p: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (p: { class?: string }) => JSX.Element;
const MessageCircle = MessageCircleIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;
const Github = GithubIcon as unknown as (p: { class?: string }) => JSX.Element;
const Instagram = InstagramIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;
const Linkedin = LinkedinIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;
const Youtube = YoutubeIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;
const X = XIcon as unknown as (p: { class?: string }) => JSX.Element;
const Twitter = TwitterIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;
const Twitch = TwitchIcon as unknown as (p: { class?: string }) => JSX.Element;
const Bluesky = BlueskyIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;
const Discord = DiscordIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;

export function ProfileAboutSection(props: {
  profile: UniventsProfile;
  ownProfile: boolean;
}) {
  const socials = () => Object.entries(props.profile.socials ?? {});

  return (
    <div class="mx-auto mt-4 grid max-w-7xl gap-4 px-4 md:grid-cols-[minmax(0,1fr)_280px] md:gap-5">
      <div class="space-y-5">
        <Show
          when={props.ownProfile && profileCompleteness(props.profile) < 100}
        >
          <ProfileCard title="Integridade do Perfil">
            <div class="flex items-center justify-between text-sm">
              <span>Complete seu perfil</span>
              <span class="rounded-full bg-primary/10 px-3 py-1 text-xs text-primary">
                {profileCompleteness(props.profile)}% completo
              </span>
            </div>
            <div class="mt-3 h-2 overflow-hidden rounded-full bg-muted">
              <div
                class="h-full bg-primary"
                style={{ width: `${profileCompleteness(props.profile)}%` }}
              />
            </div>
            <p class="mt-3 text-sm italic text-muted-foreground">
              {completenessHint(props.profile)}
            </p>
          </ProfileCard>
        </Show>
        <ProfileCard title="Sobre mim">
          <Show
            when={props.profile.aboutMe}
            fallback={
              <ProfileEmptyState message="Nada aqui ainda. Conte um pouco sobre você!" />
            }
          >
            <p class="whitespace-pre-wrap text-[15px] leading-[1.7] text-muted-foreground">
              {props.profile.aboutMe}
            </p>
          </Show>
        </ProfileCard>
      </div>
      <div class="space-y-5">
        <ProfileCard title="Idiomas">
          <Show
            when={props.profile.languages?.length}
            fallback={
              <ProfileEmptyState message="Nenhum idioma informado ainda." />
            }
          >
            <div class="flex flex-wrap gap-2">
              <For each={props.profile.languages}>
                {(item) => (
                  <span class="rounded-md border border-border bg-muted/50 px-2.5 py-1 text-xs text-muted-foreground">
                    {item}
                  </span>
                )}
              </For>
            </div>
          </Show>
        </ProfileCard>
        <ProfileCard title="Contato">
          <Show
            when={
              props.profile.website ||
              props.profile.contactEmail ||
              socials().length
            }
            fallback={
              <ProfileEmptyState message="Nenhuma informação de contato compartilhada." />
            }
          >
            <div class="grid grid-cols-2 gap-1">
              <Show when={props.profile.website}>
                {(site) => (
                  <ProfileSocial
                    href={site()}
                    label="Website"
                    network="website"
                  />
                )}
              </Show>
              <Show when={props.profile.contactEmail}>
                {(email) => (
                  <div class="hidden min-w-0 md:block">
                    <ProfileSocial
                      href={`mailto:${email()}`}
                      label="E-mail"
                      network="email"
                    />
                  </div>
                )}
              </Show>
              <For each={socials()}>
                {([network, value]) => (
                  <ProfileSocial
                    href={socialHref(network, String(value ?? ""))}
                    network={network}
                    label={
                      network === "x"
                        ? "X"
                        : network[0].toUpperCase() + network.slice(1)
                    }
                  />
                )}
              </For>
            </div>
          </Show>
        </ProfileCard>
      </div>
    </div>
  );
}

function ProfileSocial(props: {
  href: string;
  label: string;
  network?: string;
  class?: string;
}) {
  return (
    <a
      href={props.href}
      target={props.href.startsWith("mailto:") ? undefined : "_blank"}
      rel="noreferrer"
      class={cn(
        "flex min-w-0 items-center gap-2.5 rounded-md p-2 text-sm text-card-foreground transition-colors hover:bg-muted",
        props.class,
      )}
    >
      <span class="flex size-7 shrink-0 items-center justify-center rounded-md bg-muted text-foreground">
        <ProfileSocialIcon label={props.label} network={props.network} />
      </span>
      <span class="truncate font-medium">{props.label}</span>
    </a>
  );
}

function ProfileSocialIcon(props: {
  label: string;
  network?: string;
}): JSX.Element {
  const key = untrack(() => (props.network || props.label).toLowerCase());
  switch (key) {
    case "website":
      return <Globe class="size-4" />;
    case "e-mail":
    case "email":
      return <Mail class="size-4" />;
    case "github":
      return <Github class="size-4" />;
    case "instagram":
      return <Instagram class="size-4" />;
    case "linkedin":
      return <Linkedin class="size-4" />;
    case "youtube":
      return <Youtube class="size-4" />;
    case "x":
      return <X class="size-4" />;
    case "twitter":
      return <Twitter class="size-4" />;
    case "twitch":
      return <Twitch class="size-4" />;
    case "bluesky":
      return <Bluesky class="size-4" />;
    case "discord":
      return <Discord class="size-4" />;
    default:
      return <MessageCircle class="size-4" />;
  }
}

function completenessHint(profile: UniventsProfile) {
  if (!profile.preferredName && !profile.legalName)
    return "Adicione seu nome completo ou como prefere ser chamado.";
  if (!profile.aboutMe) return "Conte um pouco sobre você na bio.";
  if (!profile.languages?.length) return "Adicione seus idiomas de domínio.";
  return "Seu perfil está quase completo!";
}
