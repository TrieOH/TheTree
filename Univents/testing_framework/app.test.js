import { describe, test, beforeAll, expect } from "vitest";
import { loginAs, post, validate } from "./helpers.js";
import { createEvent } from "./fixtures/events/create.js";
import { EventSchema } from "./schemas/event.js";
import {createEdition} from "./fixtures/editions/create.js";
import {EditionSchema} from "./schemas/edition.js";
import {createKubernetesActivity, createRustActivity} from "./fixtures/activities/create.js";
import {ActivitySchema} from "./schemas/activity.js";

let owner;
beforeAll(async () => {
    owner = await loginAs(process.env.OWNER_EMAIL, process.env.OWNER_PASSWORD);
});

let event
describe("events", () => {
    test("create event", async () => {
        event = await post(owner, "/events", createEvent);
        expect(validate(EventSchema, event)).toBe(true);
        expect(event.name).toBe(createEvent.name);
    });
    test("publish event", async () => {
        await post(owner, `/events/${event.id}/publish`)
    })
});

let edition
describe("editions", () => {
    test("create edition", async () => {
        let toCreate = createEdition
        toCreate.go_auth_event_scope_id = event.goauth_scope_id
        edition = await post(owner, `/events/${event.id}/editions`, toCreate)
        expect(validate(EditionSchema, edition)).toBe(true)
        expect(edition.event_id).toBe(event.id)
    })
    test("announce edition", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/announce`)
    })
})

let rustActivity
let kubernetesActivity
describe("activities", () => {
    test("create rust activity", async () => {
        let toCreate = createRustActivity
        toCreate.edition_scope_id = edition.goauth_scope_id
        rustActivity = await post(owner, `/events/${event.id}/editions/${edition.id}/activities`, toCreate)
        expect(validate(ActivitySchema, rustActivity)).toBe(true)
        expect(rustActivity.edition_id).toBe(edition.id)
    })
    test("create kubernetes activity", async () => {
        let toCreate = createKubernetesActivity
        toCreate.edition_scope_id = edition.goauth_scope_id
        kubernetesActivity = await post(owner, `/events/${event.id}/editions/${edition.id}/activities`, toCreate)
        expect(validate(ActivitySchema, kubernetesActivity)).toBe(true)
        expect(kubernetesActivity.edition_id).toBe(edition.id)
    })
    test("publish rust activity", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${rustActivity.id}/publish`)
    })
    test("publish kubernetes activity", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${kubernetesActivity.id}/publish`)
    })
})