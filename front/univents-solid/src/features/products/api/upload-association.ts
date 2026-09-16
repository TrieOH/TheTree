import { orvalData } from "@trieoh/api-client";
import { listProductVariants } from "@trieoh/univents-api";
import {
  registerUploadAssociationHandler,
  UploadAssociationError,
} from "@/features/upload-queue";
import { appQueryClient } from "@/shared/lib/query-client";
import type { VariantI } from "../model";
import { patchVariantFn } from "./index";
import { productKeys } from "./query-keys";

registerUploadAssociationHandler("variant-gallery", async (task, url) => {
  const productId = task.association?.input?.productId;
  if (typeof productId !== "string") {
    throw new UploadAssociationError("Produto da variante não encontrado.", {
      status: 400,
    });
  }

  const variants = await listProductVariants(productId, { public: true }).then(
    orvalData<VariantI[]>,
  );
  const variant = variants.find((item) => item.id === task.owner.id);
  if (!variant) {
    throw new UploadAssociationError("Variante não encontrada.", {
      status: 404,
    });
  }

  await patchVariantFn(variant.id, {
    vendor_code: variant.vendor_code,
    name: variant.name,
    description: variant.description,
    price: variant.price,
    stock: variant.stock,
    gallery_urls: [...(variant.gallery_urls ?? []), url],
  });

  void appQueryClient.invalidateQueries({
    queryKey: productKeys.variants(productId),
  });
});
