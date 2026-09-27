import { animate } from "motion/mini";
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
  let content!: HTMLDivElement;
  let activeIndicator!: HTMLDivElement;
  let guestButton!: HTMLButtonElement;
  let hostButton!: HTMLButtonElement;
  let resizeObserver: ResizeObserver | null = null;

  const moveIndicator = (mode: Mode, immediate = false) => {
    const button = mode === "guest" ? guestButton : hostButton;
    if (!button || !activeIndicator) return;

    const keyframes = {
      left: `${button.offsetLeft}px`,
      width: `${button.offsetWidth}px`,
    };

    if (immediate) {
      Object.assign(activeIndicator.style, keyframes);
      return;
    }

    animate(activeIndicator, keyframes, {
      duration: 0.3,
      ease: [0.22, 1, 0.36, 1],
    });
  };

  const changeMode = (mode: Mode) => {

    if (mode === props.current) return;

    props.onChange(mode);
    requestAnimationFrame(() => {
      moveIndicator(mode);
      if (
        content &&
        !window.matchMedia("(prefers-reduced-motion: reduce)").matches
      ) {
        animate(
          content,
          { opacity: [0, 1], transform: ["translateY(8px)", "translateY(0)"] },
          { duration: 0.2, ease: "easeOut" },
        );
      }
    });
  };
  return (
    <div
      ref={(element) => {
        content = element;
      }}
      class="flex flex-col items-center gap-6 md:gap-8"
    >
      <div
        ref={(element) => {
          requestAnimationFrame(() => {
            moveIndicator(props.current, true);
            resizeObserver?.disconnect();
            resizeObserver = new ResizeObserver(() =>
              moveIndicator(props.current, true),
            );
            resizeObserver.observe(element);
            resizeObserver.observe(guestButton);
            resizeObserver.observe(hostButton);
          });
        }}
        class="relative inline-flex rounded-full bg-muted p-1"
      >
        <div
          ref={(element) => {
            activeIndicator = element;
          }}
          class="pointer-events-none absolute top-1 bottom-1 left-1 rounded-full bg-primary shadow-sm"
          aria-hidden="true"
        />
        <button
          ref={(element) => {
            guestButton = element;
          }}
          class={`relative z-10 overflow-hidden rounded-full border-0 bg-transparent px-4 py-2 text-xs shadow-none transition-colors duration-200 md:px-6 md:py-2.5 md:text-sm ${props.current === "guest" ? "text-primary-foreground" : "text-muted-foreground hover:bg-transparent"}`}
          onClick={() => changeMode("guest")}
        >
          <span
            class={`relative z-10 transition-transform duration-200 ${props.current === "guest" ? "scale-100" : "scale-[0.98]"}`}
          >
            Quero Participar
          </span>
        </button>
        <button
          ref={(element) => {
            hostButton = element;
          }}
          class={`relative z-10 overflow-hidden rounded-full border-0 bg-transparent px-4 py-2 text-xs shadow-none transition-colors duration-200 md:px-6 md:py-2.5 md:text-sm ${props.current === "host" ? "text-primary-foreground" : "text-muted-foreground hover:bg-transparent"}`}
          onClick={() => changeMode("host")}
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
