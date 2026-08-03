import { Link } from "@tanstack/react-router";
import type { LucideIcon } from "lucide-react";
import {
  ArrowLeft,
  Award,
  LogOut,
  Palette,
  Pencil,
  Settings2,
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
    <main className="min-h-dvh bg-background pb-28 lg:flex lg:h-dvh lg:flex-col lg:overflow-hidden lg:pb-0">
      <header className="relative shrink-0 overflow-hidden border-b border-border bg-card shadow-sm">
        <div className="absolute inset-0 bg-linear-to-r from-primary/10 via-transparent to-primary/5" />
        <div className="relative mx-auto flex max-w-7xl items-center gap-4 px-4 py-5 md:px-6">
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
          <div className="flex size-11 shrink-0 items-center justify-center rounded-xl bg-primary text-primary-foreground shadow-md shadow-primary/20">
            <Settings2 className="size-5" />
          </div>
          <div className="min-w-0">
            <h1 className="text-xl font-bold sm:text-2xl">
              Configurações do perfil
            </h1>
            <p className="text-sm text-muted-foreground">
              Personalize sua conta, aparência e preferências.
            </p>
          </div>
        </div>
      </header>
      <div className="mx-auto grid w-full max-w-7xl gap-4 px-4 py-5 md:px-6 lg:min-h-0 lg:flex-1 lg:grid-cols-[360px_minmax(0,1fr)]">
        <aside className="rounded-xl border border-border bg-card p-4 shadow-md shadow-foreground/5 lg:min-h-0 lg:overflow-y-auto">
          <SectionHeading
            icon={UserRound}
            title="Conta"
            description="Sessão, dados pessoais e acesso."
          />
          <Accordion multiple defaultValue={LEFT_PROFILE_OPTIONS}>
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
        <section className="rounded-xl border border-border bg-card p-4 shadow-md shadow-foreground/5 lg:min-h-0 lg:overflow-y-auto">
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
              <div className="flex items-center justify-between gap-4 rounded-lg border border-border bg-background p-3">
                <div>
                  <p className="text-sm font-medium">Edição rápida</p>
                  <p className="text-xs text-muted-foreground">
                    Ativa os atalhos de edição disponíveis nas áreas do sistema.
                  </p>
                </div>
                <Button
                  type="button"
                  variant={inplaceEditEnabled ? "default" : "outline"}
                  size="sm"
                  onClick={() =>
                    setInplaceEditEnabled((enabled) => {
                      const next = !enabled;
                      saveInplaceEditPreference(next);
                      return next;
                    })
                  }
                >
                  {inplaceEditEnabled ? "Ativado" : "Ativar"}
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
      className="mb-2 rounded-lg border border-border bg-background px-3 last:mb-0"
    >
      <AccordionTrigger className="py-3 hover:no-underline">
        <span className="flex items-center gap-2">
          <Icon className="size-4 text-muted-foreground" />
          {title}
        </span>
      </AccordionTrigger>
      <AccordionContent className="pb-3">{children}</AccordionContent>
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
    <div className="mb-4 flex items-center gap-3 border-b border-border pb-4">
      <div className="flex size-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
        <Icon className="size-4" />
      </div>
      <div>
        <h2 className="font-semibold">{title}</h2>
        <p className="text-xs text-muted-foreground">{description}</p>
      </div>
    </div>
  );
}
