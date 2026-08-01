import { Check } from "lucide-react";
import { cn } from "@/shared/lib/utils";

export function SignatureRequestConfirmation({
  timestamp,
}: {
  timestamp: string;
}) {
  const formatted = new Date(timestamp).toLocaleString("pt-BR", {
    month: "short",
    day: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    timeZone: "UTC",
    timeZoneName: "short",
  });

  return (
    <main className="flex min-h-screen items-center justify-center bg-muted px-3 pb-24 sm:px-4">
      <section
        className={cn(
          "flex w-full min-w-0 max-w-sm flex-col items-center gap-4",
          "rounded-2xl border border-border bg-card px-6 py-8 text-center shadow-sm",
          "sm:px-8 sm:py-10",
        )}
      >
        <span className="flex size-14 items-center justify-center rounded-2xl bg-primary text-primary-foreground">
          <Check className="size-7" />
        </span>
        <h1 className="text-xl font-semibold text-card-foreground">
          Documento assinado
        </h1>
        <p className="text-sm text-muted-foreground">
          Sua assinatura foi adicionada com segurança ao contrato do evento.
        </p>
        <div className="w-full rounded-lg bg-muted px-4 py-2 text-xs text-muted-foreground">
          Data e hora: {formatted}
        </div>
        <button
          type="button"
          onClick={() => window.close()}
          className="text-sm font-medium text-primary hover:underline"
        >
          Você pode fechar esta aba
        </button>
      </section>
    </main>
  );
}
