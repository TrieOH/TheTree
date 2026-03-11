import z from "zod";

export const oauthSetupSchema = z.object({
  fee_percent: z.coerce.number({
    error: "Fee is required",
  })
    .min(0, "Fee must be at least 0%")
    .max(100, "Fee must be at most 100%"),
});

export type OauthSetupI = z.infer<typeof oauthSetupSchema>;

export interface OauthSetupResponseI {
  redirect_url: string;
}

export interface OauthCallbackResponseI {
  url: string;
}