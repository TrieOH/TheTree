import z from "zod";

export const signatureTokenSearchSchema = z.object({
  token: z.string().default("").catch(""),
});

export type SignatureTokenSearchI = z.infer<typeof signatureTokenSearchSchema>;

export const signatureRequestTokenClaimsSchema = z.object({
  request_id: z.string().min(1),
  edition_id: z.string().min(1),
  exp: z.number().optional(),
  iat: z.number().optional(),
});

export type SignatureRequestTokenClaims = z.infer<
  typeof signatureRequestTokenClaimsSchema
>;

export const signatureRevocationTokenClaimsSchema = z.object({
  signature_id: z.string().min(1),
  edition_id: z.string().min(1),
  exp: z.number().optional(),
  iat: z.number().optional(),
});

export type SignatureRevocationTokenClaims = z.infer<
  typeof signatureRevocationTokenClaimsSchema
>;

export function parseJwtPayload<T>(
  token: string,
  schema?: z.ZodType<T>,
): T | null {
  if (!token || typeof token !== "string") return null;
  try {
    const parts = token.split(".");
    if (parts.length < 2) return null;
    const base64Url = parts[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join(""),
    );
    const parsed = JSON.parse(jsonPayload);
    if (schema) {
      const parsedWithSchema = schema.safeParse(parsed);
      return parsedWithSchema.success ? parsedWithSchema.data : null;
    }
    return parsed as T;
  } catch {
    return null;
  }
}
