import { cn } from "#/shared/lib/utils";
import { Box, Ellipsis, ExternalLink, Trash2, User2 } from "lucide-react";
import type { NamespaceI } from "../model";
import { Link } from "@tanstack/react-router";
import { timeAgo } from "#/shared/lib/helpers/date-utils";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "#/shared/ui/shadcn/dropdown-menu";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from "#/shared/ui/shadcn/context-menu";
import { Button } from "#/shared/ui/shadcn/button";

interface PropsI {
  data: NamespaceI;
}

function MenuItems({ isContext = false }: { isContext?: boolean }) {
  const Item = isContext ? ContextMenuItem : DropdownMenuItem;
  const Separator = isContext ? ContextMenuSeparator : DropdownMenuSeparator;

  return (
    <>
      <Item>
        <ExternalLink className="mr-2 size-4" />
        View Forms
      </Item>
      <Item>
        <User2 className="mr-2 size-4" />
        View Members
      </Item>
      <Separator />
      <Item variant="destructive">
        <Trash2 className="mr-2 size-4" />
        Delete (Temp)
      </Item>
    </>
  );
}

export function NamespaceCard({ data }: PropsI) {
  return (
    <ContextMenu>
      <ContextMenuTrigger
        render={
          <Link
            className={cn(
              "bg-card rounded-sm w-72 cursor-pointer",
              "ring-1 ring-foreground/10 shadow-xs",
              "relative py-4 hover:ring-primary hover:shadow-primary duration-150"
            )}
            to="/"
          />
        }
      >
        <div className="px-4 space-y-2">
          <Box className="bg-primary/80 text-primary-foreground p-1.5 rounded-sm size-8" />
          <span className="text-sm font-bold truncate">{data.name}</span>
        </div>
        <hr className="border-muted-foreground/40 mt-2" />
        <div className="flex flex-col gap-1 px-4 mt-2">
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Members</span>
            <span>40 Members</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Updated</span>
            <span>{timeAgo(data.updated_at)}</span>
          </div>
        </div>
        <div className="absolute right-4 top-2">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <Button
                  variant="ghost"
                  size="icon"
                  className={cn(
                    "text-muted-foreground hover:text-foreground/40",
                    "duration-150 cursor-pointer outline-0 select-none"
                  )}
                >
                  <Ellipsis
                    onClick={(e) => {
                      e.preventDefault();
                      e.stopPropagation();
                    }}
                  />
                </Button>
              }
            />
            <DropdownMenuContent align="end" className="w-56">
              <MenuItems />
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </ContextMenuTrigger>
      <ContextMenuContent className="w-56">
        <MenuItems isContext />
      </ContextMenuContent>
    </ContextMenu>
  );
}