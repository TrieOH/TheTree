import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  createBadgeTemplate,
  deleteBadgeTemplate,
  getBadgeTemplate,
  listBadgeTemplates,
} from "@trieoh/univents-api";
import type { BadgeTemplate, BadgeTemplateCreate } from "../model";
import { badgeKeys } from "./query-keys";

export const listBadgeTemplatesFn = createClientOnlyFn((editionId: string) =>
  listBadgeTemplates(editionId).then(orvalData<BadgeTemplate[]>),
);
export const badgeTemplatesQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: badgeKeys.byEdition(editionId),
    queryFn: () => listBadgeTemplatesFn(editionId),
  });
export const getBadgeTemplateFn = createClientOnlyFn((templateId: string) =>
  getBadgeTemplate(templateId).then(orvalData<BadgeTemplate>),
);
export const badgeTemplateQueryOptions = (templateId: string) =>
  queryOptions({
    queryKey: badgeKeys.detail(templateId),
    queryFn: () => getBadgeTemplateFn(templateId),
  });
export const createBadgeTemplateFn = createClientOnlyFn(
  (editionId: string, data: BadgeTemplateCreate) =>
    createBadgeTemplate(editionId, data).then(orvalData<BadgeTemplate>),
);
export const deleteBadgeTemplateFn = createClientOnlyFn((templateId: string) =>
  deleteBadgeTemplate(templateId).then(orvalData<null>),
);
