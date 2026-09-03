import { createFileRoute } from "@tanstack/react-router";
import { AnimatePresence, motion } from "motion/react";
import { z } from "zod";
import { useHomeEasterEgg } from "@/features/home-easter-egg/use-home-easter-egg";
import { cn } from "@/shared/lib/utils";
import { Logo } from "@/shared/ui/logo";
import { Footer } from "@/widgets/landing/ui/Footer";
import { ModeSelector } from "@/widgets/landing/ui/ModeSelector";
import { OrganizerView } from "@/widgets/landing/ui/OrganizerView";
import { ParticipantView } from "@/widgets/landing/ui/ParticipantView";

const searchSchema = z.object({
  as: z.enum(["guest", "host"]).optional().default("guest"),
});

export const Route = createFileRoute("/")({
  component: Index,
  validateSearch: searchSchema,
});

export type Mode = "guest" | "host";

function Index() {
  const { as } = Route.useSearch();
  const navigate = Route.useNavigate();
  const isAuthenticated =
    Route.useRouteContext().auth?.isAuthenticated ?? false;
  const handleLogoClick = useHomeEasterEgg(() => {
    window.location.assign("/easter-egg-slide/");
  });

  const setMode = (mode: Mode) => {
    if (mode === as) return;
    void navigate({
      search: (prev) => ({ ...prev, as: mode }),
      replace: true,
    });
  };

  return (
    <div
      className={cn(
        "min-h-screen antialiased selection:bg-primary/10 selection:text-primary",
        "bg-background text-foreground overflow-x-hidden relative",
        "pt-6 md:pb-12 pb-16",
      )}
    >
      <div className="px-4 sm:px-6 lg:px-8 relative z-10">
        <div className="pb-4 md:pb-10">
          <div className="max-w-5xl mx-auto flex flex-col items-center">
            <div className="mb-8 md:mb-12 mt-4 md:mt-8 w-32 md:w-48">
              <Logo
                variant="complete"
                className="select-none"
                onMouseDown={(event) => event.preventDefault()}
                onClick={handleLogoClick}
              />
            </div>
            <ModeSelector
              current={as}
              onChange={setMode}
              isAuthenticated={isAuthenticated}
            />
          </div>
        </div>

        <main className="pb-24 md:pb-32">
          <AnimatePresence mode="wait">
            <motion.div
              key={as}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3, ease: [0.25, 0.1, 0.25, 1] }}
            >
              {as === "guest" ? <ParticipantView /> : <OrganizerView />}
            </motion.div>
          </AnimatePresence>
        </main>
      </div>

      <Footer />
    </div>
  );
}
