import { createCrudActions, createCrudStore, type CrudActions, type CrudState } from "@/shared/lib/store/crudStore";
import type { SchemaCRUD } from "./model/types";

export interface SchemaStoreState extends CrudState<SchemaCRUD> {}

export interface SchemaStoreActions extends CrudActions<SchemaCRUD> {}

// Store Instance
export const schemaStore = createCrudStore<SchemaCRUD>();

// Actions
export const schemaActions = createCrudActions(schemaStore);

// initial state can be used to reset
export const initialSchemaState = schemaStore.state;