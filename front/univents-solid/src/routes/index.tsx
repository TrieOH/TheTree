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

  const revealSections = () => {
    if (!content) return;
    const sections = content.querySelectorAll<HTMLElement>("section");
    const items = content.querySelectorAll<HTMLElement>(".reveal-item");
    const reducedMotion = window.matchMedia(
      "(prefers-reduced-motion: reduce)",
    ).matches;

    if (reducedMotion) return;

    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (!entry.isIntersecting) continue;
          const section = entry.target as HTMLElement;
          observer.unobserve(section);
          animate(
            section,
            { opacity: [0, 1], transform: ["translateY(20px)", "translateY(0)"] },
            { duration: 0.45, ease: [0.22, 1, 0.36, 1] },
          );
        }
      },
      { rootMargin: "0px 0px -8% 0px", threshold: 0.08 },
    );

    for (const section of sections) {
      section.style.opacity = "0";
      observer.observe(section);
    }

    for (const [index, item] of items.entries()) {
      item.style.opacity = "0";
      const itemObserver = new IntersectionObserver(
        ([entry]) => {
          if (!entry?.isIntersecting) return;
          itemObserver.disconnect();
          animate(
            item,
            { opacity: [0, 1], transform: ["translateY(16px)", "translateY(0)"] },
            { delay: (index % 6) * 0.08, duration: 0.4, ease: "easeOut" },
          );
        },
        { rootMargin: "0px 0px -8% 0px", threshold: 0.08 },
      );
      itemObserver.observe(item);
    }
  };

  createEffect(
    () => mode(),
    () => {
      if (
        content &&
        !window.matchMedia("(prefers-reduced-motion: reduce)").matches
      ) {
        animate(
          content,
          { opacity: [0, 1], transform: ["translateY(10px)", "translateY(0)"] },
          { duration: 0.3, ease: "easeOut" },
        );
      }
      revealSections();
    },
  );

  const setMode = (next: Mode) => {
    if (next !== mode())
      void navigate({
        search: (previous) => ({ ...previous, as: next }),
        replace: true,
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
