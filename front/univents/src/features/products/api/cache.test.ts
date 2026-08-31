import { QueryClient } from "@tanstack/react-query";
import { describe, expect, it } from "vitest";
import type { ProductI, VariantI } from "../model";
import {
  removeProductCaches,
  removeVariantCaches,
  syncProductCaches,
  syncVariantCaches,
} from "./cache";
import { productKeys } from "./query-keys";

const product = (overrides: Partial<ProductI> = {}): ProductI => ({
  id: "product-1",
  edition_id: "edition-1",
  vendor_code: "PRODUCT",
  requires_registration: false,
  created_at: "2026-01-01T00:00:00Z",
  updated_at: null,
  deleted_at: null,
  ...overrides,
});

const variant = (overrides: Partial<VariantI> = {}): VariantI => ({
  id: "variant-1",
  edition_id: "edition-1",
  product_id: "product-1",
  vendor_code: "VARIANT",
  name: "Variant",
  description: null,
  price: 1000,
  stock: 10,
  gallery_urls: [],
  created_at: "2026-01-01T00:00:00Z",
  updated_at: null,
  deleted_at: null,
  ...overrides,
});

describe("product cache synchronization", () => {
  it("does not create partial product or variant lists", () => {
    const queryClient = new QueryClient();

    syncProductCaches(queryClient, product());
    syncVariantCaches(queryClient, variant());

    expect(
      queryClient.getQueryData(productKeys.editionList("edition-1")),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(productKeys.variants.byProduct("product-1")),
    ).toBeUndefined();
  });

  it("updates loaded product and variant lists", () => {
    const queryClient = new QueryClient();
    const previousProduct = product();
    const previousVariant = variant();
    const updatedProduct = product({ vendor_code: "UPDATED" });
    const updatedVariant = variant({ name: "Updated variant" });

    queryClient.setQueryData(productKeys.editionList("edition-1"), [
      previousProduct,
    ]);
    queryClient.setQueryData(productKeys.variants.byProduct("product-1"), [
      previousVariant,
    ]);

    syncProductCaches(queryClient, updatedProduct);
    syncVariantCaches(queryClient, updatedVariant);

    expect(
      queryClient.getQueryData(productKeys.editionList("edition-1")),
    ).toEqual([updatedProduct]);
    expect(
      queryClient.getQueryData(productKeys.variants.byProduct("product-1")),
    ).toEqual([updatedVariant]);
  });

  it("removes products, variants and product-owned detail caches", () => {
    const queryClient = new QueryClient();
    queryClient.setQueryData(productKeys.editionList("edition-1"), [product()]);
    queryClient.setQueryData(productKeys.detail("product-1"), product());
    queryClient.setQueryData(productKeys.variants.byProduct("product-1"), [
      variant(),
    ]);

    removeVariantCaches(queryClient, "edition-1", "product-1", "variant-1");
    expect(
      queryClient.getQueryData(productKeys.variants.byProduct("product-1")),
    ).toEqual([]);

    removeProductCaches(queryClient, "edition-1", "product-1");
    expect(
      queryClient.getQueryData(productKeys.editionList("edition-1")),
    ).toEqual([]);
    expect(
      queryClient.getQueryData(productKeys.detail("product-1")),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(productKeys.variants.byProduct("product-1")),
    ).toBeUndefined();
  });
});
