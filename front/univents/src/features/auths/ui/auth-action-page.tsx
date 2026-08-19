import { Link } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";

export function AuthActionPage({
  title,
  description,
  children,
}: {
  title?: string;
  description?: string;
  children: React.ReactNode;
}) {
  return (
    <main className="flex min-h-dvh items-center justify-center bg-background px-4 pb-28 pt-12">
      <section className="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-lg">
        <Link
          to="/auth"
          search={{ redirect: "" }}
          className="mb-6 inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="size-4" />
          Voltar ao login
        </Link>
        {title && (
          <div className="mb-6 text-center">
            <h1 className="font-heading text-3xl font-bold tracking-tight">
              {title}
            </h1>
            {description && (
              <p className="mt-2 text-sm text-muted-foreground">
                {description}
              </p>
            )}
          </div>
        )}
        {children}
      </section>
    </main>
  );
}
