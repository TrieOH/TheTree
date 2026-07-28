import { Check, ChevronsUpDown, Search } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";

export interface CalendarComboboxOption {
  value: string;
  label: string;
}

export function CalendarCombobox({
  value,
  options,
  placeholder,
  onChange,
  disabled = false,
  compact = false,
}: {
  value: string;
  options: CalendarComboboxOption[];
  placeholder: string;
  onChange: (value: string) => void;
  disabled?: boolean;
  compact?: boolean;
}) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);
  const selected = options.find((option) => option.value === value);
  const filtered = useMemo(
    () =>
      options.filter((option) =>
        option.label.toLowerCase().includes(query.toLowerCase()),
      ),
    [options, query],
  );
  useEffect(() => {
    if (!open) return;
    const close = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, [open]);
  return (
    <div
      ref={containerRef}
      className={cn("relative", compact ? "w-auto" : "w-full")}
    >
      <button
        type="button"
        disabled={disabled}
        aria-expanded={open}
        className={cn(
          "flex items-center justify-between gap-2 rounded-xl border border-input bg-background px-3 text-left text-sm shadow-sm hover:bg-muted/30 disabled:opacity-60",
          compact ? "h-9 w-auto min-w-28" : "h-10 w-full",
        )}
        onClick={() => setOpen(!open)}
      >
        <span className={cn("truncate", !selected && "text-muted-foreground")}>
          {selected?.label ?? placeholder}
        </span>
        <ChevronsUpDown className="size-4 shrink-0 text-muted-foreground" />
      </button>
      {open && (
        <div
          className={cn(
            "absolute top-full z-9999 mt-2 overflow-hidden rounded-xl border border-border bg-popover shadow-xl",
            compact ? "right-0 w-max min-w-full" : "left-0 w-full",
          )}
        >
          <div className="flex items-center gap-2 border-b border-border px-3 py-2">
            <Search className="size-4 text-muted-foreground" />
            <input
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="Buscar..."
              className="w-full bg-transparent text-sm outline-none"
            />
          </div>
          <div className="max-h-52 overflow-y-auto p-1">
            {filtered.map((option) => (
              <button
                type="button"
                key={option.value}
                className="flex w-full items-center justify-between rounded-lg px-3 py-2 text-left text-sm hover:bg-accent"
                onClick={() => {
                  onChange(option.value);
                  setOpen(false);
                  setQuery("");
                }}
              >
                {option.label}
                {option.value === value && <Check className="size-4" />}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
