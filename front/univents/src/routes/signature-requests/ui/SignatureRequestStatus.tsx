import { AlertCircle } from "lucide-react";
import { cn } from "@/shared/lib/utils";

interface SignatureRequestStatusProps {
  title: string;
  message: string;
  loading?: boolean;
}

export function SignatureRequestStatus({
  title,
  message,
  loading = false,
}: SignatureRequestStatusProps) {
  return (
    <main className="flex min-h-screen items-center justify-center bg-muted px-3 pb-24 sm:px-4">
      <section
        className={cn(
          "flex w-full min-w-0 max-w-sm flex-col items-center gap-4",
          "rounded-2xl border border-border bg-card px-6 py-8 text-center shadow-sm",
          "sm:px-8 sm:py-10",
        )}
      >
        <span
          className={cn(
            "flex size-14 items-center justify-center rounded-2xl",
            loading
              ? "bg-primary/10 text-primary"
              : "bg-destructive/10 text-destructive",
          )}
        >
          {loading ? (
            <span className="size-7 animate-spin rounded-full border-2 border-current border-t-transparent" />
          ) : (
            <AlertCircle className="size-7" />
          )}
        </span>
        <h1 className="text-xl font-semibold text-card-foreground">{title}</h1>
        <p className="text-sm text-muted-foreground">{message}</p>
      </section>
    </main>
  );
}
