import { describe, test, beforeAll, expect } from "vitest";
import { loginAs, post, validate } from "./helpers.js";
import { createEvent } from "./fixtures/events/create.js";
import { EventSchema } from "./schemas/event.js";

let owner;
beforeAll(async () => {
    owner = await loginAs(process.env.OWNER_EMAIL, process.env.OWNER_PASSWORD);
});

describe("events", () => {
    test("create event", async () => {
        const event = await post(owner, "/events", createEvent);
        expect(validate(EventSchema, event)).toBe(true);
        expect(event.name).toBe(createEvent.name);
    });
});