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
  FilePlus,
  CheckIcon,
} from "lucide-react"
import React from "react"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { latestSchemaVersionQueryOptions, createSchemaVersionDraftFn } from "../api"
import { cn } from "@/shared/lib/utils"
import { navigationStore, navigationActions } from "@/features/navigation"
import { useStore } from "@tanstack/react-store"
import { toast } from "sonner"

export default function SchemaVersionSelector() {
  const [open, setOpen] = React.useState(false);
  const { 
    currentProjectId, 
    currentSchemaId, 
    currentSchemaVersion: storedVersion 
  } = useStore(navigationStore)
  const queryClient = useQueryClient();

  const enabled = !!currentProjectId && !!currentSchemaId;

  const { data: latestVersion, isLoading: isLoadingLatest } = useQuery({
    ...latestSchemaVersionQueryOptions(currentProjectId!, currentSchemaId!),
    enabled,
  });

  const [showCreatingLoader, setShowCreatingLoader] = React.useState(false);
  const MIN_LOADER_DISPLAY_TIME = 500; // ms

    const createSchemaVersionMutation = useMutation({
    mutationFn: createSchemaVersionDraftFn,
    onSuccess: (response) => {
      if (response.success) {
        toast.success(response.message);
        queryClient.invalidateQueries({ 
          queryKey: ["latestSchemaVersion", currentProjectId, currentSchemaId] 
        });
        navigationActions.setCurrentSchemaVersion(response.data.version_number);
        setOpen(false);
      }
    },
  });

  React.useEffect(() => {
    let timer: NodeJS.Timeout;
    if (createSchemaVersionMutation.isPending) {
      setShowCreatingLoader(true);
    } else {
      if (showCreatingLoader) { 
        timer = setTimeout(() => {
          setShowCreatingLoader(false);
        }, MIN_LOADER_DISPLAY_TIME);
      }
    }

    return () => clearTimeout(timer);
  }, [createSchemaVersionMutation.isPending, showCreatingLoader]);



  const displayVersion = React.useMemo(() => {
    if (!enabled) return "Select project...";
    if (showCreatingLoader) return "Creating...";
    if (isLoadingLatest) return "Loading...";
    
    const version = storedVersion ?? latestVersion?.version_number;
    
    if (!version) return "Empty";
    
    const isLatest = version === latestVersion?.version_number;
    return `v${version}${isLatest ? " (lts)" : ""}`;
  }, [
    enabled, 
    storedVersion, 
    latestVersion?.version_number || null, 
    isLoadingLatest, 
    showCreatingLoader
  ]);

  const previousVersions = React.useMemo(() => {
    if (!latestVersion?.version_number) return [];
    const current = storedVersion ?? latestVersion.version_number;
    const versions = [];
    for (let i = latestVersion.version_number - 1; i >= 1; i--) {
      if (i !== latestVersion.version_number) {
        versions.push({ version_number: i, isActive: i === current });
      }
    }
    return versions;
  }, [latestVersion?.version_number, storedVersion]);

  const handleSelectVersion = (versionNumber: number) => {
    navigationActions.setCurrentSchemaVersion(versionNumber);
    setOpen(false);
  }

  const handleCreateNewDraft = () => {
    createSchemaVersionMutation.mutate({
      project_id: currentProjectId!, 
      schema_id: currentSchemaId!
    });
  };



  if (!enabled) {
    return (
      <button
        type="button"
        disabled
        className="flex h-9 items-center gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm text-muted-foreground"
      >
        <GitBranchIcon size={16} />
        <span>Select project...</span>
      </button>
    );
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          type="button"
          className={cn(
            "flex h-9 items-center justify-between gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm",
            "hover:bg-accent hover:text-accent-foreground",
            "focus:outline-none focus:ring-1 focus:ring-ring",
            "disabled:cursor-not-allowed disabled:opacity-50",
            "w-fit text-xs sm:text-sm"
          )}
          disabled={createSchemaVersionMutation.isPending}
        >
          <div className="flex items-center gap-2 overflow-hidden">
            {showCreatingLoader ? (
              <Loader2 className="h-4 w-4 animate-spin shrink-0" />
            ) : (
              <GitBranchIcon size={16} className="shrink-0" />
            )}
            <span className="truncate">{displayVersion}</span>
          </div>
          <ChevronDownIcon className="h-4 w-4 shrink-0 opacity-50" />
        </button>
      </PopoverTrigger>
      
      <PopoverContent className="w-70 p-0" align="start">
        <Command>
          <CommandInput placeholder="Search versions..." />
          <CommandList>
            <CommandEmpty>No versions found.</CommandEmpty>
            
            <CommandGroup heading="Actions">
              <CommandItem
                onSelect={handleCreateNewDraft}
                disabled={createSchemaVersionMutation.isPending}
              >
                {showCreatingLoader ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <FilePlus className="mr-2 h-4 w-4" />
                )}
                <span>Create New Draft</span>
              </CommandItem>
            </CommandGroup>

            <CommandSeparator />

            {/* Latest Version */}
            {latestVersion && (
              <CommandGroup heading="Latest Version">
                <CommandItem
                  onSelect={() => handleSelectVersion(latestVersion.version_number)}
                  className="justify-between"
                >
                  <div className="flex items-center gap-2">
                    <GitBranchIcon className="h-4 w-4" />
                    <span>v{latestVersion.version_number}</span>
                  </div>
                  {(storedVersion ?? latestVersion.version_number) === latestVersion.version_number && (
                    <CheckIcon className="h-4 w-4" />
                  )}
                </CommandItem>
              </CommandGroup>
            )}



            {previousVersions.length > 0 && (
              <>
                <CommandSeparator />
                <CommandGroup heading={`Previous Versions (${previousVersions.length})`}>
                  {previousVersions.map((version) => (
                    <CommandItem
                      key={version.version_number}
                      onSelect={() => handleSelectVersion(version.version_number)}
                      className="justify-between"
                    >
                      <div className="flex items-center gap-2">
                        <GitBranchIcon className="h-4 w-4 opacity-50" />
                        <span>v{version.version_number}</span>
                      </div>
                      {version.isActive && <CheckIcon className="h-4 w-4" />}
                    </CommandItem>
                  ))}
                </CommandGroup>
              </>
            )}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}