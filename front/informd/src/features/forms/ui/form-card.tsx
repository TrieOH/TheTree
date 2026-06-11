import { cn } from "#/shared/lib/utils";
import {
  Archive,
  ClipboardCheck,
  Copy,
  Ellipsis,
  ExternalLink,
  FileText,
  LayoutList,
  User2
} from "lucide-react";
import { Link, useNavigate } from "@tanstack/react-router";
import { timeAgo } from "#/shared/lib/helpers/date-utils";
import { toast } from "sonner";
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
import {
  FormStatusArchived,
  FormStatusClosed,
  FormStatusDraft,
  FormStatusOpen
} from "@trieoh/informd-models";
import type { FormI, FormStatusI } from "../model";

interface StatusMeta {
  /** Tailwind bg class for the card header area */
  headerBg: string;
  /** Tailwind bg class for the icon circle */
  iconBg: string;
  /** Tailwind text-color class for the icon */
  iconColor: string;
  /** Icon component */
  Icon: React.ElementType;
  /** Badge bg + text classes */
  badgeCn: string;
  /** Human-readable label */
  label: string;
}

const STATUS_META: Record<string, StatusMeta> = {
  [FormStatusOpen]: {
    headerBg: "bg-blue-50 dark:bg-blue-950/40",
    iconBg: "bg-blue-100 dark:bg-blue-900/60",
    iconColor: "text-blue-600 dark:text-blue-400",
    Icon: LayoutList,
    badgeCn:
      "bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300",
    label: "Open",
  },
  [FormStatusDraft]: {
    headerBg: "bg-muted/60",
    iconBg: "bg-muted",
    iconColor: "text-muted-foreground",
    Icon: FileText,
    badgeCn:
      "bg-secondary text-secondary-foreground",
    label: "Draft",
  },
  [FormStatusClosed]: {
    headerBg: "bg-red-50 dark:bg-red-950/40",
    iconBg: "bg-red-100 dark:bg-red-900/60",
    iconColor: "text-red-600 dark:text-red-400",
    Icon: ClipboardCheck,
    badgeCn:
      "bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300",
    label: "Closed",
  },
  [FormStatusArchived]: {
    headerBg: "bg-amber-50 dark:bg-amber-950/40",
    iconBg: "bg-amber-100 dark:bg-amber-900/60",
    iconColor: "text-amber-600 dark:text-amber-400",
    Icon: Archive,
    badgeCn:
      "bg-amber-100 text-amber-800 dark:bg-amber-900/50 dark:text-amber-300",
    label: "Archived",
  },
};

function getStatusMeta(status: FormStatusI): StatusMeta {
  return STATUS_META[status];
}

function StatusBadge({ status }: { status: FormStatusI }) {
  const meta = getStatusMeta(status);
  return (
    <span
      className={cn(
        "ml-auto shrink-0 rounded-full px-2 py-0.5 text-[10px] font-medium leading-none",
        meta.badgeCn
      )}
    >
      {meta.label}
    </span>
  );
}


interface MenuItemsProps {
  data: FormI;
  isContext?: boolean;
}

function MenuItems({ data, isContext = false }: MenuItemsProps) {
  const navigate = useNavigate();
  const Item = isContext ? ContextMenuItem : DropdownMenuItem;
  const Separator = isContext ? ContextMenuSeparator : DropdownMenuSeparator;

  return (
    <>
      {data.status === FormStatusOpen && (
        <>
          <Item
            onClick={() => {
              const url = `${window.location.origin}/view/${data.id}${data.namespace_id ? `?namespace_id=${data.namespace_id}` : ""}`
              navigator.clipboard.writeText(url)
              toast.success("Link copied to clipboard")
            }}
          >
            <Copy className="mr-2 size-4" />
            Copy Link
          </Item>
          <Item
            onClick={() => {
              navigate({
                to: '/view/$formID',
                params: { formID: data.id },
                search: { namespace_id: data.namespace_id || undefined }
              })
            }}
          >
            <ExternalLink className="mr-2 size-4" />
            View Form
          </Item>
          <Separator />
        </>
      )}

      <Item
        onClick={() => {
          navigate({
            to: '/admin/form/$formID',
            params: { formID: data.id },
            search: { namespaceID: data.namespace_id || undefined }
          })
        }}
      >
        <ExternalLink className="mr-2 size-4" />
        View Steps
      </Item>

      <Item
        onClick={() => {
          navigate({
            to: '/admin/form/$formID/submissions',
            params: { formID: data.id },
            search: { namespaceID: data.namespace_id || undefined }
          })
        }}
      >
        <ClipboardCheck className="mr-2 size-4" />
        View Submissions
      </Item>

      <Separator />
      <Item
        onClick={() => {
          navigate({
            to: '/admin/form/$formID/members',
            params: { formID: data.id },
            search: { namespaceID: data.namespace_id || undefined }
          })
        }}
      >
        <User2 className="mr-2 size-4" />
        View Members
      </Item>
    </>
  );
}


interface FormCardProps {
  data: FormI;
}

export function FormCard({ data }: FormCardProps) {
  const meta = getStatusMeta(data.status);
  const { Icon } = meta;
  const isNamespaceForm = Boolean(data.namespace_id);

  return (
    <ContextMenu>
      <ContextMenuTrigger
        render={
          <Link
            className={cn(
              "group relative flex w-72 max-w-full flex-col overflow-hidden rounded-sm",
              "bg-card ring-1 ring-foreground/10 shadow-xs",
              "cursor-pointer select-none",
              "transition-all duration-150",
              "hover:ring-primary hover:shadow-primary/20 hover:shadow-md"
            )}
            to="/admin/form/$formID"
            params={{ formID: data.id }}
            search={{ namespaceID: data.namespace_id || undefined }}
          />
        }
      >
        <div
          className={cn(
            "flex h-32 items-center justify-center relative",
            meta.headerBg
          )}
        >
          <span
            className={cn(
              'rounded-full px-1.5 py-0.5 text-[10px] font-semibold',
              'absolute left-2.5 top-2.5 uppercase tracking-wider',
              isNamespaceForm
                ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/60 dark:text-amber-200'
                : 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200'
            )}
          >
            {isNamespaceForm ? 'Namespace' : 'Personal'}
          </span>
          <div
            className={cn(
              "flex size-12 items-center justify-center rounded-full",
              meta.iconBg
            )}
          >
            <Icon className={cn("size-6", meta.iconColor)} strokeWidth={1.5} />
          </div>
        </div>

        <div className="border-t border-foreground/10 px-3 py-2.5">
          <p className="truncate text-sm font-medium leading-snug">
            {data.title}
          </p>
          <div className="mt-1.5 flex items-center gap-1 text-xs text-muted-foreground">
            <span>{timeAgo(data.updated_at)}</span>
            <StatusBadge status={data.status} />
          </div>
        </div>

        <div className="absolute right-1.5 top-1.5">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <Button
                  variant="ghost"
                  size="icon"
                  className={cn(
                    "size-7 rounded-sm",
                    "text-foreground/50 hover:text-foreground",
                    "transition-opacity duration-150",
                    "cursor-pointer opacity-100 sm:opacity-0 sm:group-hover:opacity-100"
                  )}
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                  }}
                >
                  <Ellipsis className="size-4" />
                </Button>
              }
            />
            <DropdownMenuContent align="end" className="w-48">
              <MenuItems data={data} />
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </ContextMenuTrigger>

      <ContextMenuContent className="w-48">
        <MenuItems isContext data={data} />
      </ContextMenuContent>
    </ContextMenu>
  );
}