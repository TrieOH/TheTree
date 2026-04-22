/* eslint-disable @typescript-eslint/only-throw-error */
import { redirect } from '@tanstack/react-router';
import type { AnySchema, ParsedLocation } from '@tanstack/react-router';
import type { useAuth } from '@soramux/identityx-sdk-ts/react';
// import type { BuilderMethods } from '@soramux/identityx-sdk-ts';
// import { checkAdminPermissionFn } from '@/features/events/api';

// type ReadyPermission = Pick<BuilderMethods<'object' | 'project' | 'action'>, 'user'>

interface BeforeLoadArgs {
  location: ParsedLocation<AnySchema>;
  context: { auth?: ReturnType<typeof useAuth> }
}

export const requireAuth = ({ context, location }: BeforeLoadArgs) => {
  if (context.auth?.isAuthenticated === false) {
    throw redirect({
      to: '/auth',
      search: { redirect: location.pathname, }
    });
  }
};

export const requireGuest = ({ context }: BeforeLoadArgs) => {
  if (context.auth?.isAuthenticated === true) throw redirect({ to: '/' });
};

// Maybe i will handle this in the future
// /**
//  * Creates a route guard that checks for specific user permissions.
//  * @param permissionsToCheck An array of permission objects to check against.
//  * @param logic Specifies how to combine the results: 'all' requires all permissions to be true (AND), 'any' requires at least one to be true (OR). Defaults to 'all'.
//  * @param redirectTo The path to redirect to if the user does not meet the permission criteria. Defaults to '/'.
//  * @returns A route guard function suitable for use in beforeLoad.
//  */
// export const createPermissionGuard = (
//   permissionsToCheck: ReadyPermission[],
//   logic: 'all' | 'any' = 'all', // Default to 'all' (AND logic)
//   redirectTo: string = '/'
// ) => {
//   return async ({ context, location }: BeforeLoadArgs) => {
//     requireAuth({ context, location });
//     const profile = context.auth?.auth.profile()
//     console.log(context.auth?.auth)

//     if (!profile) {
//       console.error("Auth context or hasPermission method is not available.");
//       throw redirect({ to: redirectTo });
//     }

//     let allPermissionsGranted: boolean;

//     if (logic === 'all') {
//       // Check if ALL permissions are granted
//       const results = await Promise.all(
//         permissionsToCheck.map(perm => checkAdminPermissionFn({ data: perm.user(profile.id).build() }))
//         // permissionsToCheck.map(permission => context.auth!.hasPermission(permission))
//       );
//       allPermissionsGranted = results.every(result => result.success && result.data.allowed);
//     } else {
//       // Check if ANY permission is granted
//       const results = await Promise.all(
//         permissionsToCheck.map(perm => checkAdminPermissionFn({ data: perm.user(profile.id).build() }))
//         // permissionsToCheck.map(permission => context.auth!.hasPermission(permission))
//       );
//       allPermissionsGranted = results.some(result => result.success && result.data.allowed);
//     }

//     if (!allPermissionsGranted) throw redirect({ to: redirectTo });
//   };
// };
