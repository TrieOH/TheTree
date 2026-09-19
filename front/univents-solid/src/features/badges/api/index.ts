import { orvalData } from "@trieoh/api-client";
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
} from "../model";
import { badgeKeys } from "./query-keys";

export const badgeTemplatesQueryOptions = (editionId: string) => ({
  queryKey: badgeKeys.byEdition(editionId),
  queryFn: () => listBadgeTemplates(editionId).then(orvalData<BadgeTemplate[]>),
});

export const badgeTemplateQueryOptions = (templateId: string) => ({
  queryKey: badgeKeys.detail(templateId),
  queryFn: () => getBadgeTemplate(templateId).then(orvalData<BadgeTemplate>),
});

export const userBadgesQueryOptions = (actorId: string) => ({
  queryKey: badgeKeys.user(actorId),
  queryFn: () =>
    listUserBadges(actorId, { public: true }).then(
      orvalData<BadgeProfileGroups>,
    ),
});

export const badgeEmissionsQueryOptions = (editionId: string) => ({
  queryKey: badgeKeys.emissions(editionId),
  queryFn: () =>
    listEditionBadgeEmissions(editionId).then(
      orvalData<BadgeEditionEmission[]>,
    ),
});

export const badgePrintQueryOptions = (
  editionId: string,
  emissionIds?: string[],
) => ({
  queryKey: badgeKeys.print(editionId, emissionIds),
  queryFn: () =>
    getEditionBadgesPrint(editionId, { emission_ids: emissionIds }).then(
      orvalData<BadgePrintItem[]>,
    ),
});

export * from "./query-keys";

export {
  createBadgeTemplate,
  deleteBadgeTemplate,
  getBadgeTemplate,
  getEditionBadgesPrint,
  listBadgeTemplates,
  listEditionBadgeEmissions,
  listUserBadges,
  updateBadgeTemplate,
};
