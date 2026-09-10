import { orvalData } from "@trieoh/api-client";
import { listMyCertifications } from "@trieoh/univents-api";
import type { Certification } from "@trieoh/univents-api/schemas";
import { certificationKeys } from "./query-keys";

export const myCertificationsQueryOptions = () => ({
  queryKey: certificationKeys.mine(),
  queryFn: () => listMyCertifications().then(orvalData<Certification[]>),
});
