import { createCrudActions, type CrudActions, type CrudState } from "@/shared/lib/store/crudStore";
import type { ProjectCRUD } from "./model/types";
import { Store } from "@tanstack/react-store";

export interface ProjectStoreState extends CrudState<ProjectCRUD> {
  currentProjectId: string | null;
}

export interface ProjectStoreActions extends CrudActions<ProjectCRUD> {
  setCurrentProjectId: (projectId: string | null) => void;
}

// Initialize currentProjectId from localStorage
const storedProjectId = typeof window !== 'undefined' ? localStorage.getItem('currentProjectId') : null;

// Store Instance
export const projectStore = new Store<ProjectStoreState>({
  mode: null,
  selectedItem: null,
  isLoading: false,
  isOpen: false,
  formData: null,
  currentProjectId: storedProjectId, // Initialize currentProjectId from localStorage
});

// Actions
export const projectActions: ProjectStoreActions = {
  ...createCrudActions(projectStore),
  setCurrentProjectId: (projectId: string | null) => {
    projectStore.setState((state) => ({
      ...state,
      currentProjectId: projectId,
    }));
  },
};

// Subscribe to store changes to persist currentProjectId to localStorage
projectStore.subscribe((state) => {
  if (typeof window !== 'undefined') {
    if (state.currentVal.currentProjectId)
      localStorage.setItem('currentProjectId', state.currentVal.currentProjectId);
    else localStorage.removeItem('currentProjectId');
  }
});

// initial state can be used to reset
export const initialProjectState = projectStore.state;