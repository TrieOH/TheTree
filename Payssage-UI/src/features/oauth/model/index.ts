import z from "zod";


const oauthSetupSchema = z.object({
  fee_bps: z.number({ error: "Fee bps is required" })
    .min(0, "Fee bps must be at least 0")
    .max(10000, "Fee bps must be at most 10000"),
});

export type OauthSetupI = z.infer<typeof oauthSetupSchema>;

export interface OauthSetupResponseI {
  redirect_url: string;
}

export interface OauthCallbackResponseI {
  url: string;
}