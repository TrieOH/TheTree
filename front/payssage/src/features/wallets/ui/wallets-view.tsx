import { PaginatedContainer, useLayoutHeader } from "@trieoh/ui-base";
import { Plus } from "lucide-react";
import { useMemo, useState } from "react";
import { bpsToPercentage, percentageToBps } from "#/shared/lib/utils";
import { Button } from "#/shared/ui/shadcn/button";
import FormModal from "#/widgets/modal/form-modal";
import type {
  WalletCreateI,
  WalletI,
  WalletSetFeeBpsI,
  WalletSetSandboxI,
} from "../model";
import { walletCreateSchema, walletSetFeeBpsSchema } from "../model";
import WalletCard from "./wallet-card";

interface WalletsViewProps {
  wallets: WalletI[];
  organizationId?: string;
  title: string;
  description: string;
  onCreate: (data: WalletCreateI) => void;
  isCreating: boolean;
  onSetFee: (walletId: string, data: WalletSetFeeBpsI) => void;
  onSetSandbox: (walletId: string, data: WalletSetSandboxI) => void;
  isSettingFee?: boolean;
  isSettingSandbox?: boolean;
}

export function WalletsView({
  wallets,
  organizationId,
  title,
  description,
  onCreate,
  isCreating,
  onSetFee,
  onSetSandbox,
  isSettingFee,
  isSettingSandbox,
}: WalletsViewProps) {
  const [filter, setFilter] = useState("");
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [feeWallet, setFeeWallet] = useState<WalletI | null>(null);
  const count = wallets.length;

  const header = useMemo(
    () => (
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-lg font-semibold tracking-tight">{title}</h1>
          <p className="text-sm text-muted-foreground">
            {count === 0
              ? `No wallets yet ${description}`
              : `${count} wallet${count !== 1 ? "s" : ""} ${description}`}
          </p>
        </div>
      </div>
    ),
    [count, description, title],
  );

  useLayoutHeader(header);

  const filteredWallets = wallets.filter((wallet) => {
    const search = filter.toLowerCase().trim();
    if (!search) return true;

    return (
      wallet.name.toLowerCase().includes(search) ||
      wallet.owner_id.toLowerCase().includes(search) ||
      wallet.fee_bps.toString().includes(search)
    );
  });

  return (
    <div className="space-y-6">
      <PaginatedContainer<WalletI>
        items={filteredWallets}
        layout="grid"
        minItemWidth="15rem"
        gap="6"
        pageSize={10}
        sortFields={[
          { key: "name", label: "Name" },
          { key: "fee_bps", label: "Fee" },
          { key: "created_at", label: "Created At" },
        ]}
        filterValue={filter}
        onFilterChange={setFilter}
        filterPlaceholder="Filter by name, owner or fee..."
        itemLabel="wallets"
        headerActions={
          <Button
            onClick={() => setIsCreateOpen(true)}
            variant="outline"
            className="h-9 sm:w-auto px-3 rounded-sm"
          >
            <Plus size={16} />
            Create Wallet
          </Button>
        }
        renderItems={(slice) =>
          slice.map((item) => (
            <WalletCard
              key={item.id}
              data={item}
              onEditFee={setFeeWallet}
              onSetSandbox={onSetSandbox}
              isSettingSandbox={isSettingSandbox}
            />
          ))
        }
      />

      <FormModal<WalletCreateI>
        title="Create Wallet"
        description="Give your wallet a name to identify it."
        schema={walletCreateSchema}
        formId="create-wallet-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={(data) => {
          onCreate({
            ...data,
            organization_id: organizationId,
          });
          setIsCreateOpen(false);
        }}
        defaultValues={{ name: "", organization_id: organizationId }}
        disabled={isCreating}
        buttonTitle="Create Wallet"
        fields={[
          {
            name: "name",
            label: "Wallet Name",
            placeholder: "e.g. Main Checkout Wallet",
            type: "text",
          },
        ]}
      />

      <FormModal<WalletSetFeeBpsI>
        key={feeWallet?.id ?? "fee-wallet"}
        title="Set Wallet Fee"
        description="Update the fee percentage for this wallet."
        buttonTitle={isSettingFee ? "Saving..." : "Save Fee"}
        schema={walletSetFeeBpsSchema}
        formId={feeWallet ? `set-fee-${feeWallet.id}` : "set-fee-wallet"}
        isOpen={feeWallet !== null}
        onClose={() => setFeeWallet(null)}
        onSubmit={(payload) => {
          if (!feeWallet) return;
          onSetFee(feeWallet.id, {
            fee_bps: percentageToBps(Number(payload.fee_bps)),
            organization_id: feeWallet.organization_id,
          });
          setFeeWallet(null);
        }}
        defaultValues={{
          fee_bps: bpsToPercentage(feeWallet?.fee_bps ?? 0),
          organization_id: feeWallet?.organization_id,
        }}
        disabled={isSettingFee}
        fields={[
          {
            name: "fee_bps",
            label: "Fee (%)",
            type: "percentage",
            placeholder: "e.g. 2.5",
          },
        ]}
      />
    </div>
  );
}
