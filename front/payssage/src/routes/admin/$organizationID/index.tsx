import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { toast } from "sonner";
import {
  allWalletsQueryOptions,
  createWalletFn,
  setWalletFeeBPSFn,
  setWalletSandboxFn,
} from "#/features/wallets/api";
import type {
  WalletCreateI,
  WalletI,
  WalletSetFeeBpsI,
  WalletSetSandboxI,
} from "#/features/wallets/model";
import { WalletsView } from "#/features/wallets/ui/wallets-view";

export const Route = createFileRoute("/admin/$organizationID/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { organizationID } = Route.useParams();
  const queryClient = useQueryClient();

  const { data: wallets = [] } = useQuery(
    allWalletsQueryOptions(organizationID),
  );

  const { mutate: createWallet, isPending: isCreating } = useMutation({
    mutationFn: (data: WalletCreateI) => createWalletFn(data),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allWalletsQueryOptions(organizationID).queryKey,
          (old: WalletI[] = []) => [response.data, ...old],
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

  return (
    <WalletsView
      wallets={wallets}
      organizationId={organizationID}
      title="Wallets"
      description="for this organization"
      onCreate={createWallet}
      isCreating={isCreating}
      onSetFee={(walletId, data) => setWalletFee({ walletId, data })}
      onSetSandbox={(walletId, data) => setWalletSandbox({ walletId, data })}
      isSettingFee={isSettingFee}
      isSettingSandbox={isSettingSandbox}
    />
  );
}
