import { QueryClient } from "@tanstack/react-query";
import { describe, expect, it } from "vitest";
import type { SignatureI, SignatureRequestI } from "../model";
import {
  invalidateSignatureRequestCache,
  removeSignatureCache,
  syncSignatureCache,
  syncSignatureRequestCache,
} from "./cache";
import { signatureKeys } from "./query-keys";

const signature = (overrides: Partial<SignatureI> = {}): SignatureI => ({
  id: "signature-1",
  edition_id: "edition-1",
  created_by: "actor-1",
  signatory_name: "Signer",
  image_url: "https://example.com/signature.png",
  created_at: "2026-01-01T00:00:00Z",
  ...overrides,
});

const request = (
  overrides: Partial<SignatureRequestI> = {},
): SignatureRequestI => ({
  id: "request-1",
  edition_id: "edition-1",
  created_by: "actor-1",
  signatory_name: "Signer",
  signatory_email: "signer@example.com",
  idempotency_key: "key-1",
  status: "pending",
  expires_at: "2026-02-01T00:00:00Z",
  created_at: "2026-01-01T00:00:00Z",
  ...overrides,
});

describe("signature cache synchronization", () => {
  it("does not create partial signature or request caches", () => {
    const queryClient = new QueryClient();

    syncSignatureCache(queryClient, signature());
    syncSignatureRequestCache(queryClient, request());

    expect(
      queryClient.getQueryData(signatureKeys.byEdition("edition-1")),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(signatureKeys.byId("signature-1")),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(signatureKeys.requestsByEdition("edition-1")),
    ).toBeUndefined();
  });

  it("updates and removes loaded signatures", () => {
    const queryClient = new QueryClient();
    const updated = signature({ signatory_name: "Updated signer" });
    queryClient.setQueryData(signatureKeys.byEdition("edition-1"), [
      signature(),
    ]);
    queryClient.setQueryData(signatureKeys.byId("signature-1"), signature());

    syncSignatureCache(queryClient, updated);

    expect(
      queryClient.getQueryData(signatureKeys.byEdition("edition-1")),
    ).toEqual([updated]);
    expect(queryClient.getQueryData(signatureKeys.byId("signature-1"))).toEqual(
      updated,
    );

    removeSignatureCache(queryClient, "edition-1", "signature-1");

    expect(
      queryClient.getQueryData(signatureKeys.byEdition("edition-1")),
    ).toEqual([]);
    expect(
      queryClient.getQueryData(signatureKeys.byId("signature-1")),
    ).toBeUndefined();
  });

  it("updates loaded requests and invalidates cancellation queries", () => {
    const queryClient = new QueryClient();
    const listKey = signatureKeys.requestsByEdition("edition-1");
    const detailKey = signatureKeys.requestById("request-1");
    queryClient.setQueryData(listKey, [request()]);
    queryClient.setQueryData(detailKey, request());
    const completed = request({ status: "completed" });

    syncSignatureRequestCache(queryClient, completed);

    expect(queryClient.getQueryData(listKey)).toEqual([completed]);
    expect(queryClient.getQueryData(detailKey)).toEqual(completed);

    invalidateSignatureRequestCache(queryClient, "edition-1", "request-1");

    expect(queryClient.getQueryState(listKey)?.isInvalidated).toBe(true);
    expect(queryClient.getQueryState(detailKey)?.isInvalidated).toBe(true);
  });
});
