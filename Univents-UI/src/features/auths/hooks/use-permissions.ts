import { useQueries } from "@tanstack/react-query";
import { useMemo } from "react";
import type { Permission } from "@soramux/node-perm-sdk";
import { checkAdminPermissionFn } from "@/features/events/api";

type ReadyPermission = Pick<Permission.BuilderMethods<'resource' | 'permission'>, 'subject'>
type PermissionMap<TKey extends string> = Record<TKey, ReadyPermission>

type PermissionResult<TKey extends string> = Record<TKey, boolean> & {
  some: (...keys: TKey[]) => boolean;
  every: (...keys: TKey[]) => boolean;
  isLoading: boolean;
  isFetching: boolean;
}

export function usePermissions<TKey extends string>(
  permissions: PermissionMap<TKey>,
  userId?: string
): PermissionResult<TKey> {
  const entries = useMemo(
    () => {
      if (!userId) return [];
      return (Object.entries<ReadyPermission>(permissions))
        .map(([key, permission]) => ({ key, built: permission.subject("user", userId).build() }))
    },
    [userId]
  );

  const queries = useQueries({
    queries: entries.map(({ built }) => ({
      queryKey: ['permission', userId, built.subject, built.resource] as const,
      queryFn: async (): Promise<boolean> => {
        const result = await checkAdminPermissionFn({ data: built });
        return result.success ? result.data.permissionship === "PERMISSIONSHIP_HAS_PERMISSION" : false;
      },
      placeholderData: false
    })),
  });

  const falseMap = useMemo(
    () => Object.fromEntries(Object.keys(permissions).map(k => [k, false])) as Record<TKey, boolean>,
    []
  );

  const data = useMemo(() => {
    if (!userId || entries.length === 0) return falseMap;

    return Object.fromEntries(
      entries.map(({ key }, i) => [key, queries[i].data])
    ) as Record<TKey, boolean>;
  }, [userId, entries, queries, falseMap]);

  return {
    ...data,
    isLoading: queries.some(q => q.isLoading),
    isFetching: queries.some(q => q.isFetching),
    some: (...keys: TKey[]) => keys.some(k => data[k]),
    every: (...keys: TKey[]) => keys.every(k => data[k]),
  };
}