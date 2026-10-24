import { createFileRoute } from "@tanstack/solid-router";
import { animate } from "motion/mini";
import { createEffect } from "solid-js";
import z from "zod";
import Logo from "@/shared/ui/Logo";
import { Footer } from "@/widgets/landing/ui/Footer";
import { ModeSelector } from "@/widgets/landing/ui/ModeSelector";
import { OrganizerView } from "@/widgets/landing/ui/OrganizerView";
import { ParticipantView } from "@/widgets/landing/ui/ParticipantView";

export type Mode = "guest" | "host";
const searchSchema = z.object({
  as: z.enum(["guest", "host"]).optional().default("guest"),
});

export const Route = createFileRoute("/")({
  head: () => ({ meta: [{ title: "Univents" }] }),
  validateSearch: (search) => searchSchema.parse(search),
  component: HomePage,
});

function HomePage() {
  const search = Route.useSearch();
  const navigate = Route.useNavigate();
  const mode = () => search().as;
  let content!: HTMLDivElement;

  createEffect(
    () => mode(),
    (_mode, previousMode) => {
      if (
        previousMode !== undefined && content &&
        !window.matchMedia("(prefers-reduced-motion: reduce)").matches
      ) {
        animate(
          content,
          { opacity: [0, 1], transform: ["translateY(10px)", "translateY(0)"] },
          { duration: 0.3, ease: "easeOut" },
        );
      }
    },
  );

  const setMode = (next: Mode) => {
    if (next !== mode())
      void navigate({
        search: (previous) => ({ ...previous, as: next }),
        replace: true,
        resetScroll: false,
      });
  };

  return (
    <div class="min-h-screen antialiased bg-background text-foreground overflow-x-hidden relative pt-6 md:pb-12 pb-16">
      <div class="px-4 sm:px-6 lg:px-8 relative z-10">
        <div class="pb-4 md:pb-10">
          <div class="mx-auto flex max-w-5xl flex-col items-center">
            <div class="mb-8 mt-4 w-32 md:mb-12 md:w-48">
              <Logo variant="complete" priority class="select-none" />
            </div>
            <ModeSelector current={mode()} onChange={setMode} />
          </div>
        </div>
        <main class="pb-24 md:pb-32">
          <div
            ref={(element) => {
              content = element;
            }}
          >
            {mode() === "guest" ? <ParticipantView /> : <OrganizerView />}
          </div>
        </main>
      </div>
      <Footer />
    </div>
  );
}
