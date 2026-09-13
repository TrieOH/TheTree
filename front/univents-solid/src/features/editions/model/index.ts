import type { Edition } from "@trieoh/univents-api/schemas";

export type EditionStatus = "draft" | "future" | "active" | "past";

export interface EditionI extends Edition {
  status: EditionStatus;
}

export function inferEditionStatus(
  edition: Pick<Edition, "is_draft" | "starts_at" | "ends_at">,
  now = new Date(),
): EditionStatus {
  if (edition.is_draft) return "draft";

  const currentTime = now.getTime();
  if (currentTime < new Date(edition.starts_at).getTime()) return "future";
  if (currentTime > new Date(edition.ends_at).getTime()) return "past";
  return "active";
}

export function normalizeEdition(edition: Edition): EditionI {
  return {
    ...edition,
    status: inferEditionStatus(edition),
  };
}
