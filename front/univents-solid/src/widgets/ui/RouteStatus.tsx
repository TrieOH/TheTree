import type { JSX } from "@solidjs/web";
import { Link, useRouter, type ErrorComponentProps } from "@tanstack/solid-router";
import { Button } from "@trieoh/ui-solid";
import { classifyError } from "@/shared/lib/service-error";
import { serviceHealth } from "@/shared/lib/service-health";

function StatusPage(props: {
  title: string;
  description: string;
  actions: JSX.Element;
}) {
  return (
    <section class="mx-auto flex min-h-[60vh] max-w-md flex-col items-center justify-center gap-4 px-6 py-16 text-center">
      <h1 class="text-2xl font-semibold tracking-tight">{props.title}</h1>
      <p class="text-sm text-muted-foreground">{props.description}</p>
      <div class="flex flex-wrap justify-center gap-2 pt-2">{props.actions}</div>
    </section>
  );
}

function HomeLink() {
  return (
    <Link
      to="/"
      search={{ as: "guest" }}
      class="inline-flex h-9 items-center rounded-md border border-border bg-background px-4 text-sm font-medium hover:bg-muted"
    >
      Ir para o início
    </Link>
  );
}

/**
 * A route that failed to load or render. It takes the place of that route only
 * (inside the app layout), so the dock and everything around it stay usable.
 * Failures from a backend read as "try again"; anything else is a bug, and says
 * so instead of showing a raw message.
 */
export function RouteError(props: ErrorComponentProps) {
  const router = useRouter();
  const failure = () => classifyError(props.error);
  const fromService = () => failure().transient || failure().kind === "unauthorized";

  const retry = () => {
    serviceHealth.requestRetry();
    props.reset();
    void router.invalidate();
  };

  return (
    <StatusPage
      title={fromService() ? "Não foi possível abrir esta página" : "Algo deu errado"}
      description={
        fromService()
          ? failure().message
          : "Ocorreu um erro inesperado ao exibir esta página. Tente novamente."
      }
      actions={
        <>
          <Button onClick={retry}>Tentar novamente</Button>
          <HomeLink />
        </>
      }
    />
  );
}

/** The root's own error boundary: nothing above it provides the page frame. */
export function RootError(props: ErrorComponentProps) {
  return (
    <main id="main-content" tabindex="-1">
      <RouteError {...props} />
    </main>
  );
}

export function NotFoundPage() {
  return (
    <StatusPage
      title="Página não encontrada"
      description="O endereço que você abriu não existe ou foi movido."
      actions={<HomeLink />}
    />
  );
}
