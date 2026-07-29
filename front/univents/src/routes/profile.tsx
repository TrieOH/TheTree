import { createFileRoute } from "@tanstack/react-router";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { AccountSessionContent } from "@/features/profile/ui/account-session-content";
import { AppearancePreferencesContent } from "@/features/profile/ui/appearance-preferences-content";
import { LogoutCard } from "@/features/profile/ui/logout-card";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/shared/ui/shadcn/accordion";

export const Route = createFileRoute("/profile")({
  beforeLoad: requireAuth,
  component: ProfilePage,
});

function ProfilePage() {
  return (
    <main className="h-dvh overflow-hidden bg-background">
      <aside className="mx-auto h-full max-w-7xl px-4 py-6 md:px-6 md:py-8 lg:grid lg:h-full lg:grid-cols-[360px_minmax(0,1fr)] lg:gap-6">
        {/* left */}
        <Accordion className="space-y-0 lg:sticky lg:top-6 lg:self-start lg:pr-1">
          <AccordionItem
            value="account"
            className="border-b border-border last:border-b-0"
          >
            <AccordionTrigger className="px-0 hover:no-underline">
              Conta e dados da sessão
            </AccordionTrigger>
            <AccordionContent className="px-0">
              <AccountSessionContent />
            </AccordionContent>
          </AccordionItem>

          <AccordionItem
            value="logout"
            className="border-b border-border last:border-b-0"
          >
            <AccordionTrigger className="px-0 hover:no-underline">
              Logout
            </AccordionTrigger>
            <AccordionContent className="px-0">
              <LogoutCard />
            </AccordionContent>
          </AccordionItem>
        </Accordion>

        {/* right */}
        <section className="min-h-0 lg:overflow-y-auto lg:pl-1">
          <Accordion className="space-y-0">
            <AccordionItem
              value="appearance"
              className="border-b border-border last:border-b-0"
            >
              <AccordionTrigger className="px-0 hover:no-underline">
                Aparência e Preferências
              </AccordionTrigger>
              <AccordionContent className="px-0">
                <AppearancePreferencesContent />
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </section>
      </aside>
    </main>
  );
}
