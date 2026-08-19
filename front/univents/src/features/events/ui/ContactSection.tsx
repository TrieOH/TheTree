import { Mail } from "lucide-react";
import type { EventI } from "@/features/events/model";

interface ContactSectionProps {
  event: EventI;
}

export function ContactSection({ event }: ContactSectionProps) {
  if (!event.contact_email) return null;

  return (
    <section className="w-full py-10 sm:py-14 md:py-16">
      <div className="mx-auto max-w-xl px-4 text-center sm:px-6">
        <h2 className="text-2xl font-semibold tracking-tight text-foreground sm:text-3xl md:text-4xl">
          Ficou com alguma dúvida?
        </h2>

        <p className="mt-3 text-sm leading-relaxed text-muted-foreground sm:mt-4 sm:text-base">
          A organização do evento pode ajudar com informações sobre inscrições,
          programação, local e outras dúvidas.
        </p>

        <a
          href={`mailto:${event.contact_email}`}
          className="mx-auto mt-3 inline-flex max-w-full items-center gap-3 rounded-md border border-border bg-card px-5 py-3 text-left shadow-sm transition-colors hover:bg-muted"
        >
          <Mail className="h-5 w-5 shrink-0 text-primary" />
          <span className="truncate text-sm font-medium text-foreground sm:text-base">
            {event.contact_email}
          </span>
        </a>
      </div>
    </section>
  );
}
