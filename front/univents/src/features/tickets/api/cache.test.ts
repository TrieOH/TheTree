import { QueryClient } from "@tanstack/react-query";
import { describe, expect, it } from "vitest";
import type { TicketI } from "../model";
import { syncTicketCaches } from "./cache";
import { ticketKeys } from "./query-keys";

const ticket = (overrides: Partial<TicketI> = {}): TicketI => ({
  id: "ticket-1",
  edition_id: "edition-1",
  name: "Ticket",
  access_level: 1,
  price_cents: 1000,
  created_at: "2026-01-01T00:00:00Z",
  ...overrides,
});

describe("ticket cache synchronization", () => {
  it("does not create partial caches when data has not been loaded", () => {
    const queryClient = new QueryClient();

    syncTicketCaches(queryClient, ticket());

    expect(
      queryClient.getQueryData(ticketKeys.listByEdition("edition-1")),
    ).toBeUndefined();
    expect(
      queryClient.getQueryData(ticketKeys.detail.byId("ticket-1")),
    ).toBeUndefined();
  });

  it("updates a loaded list and detail", () => {
    const queryClient = new QueryClient();
    const previous = ticket({ name: "Old ticket" });
    const updated = ticket({ name: "Updated ticket" });

    queryClient.setQueryData(ticketKeys.listByEdition("edition-1"), [previous]);
    queryClient.setQueryData(ticketKeys.detail.byId("ticket-1"), previous);

    syncTicketCaches(queryClient, updated);

    expect(
      queryClient.getQueryData(ticketKeys.listByEdition("edition-1")),
    ).toEqual([updated]);
    expect(
      queryClient.getQueryData(ticketKeys.detail.byId("ticket-1")),
    ).toEqual(updated);
  });
});
