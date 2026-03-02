import { createCrudActions, createCrudStore } from "@/shared/lib/store/crudStore";
import type { RoleCRUD } from "./model/types";

// Store Instance
export const roleStore = createCrudStore<RoleCRUD>();

// Actions
export const roleActions = createCrudActions(roleStore);

// initial state can be used to reset
export const initialRoleState = roleStore.state;