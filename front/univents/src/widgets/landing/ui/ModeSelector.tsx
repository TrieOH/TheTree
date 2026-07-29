import { motion } from "motion/react";
import type { Mode } from "@/routes/index";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";

interface Props {
  current: Mode;
  onChange: (mode: Mode) => void;
  isAuthenticated?: boolean;
}

export function ModeSelector({
  current,
  onChange,
  isAuthenticated = false,
}: Props) {
  const ctaLabel =
    current === "guest"
      ? isAuthenticated
        ? "Continuar explorando"
        : "Criar conta grátis"
      : "Começar a organizar";

  return (
    <div className="flex flex-col items-center gap-6 md:gap-8">
      {/* Toggle */}
      <div className="inline-flex p-1 bg-muted rounded-full">
        <Button
          type="button"
          onClick={() => {
            onChange("guest");
          }}
          className={cn(
            "relative z-10 rounded-full px-4 py-2 md:px-6 md:py-2.5 text-xs md:text-sm",
            "border-0 bg-transparent shadow-none hover:bg-transparent",
            current === "guest"
              ? "text-primary-foreground"
              : "text-muted-foreground",
          )}
        >
          {current === "guest" && (
            <motion.div
              layoutId="activeTab"
              className="absolute inset-0 bg-primary rounded-full shadow-sm"
              transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
            />
          )}
          <span className="relative z-10">Quero Participar</span>
        </Button>
        <Button
          type="button"
          onClick={() => {
            onChange("host");
          }}
          className={cn(
            "relative z-10 rounded-full px-4 py-2 md:px-6 md:py-2.5 text-xs md:text-sm",
            "border-0 bg-transparent shadow-none hover:bg-transparent",
            current === "host"
              ? "text-primary-foreground"
              : "text-muted-foreground",
          )}
        >
          {current === "host" && (
            <motion.div
              layoutId="activeTab"
              className="absolute inset-0 bg-primary rounded-full shadow-sm"
              transition={{ type: "spring", bounce: 0.2, duration: 0.6 }}
            />
          )}
          <span className="relative z-10">Quero Organizar</span>
        </Button>
      </div>

      {/* Headline */}
      <h1 className="text-center px-2 font-heading">
        {current === "guest" ? (
          <span
            className={cn(
              "block text-3xl sm:text-4xl md:text-6xl lg:text-7xl",
              "font-semibold tracking-tight text-foreground leading-[1.1]",
            )}
          >
            Descubra eventos,
            <br />
            <span className="text-muted-foreground">viva experiências.</span>
          </span>
        ) : (
          <span
            className={cn(
              "block text-3xl sm:text-4xl md:text-6xl lg:text-7xl",
              "font-semibold tracking-tight text-foreground leading-[1.1]",
            )}
          >
            Seus eventos,
            <br />
            <span className="text-muted-foreground">sob controle total.</span>
          </span>
        )}
      </h1>
      <p className="text-center max-w-2xl text-sm md:text-base text-muted-foreground leading-relaxed px-4">
        {current === "guest"
          ? "Veja eventos em destaque, descubra o que está acontecendo perto de você e entre direto no fluxo certo."
          : "Tenha visão clara do evento, da operação e da receita em um só lugar."}
      </p>
      <div className="flex flex-wrap justify-center gap-2 px-4">
        <span className="rounded-full border border-border bg-background px-3 py-1 text-xs text-muted-foreground">
          {ctaLabel}
        </span>
        <span className="rounded-full border border-border bg-background px-3 py-1 text-xs text-muted-foreground">
          {current === "guest" ? "Eventos em destaque" : "Gestão centralizada"}
        </span>
        <span className="rounded-full border border-border bg-background px-3 py-1 text-xs text-muted-foreground">
          {current === "guest" ? "Compra rápida" : "Operação ao vivo"}
        </span>
      </div>
    </div>
  );
}
