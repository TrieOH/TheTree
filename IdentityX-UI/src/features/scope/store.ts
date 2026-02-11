import { createCrudActions, createCrudStore, type CrudActions, type CrudState } from "@/shared/lib/store/crudStore";
import type { ScopeCRUD } from "./model/types";

export interface ProjectStoreState extends CrudState<ScopeCRUD> {}

export interface ProjectStoreActions extends CrudActions<ScopeCRUD> {}

// Store Instance
export const scopeStore = createCrudStore<ScopeCRUD>();

// Actions
export const scopeActions = createCrudActions(scopeStore);

// initial state can be used to reset
export const initialScopeState = scopeStore.state;