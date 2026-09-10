import { createFileRoute, Link } from "@tanstack/solid-router";
import { createSignal, onSettled } from "solid-js";
import type { JSX } from "@solidjs/web";

import ArrowLeftIcon from "~icons/lucide/arrow-left";
import LogOutIcon from "~icons/lucide/log-out";
import PaletteIcon from "~icons/lucide/palette";
import PencilIcon from "~icons/lucide/pencil";
import SlidersHorizontalIcon from "~icons/lucide/sliders-horizontal";
import UserIcon from "~icons/lucide/user-round";

import { requireAuth } from "@/features/auths/lib/route-guard";

import {
  readInplaceEditPreference,
  saveInplaceEditPreference,
} from "@/features/profile/lib/preferences";

import { Accordion } from "@/shared/ui/Accordion";

// import { LogoutCard } from "./logout-card";
import { AccountSessionContent } from "@/features/profile/ui/config/AccountSessionContent";
import { AppearancePreferencesContent } from "@/features/profile/ui/config/AppearancePreferencesContent";
import { LogoutCard } from "@/features/profile/ui/config/LogoutCard";

export const Route = createFileRoute("/profile/config")({
  beforeLoad: requireAuth,
  head: () => ({
    meta: [
      {
        title: "Configurações do perfil - Univents",
      },
    ],
  }),
  component: ProfileConfigPage,
});

type IconComponent = () => JSX.Element;

const ArrowLeft =
  ArrowLeftIcon as unknown as IconComponent;

const LogOut =
  LogOutIcon as unknown as IconComponent;

const Palette =
  PaletteIcon as unknown as IconComponent;

const Pencil =
  PencilIcon as unknown as IconComponent;

const SlidersHorizontal =
  SlidersHorizontalIcon as unknown as IconComponent;

const User =
  UserIcon as unknown as IconComponent;

function ProfileConfigPage() {
  const [
    inplaceEditEnabled,
    setInplaceEditEnabled,
  ] = createSignal(false);

  onSettled(() => {
    setInplaceEditEnabled(readInplaceEditPreference());
  });

  const toggleInplaceEditing = () => {
    setInplaceEditEnabled((enabled) => {
      const next = !enabled;

      saveInplaceEditPreference(next);

      return next;
    });
  };

  const accountItems = () => [
    {
      value: "account",
      title: <AccordionTitle icon={User} title="Conta e dados da sessão" />,
      content: <AccountSessionContent />,
    },
    {
      value: "logout",
      title: <AccordionTitle icon={LogOut} title="Logout" />,
      content: <LogoutCard />,
    },
  ];

  const preferenceItems = () => [
    {
      value: "appearance",
      title: <AccordionTitle icon={Palette} title="Aparência e Preferências" />,
      content: <AppearancePreferencesContent />,
    },
    {
      value: "editing",
      title: <AccordionTitle icon={Pencil} title="Modo de edição" />,
      content: (
        <div class="flex min-w-0 flex-col gap-3 rounded-sm border border-border bg-background p-3 sm:flex-row sm:items-center sm:justify-between sm:gap-4">
          <div class="min-w-0">
            <p class="text-xs leading-relaxed text-muted-foreground">
              Ativa os atalhos de edição disponíveis
              diretamente nas páginas compatíveis.
            </p>
          </div>
          <button
            type="button"
            aria-pressed={inplaceEditEnabled() ? "true" : "false"}
            class={`h-9 w-full shrink-0 rounded-md border px-4 text-sm font-medium transition-colors sm:w-auto ${inplaceEditEnabled()
              ? "border-primary bg-primary text-primary-foreground"
              : "border-border bg-background text-foreground hover:bg-muted"
              }`}
            onClick={toggleInplaceEditing}
          >
            {inplaceEditEnabled() ? "Desativar" : "Ativar"}
          </button>
        </div>
      ),
    },
  ];

  return (
    <main class="min-h-dvh min-w-0 overflow-x-clip bg-background pb-16 lg:flex lg:h-dvh lg:flex-col lg:overflow-hidden">
      {/* HEADER */}

      <header class="sticky top-0 z-40 shrink-0 overflow-hidden border-b border-border bg-card/95 shadow-sm backdrop-blur">
        <div class="relative mx-auto flex min-w-0 max-w-7xl items-center gap-3 px-3 py-4 sm:gap-4 sm:px-4 md:px-6 md:py-5">
          <Link
            to="/profile"
            search={{
              tab: "about",
            }}
            aria-label="Voltar ao perfil"
            class="flex size-10 shrink-0 items-center justify-center rounded-full border border-border bg-background shadow-sm transition-colors hover:bg-muted [&>svg]:size-4"
          >
            <ArrowLeft />
          </Link>

          <div class="min-w-0">
            <h1 class="truncate text-lg font-bold sm:text-xl">
              Configurações do perfil
            </h1>

            <p class="truncate text-xs text-muted-foreground sm:text-sm">
              Personalize sua conta, aparência e
              preferências.
            </p>
          </div>
        </div>
      </header>

      {/* CONTENT */}

      <div class="mx-auto grid w-full min-w-0 max-w-7xl gap-4 px-3 py-4 sm:px-4 md:px-6 md:py-5 lg:min-h-0 lg:flex-1 lg:grid-cols-[minmax(280px,360px)_minmax(0,1fr)] lg:py-4">
        {/* ACCOUNT */}
        <aside class="min-w-0 rounded-sm border border-border bg-card p-4 shadow-md shadow-foreground/5 lg:self-start">
          <SectionHeading
            icon={User}
            title="Conta"
            description="Sessão, dados pessoais e acesso."
          />
          <Accordion
            items={accountItems()}
            defaultValue="account"
          />
        </aside>

        {/* PREFERENCES */}

        <section class="profile-preferences-scroll min-w-0 rounded-sm border border-border bg-card p-4 pb-14 shadow-md shadow-foreground/5 lg:h-full lg:overflow-y-auto">
          <SectionHeading
            icon={SlidersHorizontal}
            title="Preferências"
            description="Ajuste como o Univents funciona para você."
          />
          <Accordion
            items={preferenceItems()}
            multiple
          />
        </section>
      </div>
    </main>
  );
}

function AccordionTitle(props: {
  icon: IconComponent;
  title: string;
}) {
  const Icon = props.icon;

  return (
    <span class="flex min-w-0 items-center gap-2">
      <span class="flex shrink-0 items-center justify-center text-muted-foreground [&>svg]:size-4">
        <Icon />
      </span>

      <span class="truncate">
        {props.title}
      </span>
    </span>
  );
}

function SectionHeading(props: {
  icon: IconComponent;
  title: string;
  description: string;
}) {
  const Icon = props.icon;

  return (
    <div class="mb-4 flex min-w-0 items-center gap-3 border-b border-border pb-4">
      <div class="flex size-9 shrink-0 items-center justify-center rounded-sm bg-primary/10 text-primary [&>svg]:size-4">
        <Icon />
      </div>

      <div class="min-w-0">
        <h2 class="truncate font-semibold">
          {props.title}
        </h2>
        <p class="truncate text-xs text-muted-foreground">
          {props.description}
        </p>
      </div>
    </div>
  );
}