import type { QueryClient } from "@tanstack/react-query";
import type { BadgeTemplate } from "../model";
import { badgeKeys } from "./query-keys";

function upsertTemplate(
  templates: BadgeTemplate[] | undefined,
  template: BadgeTemplate,
) {
  if (!templates) return templates;
  const index = templates.findIndex((item) => item.id === template.id);
  if (index === -1) return [template, ...templates];
  const next = [...templates];
  next[index] = template;
  return next;
}

function invalidateRenderedBadges(queryClient: QueryClient, editionId: string) {
  void queryClient.invalidateQueries({
    queryKey: badgeKeys.printByEdition(editionId),
  });
  void queryClient.invalidateQueries({ queryKey: badgeKeys.users() });
}

export function syncBadgeTemplateCache(
  queryClient: QueryClient,
  template: BadgeTemplate,
) {
  queryClient.setQueryData<BadgeTemplate[]>(
    badgeKeys.byEdition(template.edition_id),
    (old) => upsertTemplate(old, template),
  );
  queryClient.setQueryData<BadgeTemplate>(
    badgeKeys.detail(template.id),
    (old) => (old ? template : old),
  );
  invalidateRenderedBadges(queryClient, template.edition_id);
}

export function removeBadgeTemplateCache(
  queryClient: QueryClient,
  editionId: string,
  templateId: string,
) {
  queryClient.setQueryData<BadgeTemplate[]>(
    badgeKeys.byEdition(editionId),
    (old) => old?.filter((template) => template.id !== templateId),
  );
  queryClient.removeQueries({ queryKey: badgeKeys.detail(templateId) });
  invalidateRenderedBadges(queryClient, editionId);
}
