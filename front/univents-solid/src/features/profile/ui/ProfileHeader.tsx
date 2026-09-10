import type { JSX } from "@solidjs/web";
import { untrack } from "solid-js";
import { Link } from "@tanstack/solid-router";
import CalendarIcon from "~icons/lucide/calendar";
import MailIcon from "~icons/lucide/mail";
import PencilIcon from "~icons/lucide/pencil";
import SettingsIcon from "~icons/lucide/settings";
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
  return (
    <section class="relative w-full border-b border-border bg-card shadow-md">
      <div class="h-40 w-full overflow-hidden bg-background bg-linear-to-br from-primary/40 via-primary/15 to-muted sm:h-48 md:h-56">
        {props.profile.bannerUrl && (
          <img
            src={props.profile.bannerUrl}
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
            class="absolute right-4 top-4 flex!"
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
          {tabs(props.ownProfile).map(([tab, label]) => (
            <button
              type="button"
              class={`shrink-0 whitespace-nowrap border-b-2 px-4 py-3 text-sm font-medium transition-colors ${props.activeTab === tab ? "border-primary text-primary" : "border-transparent text-muted-foreground hover:text-foreground"}`}
              onClick={() => props.onTabChange(tab)}
            >
              {label}
            </button>
          ))}
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
    <div class={props.centered ? "text-center" : "text-left"}>
      <div
        class={`flex items-baseline gap-2 ${props.centered ? "justify-center" : ""}`}
      >
        <h1
          class={`${props.centered ? "text-2xl" : "text-3xl"} font-semibold tracking-tight text-card-foreground`}
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
          class={
            props.centered
              ? "mt-1 text-sm text-muted-foreground"
              : "text-[15px] text-muted-foreground"
          }
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
      class={`mt-1 flex items-center gap-1.5 text-sm text-muted-foreground ${props.centered ? "justify-center" : ""}`}
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
  return (
    <div
      class={`relative shrink-0 ${props.size === "mobile" ? "size-24" : "size-32"} flex items-center justify-center overflow-hidden rounded-full border-4 border-background bg-background text-2xl font-semibold shadow-xl`}
    >
      {props.imageUrl ? (
        <img
          src={props.imageUrl}
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
          class="inline-flex h-11 flex-1 items-center justify-center gap-2 rounded-lg bg-primary px-4 text-sm text-primary-foreground shadow-sm"
        >
          <Pencil class="size-4" />
          Editar perfil
        </Link>
        <Link
          to="/profile/config"
          class="inline-flex size-11 shrink-0 items-center justify-center rounded-lg border border-border bg-background shadow-sm"
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
        class="flex! size-10 rounded-md shadow-sm"
      />
      <Link
        to="/profile/edit"
        class="inline-flex h-10 items-center gap-2 rounded-md bg-primary px-4 text-sm text-primary-foreground shadow-sm"
      >
        <Pencil class="size-4" />
        Editar perfil
      </Link>
      <Link
        to="/profile/config"
        class="inline-flex size-10 items-center justify-center rounded-md border border-border bg-background shadow-sm"
        aria-label="Configurações"
      >
        <Settings class="size-4" />
      </Link>
    </div>
  );
}
