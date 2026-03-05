const now = new Date();
const inMinutes = (offset) => new Date(now.getTime() + offset * 60 * 1000).toISOString();

export const createShirt = {
    edition_scope_id: "",
    name: "SCTI 2026 T-Shirt",
    description: "Official SCTI 2026 event t-shirt",
    type: "merchandise",
    ticket_id: null,
    price_cents: 5000,
    available_from: null,
    available_until: null,
    has_inventory: true,
    inventory_quantity: 100,
};

export const createMug = {
    edition_scope_id: "",
    name: "SCTI 2026 Mug",
    description: "Official SCTI 2026 event mug",
    type: "merchandise",
    ticket_id: null,
    price_cents: 3500,
    available_from: inMinutes(1),
    available_until: inMinutes(7),
    has_inventory: false,
    inventory_quantity: 0,
};

export const createTicketProduct = (scope, name, description, ticket_id, price_cents, from, to, has_inventory, quantity)=> ({
    edition_scope_id: scope,
    name: name,
    description: description,
    type: "ticket",
    ticket_id: ticket_id,
    price_cents: price_cents,
    available_from: inMinutes(from),
    available_until: inMinutes(to),
    has_inventory: has_inventory,
    inventory_quantity: quantity,
});