import { Mail } from "lucide-react";
import type { EventI } from "@/features/events/model";

interface ContactSectionProps {
  event: EventI;
}

export function ContactSection({ event }: ContactSectionProps) {
  if (!event.contact_email) return null;

  return (
    <section className="w-full py-10 sm:py-14 md:py-16">
      <div className="max-w-xl mx-auto px-4 sm:px-6 text-center">
        {/* Title */}
        <h2 className="text-2xl sm:text-3xl md:text-4xl font-semibold text-foreground tracking-tight">
          Tem alguma dúvida?
        </h2>

        {/* Description */}
        <p className="mt-3 sm:mt-4 text-sm sm:text-base text-muted-foreground leading-relaxed">
          Nossa equipe organizadora está à disposição para ajudar com dúvidas
          sobre inscrições, palestras ou patrocínios.
        </p>

        {/* Email button */}
        <a
          href={`mailto:${event.contact_email}`}
          className="mt-6 sm:mt-8 inline-flex items-center gap-3 px-6 sm:px-8 py-3.5 sm:py-4 rounded-2xl bg-muted/80 hover:bg-muted transition-colors duration-200"
        >
          <Mail className="w-5 h-5 sm:w-6 sm:h-6 text-foreground" />
          <span className="text-sm sm:text-base font-medium text-foreground">
            {event.contact_email}
          </span>
        </a>
      </div>
    </section>
  );
}
