import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { sellersQueryOptions } from "#/features/sellers/api";
import { cn } from "#/shared/lib/utils";
import { Button } from "#/shared/ui/shadcn/button";
import { Input } from "#/shared/ui/shadcn/input";
import { Label } from "#/shared/ui/shadcn/label";
import { createWalletIntentFn } from "../api";
import type { CreateIntentFormValues } from "../model";
import { createIntentSchema } from "../model";

const fieldClassName =
  "rounded-none border-border font-mono focus-visible:border-primary focus-visible:ring-0";

export function CreateIntentForm({ walletId }: { walletId: string }) {
  const { data: sellers = [], isLoading: isLoadingSellers } = useQuery(
    sellersQueryOptions(walletId),
  );
  const activeSellers = sellers.filter((seller) => !seller.revoked_at);
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateIntentFormValues>({
    resolver: zodResolver(createIntentSchema),
    defaultValues: {
      seller_id: "",
      currency: "BRL",
      amount_cents: 0,
      checkout_provider_data: "{}",
    },
  });

  const { mutate, isPending } = useMutation({
    mutationFn: (values: CreateIntentFormValues) =>
      createWalletIntentFn(walletId, {
        seller_id: values.seller_id,
        currency: values.currency,
        amount_cents: values.amount_cents,
        checkout_provider_data: JSON.parse(
          values.checkout_provider_data,
        ) as Record<string, unknown>,
      }),
    onSuccess: (response) => {
      if (!response.success) {
        toast.error(response.message || "Failed to create intent");
        return;
      }
      toast.success("Payment intent created");
      reset({
        seller_id: "",
        currency: "BRL",
        amount_cents: 0,
        checkout_provider_data: "{}",
      });
    },
    onError: () => toast.error("Failed to create payment intent"),
  });

  return (
    <section className="space-y-4 rounded-sm border bg-card p-4">
      <div>
        <h2 className="text-sm font-semibold">Create payment intent</h2>
        <p className="text-xs text-muted-foreground">
          Create an intent using this wallet.
        </p>
      </div>

      <form
        onSubmit={handleSubmit((values) => mutate(values))}
        className="space-y-4"
      >
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label
              htmlFor="seller_id"
              className="text-[10px] font-black uppercase tracking-[0.2em]"
            >
              Seller
            </Label>
            <select
              id="seller_id"
              disabled={isLoadingSellers || activeSellers.length === 0}
              className={cn(
                "h-9 w-full border bg-background px-2.5 text-sm outline-none disabled:cursor-not-allowed disabled:opacity-50",
                fieldClassName,
                errors.seller_id && "border-destructive",
              )}
              {...register("seller_id")}
            >
              <option value="">
                {isLoadingSellers
                  ? "Loading sellers..."
                  : activeSellers.length
                    ? "Select a seller"
                    : "No active sellers available"}
              </option>
              {activeSellers.map((seller) => (
                <option key={seller.id} value={seller.id}>
                  {seller.provider.replaceAll("_", " ")} —{" "}
                  {seller.provider_user_id}
                </option>
              ))}
            </select>
            {errors.seller_id && (
              <p className="text-xs text-destructive">
                {errors.seller_id.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label
              htmlFor="currency"
              className="text-[10px] font-black uppercase tracking-[0.2em]"
            >
              Currency
            </Label>
            <Input
              id="currency"
              placeholder="BRL"
              className={fieldClassName}
              {...register("currency")}
            />
            {errors.currency && (
              <p className="text-xs text-destructive">
                {errors.currency.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label
              htmlFor="amount_cents"
              className="text-[10px] font-black uppercase tracking-[0.2em]"
            >
              Amount (cents)
            </Label>
            <Input
              id="amount_cents"
              type="number"
              min={1}
              step={1}
              placeholder="1000"
              className={fieldClassName}
              {...register("amount_cents", { valueAsNumber: true })}
            />
            {errors.amount_cents && (
              <p className="text-xs text-destructive">
                {errors.amount_cents.message}
              </p>
            )}
          </div>
        </div>

        <div className="space-y-2">
          <Label
            htmlFor="checkout_provider_data"
            className="text-[10px] font-black uppercase tracking-[0.2em]"
          >
            Provider data (JSON)
          </Label>
          <textarea
            id="checkout_provider_data"
            rows={8}
            spellCheck={false}
            placeholder={'{\n  "key": "value"\n}'}
            className={cn(
              "w-full resize-y border bg-transparent p-3 text-sm outline-none transition-colors",
              fieldClassName,
              errors.checkout_provider_data && "border-destructive",
            )}
            {...register("checkout_provider_data")}
          />
          {errors.checkout_provider_data && (
            <p className="text-xs text-destructive">
              {errors.checkout_provider_data.message}
            </p>
          )}
        </div>

        <div className="flex justify-end">
          <Button
            type="submit"
            disabled={isPending || isLoadingSellers || !activeSellers.length}
            className="rounded-none font-black uppercase tracking-widest"
          >
            {isPending ? "Creating..." : "Create intent"}
          </Button>
        </div>
      </form>
    </section>
  );
}
