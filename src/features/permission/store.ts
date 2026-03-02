import { createCrudActions, createCrudStore, type CrudActions, type CrudState } from "@/shared/lib/store/crudStore";
import type { PermissionCRUD } from "./model/types";

export interface PermissionStoreState extends CrudState<PermissionCRUD> {}

export interface PermissionStoreActions extends CrudActions<PermissionCRUD> {}

// Store Instance
export const permissionStore = createCrudStore<PermissionCRUD>();

// Actions
export const permissionActions = createCrudActions(permissionStore);

// initial state can be used to reset
export const initialPermissionState = permissionStore.state;