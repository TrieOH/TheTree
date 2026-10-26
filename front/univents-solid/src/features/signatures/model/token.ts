import {
  Schema,
  createSchema,
  minLength,
  type Infer,
} from "@/shared/lib/effect-schema";

export const signatureTokenSearchSchema = createSchema(
  Schema.Struct({
    token: Schema.optional(Schema.String),
  }),
);

export type SignatureTokenSearchI = { token: string };

export const signatureRequestTokenClaimsSchema = createSchema(
  Schema.Struct({
    request_id: Schema.String.pipe(minLength(1)),
    edition_id: Schema.String.pipe(minLength(1)),
    exp: Schema.optional(Schema.Number),
    iat: Schema.optional(Schema.Number),
  }),
);

export type SignatureRequestTokenClaims = Infer<
  typeof signatureRequestTokenClaimsSchema
>;

export const signatureRevocationTokenClaimsSchema = createSchema(
  Schema.Struct({
    signature_id: Schema.String.pipe(minLength(1)),
    edition_id: Schema.String.pipe(minLength(1)),
    exp: Schema.optional(Schema.Number),
    iat: Schema.optional(Schema.Number),
  }),
);

export type SignatureRevocationTokenClaims = Infer<
  typeof signatureRevocationTokenClaimsSchema
>;

export function parseJwtPayload<T>(
  token: string,
  schema?: { safeParse(input: unknown): { success: boolean; data?: unknown } },
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
      return parsedWithSchema.success ? (parsedWithSchema.data as T) : null;
    }
    return parsed as T;
  } catch {
    return null;
  }
}
