export const INPLACE_EDIT_PREFERENCE_KEY = "univents.inplace-edit-enabled";
export const DASHBOARD_SHOW_MONETARY_KEY = "univents.dashboard-show-monetary-default";

export function readInplaceEditPreference() {
  if (typeof window === "undefined") return false;
  return window.localStorage.getItem(INPLACE_EDIT_PREFERENCE_KEY) === "true";
}

export function saveInplaceEditPreference(enabled: boolean) {
  window.localStorage.setItem(INPLACE_EDIT_PREFERENCE_KEY, String(enabled));
  window.dispatchEvent(new CustomEvent("univents:preferences-changed"));
}

export function readDashboardMonetaryPreference(): boolean {
  if (typeof window === "undefined") return false;
  return window.localStorage.getItem(DASHBOARD_SHOW_MONETARY_KEY) === "true";
}

export function saveDashboardMonetaryPreference(enabled: boolean) {
  window.localStorage.setItem(DASHBOARD_SHOW_MONETARY_KEY, String(enabled));
  window.dispatchEvent(new CustomEvent("univents:preferences-changed"));
}
