import { createEffect, createSignal } from "solid-js";
import { Combobox } from "../combobox";
import { Dialog } from "../dialog";
import { Button } from "../button";
import { Field, Input } from "../input";
import { PRESETS, computeRangeWindow } from "./range";
import type { CustomRange, RangeKey } from "./types";

const inputDate = (date: Date) =>
  `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`;
const parseDate = (value: string) => {
  const parts = value.match(/^(\d{4})-(\d{2})-(\d{2})$/);
  if (!parts) return null;
  const date = new Date(
    Number(parts[1]),
    Number(parts[2]) - 1,
    Number(parts[3]),
  );
  return inputDate(date) === value ? date : null;
};
export function DateRangeControl(props: {
  range: RangeKey;
  customRange: CustomRange;
  onRangeChange: (range: RangeKey) => void;
  onCustomRangeChange: (range: CustomRange) => void;
  dateFormatter?: (date: Date) => string;
  disabled?: boolean;
}) {
  const [open, setOpen] = createSignal(false),
    [from, setFrom] = createSignal(""),
    [to, setTo] = createSignal("");
  createEffect(
    () => (open() ? { range: props.range, custom: props.customRange } : null),
    (snapshot) => {
      if (!snapshot) return;
      const window = computeRangeWindow(snapshot.range, snapshot.custom);
      setFrom(inputDate(window.windowStart ?? new Date()));
      setTo(inputDate(window.windowEnd ?? new Date()));
    },
  );
  const customLabel = () => {
    const from = props.customRange.from,
      to = props.customRange.to;
    const format =
      props.dateFormatter ??
      ((date: Date) =>
        new Intl.DateTimeFormat("pt-BR", {
          day: "2-digit",
          month: "short",
        }).format(date));
    return from
      ? `${format(from)}${to ? ` – ${format(to)}` : ""}`
      : "Personalizado";
  };
  const valid = () => {
    const a = parseDate(from()),
      b = parseDate(to());
    return !!a && !!b && a.getTime() <= b.getTime();
  };
  return (
    <>
      <Combobox
        class="w-44 max-w-full"
        size="sm"
        shape="soft"
        searchable={false}
        aria-label="Período do gráfico"
        placeholder="Período"
        value={props.range}
        disabled={props.disabled}
        options={[
          ...PRESETS.map((preset) => ({
            value: preset.key,
            label: preset.label,
          })),
          {
            value: "custom",
            label: props.range === "custom" ? customLabel() : "Personalizado",
          },
        ]}
        onChange={(value) => {
          if (value === "custom") setOpen(true);
          else {
            props.onRangeChange(value as RangeKey);
            props.onCustomRangeChange({ from: null, to: null });
          }
        }}
      />
      <Dialog
        open={open()}
        onOpenChange={setOpen}
        title="Período personalizado"
        description="Escolha a data inicial e final do gráfico."
        size="sm"
        footer={
          <>
            <Button size="sm" variant="outline" onClick={() => setOpen(false)}>
              Cancelar
            </Button>
            <Button
              size="sm"
              disabled={!valid()}
              onClick={() => {
                if (!valid()) return;
                props.onCustomRangeChange({
                  from: parseDate(from()),
                  to: parseDate(to()),
                });
                props.onRangeChange("custom");
                setOpen(false);
              }}
            >
              Aplicar
            </Button>
          </>
        }
      >
        <div class="grid min-w-0 grid-cols-1 gap-4">
          <Field label="De">
            {(ids) => (
              <Input
                id={ids.id}
                type="date"
                value={from()}
                max={to() || undefined}
                onInput={(event) => setFrom(event.currentTarget.value)}
              />
            )}
          </Field>
          <Field label="Até">
            {(ids) => (
              <Input
                id={ids.id}
                type="date"
                value={to()}
                min={from() || undefined}
                onInput={(event) => setTo(event.currentTarget.value)}
              />
            )}
          </Field>
        </div>
      </Dialog>
    </>
  );
}
