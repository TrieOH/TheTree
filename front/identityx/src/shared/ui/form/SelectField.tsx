import { Check, ChevronsUpDown, Search } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useFieldContext } from "@/shared/lib/forms";
import { cn } from "@/shared/lib/utils";
import type { FieldOption, RuleStatus } from "./types";

interface PropsI {
  label: string;
  placeholder: string;
  options: FieldOption[];
  required?: boolean;
  getRulesStatus?: (value: unknown) => RuleStatus[];
  submitted?: boolean;
}

export default function SelectField({
  label,
  placeholder,
  options,
  required,
  getRulesStatus,
  submitted,
}: PropsI) {
  const field = useFieldContext<string | undefined>();
  const ref = useRef<HTMLDivElement>(null);
  const searchRef = useRef<HTMLInputElement>(null);
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const value = field.state.value ?? "";
  const selected = options.find((option) => option.value === value);
  const filtered = useMemo(() => {
    const term = query.toLocaleLowerCase();
    return options.filter((option) =>
      `${option.label} ${option.value}`.toLocaleLowerCase().includes(term),
    );
  }, [options, query]);
  const hasError =
    submitted &&
    getRulesStatus?.(field.state.value).some((rule) => !rule.passed);

  useEffect(() => {
    if (open) searchRef.current?.focus();
    const close = (event: MouseEvent) => {
      if (open && !ref.current?.contains(event.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, [open]);

  return (
    <div
      ref={ref}
      className="relative mt-2 flex min-w-0 w-full max-w-full self-stretch flex-col gap-1"
    >
      <label
        htmlFor={field.name}
        className="text-base font-semibold text-foreground"
      >
        {required ? `${label} *` : label}
      </label>
      <button
        type="button"
        id={field.name}
        aria-haspopup="listbox"
        aria-expanded={open}
        onClick={() => setOpen((current) => !current)}
        className={cn(
          "box-border grid h-10 min-w-0 w-full max-w-full grid-cols-[minmax(0,1fr)_auto] items-center gap-2 overflow-hidden rounded-sm",
          "border bg-card px-3 text-left text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring",
          "max-w-full!",
          hasError ? "border-destructive" : "border-input",
        )}
      >
        <span
          className={cn(
            "min-w-0 overflow-hidden truncate",
            !selected && "text-muted-foreground",
          )}
        >
          {selected?.label ?? placeholder}
        </span>
        <ChevronsUpDown className="size-4 shrink-0 text-muted-foreground" />
      </button>
      {open && (
        <div className="absolute inset-x-0 top-full z-30 mt-1 overflow-hidden rounded-sm border border-border bg-popover shadow-xl">
          <div className="flex items-center gap-2 border-b border-border px-3 py-2">
            <Search className="size-4 text-muted-foreground" />
            <input
              ref={searchRef}
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="Search actors..."
              className="min-w-0 flex-1 bg-transparent text-sm outline-none"
            />
          </div>
          <div className="max-h-60 overflow-y-auto p-1" role="listbox">
            <button
              type="button"
              role="option"
              aria-selected={!value}
              className="flex w-full items-center justify-between rounded-sm px-3 py-2 text-left text-sm hover:bg-accent"
              onClick={() => {
                field.handleChange(undefined);
                setOpen(false);
                setQuery("");
              }}
            >
              {" "}
              {placeholder} {!value && <Check className="size-4" />}{" "}
            </button>
            {filtered.map((option) => (
              <button
                type="button"
                role="option"
                aria-selected={option.value === value}
                key={option.value}
                className="flex w-full items-center justify-between rounded-sm px-3 py-2 text-left text-sm hover:bg-accent"
                onClick={() => {
                  field.handleChange(option.value);
                  setOpen(false);
                  setQuery("");
                }}
              >
                <span className="truncate">{option.label}</span>
                {option.value === value && <Check className="size-4" />}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
