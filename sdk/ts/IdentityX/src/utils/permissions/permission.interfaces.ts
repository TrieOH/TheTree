import type { ValidateNamespace, ValidateSpecifier } from "../../types/permission-types";
import type { PermissionResult } from "./permission.result";

export interface PermissionRoot {
  on<NS extends string, SP extends string>(
    namespace: ValidateNamespace<NS>,
    specifier: ValidateSpecifier<SP>
  ): PermissionObjectBuilder;

  onAll<NS extends string>(
    namespace: ValidateNamespace<NS>
  ): PermissionObjectBuilder;

  any(): PermissionFinal;
}

export interface PermissionObjectBuilder {
  on<NS extends string, SP extends string>(
    namespace: ValidateNamespace<NS>,
    specifier: SP
  ): PermissionObjectBuilder;

  onAll<NS extends string>(
    namespace: ValidateNamespace<NS>
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
