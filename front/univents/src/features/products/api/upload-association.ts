import {
  registerUploadAssociationHandler,
  UploadAssociationError,
  uploadAssociationErrorFromResponse,
} from "@/features/upload-queue";
import { getContext } from "@/integrations/tanstack-query/root-provider";
import { authQueryFetcher } from "@/shared/lib/api/fetch";
import type { VariantI } from "../model";
import { patchVariantFn } from "./index";
import { productKeys } from "./query-keys";

registerUploadAssociationHandler("variant-gallery", async (task, url) => {
  const productId = task.association?.input?.productId;
  if (typeof productId !== "string")
    throw new UploadAssociationError("Produto da variante não encontrado.", {
      status: 400,
    });
  const variants = await authQueryFetcher<VariantI[]>(
    `/products/${productId}/variants`,
  );
  const variant = variants.find((item) => item.id === task.owner.id);
  if (!variant)
    throw new UploadAssociationError("Variante não encontrada.", {
      status: 404,
    });
  const response = await patchVariantFn(variant.id, {
    vendor_code: variant.vendor_code,
    name: variant.name,
    description: variant.description,
    price: variant.price,
    stock: variant.stock,
    gallery_urls: [...(variant.gallery_urls ?? []), url],
  });
  if (!response.success)
    throw uploadAssociationErrorFromResponse(
      response,
      "Não foi possível associar a imagem.",
    );
  getContext().queryClient.setQueryData<VariantI[]>(
    productKeys.variants.byProduct(productId),
    (old = []) =>
      old.map((item) => (item.id === response.data.id ? response.data : item)),
  );
});
