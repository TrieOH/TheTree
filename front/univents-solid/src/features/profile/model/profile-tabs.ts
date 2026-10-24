export const PROFILE_TABS = ["about", "badges", "certificates", "purchases"] as const;
export type ProfileTab = (typeof PROFILE_TABS)[number];

export const PROFILE_TAB_LABELS: Record<ProfileTab, string> = {
  about: "Sobre",
  badges: "Crachás",
  certificates: "Certificados",
  purchases: "Compras",
};
