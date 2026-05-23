import { cn } from "#/shared/lib/utils";
import { Box, Ellipsis } from "lucide-react";
import type { NamespaceI } from "../model";
import { Link } from "@tanstack/react-router";
import { timeAgo } from "#/shared/lib/helpers/date-utils";

interface PropsI {
  data: NamespaceI
}

export function NamespaceCard({ data }: PropsI) {
  return (
    <Link
      className={cn(
        "bg-card rounded-sm w-72 cursor-pointer",
        "ring-1 ring-foreground/10 shadow-xs",
        "relative py-4 hover:ring-primary hover:shadow-primary duration-150"
      )}
      to="/"
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
      <Ellipsis
        className={cn(
          "absolute text-muted-foreground hover:text-foreground/40 cursor-pointer",
          "right-4 top-2 duration-150 active:scale-90"
        )}
      />
    </Link>
  )
}