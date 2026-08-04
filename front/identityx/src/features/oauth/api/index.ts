import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  createProjectOAuthProvider,
  deleteOAuthProvider,
  disableOAuthProvider,
  enableOAuthProvider,
  listProjectOAuthProviders,
  updateOAuthProvider,
} from "@trieoh/identityx-api";
import type {
  OAuthProviderCreateI,
  OAuthProviderI,
  OAuthProviderUpdateI,
} from "../model";

export const allOAuthProvidersQueryOptions = (projectId: string) =>
  queryOptions({
    queryKey: ["projects", projectId, "oauth-providers"],
    queryFn: () =>
      listProjectOAuthProviders(projectId).then(orvalData<OAuthProviderI[]>),
  });

export const createOAuthProviderFn = createClientOnlyFn(
  (projectId: string, data: OAuthProviderCreateI) =>
    createProjectOAuthProvider(projectId, data).then(orvalData<OAuthProviderI>),
);

export const updateOAuthProviderFn = createClientOnlyFn(
  (id: string, data: OAuthProviderUpdateI) =>
    updateOAuthProvider(id, data).then(orvalData<OAuthProviderI>),
);

export const deleteOAuthProviderFn = createClientOnlyFn((id: string) =>
  deleteOAuthProvider(id).then(orvalData<void>),
);

export const setOAuthProviderEnabledFn = createClientOnlyFn(
  (id: string, enabled: boolean) =>
    (enabled ? enableOAuthProvider(id) : disableOAuthProvider(id)).then(
      orvalData<OAuthProviderI>,
    ),
);
