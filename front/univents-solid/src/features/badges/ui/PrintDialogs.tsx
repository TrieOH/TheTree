import type { JSX } from "@solidjs/web";
import { Button, Dialog, Field, Input, Label } from "@trieoh/ui-solid";
import { For, createEffect, createSignal } from "solid-js";

export interface DateFilterDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  printedAfter: string;
  onApply: (date: string) => void;
}

export function DateFilterDialog(props: DateFilterDialogProps): JSX.Element {
  const [localDate, setLocalDate] = createSignal("");

  createEffect(
    () => props.printedAfter,
    (val) => {
      setLocalDate(val ?? "");
    },
  );

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title="Filtrar crachás por data"
      description="Exiba e imprima somente os crachás gerados a partir da data e hora escolhidas."
      class="sm:max-w-sm"
    >
      <div class="space-y-4 pt-2">
        <Field label="Gerados a partir de">
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              type="datetime-local"
              value={localDate()}
              onInput={(e) => setLocalDate(e.currentTarget.value)}
              class="mt-1.5 h-10 w-full"
            />
          )}
        </Field>

        <div class="flex gap-2 pt-2">
          <Button
            type="button"
            variant="outline"
            class="flex-1"
            onClick={() => {
              setLocalDate("");
              props.onApply("");
              props.onOpenChange(false);
            }}
          >
            Limpar
          </Button>
          <Button
            type="button"
            class="flex-1"
            onClick={() => {
              props.onApply(localDate());
              props.onOpenChange(false);
            }}
          >
            Aplicar
          </Button>
        </div>
      </div>
    </Dialog>
  );
}

export interface QrPrintDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirmPrint: (sizeMm: number) => void;
  loading?: boolean;
}

const PRESET_QR_SIZES = [
  { size: 30, label: "30 mm (Pequeno)" },
  { size: 48, label: "48 mm (Padrão)" },
  { size: 60, label: "60 mm (Médio)" },
  { size: 210, label: "210 mm (A4 Folha Cheia)" },
];

export function QrPrintDialog(props: QrPrintDialogProps): JSX.Element {
  const [selectedSize, setSelectedSize] = createSignal(48);

  const handlePrint = () => {
    props.onConfirmPrint(selectedSize());
    props.onOpenChange(false);
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title="Imprimir QR Codes de Check-in"
      description="Selecione o tamanho para impressão dos códigos de cada participante."
      class="sm:max-w-sm"
    >
      <div class="space-y-4 pt-2">
        <div class="space-y-2">
          <Label>Tamanhos sugeridos</Label>
          <div class="grid grid-cols-2 gap-2">
            <For each={PRESET_QR_SIZES}>
              {(preset) => (
                <button
                  type="button"
                  onClick={() => setSelectedSize(preset.size)}
                  class={`rounded-lg border p-2.5 text-left text-xs font-medium transition-colors ${
                    selectedSize() === preset.size
                      ? "border-primary bg-primary/10 text-primary"
                      : "border-border bg-card text-foreground hover:bg-muted"
                  }`}
                >
                  {preset.label}
                </button>
              )}
            </For>
          </div>
        </div>

        <Field label="Tamanho customizado (mm)">
          {(ids) => (
            <Input
              id={ids.id}
              aria-describedby={ids.describedBy}
              type="number"
              min={20}
              max={210}
              value={selectedSize()}
              onInput={(e) => {
                const val = Number(e.currentTarget.value);
                if (Number.isFinite(val) && val >= 20 && val <= 210) {
                  setSelectedSize(val);
                }
              }}
              class="mt-1"
            />
          )}
        </Field>

        <div class="flex justify-end gap-2 pt-2">
          <Button
            type="button"
            variant="outline"
            disabled={props.loading}
            onClick={() => props.onOpenChange(false)}
          >
            Cancelar
          </Button>
          <Button
            type="button"
            disabled={props.loading}
            onClick={handlePrint}
          >
            Imprimir
          </Button>
        </div>
      </div>
    </Dialog>
  );
}
