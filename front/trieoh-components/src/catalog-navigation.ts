import type { CatalogNavSection } from "@trieoh/ui-solid";
import type { Entry } from "./catalog";

export const ALL_GROUPS = "Todos";
export const groups = [
  "Ações",
  "Formulários",
  "Navegação",
  "Overlays",
  "Feedback",
  "Layout",
];
export const groupHref = (group: string) =>
  group === ALL_GROUPS ? "#/" : `#/group/${encodeURIComponent(group)}`;
export const catalogTop = {
  id: ALL_GROUPS,
  label: "Todos os componentes",
  href: groupHref(ALL_GROUPS),
};

export function catalogSections(
  entries: readonly Entry[],
): CatalogNavSection[] {
  return groups.map((title) => ({
    id: title,
    title,
    href: groupHref(title),
    items: entries
      .filter((entry) => entry.group === title)
      .map((entry) => ({
        id: entry.name,
        label: entry.name,
        href: `#/components/${entry.name}`,
      })),
  }));
}
