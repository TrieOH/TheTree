import { Link, useNavigate } from "@tanstack/react-router";
import { timeAgo, truncateString } from "@trieoh/shared-utils";
import {
  Check,
  Copy,
  Ellipsis,
  ExternalLink,
  Pencil,
  Sparkles,
  Wallet,
} from "lucide-react";
import type { MouseEvent } from "react";
import { useState } from "react";
import { toast } from "sonner";
import { bpsToPercentage, cn } from "#/shared/lib/utils";
import { Button } from "#/shared/ui/shadcn/button";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from "#/shared/ui/shadcn/context-menu";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "#/shared/ui/shadcn/dropdown-menu";
import type { WalletI, WalletSetSandboxI } from "../model";

interface Props {
  data: WalletI;
  onEditFee: (wallet: WalletI) => void;
  onSetSandbox: (walletId: string, data: WalletSetSandboxI) => void;
  isSettingSandbox?: boolean;
}

export default function WalletCard({
  data,
  onEditFee,
  onSetSandbox,
  isSettingSandbox = false,
}: Props) {
  const navigate = useNavigate();
  const [copied, setCopied] = useState(false);
  const openWallet = () =>
    navigate({ to: "/admin/wallets/$walletID", params: { walletID: data.id } });
  const copyId = async () => {
    await navigator.clipboard.writeText(data.id);
    setCopied(true);
    toast.success("Wallet ID copied");
    setTimeout(() => setCopied(false), 1500);
  };
  const runMenuAction = (event: MouseEvent, action: () => void) => {
    event.preventDefault();
    event.stopPropagation();
    action();
  };

  const menuItems = (context = false) => {
    const Item = context ? ContextMenuItem : DropdownMenuItem;
    const Separator = context ? ContextMenuSeparator : DropdownMenuSeparator;
    return (
      <>
        <Item onClick={(event) => runMenuAction(event, openWallet)}>
          <ExternalLink /> Open wallet
        </Item>
        <Item onClick={(event) => runMenuAction(event, () => onEditFee(data))}>
          <Pencil /> Edit fee
        </Item>
        <Item
          disabled={isSettingSandbox}
          aria-label={`Switch environment to ${data.sandbox ? "production" : "sandbox"}`}
          onClick={(event) =>
            runMenuAction(event, () =>
              onSetSandbox(data.id, {
                sandbox: !data.sandbox,
                organization_id: data.organization_id,
              }),
            )
          }
        >
          <Sparkles /> Switch to {data.sandbox ? "production" : "sandbox"}
        </Item>
        <Separator />
        <Item onClick={(event) => runMenuAction(event, () => void copyId())}>
          {copied ? <Check /> : <Copy />} Copy wallet ID
        </Item>
      </>
    );
  };

  return (
    <ContextMenu>
      <ContextMenuTrigger
        render={
          <Link
            to="/admin/wallets/$walletID"
            params={{ walletID: data.id }}
            className={cn(
              "relative w-full cursor-pointer rounded-sm bg-card py-4 ring-1 ring-foreground/10 shadow-xs duration-150",
              "hover:ring-primary hover:shadow-primary",
            )}
          />
        }
      >
        <div className="space-y-2 px-4 pr-12">
          <Wallet className="size-8 rounded-sm bg-primary/80 p-1.5 text-primary-foreground" />
          <div className="space-y-0.5">
            <span className="block truncate text-sm font-bold">
              {data.name}
            </span>
            <span className="block truncate font-mono text-xs text-muted-foreground">
              {truncateString(data.id, 8, 4)}
            </span>
          </div>
        </div>
        <hr className="mt-2 border-muted-foreground/40" />
        <div className="mt-2 flex flex-col gap-1 px-4 text-sm">
          <div className="flex justify-between gap-3">
            <span className="text-muted-foreground">Owner</span>
            <span>{data.organization_id ? "Organization" : "Personal"}</span>
          </div>
          <div className="flex justify-between gap-3">
            <span className="text-muted-foreground">Fee</span>
            <span>{bpsToPercentage(data.fee_bps).toFixed(2)}%</span>
          </div>
          <div className="flex justify-between gap-3">
            <span className="text-muted-foreground">Collector</span>
            <span>{data.collector_id ? "Connected" : "Not connected"}</span>
          </div>
          <div className="flex justify-between gap-3">
            <span className="text-muted-foreground">Environment</span>
            <span
              className={cn(
                "inline-flex items-center gap-1.5 font-medium",
                data.sandbox ? "text-amber-600" : "text-emerald-600",
              )}
            >
              <span
                className="size-1.5 rounded-full bg-current"
                aria-hidden="true"
              />
              {data.sandbox ? "Sandbox" : "Production"}
            </span>
          </div>
          <div className="flex justify-between gap-3">
            <span className="text-muted-foreground">Created</span>
            <span>{timeAgo(data.created_at)}</span>
          </div>
        </div>
        <div className="absolute right-3 top-2">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <Button
                  variant="ghost"
                  size="icon-sm"
                  onClick={(event) => {
                    event.preventDefault();
                    event.stopPropagation();
                  }}
                />
              }
            >
              <Ellipsis />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-52">
              {menuItems()}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </ContextMenuTrigger>
      <ContextMenuContent className="w-52">
        {menuItems(true)}
      </ContextMenuContent>
    </ContextMenu>
  );
}
