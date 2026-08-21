import type { EditionI } from "./index";

type VisualSource = Partial<Pick<EditionI, "logo_url" | "banner_url">>;

export function resolveEditionVisuals(
  event: VisualSource,
  activeEdition: VisualSource | null | undefined,
  upcomingEditions: VisualSource[],
  pastEditions: VisualSource[],
) {
  const sources = [activeEdition, upcomingEditions[0], event, pastEditions[0]];

  return {
    logo_url: sources.find((source) => source?.logo_url)?.logo_url ?? null,
    banner_url:
      sources.find((source) => source?.banner_url)?.banner_url ?? null,
  };
}
