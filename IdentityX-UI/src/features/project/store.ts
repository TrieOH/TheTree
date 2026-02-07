import { createCrudActions, createCrudStore } from "@/shared/lib/store/crudStore";
import type { ProjectCRUD } from "./model/types";

// Store Instance
export const projectStore = createCrudStore<ProjectCRUD>();

// Actions
export const projectActions = createCrudActions(projectStore);

// initial state can be used to reset
export const initialProjectState = projectStore.state;