const now = new Date();
const inMinutes = (offset) => new Date(now.getTime() + offset * 60 * 1000).toISOString();

export const createCoffeeBreak = {
    name: "Coffee Break",
    starts_at: inMinutes(10),
    ends_at: inMinutes(12),
    type: "amenity",
    access_mode: "ticket",
};

export const createCheckInArea = {
    name: "CheckIn",
    starts_at: inMinutes(10),
    ends_at: inMinutes(14),
    type: "zone",
    access_mode: "open",
};