import type { PermissionResult } from "./permission.result";
import type { ValidateAction, ValidateNamespace, ValidateSpecifier } from "./permission.validators";

export interface PermissionRoot {
  on<NS extends string, SP extends string>(
    namespace: ValidateNamespace<NS>,
    specifier: ValidateSpecifier<SP>
  ): PermissionObjectBuilder; // namespace:specifier/

  onAll<NS extends string>(
    namespace: ValidateNamespace<NS>
  ): PermissionObjectBuilder; // namespace:*/

  any(): PermissionFinal; // *
}

export interface PermissionObjectBuilder {
  on<NS extends string, SP extends string>(
    namespace: ValidateNamespace<NS>,
    specifier: SP
  ): PermissionObjectBuilder; // namespace:specifier/

  onAll<NS extends string>(
    namespace: ValidateNamespace<NS>
  ): PermissionObjectBuilder; // namespace:*/
  forAnyChild(): PermissionActionBuilder; // ...namespace:specifier.../*
  forAnyDescendant(): PermissionActionBuilder; // ...namespace:specifier.../**

  // Start of Actions
  in<Token extends string>(action: ValidateAction<Token>): PermissionActionChain; // token
  inAny(): PermissionFinal; // *
  can<Token extends string>(action: ValidateAction<Token>): PermissionFinal; // token
}

// Start of Actions
export interface PermissionActionBuilder {
  in<Token extends string>(action: ValidateAction<Token>): PermissionActionChain; // token
  inAny(): PermissionFinal; // *
  can<Token extends string>(action: ValidateAction<Token>): PermissionFinal; // token
}

export interface PermissionActionChain {
  in<Token extends string>(action: ValidateAction<Token>): PermissionActionChain; // token:...:stage
  inAny(): PermissionActionChain; // token:...:*
  can<Token extends string>(action: ValidateAction<Token>): PermissionFinal; // token:...:stage
  forAnyChild(): PermissionFinal // // token:...:*
  forAnyDescendant(): PermissionFinal; // token:...:**
}

export interface PermissionFinal {
  finish(): PermissionResult; // end
}
