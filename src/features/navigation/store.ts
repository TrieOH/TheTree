import { Store } from "@tanstack/react-store";
import type { NavigationActions, NavigationStoreState } from "./model/types";

// Initialize currentProjectId from localStorage
const storedProjectId = typeof window !== 'undefined' ? localStorage.getItem('currentProjectId') : null;
// Store Instance
export const navigationStore = new Store<NavigationStoreState>({
  currentProjectId: storedProjectId,
});

// Actions
export const navigationActions: NavigationActions = {
  setCurrentProjectId: (projectId: string | null) => {
    navigationStore.setState((state) => ({
      ...state,
      currentProjectId: projectId,
    }));
  },
};

// Subscribe to store changes to persist currentProjectId to localStorage
navigationStore.subscribe((state) => {
  if (typeof window !== 'undefined') {
    if (state.currentVal.currentProjectId)
      localStorage.setItem('currentProjectId', state.currentVal.currentProjectId);
    else localStorage.removeItem('currentProjectId');
  }
});

// initial state can be used to reset
export const initialNavigationState = navigationStore.state;
