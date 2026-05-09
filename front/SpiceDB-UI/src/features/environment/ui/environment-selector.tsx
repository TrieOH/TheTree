import { ChevronsUpDown, Globe } from "lucide-react"
import { Button } from "#/shared/ui/shadcn/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "#/shared/ui/shadcn/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "#/shared/ui/shadcn/popover"
import { useEnvironment } from "../hooks/use-environment"
import { useState } from "react"

export function EnvironmentSelector() {
  const [open, setOpen] = useState(false)
  const { environments, currentEnvironment, navigateToEnvironment } = useEnvironment()

  if (environments.length === 0) return null

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger render={
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-56 justify-between h-9"
        >
          <div className="flex items-center gap-2 truncate">
            <Globe className="w-4 h-4 text-primary" />
            <span className="truncate">
              {currentEnvironment?.name || "Select environment..."}
            </span>
          </div>
          <ChevronsUpDown className="w-4 h-4 ml-2 opacity-50 shrink-0" />
        </Button>
      } />
      <PopoverContent className="w-56 p-0" align="start">
        <Command>
          <CommandInput placeholder="Search environments..." />
          <CommandList>
            <CommandEmpty>No environment found.</CommandEmpty>
            <CommandGroup>
              {environments.map((env) => (
                <CommandItem
                  key={env.name}
                  value={env.name}
                  data-checked={currentEnvironment?.name === env.name}
                  onSelect={() => {
                    navigateToEnvironment(env.name)
                    setOpen(false)
                  }}
                >
                  {env.name}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
