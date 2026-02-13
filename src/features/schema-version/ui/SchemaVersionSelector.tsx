import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator } from "@/shared/ui/shadcn/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/shared/ui/shadcn/popover"
import {
  GitBranchIcon,
  ChevronDownIcon,
  Loader2,
} from "lucide-react"
import React, { useEffect } from "react"
import { useQuery } from "@tanstack/react-query"
import { latestSchemaVersionQueryOptions, currentSchemaVersionQueryOptions } from "../api"
import { cn } from "@/shared/lib/utils"
import { navigationStore, navigationActions } from "@/features/navigation"
import { useStore } from "@tanstack/react-store"


export default function SchemaVersionSelector() {
  const [open, setOpen] = React.useState(false);
  const { currentProjectId, currentSchemaId, currentSchemaVersion: storedCurrentSchemaVersion } = useStore(navigationStore)

  const { data: latestVersion, isLoading: isLoadingLatest, isError: isErrorLatest } = useQuery(
    latestSchemaVersionQueryOptions(currentProjectId || "", currentSchemaId || "")
  );

  const { data: currentVersion, isLoading: isLoadingCurrent, isError: isErrorCurrent } = useQuery(
    currentSchemaVersionQueryOptions(currentProjectId || "", currentSchemaId || "")
  );

  useEffect(() => {
    if (currentVersion && currentVersion.version_number !== storedCurrentSchemaVersion) {
      navigationActions.setCurrentSchemaVersion(currentVersion.version_number);
    }
  }, [currentVersion, storedCurrentSchemaVersion]);

  const isLoading = isLoadingLatest || isLoadingCurrent;
  const isError = isErrorLatest || isErrorCurrent;

  const previousVersions = React.useMemo(() => {
    if (!currentVersion || currentVersion.version_number <= 1) return [];
    const versions = [];
    for (let i = currentVersion.version_number - 1; i >= 1; i--) versions.push({ version_number: i });
    return versions;
  }, [currentVersion]);

  const handleSelectVersion = (versionNumber: number) => {
    navigationActions.setCurrentSchemaVersion(versionNumber);
    setOpen(false);
  }

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
                {storedCurrentSchemaVersion ? `v${storedCurrentSchemaVersion}` : (currentVersion ? `v${currentVersion.version_number}` : (latestVersion ? `v${latestVersion.version_number} (Latest)` : "Select version..."))}
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
            <CommandGroup heading="Latest Version">
              <CommandItem
                key={latestVersion?.id || "latest"}
                onSelect={() => handleSelectVersion(latestVersion?.version_number || 1)}
                disabled={!latestVersion}
              >
                <GitBranchIcon />
                <span>{latestVersion ? `v${latestVersion.version_number}` : "N/A"}</span>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading="Current Active">
              <CommandItem
                key={currentVersion?.id || "current"}
                onSelect={() => handleSelectVersion(currentVersion?.version_number || 1)}
                disabled={!currentVersion}
              >
                <GitBranchIcon />
                <span>{currentVersion ? `v${currentVersion.version_number}` : "N/A"}</span>
              </CommandItem>
            </CommandGroup>
            <CommandSeparator />
            <CommandGroup heading={`Previous Versions (${previousVersions.length})`}>
              {previousVersions.map((version) => (
                <CommandItem
                  key={version.version_number}
                  onSelect={() => handleSelectVersion(version.version_number)}
                >
                  <GitBranchIcon />
                  <span>{`v${version.version_number}`}</span>
                </CommandItem>
              ))}
              {!previousVersions.length && (
                <CommandItem disabled>
                  <span>No previous versions</span>
                </CommandItem>
              )}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
