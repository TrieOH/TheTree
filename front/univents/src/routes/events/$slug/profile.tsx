import { createFileRoute } from "@tanstack/react-router";
import {
  Award,
  LogOut,
  type LucideIcon,
  Palette,
  Pencil,
  UserRound,
} from "lucide-react";
import { type ReactNode, useEffect, useState } from "react";
import { UserCertificationsSection } from "@/features/certifications/ui/UserCertificationsSection";
import {
  readInplaceEditPreference,
  saveInplaceEditPreference,
} from "@/features/profile/lib/preferences";
import { AccountSessionContent } from "@/features/profile/ui/account-session-content";
import { AppearancePreferencesContent } from "@/features/profile/ui/appearance-preferences-content";
import { LogoutCard } from "@/features/profile/ui/logout-card";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/shared/ui/shadcn/accordion";
import { Button } from "@/shared/ui/shadcn/button";

export const Route = createFileRoute("/events/$slug/profile")({
  component: EventProfilePage,
});

const LEFT_PROFILE_OPTIONS = ["account", "logout"];

function EventProfilePage() {
  const [inplaceEditEnabled, setInplaceEditEnabled] = useState(false);

  useEffect(() => {
    setInplaceEditEnabled(readInplaceEditPreference());
  }, []);

  const handleInplaceEditPreference = () => {
    setInplaceEditEnabled((enabled) => {
      const next = !enabled;
      saveInplaceEditPreference(next);
      return next;
    });
  };

  return (
    <main className="h-full min-h-dvh bg-background pb-28!">
      <aside className="mx-auto h-full max-w-7xl px-4 py-5 md:px-6 md:py-6 lg:grid lg:h-full lg:grid-cols-[360px_minmax(0,1fr)] lg:gap-4">
        {/* left */}
        <div className="lg:sticky lg:top-5 lg:self-start lg:pr-1">
          <Accordion
            multiple
            defaultValue={LEFT_PROFILE_OPTIONS}
            className="space-y-0"
          >
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
        </div>

        {/* right */}
        <section className="min-h-0 lg:overflow-y-auto">
          <Accordion className="space-y-0" multiple>
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
              <div className="flex items-center justify-between gap-4 rounded-lg border border-border/40 bg-card p-3">
                <div className="space-y-4!">
                  <p className="text-sm font-medium m-0!">Edição rápida</p>
                  <p className="mt-1 text-xs text-muted-foreground">
                    Ativa os atalhos de edição disponíveis nas áreas do sistema.
                  </p>
                </div>
                <Button
                  type="button"
                  variant={inplaceEditEnabled ? "default" : "outline"}
                  size="sm"
                  onClick={handleInplaceEditPreference}
                >
                  {inplaceEditEnabled ? "Ativado" : "Ativar"}
                </Button>
              </div>
            </ProfileOption>
          </Accordion>
        </section>
      </aside>
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
      className="border-b border-border/50 last:border-b-0"
    >
      <AccordionTrigger className="px-0 py-3 hover:no-underline focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0">
        <span className="flex items-center gap-2">
          <Icon className="size-4 text-muted-foreground" />
          {title}
        </span>
      </AccordionTrigger>
      <AccordionContent className="px-0">{children}</AccordionContent>
    </AccordionItem>
  );
}
