import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator, CommandShortcut } from "@/shared/ui/shadcn/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/shared/ui/shadcn/popover"
import {
  CalendarIcon,
  CreditCardIcon,
  SettingsIcon,
  UserIcon,
  GitBranchIcon,
  ChevronDownIcon,
  Loader2,
} from "lucide-react"
import React from "react"
import { useQuery } from "@tanstack/react-query"
import { latestSchemaVersionQueryOptions } from "../api"
import { cn } from "@/shared/lib/utils"
import { navigationStore } from "@/features/navigation"
import { useStore } from "@tanstack/react-store"


export default function SchemaVersionSelector() {
  const [open, setOpen] = React.useState(false);
  const { currentProjectId, currentSchemaId } = useStore(navigationStore)

  const { data: latestVersion, isLoading, isError } = useQuery(
    latestSchemaVersionQueryOptions(currentProjectId || "", currentSchemaId || "")
  );

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          type="button"
          className={cn(
            "flex h-9 items-center justify-between whitespace-nowrap rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
            "w-fit sm:w-37.5"
          )}
          role="combobox"
          aria-expanded={open}
          aria-label="Select schema version"
          disabled={isLoading || isError}
        >
          {isLoading ? (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          ) : isError ? (
            "Error"
          ) : (
            <div className="flex items-center gap-1.5">
              <GitBranchIcon size={16} />
              <span className="truncate mb-0.5">
                {latestVersion ? `v${latestVersion.version_number}` : "Select version..."}
              </span>
            </div>
          )}
          <ChevronDownIcon className="ml-2 h-4 w-4 shrink-0 opacity-50 hidden sm:inline" />
        </button>
      </PopoverTrigger>
      <PopoverContent className="w-[calc(100vw-2rem)] sm:w-75 p-0">
        <Command>
          <CommandInput placeholder="Search versions..." />
          <CommandList>
            <CommandEmpty>No versions found.</CommandEmpty>
            <CommandGroup heading="Current Active">
              <CommandItem>
                <CalendarIcon />
                <span>Calendar</span>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading="Previous Versions (10)">
              <CommandItem>
                <UserIcon />
                <span>Profile</span>
                <CommandShortcut>⌘P</CommandShortcut>
              </CommandItem>
              <CommandItem>
                <CreditCardIcon />
                <span>Billing</span>
                <CommandShortcut>⌘B</CommandShortcut>
              </CommandItem>
              <CommandItem>
                <SettingsIcon />
                <span>Settings</span>
                <CommandShortcut>⌘S</CommandShortcut>
              </CommandItem>
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
