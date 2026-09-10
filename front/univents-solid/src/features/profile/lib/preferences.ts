export const INPLACE_EDIT_PREFERENCE_KEY = "univents.inplace-edit-enabled";

export function readInplaceEditPreference() {
  if (typeof window === "undefined") return false;
  return window.localStorage.getItem(INPLACE_EDIT_PREFERENCE_KEY) === "true";
}

export function saveInplaceEditPreference(enabled: boolean) {
  window.localStorage.setItem(INPLACE_EDIT_PREFERENCE_KEY, String(enabled));
  window.dispatchEvent(new CustomEvent("univents:preferences-changed"));
}
