import { createEffect, createSignal } from "solid-js";
import CalendarIcon from "~icons/lucide/calendar-days";
import { Button } from "./button";
import { Popover } from "./popover";
import {
  Calendar,
  dayTime,
  formatDate,
  validDate,
} from "./date-picker/calendar";
import type { Shape } from "./lib/shape";
export interface DateRange {
  from: Date | null;
  to: Date | null;
}
export interface DateRangePickerProps {
  value?: DateRange;
  onChange: (range: DateRange) => void;
  min?: Date;
  max?: Date;
  isDateDisabled?: (date: Date) => boolean;
  disabled?: boolean;
  placeholder?: string;
  ariaLabel?: string;
  clearable?: boolean;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
  shape?: Shape;
  class?: string;
}
export function DateRangePicker(props: DateRangePickerProps) {
  const [opened, setOpened] = createSignal(false);
  const [draft, setDraft] = createSignal<DateRange>({ from: null, to: null });
  const open = () => !props.disabled && (props.open ?? opened());
  const changeOpen = (value: boolean) => {
    setOpened(value);
    props.onOpenChange?.(value);
  };
  createEffect(
    () => ({ open: open(), value: props.value }),
    (snapshot) => {
      if (snapshot.open) setDraft(snapshot.value ?? { from: null, to: null });
    },
  );
  const valid = () => validDate(draft().from) && validDate(draft().to);
  return (
    <Popover
      open={open()}
      onOpenChange={changeOpen}
      disabled={props.disabled}
      ariaLabel={props.ariaLabel ?? "Selecionar período"}
      width={304}
      shape={props.shape}
      trigger={(trigger) => (
        <Button
          {...trigger}
          variant="outline"
          class={props.class}
          shape={props.shape}
          aria-label={props.ariaLabel ?? "Selecionar período"}
        >
          <CalendarIcon class="size-4 shrink-0" />
          <span class="min-w-0 truncate">
            {validDate(props.value?.from)
              ? `${formatDate(props.value.from)}${validDate(props.value.to) ? ` – ${formatDate(props.value.to)}` : ""}`
              : (props.placeholder ?? "Selecionar período")}
          </span>
        </Button>
      )}
    >
      {(close) => (
        <>
          <p aria-live="polite" class="mb-2 text-xs text-muted-foreground">
            {!draft().from
              ? "Escolha a data inicial"
              : !draft().to
                ? "Escolha a data final"
                : "Período selecionado"}
          </p>
          <Calendar
            from={draft().from}
            to={draft().to}
            min={props.min}
            max={props.max}
            isDateDisabled={props.isDateDisabled}
            onSelect={(date) =>
              setDraft((previous) =>
                !previous.from || previous.to
                  ? { from: date, to: null }
                  : dayTime(date) < dayTime(previous.from)
                    ? { from: date, to: previous.from }
                    : { from: previous.from, to: date },
              )
            }
          />
          <div class="mt-3 flex flex-wrap justify-end gap-2 border-t border-border pt-3">
            {props.clearable !== false && (
              <Button
                size="sm"
                variant="ghost"
                onClick={() => {
                  props.onChange({ from: null, to: null });
                  close();
                }}
              >
                Limpar
              </Button>
            )}
            <Button size="sm" variant="outline" onClick={close}>
              Cancelar
            </Button>
            <Button
              size="sm"
              disabled={!valid()}
              onClick={() => {
                if (valid()) {
                  props.onChange(draft());
                  close();
                }
              }}
            >
              Aplicar
            </Button>
          </div>
        </>
      )}
    </Popover>
  );
}
