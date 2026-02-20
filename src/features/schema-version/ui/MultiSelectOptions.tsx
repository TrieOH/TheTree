import type * as React from "react";
import { useState } from "react";
import { ChevronDown, XIcon } from "lucide-react";

import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import {
  Command,
  CommandGroup,
  CommandItem,
  CommandList,
  CommandEmpty,
} from "@/shared/ui/shadcn/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/shared/ui/shadcn/popover";
import { Checkbox } from "@/shared/ui/shadcn/checkbox";

import type { OptionResultI } from '../model/types';

interface MultiSelectOptionsProps {
  id: string;
  options: OptionResultI[];
  value?: string[] | null;
  onChange: (value: string[] | undefined) => void;
  placeholder?: string;
  className?: string;
}

export const MultiSelectOptions: React.FC<MultiSelectOptionsProps> = ({
  id,
  options,
  value,
  onChange,
  placeholder = "Select...",
  className,
}) => {
  const [open, setOpen] = useState(false);

  const selectedValues = Array.isArray(value) ? value : [];

  const handleSelect = (optionValue: string) => {
    const isSelected = selectedValues.includes(optionValue);
    let newSelectedValues: string[];

    if (isSelected) {
      newSelectedValues = selectedValues.filter((val) => val !== optionValue);
    } else {
      newSelectedValues = [...selectedValues, optionValue];
    }
    onChange(newSelectedValues.length > 0 ? newSelectedValues : undefined);
  };

  const handleRemove = (val: string) => {
    const newSelectedValues = selectedValues.filter((v) => v !== val);
    onChange(newSelectedValues.length > 0 ? newSelectedValues : undefined);
  };

  return (
    <div className={cn("flex flex-col gap-1.5 w-full", className)}>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <div
            className="flex h-8 w-full min-w-0 items-center justify-between rounded-md border border-input bg-background px-2 sm:px-3 text-sm shadow-[1px_1px_0_0_var(--color-input)] focus-within:shadow-[2px_2px_0_0_var(--color-input)] transition-all duration-300 ease-out cursor-pointer"
            id={id}
          >
            {selectedValues.length > 0 ? (
              <span className="truncate text-xs sm:text-sm">
                {selectedValues.length} selected
              </span>
            ) : (
              <span className="text-muted-foreground text-xs sm:text-sm truncate">{placeholder}</span>
            )}
            <ChevronDown className="h-4 w-4 shrink-0 opacity-50 ml-1 sm:ml-2" />
          </div>
        </PopoverTrigger>
        <PopoverContent className="w-[--radix-popover-trigger-width] p-0" align="start">
          <Command>
            <CommandList className="max-h-48 overflow-auto">
              <CommandEmpty className="text-xs py-2">No results found.</CommandEmpty>
              <CommandGroup>
                {options.map((option) => (
                  <CommandItem
                    key={option.value}
                    onSelect={() => handleSelect(option.value)}
                    value={option.value}
                    className="text-xs cursor-pointer"
                  >
                    <Checkbox
                      checked={selectedValues.includes(option.value)}
                      onCheckedChange={() => handleSelect(option.value)}
                      className="mr-2 h-3.5 w-3.5"
                    />
                    <span className="truncate">{option.label}</span>
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      {selectedValues.length > 0 && (
        <div className="flex flex-wrap items-center gap-1 p-1.5 border rounded-md w-full min-w-0 max-h-24 overflow-y-auto bg-muted/30">
          {selectedValues.map((val) => {
            const option = options.find(o => o.value === val);
            return (
              <Badge key={val} variant="secondary" className="text-xs px-1.5 py-0.5 shrink-0">
                <span className="truncate max-w-25 sm:max-w-37.5">
                  {option ? option.label : val}
                </span>
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleRemove(val);
                  }}
                  className="ml-1 rounded-full p-0.5 hover:bg-secondary-foreground/20 shrink-0"
                >
                  <XIcon className="h-3 w-3" />
                </button>
              </Badge>
            );
          })}
        </div>
      )}
    </div>
  );
};