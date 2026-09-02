import { orvalData } from "@trieoh/api-client";
import { listProductVariants } from "@trieoh/univents-api";
import {
  registerUploadAssociationHandler,
  UploadAssociationError,
} from "@/features/upload-queue";
import { getContext } from "@/integrations/tanstack-query/root-provider";
import type { VariantI } from "../model";
import { syncVariantCaches } from "./cache";
import { patchVariantFn } from "./index";

registerUploadAssociationHandler("variant-gallery", async (task, url) => {
  const productId = task.association?.input?.productId;
  if (typeof productId !== "string")
    throw new UploadAssociationError("Produto da variante não encontrado.", {
      status: 400,
    });
  const variants = await listProductVariants(productId, { public: true }).then(
    orvalData<VariantI[]>,
  );
  const variant = variants.find((item) => item.id === task.owner.id);
  if (!variant)
    throw new UploadAssociationError("Variante não encontrada.", {
      status: 404,
    });
  const updated = await patchVariantFn(variant.id, {
    vendor_code: variant.vendor_code,
    name: variant.name,
    description: variant.description,
    price: variant.price,
    stock: variant.stock,
    gallery_urls: [...(variant.gallery_urls ?? []), url],
  });
  syncVariantCaches(getContext().queryClient, updated);
});
