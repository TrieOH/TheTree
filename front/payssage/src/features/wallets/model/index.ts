import type {
  CreateWalletRequest,
  SetFeeBPSRequest,
  SetSandboxRequest,
  Wallet,
} from "@trieoh/payssage-models";
import z from "zod";

export const walletCreateSchema = z.object({
  name: z
    .string({ error: "Name is required" })
    .min(3, "Name must be at least 3 characters long"),
  organization_id: z.string().optional().catch(undefined),
}) satisfies z.ZodType<CreateWalletRequest>;

export type WalletCreateI = CreateWalletRequest;
export interface WalletI extends Wallet {
  collector_id?: string;
}

export interface WalletBindCollectorI {
  collector_id: string;
}

export const walletSetFeeBpsSchema = z.object({
  fee_bps: z
    .number({ error: "Fee is required" })
    .min(0, "Fee must be greater than or equal to 0")
    .max(10000, "Fee must be less than or equal to 10000"),
  organization_id: z.string().optional().catch(undefined),
}) satisfies z.ZodType<SetFeeBPSRequest>;

export type WalletSetFeeBpsI = SetFeeBPSRequest;

export const walletSetSandboxSchema = z.object({
  sandbox: z.boolean({ error: "Sandbox is required" }),
  organization_id: z.string().optional().catch(undefined),
}) satisfies z.ZodType<SetSandboxRequest>;

export type WalletSetSandboxI = SetSandboxRequest;
