import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useForm, useWatch } from "react-hook-form";
import { toast } from "sonner";
import { env } from "#/env";
import { collectorsQueryOptions } from "#/features/collectors/api";
import { sellersQueryOptions } from "#/features/sellers/api";
import { allWalletsQueryOptions } from "#/features/wallets/api";
import { cn } from "#/shared/lib/utils";
import { Button } from "#/shared/ui/shadcn/button";
import { Input } from "#/shared/ui/shadcn/input";
import { Label } from "#/shared/ui/shadcn/label";
import { createTestModeWalletIntentFn } from "../api";
import type { CreateTestModeIntentFormValues } from "../model";
import { createTestModeIntentSchema, intentStatuses } from "../model";

const controlClassName =
  "h-9 w-full rounded-none border border-border bg-background px-2.5 text-sm outline-none focus:border-primary disabled:cursor-not-allowed disabled:opacity-50";

export function CreateTestModeIntentForm() {
  const { data: wallets = [] } = useQuery(allWalletsQueryOptions());
  const {
    control,
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateTestModeIntentFormValues>({
    resolver: zodResolver(createTestModeIntentSchema),
    defaultValues: {
      wallet_id: "",
      seller_id: "",
      collector_id: "",
      amount_cents: 1000,
      currency: "BRL",
      sandbox: true,
      provider: env.VITE_SUPPORTED_PROVIDERS[0] ?? "mercadopago",
      status: "succeeded",
      provider_data: JSON.stringify(
        {
          installments: 1,
          token: "replace-me",
          payment_method_id: "replace-me",
          payer: {
            email: "replace-me@example.com",
            identification_type: "CPF",
            identification_number: "replace-me",
          },
        },
        null,
        2,
      ),
      metadata: "{}",
    },
  });
  const walletId = useWatch({ control, name: "wallet_id" });
  const selectedWallet = wallets.find((wallet) => wallet.id === walletId);
  const sellersOptions = sellersQueryOptions(walletId);
  const { data: sellers = [] } = useQuery({
    ...sellersOptions,
    enabled: Boolean(walletId),
  });
  const { data: collectors = [] } = useQuery({
    ...collectorsQueryOptions(selectedWallet?.organization_id),
    enabled: Boolean(walletId),
  });
  const activeSellers = sellers.filter((seller) => !seller.revoked_at);
  const activeCollectors = collectors.filter(
    (collector) => !collector.revoked_at,
  );

  const { mutate, isPending } = useMutation({
    mutationFn: (values: CreateTestModeIntentFormValues) =>
      createTestModeWalletIntentFn({
        wallet_id: values.wallet_id,
        seller_id: values.seller_id,
        ...(values.collector_id && { collector_id: values.collector_id }),
        amount_cents: values.amount_cents,
        currency: values.currency,
        sandbox: values.sandbox,
        provider: values.provider,
        status: values.status,
        provider_data: JSON.parse(values.provider_data) as Record<
          string,
          unknown
        >,
        metadata: JSON.parse(values.metadata) as Record<string, unknown>,
      }),
    onSuccess: (response) => {
      if (!response.success)
        return toast.error(response.message || "Failed to create test intent");
      toast.success("Test intent created");
      reset();
    },
    onError: () => toast.error("Failed to create test intent"),
  });

  const error = (name: keyof CreateTestModeIntentFormValues) =>
    errors[name] ? (
      <p className="text-xs text-destructive">{errors[name].message}</p>
    ) : null;

  return (
    <section className="rounded-sm border bg-card p-4">
      <form
        onSubmit={handleSubmit((values) => mutate(values))}
        className="space-y-5"
      >
        <div className="grid gap-4 md:grid-cols-2">
          <Field label="Wallet" name="wallet_id">
            <select
              id="wallet_id"
              className={controlClassName}
              {...register("wallet_id")}
            >
              <option value="">Select a wallet</option>
              {wallets.map((wallet) => (
                <option key={wallet.id} value={wallet.id}>
                  {wallet.name}
                </option>
              ))}
            </select>
            {error("wallet_id")}
          </Field>

          <Field label="Seller" name="seller_id">
            <select
              id="seller_id"
              disabled={!walletId}
              className={controlClassName}
              {...register("seller_id")}
            >
              <option value="">
                {walletId ? "Select a seller" : "Select a wallet first"}
              </option>
              {activeSellers.map((seller) => (
                <option key={seller.id} value={seller.id}>
                  {seller.provider} — {seller.provider_user_id}
                </option>
              ))}
            </select>
            {error("seller_id")}
          </Field>

          <Field label="Collector (optional)" name="collector_id">
            <select
              id="collector_id"
              disabled={!walletId}
              className={controlClassName}
              {...register("collector_id")}
            >
              <option value="">No collector</option>
              {activeCollectors.map((collector) => (
                <option key={collector.id} value={collector.id}>
                  {collector.provider} — {collector.provider_user_id}
                </option>
              ))}
            </select>
            {error("collector_id")}
          </Field>

          <Field label="Amount (cents)" name="amount_cents">
            <Input
              id="amount_cents"
              type="number"
              min={1}
              step={1}
              className={controlClassName}
              {...register("amount_cents", { valueAsNumber: true })}
            />
            {error("amount_cents")}
          </Field>

          <Field label="Currency" name="currency">
            <Input
              id="currency"
              className={controlClassName}
              {...register("currency")}
            />
            {error("currency")}
          </Field>

          <Field label="Provider" name="provider">
            <select
              id="provider"
              className={controlClassName}
              {...register("provider")}
            >
              {env.VITE_SUPPORTED_PROVIDERS.map((provider) => (
                <option key={provider} value={provider}>
                  {provider}
                </option>
              ))}
            </select>
            {error("provider")}
          </Field>

          <Field label="Status" name="status">
            <select
              id="status"
              className={controlClassName}
              {...register("status")}
            >
              {intentStatuses.map((status) => (
                <option key={status} value={status}>
                  {status}
                </option>
              ))}
            </select>
            {error("status")}
          </Field>

          <label className="flex items-center gap-3 self-end pb-2 text-sm font-medium">
            <input
              type="checkbox"
              className="size-4"
              {...register("sandbox")}
            />
            Sandbox intent
          </label>
        </div>

        <div className="grid gap-4 lg:grid-cols-2">
          <JsonField
            id="provider_data"
            label="Provider data (JSON)"
            error={errors.provider_data?.message}
            register={register("provider_data")}
          />
          <JsonField
            id="metadata"
            label="Metadata (JSON)"
            error={errors.metadata?.message}
            register={register("metadata")}
          />
        </div>

        <div className="flex justify-end">
          <Button
            type="submit"
            disabled={isPending}
            className="rounded-none font-black uppercase tracking-widest"
          >
            {isPending ? "Creating..." : "Create test intent"}
          </Button>
        </div>
      </form>
    </section>
  );
}

function Field({
  label,
  name,
  children,
}: {
  label: string;
  name: string;
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-2">
      <Label
        htmlFor={name}
        className="text-[10px] font-black uppercase tracking-[0.2em]"
      >
        {label}
      </Label>
      {children}
    </div>
  );
}

function JsonField({
  id,
  label,
  error,
  register,
}: {
  id: string;
  label: string;
  error?: string;
  register: React.TextareaHTMLAttributes<HTMLTextAreaElement>;
}) {
  return (
    <Field label={label} name={id}>
      <textarea
        id={id}
        rows={9}
        spellCheck={false}
        className={cn(
          controlClassName,
          "h-auto resize-y p-3 font-mono",
          error && "border-destructive",
        )}
        {...register}
      />
      {error && <p className="text-xs text-destructive">{error}</p>}
    </Field>
  );
}
