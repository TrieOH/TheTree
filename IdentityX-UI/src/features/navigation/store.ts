import { Store } from "@tanstack/react-store";
import type { NavigationActions, NavigationStoreState } from "./model/types";

const storedSchemaVersion = typeof window !== "undefined" ? localStorage.getItem("currentSchemaVersion") : null;

// Store Instance
export const navigationStore = new Store<NavigationStoreState>({
  currentSchemaVersion: storedSchemaVersion ? parseInt(storedSchemaVersion, 10) : null,
});

// Actions
export const navigationActions: NavigationActions = {
  setCurrentSchemaVersion: (schemaVersion: number | null) => {
    navigationStore.setState((state) => ({
      ...state,
      currentSchemaVersion: schemaVersion,
    }));
  },
};


navigationStore.subscribe((state) => {
  if (typeof window !== "undefined") {
    const current = state.currentVal ?? state;

    if (current.currentSchemaVersion) localStorage.setItem("currentSchemaVersion", current.currentSchemaVersion.toString());
    else localStorage.removeItem("currentSchemaVersion");
  }
});

// initial state can be used to reset
export const initialNavigationState = navigationStore.state;
