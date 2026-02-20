import { Store } from "@tanstack/react-store";
import type { NavigationActions, NavigationStoreState } from "./model/types";

const storedProjectId = typeof window !== 'undefined' ? localStorage.getItem('currentProjectId') : null;
const storedSchemaId = typeof window !== "undefined" ? localStorage.getItem("currentSchemaId") : null;
const storedSchemaVersion = typeof window !== "undefined" ? localStorage.getItem("currentSchemaVersion") : null;

// Store Instance
export const navigationStore = new Store<NavigationStoreState>({
  currentProjectId: storedProjectId,
  currentSchemaId: storedSchemaId,
  currentSchemaVersion: storedSchemaVersion ? parseInt(storedSchemaVersion, 10) : null,
});

// Actions
export const navigationActions: NavigationActions = {
  setCurrentProjectId: (projectId: string | null) => {
    navigationStore.setState((state) => ({
      ...state,
      currentProjectId: projectId,
    }));
  },
  setCurrentSchemaId: (schemaId: string | null) => {
    navigationStore.setState((state) => ({
      ...state,
      currentSchemaId: schemaId,
    }));
  },
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

    if (current.currentProjectId) localStorage.setItem("currentProjectId", current.currentProjectId);
    else localStorage.removeItem("currentProjectId");

    if (current.currentSchemaId) localStorage.setItem("currentSchemaId", current.currentSchemaId);
    else localStorage.removeItem("currentSchemaId");

    if (current.currentSchemaVersion) localStorage.setItem("currentSchemaVersion", current.currentSchemaVersion.toString());
    else localStorage.removeItem("currentSchemaVersion");
  }
});

// initial state can be used to reset
export const initialNavigationState = navigationStore.state;
