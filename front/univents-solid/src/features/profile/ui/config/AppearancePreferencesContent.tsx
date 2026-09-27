import { Show, untrack } from "solid-js";
import type { JSX } from "@solidjs/web";

import { useTheme } from "@/shared/lib/theme";

import CheckIcon from "~icons/lucide/check";
import MonitorIcon from "~icons/lucide/monitor";
import MoonIcon from "~icons/lucide/moon";
import SunIcon from "~icons/lucide/sun";

type IconComponent = () => JSX.Element;

const Check =
  CheckIcon as unknown as IconComponent;

const Monitor =
  MonitorIcon as unknown as IconComponent;

const Moon =
  MoonIcon as unknown as IconComponent;

const Sun =
  SunIcon as unknown as IconComponent;

type Theme = "light" | "dark" | "system";

export function AppearancePreferencesContent() {
  const { theme, setTheme } = useTheme();

  const chooseTheme = (nextTheme: Theme) => {
    setTheme(nextTheme);
  };

  return (
    <div class="min-w-0">
      <div class="mb-4">
        <p class="text-sm font-semibold">
          Tema da interface
        </p>

        <p class="mt-1 max-w-lg text-xs leading-relaxed text-muted-foreground">
          Escolha como o Univents deve aparecer para
          você.
        </p>
      </div>

      <div class="grid gap-2.5 sm:grid-cols-3 sm:gap-3">
        <ThemeOption
          selected={theme() === "light"}
          label="Claro"
          description="Interface clara"
          icon={Sun}
          onClick={() => chooseTheme("light")}
        >
          <LightPreview />
        </ThemeOption>

        <ThemeOption
          selected={theme() === "dark"}
          label="Escuro"
          description="Interface escura"
          icon={Moon}
          onClick={() => chooseTheme("dark")}
        >
          <DarkPreview />
        </ThemeOption>

        <ThemeOption
          selected={theme() === "system"}
          label="Sistema"
          description="Segue seu dispositivo"
          icon={Monitor}
          onClick={() => chooseTheme("system")}
        >
          <SystemPreview />
        </ThemeOption>
      </div>

      <div class="mt-3 flex items-start gap-2 px-1">
        <span class="mt-0.5 flex shrink-0 text-muted-foreground [&>svg]:size-3.5">
          <Monitor />
        </span>

        <p class="text-[11px] leading-relaxed text-muted-foreground">
          No modo Sistema, o Univents acompanha
          automaticamente o tema do seu dispositivo.
        </p>
      </div>
    </div>
  );
}

function ThemeOption(props: {
  selected: boolean;
  label: string;
  description: string;
  icon: IconComponent;
  children: JSX.Element;
  onClick: () => void;
}) {
  const Icon = untrack(() => props.icon);

  return (
    <button
      type="button"
      aria-pressed={
        props.selected ? "true" : "false"
      }
      onClick={() => props.onClick()}
      class={`group relative flex min-w-0 items-center gap-3 overflow-hidden rounded-xl border-2 p-2 text-left transition-all duration-200 sm:block sm:p-0 ${props.selected
        ? "border-primary bg-primary/2.5 shadow-sm shadow-primary/10"
        : "border-border bg-background hover:border-foreground/20 hover:bg-muted/20"
        }`}
    >
      <div class="w-23 shrink-0 overflow-hidden rounded-lg border border-border/60 sm:w-auto sm:rounded-b-none sm:rounded-t-[9px] sm:border-x-0 sm:border-t-0">
        {props.children}
      </div>

      <div class="flex min-w-0 flex-1 items-center gap-2.5 pr-1 sm:px-3 sm:pb-3 sm:pt-2.5">
        <span
          class={`flex size-8 shrink-0 items-center justify-center rounded-lg transition-colors [&>svg]:size-4 ${props.selected
            ? "bg-primary text-primary-foreground"
            : "bg-muted text-muted-foreground group-hover:text-foreground"
            }`}
        >
          <Icon />
        </span>

        <div class="min-w-0 flex-1">
          <p class="truncate text-xs font-semibold">
            {props.label}
          </p>

          <p class="mt-0.5 truncate text-[10px] text-muted-foreground">
            {props.description}
          </p>
        </div>

        <span
          class={`flex size-5 shrink-0 items-center justify-center rounded-full transition-all [&>svg]:size-3 ${props.selected
            ? "scale-100 bg-primary text-primary-foreground opacity-100"
            : "scale-75 border border-border opacity-40"
            }`}
        >
          <Show when={props.selected}>
            <Check />
          </Show>
        </span>
      </div>
    </button>
  );
}

function LightPreview() {
  return (
    <div class="h-15.5 bg-[#f8fafc] p-1.5 sm:h-24 sm:p-2">
      <div class="flex h-full overflow-hidden rounded-md border border-black/10 bg-white">
        <div class="w-3.5 border-r border-black/10 bg-[#f1f5f9] sm:w-5" />

        <div class="flex-1 p-1.5 sm:p-2">
          <div class="mb-1.5 h-1.5 w-1/2 rounded-full bg-[#cbd5e1] sm:mb-2 sm:h-2" />

          <div class="space-y-1 sm:space-y-1.5">
            <div class="h-2.5 rounded bg-[#f1f5f9] sm:h-4" />
            <div class="h-2.5 rounded bg-[#f1f5f9] sm:h-4" />
          </div>
        </div>
      </div>
    </div>
  );
}

function DarkPreview() {
  return (
    <div class="h-15.5 bg-[#09090b] p-1.5 sm:h-24 sm:p-2">
      <div class="flex h-full overflow-hidden rounded-md border border-white/10 bg-[#18181b]">
        <div class="w-3.5 border-r border-white/10 bg-[#27272a] sm:w-5" />

        <div class="flex-1 p-1.5 sm:p-2">
          <div class="mb-1.5 h-1.5 w-1/2 rounded-full bg-[#52525b] sm:mb-2 sm:h-2" />

          <div class="space-y-1 sm:space-y-1.5">
            <div class="h-2.5 rounded bg-[#27272a] sm:h-4" />
            <div class="h-2.5 rounded bg-[#27272a] sm:h-4" />
          </div>
        </div>
      </div>
    </div>
  );
}

function SystemPreview() {
  return (
    <div class="flex h-15.5 sm:h-24">
      <div class="w-1/2 bg-[#f8fafc] p-1.5 pr-0.5 sm:p-2 sm:pr-1">
        <div class="h-full rounded-l-md border border-r-0 border-black/10 bg-white p-1 sm:p-2">
          <div class="mb-1.5 h-1.5 w-2/3 rounded-full bg-[#cbd5e1] sm:mb-2 sm:h-2" />

          <div class="space-y-1 sm:space-y-1.5">
            <div class="h-2.5 rounded bg-[#f1f5f9] sm:h-4" />
            <div class="h-2.5 rounded bg-[#f1f5f9] sm:h-4" />
          </div>
        </div>
      </div>

      <div class="w-1/2 bg-[#09090b] p-1.5 pl-0.5 sm:p-2 sm:pl-1">
        <div class="h-full rounded-r-md border border-l-0 border-white/10 bg-[#18181b] p-1 sm:p-2">
          <div class="mb-1.5 h-1.5 w-2/3 rounded-full bg-[#52525b] sm:mb-2 sm:h-2" />

          <div class="space-y-1 sm:space-y-1.5">
            <div class="h-2.5 rounded bg-[#27272a] sm:h-4" />
            <div class="h-2.5 rounded bg-[#27272a] sm:h-4" />
          </div>
        </div>
      </div>
    </div>
  );
}