import { createCrudActions, createCrudStore, type CrudActions, type CrudState } from "@/shared/lib/store/crudStore";
import type { ProjectCRUD } from "./model/types";

export interface ProjectStoreState extends CrudState<ProjectCRUD> {}

export interface ProjectStoreActions extends CrudActions<ProjectCRUD> {}

// Store Instance
export const projectStore = createCrudStore<ProjectCRUD>();

// Actions
export const projectActions = createCrudActions(projectStore);

// initial state can be used to reset
export const initialProjectState = projectStore.state;