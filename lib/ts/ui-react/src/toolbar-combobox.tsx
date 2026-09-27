import { Check, ChevronsUpDown, Search } from "lucide-react";
import type { ReactNode } from "react";
import { useEffect, useLayoutEffect, useMemo, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { cn } from "./cn";

export interface ComboboxOption {
  value: string;
  label: string;
  description?: string;
}

export interface ComboboxProps {
  value?: string;
  options: readonly ComboboxOption[];
  placeholder: string;
  searchPlaceholder?: string;
  disabled?: boolean;
  icon?: ReactNode;
  iconOnly?: boolean;
  className?: string;
  triggerClassName?: string;
  dropdownClassName?: string;
  onChange: (value: string) => void;
}

export function Combobox({
  value,
  options,
  placeholder,
  searchPlaceholder = "Buscar…",
  disabled = false,
  icon,
  iconOnly = false,
  className,
  triggerClassName,
  dropdownClassName,
  onChange,
}: ComboboxProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const [position, setPosition] = useState({ top: 0, left: 0, minWidth: 0 });
  const [query, setQuery] = useState("");
  const [highlightedIndex, setHighlightedIndex] = useState(0);
  const selectedOption = options.find((option) => option.value === value);
  const visibleOptions = useMemo(() => {
    const normalized = query.trim().toLowerCase();
    return options.filter((option) =>
      `${option.label} ${option.description ?? ""}`
        .toLowerCase()
        .includes(normalized),
    );
  }, [options, query]);

  useEffect(() => {
    if (!open) return;
    const close = (event: MouseEvent) => {
      const target = event.target as Node;
      if (
        !containerRef.current?.contains(target) &&
        !dropdownRef.current?.contains(target)
      ) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, [open]);

  useEffect(() => setHighlightedIndex(0), [open, query]);

  useLayoutEffect(() => {
    if (!open) return;
    const placeDropdown = () => {
      const rect = containerRef.current?.getBoundingClientRect();
      if (rect) {
        setPosition({ top: rect.bottom + 4, left: rect.left, minWidth: rect.width });
      }
    };
    placeDropdown();
    window.addEventListener("resize", placeDropdown);
    window.addEventListener("scroll", placeDropdown, true);
    return () => {
      window.removeEventListener("resize", placeDropdown);
      window.removeEventListener("scroll", placeDropdown, true);
    };
  }, [open]);

  return (
    <div
      ref={containerRef}
      className={cn(
        "relative min-w-0 shrink-0 rounded-md border border-border",
        className,
      )}
    >
      <button
        type="button"
        disabled={disabled}
        aria-haspopup="listbox"
        aria-expanded={open}
        className={cn(
          "flex h-7 items-center rounded-md bg-background text-xs transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-45",
          triggerClassName,
          iconOnly
            ? "w-7 justify-center p-0"
            : "w-full justify-between gap-2 px-2",
        )}
        onClick={() => setOpen((current) => !current)}
      >
        {iconOnly ? (
          icon
        ) : (
          <>
            <span
              className={cn(
                "truncate",
                !selectedOption && "text-muted-foreground",
              )}
            >
              {selectedOption?.label ?? placeholder}
            </span>
            <ChevronsUpDown className="size-3 shrink-0 text-muted-foreground" />
          </>
        )}
      </button>
      {open ? createPortal(
        <div
          ref={dropdownRef}
          style={position}
          className={cn(
            "fixed z-70 w-max max-w-[calc(100vw-1rem)] overflow-hidden rounded-lg border border-border bg-popover text-popover-foreground shadow-lg",
            iconOnly && "w-64 max-w-none",
            dropdownClassName,
          )}
        >
          <div className="flex items-center gap-2 border-b border-border px-2 py-1.5">
            <Search className="size-3.5 shrink-0 text-muted-foreground" />
            <input
              autoFocus
              value={query}
              placeholder={searchPlaceholder}
              className="h-7 min-w-0 flex-1 bg-transparent text-sm outline-none placeholder:text-muted-foreground"
              onChange={(event) => setQuery(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === "Escape") setOpen(false);
                if (event.key === "ArrowDown") {
                  event.preventDefault();
                  setHighlightedIndex((index) =>
                    Math.min(index + 1, Math.max(visibleOptions.length - 1, 0)),
                  );
                }
                if (event.key === "ArrowUp") {
                  event.preventDefault();
                  setHighlightedIndex((index) => Math.max(index - 1, 0));
                }
                if (event.key === "Enter" && visibleOptions[highlightedIndex]) {
                  onChange(visibleOptions[highlightedIndex].value);
                  setOpen(false);
                  setQuery("");
                }
              }}
            />
          </div>
          <ul className="max-h-56 overflow-y-auto py-1">
            {visibleOptions.length === 0 ? (
              <li className="px-3 py-2 text-xs text-muted-foreground">
                Nenhum resultado encontrado
              </li>
            ) : (
              visibleOptions.map((option, index) => (
                <li key={option.value}>
                  <button
                    type="button"
                    role="option"
                    aria-selected={option.value === value}
                    className={cn(
                      "flex w-full items-center justify-between gap-2 px-3 py-2.5 text-left text-sm transition-colors",
                      index === highlightedIndex
                        ? "bg-primary/10"
                        : "hover:bg-muted",
                    )}
                    onMouseEnter={() => setHighlightedIndex(index)}
                    onClick={() => {
                      onChange(option.value);
                      setOpen(false);
                      setQuery("");
                    }}
                  >
                    <span className="min-w-0">
                      <span className="block truncate">{option.label}</span>
                      {option.description ? (
                        <span className="block truncate text-xs text-muted-foreground">
                          {option.description}
                        </span>
                      ) : null}
                    </span>
                    {option.value === value ? (
                      <Check className="size-3.5 shrink-0" />
                    ) : null}
                  </button>
                </li>
              ))
            )}
          </ul>
        </div>,
        document.body,
      ) : null}
    </div>
  );
}

// Backwards-compatible names for existing toolbar consumers.
export const ToolbarCombobox = Combobox;
export type ToolbarComboboxOption = ComboboxOption;
export type ToolbarComboboxProps = ComboboxProps;
