import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createFileRoute,
  Outlet,
  useRouterState,
} from "@tanstack/react-router";
import { LayoutContext } from "@trieoh/ui-base";
import { useState } from "react";
import { toast } from "sonner";
import z from "zod";
import {
  allWalletsQueryOptions,
  createWalletFn,
  setWalletFeeBPSFn,
  setWalletSandboxFn,
} from "#/features/wallets/api";
import type {
  WalletI,
  WalletSetFeeBpsI,
  WalletSetSandboxI,
} from "#/features/wallets/model";
import { WalletsView } from "#/features/wallets/ui/wallets-view";

const walletsSearchSchema = z.object({
  organizationID: z.string().optional(),
});

export const Route = createFileRoute("/admin/wallets")({
  validateSearch: walletsSearchSchema.parse,
  component: RouteComponent,
});

function RouteComponent() {
  const { organizationID } = Route.useSearch();
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  const queryClient = useQueryClient();

  const [headerSlot, setHeaderSlot] = useState<React.ReactNode>(null);

  const { data: wallets = [] } = useQuery(
    allWalletsQueryOptions(organizationID),
  );

  const { mutate: createWallet, isPending: isCreating } = useMutation({
    mutationFn: createWalletFn,
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allWalletsQueryOptions(organizationID).queryKey,
          (oldData = []) => {
            return [response.data, ...oldData];
          },
        );
        toast.success(response.message || "Wallet created successfully");
      } else toast.error(response.message || "Failed to create wallet");
    },
    onError: (error: Error) => toast.error(error.message),
  });

  const { mutate: setWalletFee, isPending: isSettingFee } = useMutation({
    mutationFn: ({
      walletId,
      data,
    }: {
      walletId: string;
      data: WalletSetFeeBpsI;
    }) => setWalletFeeBPSFn(walletId, data),
    onSuccess: (response, variables) => {
      if (response.success) {
        const { walletId, data } = variables;
        queryClient.setQueryData(
          allWalletsQueryOptions(organizationID).queryKey,
          (old: WalletI[] = []) =>
            old.map((wallet) =>
              wallet.id === walletId
                ? { ...wallet, fee_bps: data.fee_bps }
                : wallet,
            ),
        );
        toast.success(response.message || "Wallet fee updated successfully");
      } else toast.error(response.message || "Failed to update wallet fee");
    },
    onError: (error: Error) => toast.error(error.message),
  });

  const { mutate: setWalletSandbox, isPending: isSettingSandbox } = useMutation(
    {
      mutationFn: ({
        walletId,
        data,
      }: {
        walletId: string;
        data: WalletSetSandboxI;
      }) => setWalletSandboxFn(walletId, data),
      onSuccess: (response, variables) => {
        if (response.success) {
          const { walletId, data } = variables;
          queryClient.setQueryData(
            allWalletsQueryOptions(organizationID).queryKey,
            (old: WalletI[] = []) =>
              old.map((wallet) =>
                wallet.id === walletId
                  ? { ...wallet, sandbox: data.sandbox }
                  : wallet,
              ),
          );
          toast.success(
            response.message || "Wallet sandbox state updated successfully",
          );
        } else
          toast.error(
            response.message || "Failed to update wallet sandbox state",
          );
      },
      onError: (error: Error) => toast.error(error.message),
    },
  );

  if (pathname !== "/admin/wallets") {
    return (
      <LayoutContext.Provider value={{ setHeader: setHeaderSlot }}>
        {headerSlot && (
          <div className="border-b border-border/40 bg-background px-6 py-4">
            {headerSlot}
          </div>
        )}
        <div className="flex-1 p-6">
          <Outlet />
        </div>
      </LayoutContext.Provider>
    );
  }

  return (
    <LayoutContext.Provider value={{ setHeader: setHeaderSlot }}>
      {/* Page Header Slot */}
      {/*
          Rendered only when a child page calls useLayoutHeader().
          Sits between the tab bar and the page content.
        */}
      {headerSlot && (
        <div className="border-b border-border/40 px-6 py-4 bg-background">
          {headerSlot}
        </div>
      )}

      <div className="flex-1 p-6">
        <WalletsView
          wallets={wallets}
          organizationId={organizationID}
          title={organizationID ? "Organization Wallets" : "My Wallets"}
          description={
            organizationID
              ? "in this organization"
              : "associated with your account"
          }
          onCreate={createWallet}
          isCreating={isCreating}
          onSetFee={(walletId, data) => setWalletFee({ walletId, data })}
          onSetSandbox={(walletId, data) =>
            setWalletSandbox({ walletId, data })
          }
          isSettingFee={isSettingFee}
          isSettingSandbox={isSettingSandbox}
        />
      </div>
    </LayoutContext.Provider>
  );
}
