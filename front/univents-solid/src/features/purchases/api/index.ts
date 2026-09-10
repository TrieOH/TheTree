import { orvalData } from "@trieoh/api-client";
import { listMyPurchases } from "@trieoh/univents-api";
import type { MyPurchases } from "@trieoh/univents-api/schemas";
import { purchaseKeys } from "./query-keys";

export const myPurchasesQueryOptions = () => ({
  queryKey: purchaseKeys.mine(),
  queryFn: () => listMyPurchases().then(orvalData<MyPurchases>),
});
