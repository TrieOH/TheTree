import type { JSX } from "@solidjs/web";
import { createSignal } from "solid-js";
import CalendarIcon from "~icons/lucide/calendar";
import ChevronDownIcon from "~icons/lucide/chevron-down";
import { DateRangeModal, PRESETS } from "./date-range-modal";
import { isSameDay } from "./date-utils";
import type { CustomRange, RangeKey } from "./types";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideCalendar = CalendarIcon as unknown as IconComp;
const LucideChevronDown = ChevronDownIcon as unknown as IconComp;

const defaultTriggerFormatter = (date: Date) =>
  new Intl.DateTimeFormat("pt-BR", { day: "2-digit", month: "short" }).format(
    date,
  );

export function DateRangeControl(props: {
  range: RangeKey;
  onRangeChange: (range: RangeKey) => void;
  customRange: CustomRange;
  onCustomRangeChange: (customRange: CustomRange) => void;
  dateFormatter?: (date: Date) => string;
}): JSX.Element {
  const [open, setOpen] = createSignal(false);
  const formatDate = () => props.dateFormatter ?? defaultTriggerFormatter;

  const label = () => {
    if (props.range === "custom" && props.customRange.from) {
      if (
        props.customRange.to &&
        !isSameDay(props.customRange.from, props.customRange.to)
      ) {
        return `${formatDate()(props.customRange.from)} – ${formatDate()(props.customRange.to)}`;
      }
      return formatDate()(props.customRange.from);
    }
    return PRESETS.find((p) => p.key === props.range)?.label ?? "Período";
  };

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        class="flex h-8 items-center gap-1.5 rounded-full border border-border bg-background px-3 text-xs font-medium text-muted-foreground shadow-xs transition-colors hover:bg-muted hover:text-foreground"
      >
        <LucideCalendar class="size-3.5 text-muted-foreground" />
        <span class="max-w-[9rem] truncate sm:max-w-none">{label()}</span>
        <LucideChevronDown class="size-3 text-muted-foreground" />
      </button>

      <DateRangeModal
        open={open()}
        onClose={() => setOpen(false)}
        range={props.range}
        customRange={props.customRange}
        onApply={(r, cr) => {
          props.onRangeChange(r);
          props.onCustomRangeChange(cr);
        }}
        dateFormatter={props.dateFormatter}
      />
    </>
  );
}
