import { authFetcher, publicFetcher } from "#/shared/lib/api/fetch";
import type { ConnectRequest, RevokeRequest } from "../model";

export const getProviderCallbackUrl = (provider: string) =>
  new URL(
    `/callback/${encodeURIComponent(provider)}`,
    window.location.origin,
  ).toString();

export const connectProviderFn = (provider: string, payload: ConnectRequest) =>
  authFetcher.post<string>(`/providers/${provider}/connect`, payload);

export const revokeProviderFn = (provider: string, payload: RevokeRequest) =>
  authFetcher.post<void>(`/providers/${provider}/revoke`, payload);

export const getProviderCallbackFn = (
  code: string,
  state: string,
  provider: string,
) =>
  publicFetcher.get<string>(
    `/providers/${provider}/callback?code=${encodeURIComponent(code)}&state=${encodeURIComponent(state)}`,
  );
