import { createSignal, onSettled } from "solid-js";
import {
  readDashboardMonetaryPreference,
} from "@/features/profile/lib/preferences";

/**
 * Hook to manage monetary visibility in dashboards.
 *
 * Rules:
 * - Always starts with the user's default setting (which defaults to false / hidden).
 * - Allows the user to toggle visibility for the current dashboard session.
 * - Listens for preferences changes (e.g. if modified in another tab or in profile settings).
 */
export function useMonetaryVisibility() {
  const [showMonetary, setShowMonetary] = createSignal(
    readDashboardMonetaryPreference(),
  );

  onSettled(() => {
    setShowMonetary(readDashboardMonetaryPreference());

    const handlePrefChange = () => {
      setShowMonetary(readDashboardMonetaryPreference());
    };

    window.addEventListener("univents:preferences-changed", handlePrefChange);

    return () => {
      window.removeEventListener(
        "univents:preferences-changed",
        handlePrefChange,
      );
    };
  });

  const toggleMonetary = () => setShowMonetary((v) => !v);

  return {
    showMonetary,
    setShowMonetary,
    toggleMonetary,
  };
}
