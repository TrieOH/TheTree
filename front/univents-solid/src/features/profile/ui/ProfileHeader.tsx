import type { JSX } from "@solidjs/web";
import { For, untrack } from "solid-js";
import { Link } from "@tanstack/solid-router";
import { cn } from "@trieoh/ui-solid";
import CalendarIcon from "~icons/lucide/calendar";
import MailIcon from "~icons/lucide/mail";
import PencilIcon from "~icons/lucide/pencil";
import SettingsIcon from "~icons/lucide/settings";
import { resolveStorageUrl } from "@/shared/lib/storage-url";
import type { UniventsProfile } from "../model/profile-data";
import { ProfileShareButton } from "./ProfileShareButton";

const Calendar = CalendarIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;
const Mail = MailIcon as unknown as (p: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (p: { class?: string }) => JSX.Element;
const Settings = SettingsIcon as unknown as (p: {
  class?: string;
}) => JSX.Element;
type Tab = "about" | "badges" | "certificates" | "purchases";

export function ProfileHeader(props: {
  profile: UniventsProfile;
  name: string;
  handle?: string;
  ownProfile: boolean;
  profileUrl: string;
  activeTab: Tab;
  onTabChange: (tab: Tab) => void;
}) {
  const bannerUrl = () => resolveStorageUrl(props.profile.bannerUrl);

  return (
    <section class="relative w-full border-b border-border bg-card shadow-md">
      <div class="h-40 w-full overflow-hidden bg-background bg-linear-to-br from-primary/40 via-primary/15 to-muted sm:h-48 md:h-56">
        {bannerUrl() && (
          <img
            src={bannerUrl()}
            alt=""
            class="size-full object-cover object-center"
          />
        )}
      </div>
      <div class="md:hidden">
        {props.ownProfile && (
          <ProfileShareButton
            name={props.name}
            url={props.profileUrl}
            class="absolute right-4 top-4 size-11 rounded-full border border-border/40 bg-background/90 text-foreground shadow-lg backdrop-blur-sm"
          />
        )}
        <div class="relative mx-auto px-4">
          <div class="relative z-10 -mt-12 flex flex-col items-center">
            <Avatar
              name={props.name}
              imageUrl={props.profile.pfpUrl}
              size="mobile"
            />
          </div>
          <div class="mt-3 pb-4 text-center">
            <Identity
              profile={props.profile}
              name={props.name}
              handle={props.handle}
              centered
            />
            <MemberSince createdAt={props.profile.createdAt} centered />
            {props.profile.contactEmail && (
              <a
                href={`mailto:${props.profile.contactEmail}`}
                class="mt-1 inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
              >
                <Mail class="size-3.5" />
                {props.profile.contactEmail}
              </a>
            )}
          </div>
          {props.ownProfile && <Actions mobile />}
        </div>
      </div>
      <div class="mx-auto hidden max-w-7xl px-4 md:block">
        <div class="relative pb-6">
          <div class="absolute left-0 -top-16">
            <Avatar
              name={props.name}
              imageUrl={props.profile.pfpUrl}
              size="desktop"
            />
          </div>
          {props.ownProfile && (
            <Actions name={props.name} profileUrl={props.profileUrl} />
          )}
          <div class="pt-18">
            <Identity
              profile={props.profile}
              name={props.name}
              handle={props.handle}
            />
            <MemberSince createdAt={props.profile.createdAt} />
          </div>
        </div>
      </div>
      <div class="relative mx-auto max-w-7xl">
        <nav class="flex overflow-x-auto px-4 pr-12 sm:pr-4">
          <For each={tabs(props.ownProfile)}>
            {([tab, label]) => (
              <button
                type="button"
                class={cn(
                  "shrink-0 whitespace-nowrap border-b-2 px-4 py-3 text-sm font-medium transition-colors",
                  props.activeTab === tab
                    ? "border-primary text-primary"
                    : "border-transparent text-muted-foreground hover:text-foreground",
                )}
                onClick={() => props.onTabChange(tab)}
              >
                {label}
              </button>
            )}
          </For>
        </nav>
        {props.ownProfile && (
          <div
            aria-hidden="true"
            class="pointer-events-none absolute inset-y-0 right-0 w-12 bg-linear-to-r from-transparent to-card sm:hidden"
          />
        )}
      </div>
    </section>
  );
}

function tabs(own: boolean): [Tab, string][] {
  return [
    ["about", "Sobre"],
    ["badges", "Crachás"],
    ...(own
      ? [
        ["certificates", "Certificados"],
        ["purchases", "Compras"],
      ]
      : []),
  ] as [Tab, string][];
}

function Identity(props: {
  profile: UniventsProfile;
  name: string;
  handle?: string;
  centered?: boolean;
}) {
  return (
    <div class={cn(props.centered ? "text-center" : "text-left")}>
      <div
        class={cn(
          "flex items-baseline gap-2",
          props.centered && "justify-center",
        )}
      >
        <h1
          class={cn(
            "font-semibold tracking-tight text-card-foreground",
            props.centered ? "text-2xl" : "text-3xl",
          )}
        >
          {props.name}
        </h1>
        {props.profile.pronouns && (
          <span class="text-sm text-muted-foreground">
            {props.profile.pronouns}
          </span>
        )}
      </div>
      {props.handle && (
        <p class="text-sm text-muted-foreground">
          @{props.handle.replace(/^@/, "")}
        </p>
      )}
      {props.profile.legalName &&
        props.profile.legalName !== props.name &&
        props.profile.visibility?.hideLegalName === false && (
          <p class="text-sm text-muted-foreground">{props.profile.legalName}</p>
        )}
      {(props.profile.role || props.profile.organization) && (
        <p
          class={cn(
            "text-muted-foreground",
            props.centered ? "mt-1 text-sm" : "text-[15px]",
          )}
        >
          {[props.profile.role, props.profile.organization]
            .filter(Boolean)
            .join(" · ")}
        </p>
      )}
    </div>
  );
}

function MemberSince(props: { createdAt?: string; centered?: boolean }) {
  return (
    <p
      class={cn(
        "mt-1 flex items-center gap-1.5 text-sm text-muted-foreground",
        props.centered && "justify-center",
      )}
    >
      <Calendar class="size-3.5" />
      Membro desde{" "}
      {props.createdAt
        ? new Date(props.createdAt).toLocaleDateString("pt-BR", {
          month: "long",
          year: "numeric",
        })
        : "—"}
    </p>
  );
}

function Avatar(props: {
  name: string;
  imageUrl?: string | null;
  size: "mobile" | "desktop";
}) {
  const url = () => resolveStorageUrl(props.imageUrl);

  return (
    <div
      class={cn(
        "relative flex shrink-0 items-center justify-center overflow-hidden rounded-full border-4 border-background bg-background text-2xl font-semibold shadow-xl",
        props.size === "mobile" ? "size-24" : "size-32",
      )}
    >
      {url() ? (
        <img
          src={url()}
          alt={props.name}
          class="size-full object-cover"
        />
      ) : (
        <span class="flex size-full items-center justify-center rounded-full bg-muted text-muted-foreground">
          {props.name.slice(0, 2).toUpperCase()}
        </span>
      )}
    </div>
  );
}

function Actions(props: {
  mobile?: boolean;
  name?: string;
  profileUrl?: string;
}) {
  const mobile = untrack(() => props.mobile);
  if (mobile)
    return (
      <div class="mb-4 flex gap-2">
        <Link
          to="/profile/edit"
          class={cn(
            "inline-flex h-11 flex-1 items-center justify-center gap-2 rounded-lg bg-primary px-4 text-sm text-primary-foreground shadow-sm",
            "transition-all active:scale-95",
          )}
        >
          <Pencil class="size-4" />
          Editar perfil
        </Link>
        <Link
          to="/profile/config"
          class={cn(
            "inline-flex size-11 shrink-0 items-center justify-center rounded-lg border border-border bg-background shadow-sm",
            "transition-all active:scale-95",
          )}
          aria-label="Configurações do perfil"
        >
          <Settings class="size-4" />
        </Link>
      </div>
    );
  return (
    <div class="absolute right-0 top-4 flex gap-2">
      <ProfileShareButton
        name={props.name ?? ""}
        url={props.profileUrl ?? ""}
        class={cn(
          "size-10 rounded-md border border-border bg-background text-foreground shadow-sm hover:bg-accent",
          "transition-all active:scale-95",
        )}
      />
      <Link
        to="/profile/edit"
        class={cn(
          "inline-flex h-10 items-center gap-1.5 rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground shadow-sm hover:bg-primary/90",
          "transition-all active:scale-95",
        )}
      >
        <Pencil class="size-4" />
        Editar perfil
      </Link>
      <Link
        to="/profile/config"
        class={cn(
          "inline-flex size-10 items-center justify-center rounded-md border border-border bg-background text-foreground shadow-sm hover:bg-accent",
          "transition-all active:scale-95",
        )}
        aria-label="Configurações"
      >
        <Settings class="size-4" />
      </Link>
    </div>
  );
}
