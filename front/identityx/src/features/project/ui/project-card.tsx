import { cn } from "@/shared/lib/utils";
import {
  BadgeCheck,
  Ellipsis,
  ExternalLink,
  FolderKanban,
  Globe,
  Users2,
} from "lucide-react";
import type { ProjectI } from "../model";
import { Link, useNavigate } from "@tanstack/react-router";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/ui/shadcn/dropdown-menu";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from "@/shared/ui/shadcn/context-menu";
import { timeAgo } from "@/shared/lib/date-utils";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";

interface PropsI {
  data: ProjectI;
}

function MenuItems({
  isContext = false,
  data,
}: {
  isContext?: boolean;
  data: ProjectI;
}) {
  const navigate = useNavigate();
  const Item = isContext ? ContextMenuItem : DropdownMenuItem;
  const Separator = isContext ? ContextMenuSeparator : DropdownMenuSeparator;

  return (
    <>
      <Item
        onClick={() =>
          navigate({
            to: "/admin/projects/$projectID",
            params: { projectID: data.id },
            search: { organizationID: data.organization_id || undefined }
          })
        }
      >
        <ExternalLink className="mr-2 size-4" />
        View Project
      </Item>
      <Item
        onClick={() =>
          navigate({
            to: "/admin/projects/$projectID/members",
            params: { projectID: data.id },
            search: { organizationID: data.organization_id || undefined }
          })
        }
      >
        <Users2 className="mr-2 size-4" />
        View Members
      </Item>
      {data.domain && (
        <>
          <Separator />
          <Item
            onClick={() => window.open(`https://${data.domain}`, "_blank")}
          >
            <Globe className="mr-2 size-4" />
            Visit Domain
          </Item>
        </>
      )}
    </>
  );
}

export default function ProjectCard({ data }: PropsI) {
  return (
    <ContextMenu>
      <ContextMenuTrigger
        render={
          <Link
            className={cn(
              "bg-card rounded-sm w-full cursor-pointer",
              "ring-1 ring-foreground/10 shadow-xs",
              "relative py-4 hover:ring-primary hover:shadow-primary duration-150"
            )}
            to="/admin/projects/$projectID"
            params={{ projectID: data.id }}
            search={{ organizationID: data.organization_id || undefined }}
          />
        }
      >
        <div className="px-4 space-y-2 pr-10">
          <FolderKanban className="bg-primary/80 text-primary-foreground p-1.5 rounded-sm size-8" />
          <div className="space-y-0.5">
            <span className="text-sm font-bold truncate block">
              {data.name}
            </span>
            {data.domain ? (
              <span className="text-xs text-muted-foreground truncate flex items-center gap-1">
                {data.domain}
                {data.domain_verified_at && (
                  <BadgeCheck
                    className="size-3.5 text-primary shrink-0"
                    aria-label="Domain verified"
                  />
                )}
              </span>
            ) : (
              <span className="text-xs text-muted-foreground/60 truncate block">
                No domain set
              </span>
            )}
          </div>
        </div>

        <hr className="border-muted-foreground/40 mt-2" />

        <div className="flex flex-col gap-1 px-4 mt-2">
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Owner</span>
            <span className="truncate max-w-35" title={data.owner_id}>
              {data.owner_id}
            </span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Created</span>
            <span>{timeAgo(data.created_at)}</span>
          </div>
        </div>

        <div className="absolute right-4 top-2">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <ShadowButton
                  variant="ghost"
                  className={cn(
                    "text-muted-foreground hover:text-foreground/40",
                    "duration-150 cursor-pointer outline-0 select-none"
                  )}
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                  }}
                  leftIcon={<Ellipsis />}
                />
              }
            />
            <DropdownMenuContent align="end" className="w-56">
              <MenuItems data={data} />
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </ContextMenuTrigger>
      <ContextMenuContent className="w-56">
        <MenuItems isContext data={data} />
      </ContextMenuContent>
    </ContextMenu>
  );
}