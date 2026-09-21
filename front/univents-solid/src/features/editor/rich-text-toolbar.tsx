import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
} from "solid-js";
import AlignCenterIcon from "~icons/lucide/align-center";
import AlignJustifyIcon from "~icons/lucide/align-justify";
import AlignLeftIcon from "~icons/lucide/align-left";
import AlignRightIcon from "~icons/lucide/align-right";
import BoldIcon from "~icons/lucide/bold";
import ItalicIcon from "~icons/lucide/italic";
import MinusIcon from "~icons/lucide/minus";
import PlusIcon from "~icons/lucide/plus";
import UnderlineIcon from "~icons/lucide/underline";
import { cn } from "@trieoh/ui-solid";
import { FONT_FAMILIES, LINE_HEIGHT_OPTIONS } from "./types";
import { ToolbarCombobox } from "./toolbar-combobox";
import type { RichTextController, TextSelectionStyles } from "./types";

const Bold = BoldIcon as unknown as (props: { class?: string }) => JSX.Element;
const Italic = ItalicIcon as unknown as (props: { class?: string }) => JSX.Element;
const Underline = UnderlineIcon as unknown as (props: { class?: string }) => JSX.Element;
const AlignLeft = AlignLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const AlignCenter = AlignCenterIcon as unknown as (props: { class?: string }) => JSX.Element;
const AlignRight = AlignRightIcon as unknown as (props: { class?: string }) => JSX.Element;
const AlignJustify = AlignJustifyIcon as unknown as (props: { class?: string }) => JSX.Element;
const Minus = MinusIcon as unknown as (props: { class?: string }) => JSX.Element;
const Plus = PlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export function normalizeHexColor(color?: string | null): string {
  if (!color) return "#0f172a";
  if (/^#[0-9a-f]{6}$/i.test(color)) return color;
  if (/^#[0-9a-f]{3}$/i.test(color)) {
    return `#${color[1]}${color[1]}${color[2]}${color[2]}${color[3]}${color[3]}`;
  }
  const match = color.match(/rgba?\(([0-9]+),\s*([0-9]+),\s*([0-9]+)/i);
  if (match) {
    const [, r, g, b] = match;
    const hex = (n: string) =>
      Math.min(255, Math.max(0, parseInt(n, 10)))
        .toString(16)
        .padStart(2, "0");
    return `#${hex(r)}${hex(g)}${hex(b)}`;
  }
  return "#0f172a";
}

export interface RichTextToolbarProps {
  controller: RichTextController | null;
  selectionStyles: TextSelectionStyles | null;
  variableOptions?: Array<{
    value: string;
    label: string;
    description?: string;
  }>;
}

export function RichTextToolbar(props: RichTextToolbarProps): JSX.Element {
  const [fontSize, setFontSize] = createSignal("24");
  const [textColor, setTextColor] = createSignal("#0f172a");

  const disabled = createMemo(() => !props.controller);

  createEffect(
    () => props.selectionStyles,
    (styles) => {
      if (styles) {
        queueMicrotask(() => {
          setFontSize(styles.fontSize === null ? "—" : String(styles.fontSize));
          if (styles.color) {
            setTextColor(normalizeHexColor(styles.color));
          }
        });
      }
    },
  );

  const alignments = [
    ["left", AlignLeft, "Alinhar à esquerda"],
    ["center", AlignCenter, "Centralizar"],
    ["right", AlignRight, "Alinhar à direita"],
    ["justify", AlignJustify, "Justificar"],
  ] as const;

  function applyFontSize(value: number) {
    const nextValue = Math.min(400, Math.max(6, Math.round(value)));
    setFontSize(String(nextValue));
    props.controller?.setFontSize(nextValue);
  }

  function stepFontSize(step: number) {
    const current = Number(fontSize());
    applyFontSize(
      Number.isFinite(current)
        ? current + step
        : (props.selectionStyles?.fontSize ?? 24) + step,
    );
  }

  return (
    <div
      data-text-toolbar
      class="relative z-30 flex h-10 shrink-0 items-center justify-center gap-0.5 border-b border-border bg-muted/70 px-2 shadow-xs select-none"
      onPointerDown={(event) => event.stopPropagation()}
    >
      {/* Bold */}
      <button
        type="button"
        title="Negrito"
        disabled={disabled()}
        class={cn(
          "relative flex size-7 cursor-pointer items-center justify-center rounded-md transition-colors disabled:pointer-events-none disabled:opacity-40",
          props.selectionStyles?.bold ? "bg-accent text-accent-foreground shadow-2xs" : "hover:bg-muted text-foreground",
        )}
        onMouseDown={(event) => event.preventDefault()}
        onClick={() => props.controller?.toggleBold()}
      >
        <Bold class="size-3.5" />
        <Show when={props.selectionStyles?.bold === null}>
          <span class="absolute right-0.5 bottom-0 text-[9px] leading-none">
            −
          </span>
        </Show>
      </button>

      {/* Italic */}
      <button
        type="button"
        title="Itálico"
        disabled={disabled()}
        class={cn(
          "relative flex size-7 cursor-pointer items-center justify-center rounded-md transition-colors disabled:pointer-events-none disabled:opacity-40",
          props.selectionStyles?.italic ? "bg-accent text-accent-foreground shadow-2xs" : "hover:bg-muted text-foreground",
        )}
        onMouseDown={(event) => event.preventDefault()}
        onClick={() => props.controller?.toggleItalic()}
      >
        <Italic class="size-3.5" />
        <Show when={props.selectionStyles?.italic === null}>
          <span class="absolute right-0.5 bottom-0 text-[9px] leading-none">
            −
          </span>
        </Show>
      </button>

      {/* Underline */}
      <button
        type="button"
        title="Sublinhado"
        disabled={disabled()}
        class={cn(
          "relative flex size-7 cursor-pointer items-center justify-center rounded-md transition-colors disabled:pointer-events-none disabled:opacity-40",
          props.selectionStyles?.underline ? "bg-accent text-accent-foreground shadow-2xs" : "hover:bg-muted text-foreground",
        )}
        onMouseDown={(event) => event.preventDefault()}
        onClick={() => props.controller?.toggleUnderline()}
      >
        <Underline class="size-3.5" />
        <Show when={props.selectionStyles?.underline === null}>
          <span class="absolute right-0.5 bottom-0 text-[9px] leading-none">
            −
          </span>
        </Show>
      </button>

      <div class="mx-1 h-4 w-px bg-border" />

      {/* Alignments */}
      <For each={alignments}>
        {([align, Icon, label]) => {
          const isCurrent = () => props.selectionStyles?.align === align;

          return (
            <button
              type="button"
              title={label}
              disabled={disabled()}
              class={cn(
                "flex size-7 cursor-pointer items-center justify-center rounded-md transition-colors disabled:pointer-events-none disabled:opacity-40",
                isCurrent() ? "bg-accent text-accent-foreground shadow-2xs" : "hover:bg-muted text-foreground",
              )}
              onMouseDown={(event) => event.preventDefault()}
              onClick={() => props.controller?.setAlign(align)}
            >
              <Icon class="size-3.5" />
            </button>
          );
        }}
      </For>

      <div class="mx-1 h-4 w-px bg-border" />

      {/* Font Family */}
      <ToolbarCombobox
        value={
          props.selectionStyles?.fontFamily === null
            ? ""
            : (props.selectionStyles?.fontFamily ?? "")
        }
        options={FONT_FAMILIES}
        placeholder={
          props.selectionStyles?.fontFamily === null ? "Várias fontes" : "Fonte"
        }
        class="w-36"
        disabled={disabled()}
        onChange={(font) => props.controller?.setFontFamily(font)}
      />

      <div class="mx-1 h-4 w-px bg-border" />

      {/* Line Height */}
      <ToolbarCombobox
        value={
          props.selectionStyles?.lineHeight === null
            ? ""
            : String(props.selectionStyles?.lineHeight ?? "")
        }
        options={LINE_HEIGHT_OPTIONS}
        placeholder={
          props.selectionStyles?.lineHeight === null
            ? "Vários espaçamentos"
            : "Espaçamento"
        }
        class="w-24"
        disabled={disabled()}
        onChange={(height) => props.controller?.setLineHeight(Number(height))}
      />

      <div class="mx-1 h-4 w-px bg-border" />

      {/* Font Size Input & Step Buttons */}
      <div class="flex items-center gap-1">
        <button
          type="button"
          title="Diminuir tamanho da fonte"
          disabled={disabled()}
          class="flex size-7 cursor-pointer items-center justify-center rounded-md text-foreground transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-40"
          onMouseDown={(event) => event.preventDefault()}
          onClick={() => stepFontSize(-1)}
        >
          <Minus class="size-3" />
        </button>
        <input
          type="text"
          inputmode="numeric"
          pattern="[0-9]*"
          value={fontSize()}
          disabled={disabled()}
          aria-label="Tamanho da fonte em pixels"
          class="h-7 w-11 rounded-md border border-input bg-background px-1 text-center text-xs text-foreground focus-visible:border-ring focus-visible:outline-hidden disabled:opacity-40"
          onInput={(event) => setFontSize(event.currentTarget.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              applyFontSize(Number(event.currentTarget.value));
            }
          }}
          onBlur={(event) => applyFontSize(Number(event.currentTarget.value))}
        />
        <button
          type="button"
          title="Aumentar tamanho da fonte"
          disabled={disabled()}
          class="flex size-7 cursor-pointer items-center justify-center rounded-md text-foreground transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-40"
          onMouseDown={(event) => event.preventDefault()}
          onClick={() => stepFontSize(1)}
        >
          <Plus class="size-3" />
        </button>
      </div>

      {/* Text Color */}
      <div class="relative h-7 w-8 shrink-0">
        <input
          type="color"
          value={textColor()}
          disabled={disabled()}
          aria-label="Cor do texto"
          class="h-7 w-8 cursor-pointer rounded-md p-1 border border-input bg-background"
          onChange={(event) => {
            setTextColor(event.currentTarget.value);
            props.controller?.setColor(event.currentTarget.value);
          }}
        />
        <Show when={props.controller && props.selectionStyles?.color === null}>
          <span
            class="pointer-events-none absolute inset-1 rounded-sm border border-border"
            style={{
              background:
                "conic-gradient(#cbd5e1 0 25%, transparent 0 50%, #cbd5e1 0 75%, transparent 0)",
            }}
          />
        </Show>
      </div>

      {/* Variables Combobox (only if passed) */}
      <Show
        when={props.variableOptions && props.variableOptions.length > 0}
      >
        <div class="mx-1 h-4 w-px bg-border" />
        <ToolbarCombobox
          value=""
          options={props.variableOptions ?? []}
          placeholder="Variável"
          class="w-32"
          disabled={disabled()}
          onChange={(variable) => props.controller?.insertText(variable)}
        />
      </Show>
    </div>
  );
}
