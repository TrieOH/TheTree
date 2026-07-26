export const productKeys = {
  all: ["products"] as const,

  lists: () => [...productKeys.all, "list"] as const,
  editionList: (editionId: string) =>
    [...productKeys.lists(), "edition", editionId] as const,
  detail: (productId: string) => [...productKeys.all, productId] as const,
  byVendorCode: (editionId: string, vendorCode: string) =>
    [...productKeys.all, "vendor-code", editionId, vendorCode] as const,

  variants: {
    all: ["variants"] as const,
    byProduct: (productId: string) =>
      [...productKeys.variants.all, "product", productId] as const,
    byVendorCode: (editionId: string, vendorCode: string) =>
      [
        ...productKeys.variants.all,
        "vendor-code",
        editionId,
        vendorCode,
      ] as const,
  },
};
