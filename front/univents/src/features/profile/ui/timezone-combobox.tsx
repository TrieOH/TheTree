import { Check, ChevronsUpDown, Search } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";

interface TimezoneComboboxProps {
  id: string;
  value: string;
  options: readonly { value: string; label: string }[];
  onChange: (value: string) => void;
}

export function TimezoneCombobox({
  id,
  value,
  options,
  onChange,
}: TimezoneComboboxProps) {
  const ref = useRef<HTMLDivElement>(null);
  const searchRef = useRef<HTMLInputElement>(null);
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const selected = options.find((option) => option.value === value);
  const filtered = useMemo(() => {
    const term = query.trim().toLocaleLowerCase("pt-BR");
    return options.filter((option) =>
      `${option.label} ${option.value}`
        .toLocaleLowerCase("pt-BR")
        .includes(term),
    );
  }, [options, query]);

  useEffect(() => {
    if (!open) return;
    const close = (event: MouseEvent) => {
      if (!ref.current?.contains(event.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, [open]);

  useEffect(() => {
    if (open) searchRef.current?.focus();
  }, [open]);

  return (
    <div ref={ref} className="relative w-full">
      <button
        id={id}
        type="button"
        aria-haspopup="listbox"
        aria-expanded={open}
        className={cn(
          "flex h-10 w-full items-center justify-between gap-2 rounded-md",
          "border border-input bg-background px-3 text-left text-sm",
          "outline-none focus-visible:ring-2 focus-visible:ring-ring",
        )}
        onClick={() => setOpen((current) => !current)}
      >
        <span className={cn("truncate", !selected && "text-muted-foreground")}>
          {selected?.label ?? "Selecione o fuso horário…"}
        </span>
        <ChevronsUpDown className="size-4 shrink-0 text-muted-foreground" />
      </button>
      {open && (
        <div className="absolute inset-x-0 top-full z-30 mt-1 overflow-hidden rounded-md border border-border bg-popover shadow-xl">
          <div className="flex items-center gap-2 border-b border-border px-3 py-2">
            <Search className="size-4 text-muted-foreground" />
            <input
              ref={searchRef}
              value={query}
              placeholder="Buscar cidade ou fuso…"
              className="min-w-0 flex-1 bg-transparent text-sm outline-none"
              onChange={(event) => setQuery(event.target.value)}
            />
          </div>
          <div className="max-h-60 overflow-y-auto p-1" role="listbox">
            {filtered.length ? (
              filtered.map((option) => (
                <button
                  type="button"
                  role="option"
                  aria-selected={option.value === value}
                  key={option.value}
                  className="flex w-full items-center justify-between rounded-sm px-3 py-2 text-left text-sm hover:bg-accent"
                  onClick={() => {
                    onChange(option.value);
                    setOpen(false);
                    setQuery("");
                  }}
                >
                  <span className="truncate">{option.label}</span>
                  {option.value === value && <Check className="size-4" />}
                </button>
              ))
            ) : (
              <p className="px-3 py-4 text-center text-sm text-muted-foreground">
                Nenhum fuso encontrado.
              </p>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
