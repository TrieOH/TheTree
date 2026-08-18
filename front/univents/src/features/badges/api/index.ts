import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import {
  createBadgeTemplate,
  deleteBadgeTemplate,
  getBadgeTemplate,
  getEditionBadgesPrint,
  listBadgeTemplates,
  listEditionBadgeEmissions,
  listUserBadges,
  updateBadgeTemplate,
} from "@trieoh/univents-api";
import type {
  BadgeEditionEmission,
  BadgePrintItem,
  BadgeProfileGroups,
  BadgeTemplate,
  BadgeTemplateCreate,
} from "../model";
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
    withSpan("action:badge-template-create", () =>
      createBadgeTemplate(editionId, {
        ...data,
        origin: data.origin ?? null,
      }).then(orvalData<BadgeTemplate>),
    ),
);
export const updateBadgeTemplateFn = createClientOnlyFn(
  (templateId: string, data: BadgeTemplateCreate) =>
    withSpan("action:badge-template-update", () =>
      updateBadgeTemplate(templateId, {
        name: data.name,
        design_data: data.design_data,
      }).then(orvalData<BadgeTemplate>),
    ),
);
export const deleteBadgeTemplateFn = createClientOnlyFn((templateId: string) =>
  withSpan("action:badge-template-delete", () =>
    deleteBadgeTemplate(templateId).then(orvalData<null>),
  ),
);

export const userBadgesQueryOptions = (userId: string) =>
  queryOptions({
    queryKey: badgeKeys.user(userId),
    enabled: Boolean(userId),
    queryFn: () =>
      listUserBadges(userId, { public: true }).then(
        orvalData<BadgeProfileGroups>,
      ),
  });

export const badgeEmissionsQueryOptions = (editionId: string) =>
  queryOptions({
    queryKey: badgeKeys.emissions(editionId),
    queryFn: () =>
      listEditionBadgeEmissions(editionId).then(
        orvalData<BadgeEditionEmission[]>,
      ),
  });

export const badgePrintQueryOptions = (
  editionId: string,
  emissionIds?: string[],
) =>
  queryOptions({
    queryKey: badgeKeys.print(editionId, emissionIds),
    enabled: false,
    queryFn: () =>
      getEditionBadgesPrint(editionId, { emission_ids: emissionIds }).then(
        orvalData<BadgePrintItem[]>,
      ),
  });
