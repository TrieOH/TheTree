import type { CannotStartWithNumber, OnlyAlphanumeric } from "../../types/permission-types";
import type { PermissionResult } from "./permission.result";

export interface PermissionRoot {
  on<NS extends string, SP extends string>(
    namespace: CannotStartWithNumber<NS>,
    specifier: OnlyAlphanumeric<SP>
  ): PermissionObjectBuilder;

  onWildcard<NS extends string>(
    namespace: CannotStartWithNumber<NS>
  ): PermissionObjectBuilder;

  any(): PermissionFinal;
}

export interface PermissionObjectBuilder {
  on<NS extends string, SP extends string>(
    namespace: CannotStartWithNumber<NS>,
    specifier: SP
  ): PermissionObjectBuilder;

  onWildcard<NS extends string>(
    namespace: CannotStartWithNumber<NS>
  ): PermissionObjectBuilder;

  forAnyChild(): PermissionObjectFinal;
  forAnyDescendant(): PermissionObjectFinal;
  done(): PermissionObjectFinal;
}

export interface PermissionObjectFinal {
  can<ActionToken extends string>(action: ActionToken): PermissionActionChain;
  canAnyAction(): PermissionFinal;
}

export interface PermissionActionChain {
  and<Part extends string>(part: Part): PermissionActionChain;
  build(): PermissionResult;
}

export interface PermissionFinal {
  build(): PermissionResult;
}
