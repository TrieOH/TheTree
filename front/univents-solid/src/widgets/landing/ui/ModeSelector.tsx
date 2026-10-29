import type { Mode } from "@/routes/index";
export function ModeSelector(props: {
  current: Mode;
  onChange: (mode: Mode) => void;
  isAuthenticated?: boolean;
}) {
  const ctaLabel = () =>
    props.current === "guest"
      ? props.isAuthenticated
        ? "Continuar explorando"
        : "Criar conta grátis"
      : "Começar a organizar";
  return (
    <div class="flex flex-col items-center gap-6 md:gap-8">
      <div class="relative inline-grid grid-cols-2 rounded-full bg-muted p-1">
        <div
          class="pointer-events-none absolute top-1 bottom-1 left-1 w-[calc((100%_-_0.5rem)/2)] rounded-full bg-primary shadow-sm transition-transform duration-200 ease-out motion-reduce:transition-none"
          style={{ transform: props.current === "host" ? "translateX(100%)" : "translateX(0)" }}
          aria-hidden="true"
        />
        <button
          class={`relative z-10 overflow-hidden rounded-full border-0 bg-transparent px-4 py-2 text-xs shadow-none transition-colors duration-200 md:px-6 md:py-2.5 md:text-sm ${props.current === "guest" ? "text-primary-foreground" : "text-muted-foreground hover:bg-transparent"}`}
          onClick={() => props.onChange("guest")}
        >
          <span
            class={`relative z-10 transition-transform duration-200 ${props.current === "guest" ? "scale-100" : "scale-[0.98]"}`}
          >
            Quero Participar
          </span>
        </button>
        <button
          class={`relative z-10 overflow-hidden rounded-full border-0 bg-transparent px-4 py-2 text-xs shadow-none transition-colors duration-200 md:px-6 md:py-2.5 md:text-sm ${props.current === "host" ? "text-primary-foreground" : "text-muted-foreground hover:bg-transparent"}`}
          onClick={() => props.onChange("host")}
        >
          <span
            class={`relative z-10 transition-transform duration-200 ${props.current === "host" ? "scale-100" : "scale-[0.98]"}`}
          >
            Quero Organizar
          </span>
        </button>
      </div>
      <h1 class="px-2 text-center font-heading text-3xl font-semibold leading-[1.1] tracking-tight sm:text-4xl md:text-6xl">
        {props.current === "guest" ? (
          <>
            Descubra eventos,
            <br />
            <span class="text-muted-foreground">viva experiências.</span>
          </>
        ) : (
          <>
            Seus eventos,
            <br />
            <span class="text-muted-foreground">sob controle total.</span>
          </>
        )}
      </h1>
      <p class="max-w-2xl px-4 text-center text-sm leading-relaxed text-muted-foreground md:text-base">
        {props.current === "guest"
          ? "Veja eventos em destaque, descubra o que está acontecendo perto de você e entre direto no fluxo certo."
          : "Tenha visão clara do evento, da operação e da receita em um só lugar."}
      </p>
      <div class="flex flex-wrap justify-center gap-2 px-4">
        <span class="rounded-full border border-border bg-background px-3 py-1 text-xs text-muted-foreground">
          {ctaLabel()}
        </span>
        <span class="rounded-full border border-border bg-background px-3 py-1 text-xs text-muted-foreground">
          {props.current === "guest"
            ? "Eventos em destaque"
            : "Gestão centralizada"}
        </span>
        <span class="rounded-full border border-border bg-background px-3 py-1 text-xs text-muted-foreground">
          {props.current === "guest" ? "Compra rápida" : "Operação ao vivo"}
        </span>
      </div>
    </div>
  );
}
