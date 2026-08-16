import { Link } from "@tanstack/react-router";
import type { LucideIcon } from "lucide-react";
import {
  ArrowLeft,
  Award,
  LogOut,
  Palette,
  Pencil,
  SlidersHorizontal,
  UserRound,
} from "lucide-react";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";
import { UserCertificationsSection } from "@/features/certifications/ui/UserCertificationsSection";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/shared/ui/shadcn/accordion";
import { Button, buttonVariants } from "@/shared/ui/shadcn/button";
import {
  readInplaceEditPreference,
  saveInplaceEditPreference,
} from "../lib/preferences";
import { AccountSessionContent } from "./account-session-content";
import { AppearancePreferencesContent } from "./appearance-preferences-content";
import { LogoutCard } from "./logout-card";

const LEFT_PROFILE_OPTIONS = ["account"];

export function ProfileConfigPage() {
  const [inplaceEditEnabled, setInplaceEditEnabled] = useState(false);
  useEffect(() => setInplaceEditEnabled(readInplaceEditPreference()), []);

  return (
    <main className="min-h-dvh min-w-0 overflow-x-clip bg-background pb-16 lg:flex lg:h-dvh lg:flex-col lg:overflow-hidden">
      <header className="sticky top-0 z-40 shrink-0 overflow-hidden border-b border-border bg-card/95 shadow-sm backdrop-blur">
        <div className="relative mx-auto flex min-w-0 max-w-7xl items-center gap-3 px-3 py-4 sm:gap-4 sm:px-4 md:px-6 md:py-5">
          <Link
            to="/profile"
            className={buttonVariants({
              variant: "outline",
              size: "icon",
              className: "shrink-0 rounded-full bg-background shadow-sm",
            })}
            aria-label="Voltar ao perfil"
          >
            <ArrowLeft className="size-4" />
          </Link>
          <div className="min-w-0">
            <h1 className="truncate text-lg font-bold sm:text-xl">
              Configurações do perfil
            </h1>
            <p className="truncate text-xs text-muted-foreground sm:text-sm">
              Personalize sua conta, aparência e preferências.
            </p>
          </div>
        </div>
      </header>
      <div className="mx-auto grid w-full min-w-0 max-w-7xl gap-4 px-3 py-4 sm:px-4 md:px-6 md:py-5 lg:min-h-0 lg:flex-1 lg:grid-cols-[minmax(280px,360px)_minmax(0,1fr)] lg:py-4">
        <aside className="min-w-0 rounded-sm border border-border bg-card p-4 shadow-md shadow-foreground/5 lg:self-start">
          <SectionHeading
            icon={UserRound}
            title="Conta"
            description="Sessão, dados pessoais e acesso."
          />
          <Accordion defaultValue={LEFT_PROFILE_OPTIONS}>
            <ProfileOption
              value="account"
              icon={UserRound}
              title="Conta e dados da sessão"
            >
              <AccountSessionContent />
            </ProfileOption>
            <ProfileOption value="logout" icon={LogOut} title="Logout">
              <LogoutCard />
            </ProfileOption>
          </Accordion>
        </aside>
        <section className="profile-preferences-scroll min-w-0 rounded-sm border border-border bg-card p-4 shadow-md shadow-foreground/5 lg:h-full lg:overflow-y-auto pb-14">
          <SectionHeading
            icon={SlidersHorizontal}
            title="Preferências"
            description="Ajuste como o Univents funciona para você."
          />
          <Accordion multiple>
            <ProfileOption
              value="certificates"
              icon={Award}
              title="Meus certificados"
            >
              <UserCertificationsSection />
            </ProfileOption>
            <ProfileOption
              value="appearance"
              icon={Palette}
              title="Aparência e Preferências"
            >
              <AppearancePreferencesContent />
            </ProfileOption>
            <ProfileOption value="editing" icon={Pencil} title="Modo de edição">
              <div className="flex min-w-0 flex-col gap-3 rounded-sm border border-border bg-background p-3 sm:flex-row sm:items-center sm:justify-between sm:gap-4">
                <div className="min-w-0">
                  <p className="text-xs leading-relaxed text-muted-foreground">
                    Ativa os atalhos de edição disponíveis diretamente nas
                    páginas compatíveis
                  </p>
                </div>
                <Button
                  type="button"
                  variant={inplaceEditEnabled ? "default" : "outline"}
                  size="sm"
                  className="w-full shrink-0 h-9 sm:w-auto"
                  onClick={() =>
                    setInplaceEditEnabled((enabled) => {
                      const next = !enabled;
                      saveInplaceEditPreference(next);
                      return next;
                    })
                  }
                >
                  {inplaceEditEnabled ? "Desativar" : "Ativar"}
                </Button>
              </div>
            </ProfileOption>
          </Accordion>
        </section>
      </div>
    </main>
  );
}

function ProfileOption({
  value,
  icon: Icon,
  title,
  children,
}: {
  value: string;
  icon: LucideIcon;
  title: string;
  children: ReactNode;
}) {
  return (
    <AccordionItem
      value={value}
      className="mb-2 min-w-0 rounded-sm border border-border bg-background px-3 last:mb-0"
    >
      <AccordionTrigger className="py-3 hover:no-underline">
        <span className="flex min-w-0 items-center gap-2 text-left">
          <Icon className="size-4 text-muted-foreground" />
          {title}
        </span>
      </AccordionTrigger>
      <AccordionContent className="min-w-0 max-w-full overflow-x-auto pb-3">
        {children}
      </AccordionContent>
    </AccordionItem>
  );
}

function SectionHeading({
  icon: Icon,
  title,
  description,
}: {
  icon: LucideIcon;
  title: string;
  description: string;
}) {
  return (
    <div className="mb-4 flex min-w-0 items-center gap-3 border-b border-border pb-4">
      <div className="flex size-9 items-center justify-center rounded-sm bg-primary/10 text-primary">
        <Icon className="size-4" />
      </div>
      <div className="min-w-0">
        <h2 className="truncate font-semibold">{title}</h2>
        <p className="truncate text-xs text-muted-foreground">{description}</p>
      </div>
    </div>
  );
}
