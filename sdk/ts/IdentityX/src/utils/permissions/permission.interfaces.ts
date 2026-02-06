import type { PermissionResult } from "./permission.result";
import type { ValidateAction, ValidateNamespace, ValidateSpecifier } from "./permission.validators";

/**
 * The entry point of the Permission DSL.
 * Defines the initial scope or a global wildcard(*).
 */
export interface PermissionRoot {
  /**
   * Targets a specific resource (e.g., a table object) by namespace and unique identifier.
   * @param {ValidateNamespace<NS>} namespace - The namespace of the resource.
   * @param {ValidateSpecifier<SP>} specifier - The unique identifier/specifier for the resource.
   * @returns {PermissionObjectBuilder} A builder for configuring permissions.
   * @example
   * .on('event', '123') // Resulting selector: "event:123/"
   */
  on<NS extends string, SP extends string>(
    namespace: ValidateNamespace<NS>,
    specifier: ValidateSpecifier<SP>
  ): PermissionObjectBuilder;


  /**
   * Targets all resources (e.g., table objects) within a given namespace.
   * @param {ValidateNamespace<NS>} namespace - The namespace to target.
   * @returns {PermissionObjectBuilder} A builder for configuring permissions.
   * @example
   * .onAll('project') // Resulting selector: "project:*"
   */
  onAll<NS extends string>(
    namespace: ValidateNamespace<NS>
  ): PermissionObjectBuilder;

  /**
   * Grants access to all resources across all namespaces.
   * @returns {PermissionFinal} The final permission object.
   * @example
   * .any() // Resulting selector: "*"
   */
  any(): PermissionFinal;
}

/**
 * Builder for the resource path (Objects).
 * Handles the hierarchy of "where" the permission applies.
 */
export interface PermissionObjectBuilder {
  /**
   * Adds a nested resource level to the current path.
  * @param {ValidateNamespace<NS>} namespace - The sub-namespace.
   * @param {ValidateNamespace<SP>} specifier - The sub-specifier identifier.
   * @example
   * // Resulting selector: "org:1/project:A/"
   * .on('org', '1').on('project', 'A')
   */
  on<NS extends string, SP extends string>(
    namespace: ValidateNamespace<NS>,
    specifier: ValidateSpecifier<SP>
  ): PermissionObjectBuilder;

  /**
   * Adds a wildcard(*) namespace level to the current path.
   * @param {ValidateNamespace<NS>} namespace - The sub-namespace to target. 
   * @returns {PermissionObjectBuilder} The builder with the appended wildcard(*) path.
   * @example
   * // Resulting selector: "org:1/project:*"
   * .on('org', '1').onAll('project')
   */
  onAll<NS extends string>(
    namespace: ValidateNamespace<NS>
  ): PermissionObjectBuilder;

  /**
   * Targets any immediate child under the current resource path.
   * @returns {PermissionActionBuilder} A builder to define actions for children.
   * @example
   * .on('folder', 'work').forAnyChild() // Path: "folder:work/*"
   */
  forAnyChild(): PermissionActionBuilder;

  /**
   * Targets all descendants (recursive) under the current resource path.
   * @returns {PermissionActionBuilder} A builder to define actions for all descendants.
   * @example
   * .on('folder', 'work').forAnyDescendant() // Path: "folder:work/**"
   */
  forAnyDescendant(): PermissionActionBuilder;

  // ------------- Start of Actions -------------
  /**
   * Defines the domain or category for the action.
   * @param {ValidateAction<Token>} action - The action domain/context (e.g., 'settings').
   * @returns {PermissionActionChain} A chain to refine stages or finalize with a verb.
   * @example
   * .on('file', '123').in('metadata') // Action Domain: "metadata"
   */
  in<Token extends string>(action: ValidateAction<Token>): PermissionActionChain;
  
  /**
   * Resolves the permission by granting any action (*) on the current resource.
   * @returns {PermissionFinal} The final permission object.
   * @example
   * .on('file', '123').inAny() // Action: "*"
   */
  inAny(): PermissionFinal;

  /**
   * Resolves the permission by granting a specific action/verb.
   * @param {ValidateAction<Token>} action - The specific verb to grant (e.g., 'edit').
   * @returns {PermissionFinal} The final permission object.
   * @example
   * .on('file', '123').can('read') // Action: "read"
   */
  can<Token extends string>(action: ValidateAction<Token>): PermissionFinal;
}

/**
 * Handles action verbs and domains for resources targeted with wildcards(*).
 */
export interface PermissionActionBuilder {
  /**
   * Defines the action domain for the targeted wildcard(*) path.
   * @param {ValidateAction<Token>} action - The action domain.
   * @returns {PermissionActionChain} A chain to refine or finalize the action.
   * @example
   * .forAnyChild().in('settings') // Action Domain: "settings"
   */
  in<Token extends string>(action: ValidateAction<Token>): PermissionActionChain;

  /**
   * Resolves the permission by allowing any action on the targeted scope.
   * @returns {PermissionFinal} The final permission object.
   * @example
   * .forAnyChild().inAny() // Action: "*"
   */
  inAny(): PermissionFinal;

  /**
   * Resolves the permission by granting a specific action/verb on the targeted scope.
   * @param {ValidateAction<Token>} action - The action/verb to grant.
   * @returns {PermissionFinal} The final permission object.
   * @example
   * .forAnyChild().can('view') // Action: View
   */
  can<Token extends string>(action: ValidateAction<Token>): PermissionFinal;
}

/**
 * Handles the definition of multi-stage action domains or action wildcards(*).
 */
export interface PermissionActionChain {
  /**
   * Appends a sub-domain or refined stage to the current action context.
   * @param {ValidateAction<Token>} action - The sub-action domain.
   * @returns {PermissionActionChain} The chain with the appended domain.
   * @example
   * .in('order').in('payment') // Action Domain: "order:payment"
   */
  in<Token extends string>(action: ValidateAction<Token>): PermissionActionChain;

  /**
   * Adds a wildcard(*) to the current action domain stage.
   * @returns {PermissionActionChain} The chain with a wildcar(*) domain.
   * @example
   * .in('order').inAny() // Action: "order:*"
   */
  inAny(): PermissionActionChain;

  /**
   * Resolves the action chain with a specific final verb.
   * @param {ValidateAction<Token>} action - The final verb to grant.
   * @returns {PermissionFinal} The final permission object.
   * @example
   * .in('order').in('payment').can('approve') // Action: "order:payment:approve"
   */
  can<Token extends string>(action: ValidateAction<Token>): PermissionFinal;

  /**
   * Resolves the action by targeting any immediate child stage of the current domain.
   * @returns {PermissionFinal} The final permission object.
   * @example
   * .in('payment').forAnyChild() // Action: "payment:*"
   */
  forAnyChild(): PermissionFinal

  /**
   * Resolves the action by targeting all future sub-stages of the current domain.
   * @returns {PermissionFinal} The final permission object.
   * @example
   * .in('payment').forAnyDescendant() // Action: "payment:**"
   */
  forAnyDescendant(): PermissionFinal;
}

/**
 * Represents the completed state of the permission builder.
 */
export interface PermissionFinal {
  /**
   * Finalizes the construction process and returns the Result object.
   * @returns {PermissionResult} The result containing the serialized permission.
   */
  finish(): PermissionResult;
}
