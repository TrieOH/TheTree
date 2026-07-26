import { Check, ChevronsUpDown, Loader2, Search } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import type { FieldValues } from "react-hook-form";
import type {
  ComboboxOption,
  FieldConfig,
  FieldFormApi,
} from "../../model/types";
import { getFieldError } from "../../utils/get-field-error";

export interface ComboboxFieldRendererProps<TFieldValues extends FieldValues> {
  field: FieldConfig<TFieldValues>;
  form: FieldFormApi<TFieldValues>;
}

function useDebouncedValue(value: string, delayMs: number): string {
  const [debounced, setDebounced] = useState(value);

  useEffect(() => {
    const timer = setTimeout(() => setDebounced(value), delayMs);
    return () => clearTimeout(timer);
  }, [value, delayMs]);

  return debounced;
}

export function ComboboxFieldRenderer<TFieldValues extends FieldValues>({
  field,
  form,
}: ComboboxFieldRendererProps<TFieldValues>) {
  if (field.kind !== "combobox") return null;

  const containerRef = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [highlightedIndex, setHighlightedIndex] = useState(0);
  const [asyncOptions, setAsyncOptions] = useState<ComboboxOption[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const error = getFieldError(form.formState.errors, field.name);
  const selectedValue = form.watch(field.name);

  const isAsync = typeof field.options === "function";
  const debouncedQuery = useDebouncedValue(query, field.debounceMs ?? 250);
  const normalizedQuery = useMemo(() => query.trim().toLowerCase(), [query]);

  // Static options: filter client-side. Async options: fetched below.
  const staticOptions = useMemo(
    () => (isAsync ? [] : (field.options as ComboboxOption[])),
    [isAsync, field.options],
  );

  useEffect(() => {
    if (!isAsync || !isOpen) return;
    const loader = field.options as (
      query: string,
    ) => Promise<ComboboxOption[]>;
    let cancelled = false;

    setIsLoading(true);
    loader(debouncedQuery)
      .then((results) => {
        if (!cancelled) setAsyncOptions(results);
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [isAsync, isOpen, debouncedQuery, field.options]);

  const visibleOptions = useMemo(
    () =>
      isAsync
        ? asyncOptions
        : staticOptions.filter((option) =>
            option.label.toLowerCase().includes(normalizedQuery),
          ),
    [asyncOptions, isAsync, normalizedQuery, staticOptions],
  );

  const selectedOption = useMemo(() => {
    if (typeof selectedValue !== "string" || selectedValue.length === 0)
      return undefined;
    return (isAsync ? asyncOptions : staticOptions).find(
      (option) => option.value === selectedValue,
    );
  }, [selectedValue, staticOptions, asyncOptions, isAsync]);

  useEffect(() => {
    if (!isOpen) return;
    function handleClickOutside(event: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [isOpen]);

  useEffect(() => {
    setHighlightedIndex(0);
  }, [query, isOpen]);

  const selectOption = (option: ComboboxOption) => {
    form.setValue(field.name, option.value as never, {
      shouldDirty: true,
      shouldValidate: true,
    });
    setIsOpen(false);
    setQuery("");
  };

  return (
    <div className="space-y-1.5" ref={containerRef}>
      <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
        {field.label}
        {field.optional ? (
          <span className="ml-1 font-normal normal-case text-muted-foreground/70">
            (opcional)
          </span>
        ) : null}
      </span>

      <div className="relative">
        <button
          type="button"
          disabled={field.disabled}
          onClick={() => setIsOpen((open) => !open)}
          aria-haspopup="listbox"
          aria-expanded={isOpen}
          className={
            "flex h-10 w-full items-center justify-between rounded-md border bg-transparent px-3 text-sm " +
            (error ? "border-destructive" : "border-input") +
            (field.disabled ? " opacity-60" : "")
          }
        >
          <span className={selectedOption ? "" : "text-muted-foreground"}>
            {selectedOption?.label ?? field.placeholder ?? "Selecione..."}
          </span>
          <ChevronsUpDown className="h-4 w-4 shrink-0 text-muted-foreground" />
        </button>

        {isOpen ? (
          <div className="absolute z-10 mt-1 w-full overflow-hidden rounded-xl border border-border/60 bg-popover shadow-md">
            <div className="flex items-center gap-2 border-b border-border/60 bg-background/70 px-3 py-2.5">
              <Search className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
              <input
                value={query}
                onChange={(event) => setQuery(event.target.value)}
                placeholder="Buscar..."
                className="h-5 w-full bg-transparent text-sm outline-none placeholder:text-muted-foreground"
                onKeyDown={(event) => {
                  if (event.key === "ArrowDown") {
                    event.preventDefault();
                    setHighlightedIndex((index) =>
                      Math.min(index + 1, visibleOptions.length - 1),
                    );
                  } else if (event.key === "ArrowUp") {
                    event.preventDefault();
                    setHighlightedIndex((index) => Math.max(index - 1, 0));
                  } else if (event.key === "Enter") {
                    event.preventDefault();
                    const option = visibleOptions[highlightedIndex];
                    if (option) selectOption(option);
                  } else if (event.key === "Escape") {
                    setIsOpen(false);
                  }
                }}
              />
            </div>

            <ul className="max-h-56 overflow-y-auto py-1">
              {isLoading ? (
                <li className="flex items-center justify-center gap-2 px-3 py-3 text-xs text-muted-foreground">
                  <Loader2 className="h-3.5 w-3.5 animate-spin" />
                  Carregando...
                </li>
              ) : visibleOptions.length === 0 ? (
                <li className="px-3 py-3 text-xs text-muted-foreground">
                  {field.emptyMessage ?? "Nenhum resultado encontrado"}
                </li>
              ) : (
                visibleOptions.map((option, index) => (
                  <li key={option.value}>
                    <button
                      type="button"
                      role="option"
                      aria-selected={option.value === selectedValue}
                      onMouseEnter={() => setHighlightedIndex(index)}
                      onClick={() => selectOption(option)}
                      className={
                        "flex w-full items-center justify-between gap-2 px-3 py-2 text-left text-sm transition-colors " +
                        (index === highlightedIndex
                          ? "bg-primary/10"
                          : "hover:bg-accent/10")
                      }
                    >
                      <span>
                        <span className="block">{option.label}</span>
                        {option.description ? (
                          <span className="block text-xs text-muted-foreground">
                            {option.description}
                          </span>
                        ) : null}
                      </span>
                      {option.value === selectedValue ? (
                        <Check className="h-4 w-4 shrink-0" />
                      ) : null}
                    </button>
                  </li>
                ))
              )}
            </ul>
          </div>
        ) : null}
      </div>

      {field.description ? (
        <p className="text-xs text-muted-foreground">{field.description}</p>
      ) : null}
      {error ? <p className="text-xs text-destructive">{error}</p> : null}
    </div>
  );
}
