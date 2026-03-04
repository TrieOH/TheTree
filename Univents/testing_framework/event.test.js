import { describe, test, beforeAll, expect } from "vitest";
import {loginAs, fixture, schema, post, validate} from "./helpers.js";

const validateEvent = schema("event");

let owner;
beforeAll(async () => {
    owner = await loginAs(process.env.OWNER_EMAIL, process.env.OWNER_PASSWORD);
});

describe("events", () => {
    test("create event", async () => {
        const event = await post(owner, "/events", fixture("events/create"));
        expect(validate(validateEvent, event)).toBe(true);
        expect(event.name).toBe(fixture("events/create").name);
    });
});