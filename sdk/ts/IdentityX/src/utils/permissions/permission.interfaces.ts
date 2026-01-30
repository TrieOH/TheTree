import type { PermissionResult } from "./permission.result";
import type { ValidateAction, ValidateNamespace, ValidateSpecifier } from "./permission.validators";

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
  can<Token extends string>(action: ValidateAction<Token>): PermissionActionChain;
  canAnyAction(): PermissionFinal;
}

export interface PermissionActionChain {
  and<Token extends string>(action: ValidateAction<Token>): PermissionActionChain;
  build(): PermissionResult;
}

export interface PermissionFinal {
  build(): PermissionResult;
}
