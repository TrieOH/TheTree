import {
  Schema,
  createSchema,
  minLength,
  maxLength,
  nonNegative,
  positive,
  refine,
  type Infer,
} from "@/shared/lib/effect-schema";

export const productStep1Schema = createSchema(
  Schema.Struct({
    vendor_code: Schema.Trim.pipe(
      minLength(2, { message: "O código do produto deve ter pelo menos 2 caracteres." }),
    ),
  }),
);

export const productVariantStepSchema = createSchema(
  Schema.Struct({
    variant_vendor_code: Schema.Trim.pipe(
      minLength(2, { message: "O código da variação deve ter pelo menos 2 caracteres." }),
    ),
    variant_name: Schema.Trim.pipe(
      minLength(2, { message: "O nome da variação deve ter pelo menos 2 caracteres." }),
    ),
    variant_price_cents: Schema.Unknown.pipe(
      refine(
        (v) => typeof v === "number" && !isNaN(v) && v >= 0,
        { message: "Informe um preço válido (>= 0)." },
      ),
    ),
    variant_stock: Schema.optional(
      Schema.Unknown.pipe(
        refine(
          (v) => {
            if (v === null || v === undefined || String(v).trim() === "") return true;
            const n = Number(v);
            return !isNaN(n) && Number.isInteger(n) && n > 0;
          },
          { message: "Quantidade em estoque deve ser um número inteiro maior que 0." },
        ),
      ),
    ),
  }),
);

export const variantStep1Schema = createSchema(
  Schema.Struct({
    vendor_code: Schema.Trim.pipe(
      minLength(2, { message: "O código deve ter pelo menos 2 caracteres." }),
    ),
    name: Schema.Trim.pipe(
      minLength(2, { message: "O nome deve ter pelo menos 2 caracteres." }),
    ),
  }),
);

export const variantStep2Schema = createSchema(
  Schema.Struct({
    price_cents: Schema.Unknown.pipe(
      refine(
        (v) => typeof v === "number" && !isNaN(v) && v >= 0,
        { message: "Informe um preço válido (>= 0)." },
      ),
    ),
    stock: Schema.optional(
      Schema.Unknown.pipe(
        refine(
          (v) => {
            if (v === null || v === undefined || String(v).trim() === "") return true;
            const n = Number(v);
            return !isNaN(n) && Number.isInteger(n) && n > 0;
          },
          { message: "Quantidade em estoque deve ser um número inteiro maior que 0." },
        ),
      ),
    ),
  }),
);

export const createInitialProductSchema = createSchema(
  Schema.Struct({
    requires_registration: Schema.optional(Schema.Boolean),
    vendor_code: Schema.Trim.pipe(
      minLength(2, { message: "O código do produto deve ter pelo menos 2 caracteres." }),
      maxLength(255),
    ),
    variant_vendor_code: Schema.Trim.pipe(
      minLength(2, { message: "O código da variação deve ter pelo menos 2 caracteres." }),
      maxLength(255),
    ),
    name: Schema.Trim.pipe(
      minLength(2, { message: "O nome da variação deve ter pelo menos 2 caracteres." }),
    ),
    description: Schema.optional(
      Schema.NullOr(Schema.String.pipe(maxLength(8000))),
    ),
    price: Schema.Number.pipe(
      nonNegative({ message: "Informe um preço válido (>= 0)." }),
    ),
    stock: Schema.optional(
      Schema.NullOr(
        Schema.Number.pipe(
          positive({ message: "Quantidade em estoque deve ser um número inteiro maior que 0." }),
        ),
      ),
    ),
  }),
);

export type CreateInitialProductInputI = Infer<
  typeof createInitialProductSchema
>;
export type CreateInitialProductOutputI = Infer<
  typeof createInitialProductSchema
>;

export const productPatchSchema = createSchema(
  Schema.Struct({
    requires_registration: Schema.optional(Schema.Boolean),
    vendor_code: Schema.Trim.pipe(
      minLength(2, { message: "O código do produto deve ter pelo menos 2 caracteres." }),
      maxLength(255),
    ),
  }),
);

export type ProductPatchInputI = Infer<typeof productPatchSchema>;
export type ProductPatchOutputI = Infer<typeof productPatchSchema>;

export interface ProductI {
  id: string;
  edition_id: string;
  vendor_code: string;
  requires_registration: boolean;
  created_at: string;
  updated_at?: string | null;
  deleted_at?: string | null;
}

export const variantCreateSchema = createSchema(
  Schema.Struct({
    vendor_code: Schema.Trim.pipe(
      minLength(2, { message: "O código deve ter pelo menos 2 caracteres." }),
      maxLength(255),
    ),
    name: Schema.Trim.pipe(
      minLength(2, { message: "O nome deve ter pelo menos 2 caracteres." }),
    ),
    description: Schema.optional(
      Schema.NullOr(Schema.String.pipe(maxLength(8000))),
    ),
    price: Schema.Number.pipe(
      nonNegative({ message: "Informe um preço válido (>= 0)." }),
    ),
    stock: Schema.optional(
      Schema.NullOr(
        Schema.Number.pipe(
          positive({ message: "Quantidade em estoque deve ser um número inteiro maior que 0." }),
        ),
      ),
    ),
    gallery_urls: Schema.optional(Schema.Array(Schema.String)),
  }),
);

export type VariantCreateInputI = Infer<typeof variantCreateSchema>;
export type VariantCreateOutputI = Infer<typeof variantCreateSchema>;

export interface VariantI {
  id: string;
  edition_id: string;
  product_id: string;
  vendor_code: string;
  name: string;
  description?: string | null;
  price: number;
  stock?: number | null;
  gallery_urls: string[];
  created_at: string;
  updated_at?: string | null;
  deleted_at?: string | null;
}
