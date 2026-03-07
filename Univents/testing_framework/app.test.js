import {beforeAll, describe, expect, test} from "vitest";
import {get, loginAs, post, validate} from "./helpers.js";
import {createEvent} from "./fixtures/events/create.js";
import {EventSchema} from "./schemas/event.js";
import {createEdition} from "./fixtures/editions/create.js";
import {EditionSchema} from "./schemas/edition.js";
import {createKubernetesActivity, createPremiumWorkshop, createRustActivity} from "./fixtures/activities/create.js";
import {ActivitySchema} from "./schemas/activity.js";
import {createCheckInArea, createCoffeeBreak} from "./fixtures/checkpoints/create.js";
import {CheckpointSchema} from "./schemas/checkpoint.js";
import {createMug, createShirt, createTicketProduct} from "./fixtures/products/create.js";
import {ProductSchema} from "./schemas/product.js";
import {createFullAccessTicket, createStandardTicket, createVIPTicket} from "./fixtures/tickets/create.js";
import {TicketSchema} from "./schemas/ticket.js";
import {addActivityPermission, addCheckpointPermission} from "./fixtures/ticket_permissions/add.js";
import {TicketPermissionSchema} from "./schemas/ticket_permission.js";
import WebSocket from "ws"

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
        edition = await post(owner, `/events/${event.id}/editions`, createEdition)
        expect(validate(EditionSchema, edition)).toBe(true)
        expect(edition.event_id).toBe(event.id)
    })
    test("announce edition", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/announce`)
    })
})

let rustActivity
let kubernetesActivity
let premiumWorkshop
describe("activities", () => {
    test("create rust activity", async () => {
        rustActivity = await post(owner, `/events/${event.id}/editions/${edition.id}/activities`, createRustActivity)
        expect(validate(ActivitySchema, rustActivity)).toBe(true)
        expect(rustActivity.edition_id).toBe(edition.id)
    })
    test("create kubernetes activity", async () => {
        kubernetesActivity = await post(owner, `/events/${event.id}/editions/${edition.id}/activities`, createKubernetesActivity)
        expect(validate(ActivitySchema, kubernetesActivity)).toBe(true)
        expect(kubernetesActivity.edition_id).toBe(edition.id)
    })
    test("create premium workshop", async () => {
        premiumWorkshop = await post(owner, `/events/${event.id}/editions/${edition.id}/activities`, createPremiumWorkshop)
        expect(validate(ActivitySchema, premiumWorkshop)).toBe(true)
        expect(premiumWorkshop.edition_id).toBe(edition.id)
    })
    test("publish rust activity", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${rustActivity.id}/publish`)
    })
    test("publish kubernetes activity", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${kubernetesActivity.id}/publish`)
    })
    test("publish premium workshop", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${premiumWorkshop.id}/publish`)
    })
})

let coffeeBreak
let checkIn
describe("checkpoints", () => {
    test("create coffee break", async () => {
        coffeeBreak = await post(owner, `/events/${event.id}/editions/${edition.id}/checkpoints`, createCoffeeBreak)
        expect(validate(CheckpointSchema, coffeeBreak)).toBe(true)
        expect(coffeeBreak.edition_id).toBe(edition.id)
    })
    test("create check in", async () => {
        checkIn = await post(owner, `/events/${event.id}/editions/${edition.id}/checkpoints`, createCheckInArea)
        expect(validate(CheckpointSchema, checkIn)).toBe(true)
        expect(checkIn.edition_id).toBe(edition.id)
    })
})

let standardTicket
let vipTicket
let fullTicket
describe("tickets", () => {
    test("create standard ticket", async () => {
        standardTicket = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets`, createStandardTicket)
        expect(validate(TicketSchema, standardTicket)).toBe(true)
        expect(standardTicket.edition_id).toBe(edition.id)
    })
    test("create vip ticket", async () => {
        vipTicket = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets`, createVIPTicket)
        expect(validate(TicketSchema, vipTicket)).toBe(true)
        expect(vipTicket.edition_id).toBe(edition.id)
    })
    test("create full ticket", async () => {
        fullTicket = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets`, createFullAccessTicket)
        expect(validate(TicketSchema, fullTicket)).toBe(true)
        expect(fullTicket.edition_id).toBe(edition.id)
    })
})

describe('ticket permissions', () => {
    describe('standard ticket', () => {
        test("add rust activity permission", async () => {
            let toAdd = addActivityPermission(rustActivity.id)
            const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${standardTicket.id}/permissions`, toAdd)
            expect(validate(TicketPermissionSchema, resData)).toBe(true)
            expect(resData.ticket_id).toBe(standardTicket.id)
            expect(resData.activity_id).toBe(rustActivity.id)
        })
        test("add check in checkpoint permission", async () => {
            let toAdd = addCheckpointPermission(checkIn.id)
            const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${standardTicket.id}/permissions`, toAdd)
            expect(validate(TicketPermissionSchema, resData)).toBe(true)
            expect(resData.ticket_id).toBe(standardTicket.id)
            expect(resData.checkpoint_id).toBe(checkIn.id)
        })
    });
    describe('vip ticket', () => {
        test("add rust activity permission", async () => {
            let toAdd = addActivityPermission(rustActivity.id)
            const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${vipTicket.id}/permissions`, toAdd)
            expect(validate(TicketPermissionSchema, resData)).toBe(true)
            expect(resData.ticket_id).toBe(vipTicket.id)
            expect(resData.activity_id).toBe(rustActivity.id)
        })
        test("add kubernetes activity permission", async () => {
            let toAdd = addActivityPermission(kubernetesActivity.id)
            const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${vipTicket.id}/permissions`, toAdd)
            expect(validate(TicketPermissionSchema, resData)).toBe(true)
            expect(resData.ticket_id).toBe(vipTicket.id)
            expect(resData.activity_id).toBe(kubernetesActivity.id)
        })
        test("add check in checkpoint permission", async () => {
            let toAdd = addCheckpointPermission(checkIn.id)
            const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${vipTicket.id}/permissions`, toAdd)
            expect(validate(TicketPermissionSchema, resData)).toBe(true)
            expect(resData.ticket_id).toBe(vipTicket.id)
            expect(resData.checkpoint_id).toBe(checkIn.id)
        })
        describe('full access ticket', () => {
            test("add rust activity permission", async () => {
                let toAdd = addActivityPermission(rustActivity.id)
                const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${fullTicket.id}/permissions`, toAdd)
                expect(validate(TicketPermissionSchema, resData)).toBe(true)
                expect(resData.ticket_id).toBe(fullTicket.id)
                expect(resData.activity_id).toBe(rustActivity.id)
            })
            test("add kubernetes activity permission", async () => {
                let toAdd = addActivityPermission(kubernetesActivity.id)
                const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${fullTicket.id}/permissions`, toAdd)
                expect(validate(TicketPermissionSchema, resData)).toBe(true)
                expect(resData.ticket_id).toBe(fullTicket.id)
                expect(resData.activity_id).toBe(kubernetesActivity.id)
            })
            test("add premium workshop permission", async () => {
                let toAdd = addActivityPermission(premiumWorkshop.id)
                const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${fullTicket.id}/permissions`, toAdd)
                expect(validate(TicketPermissionSchema, resData)).toBe(true)
                expect(resData.ticket_id).toBe(fullTicket.id)
                expect(resData.activity_id).toBe(premiumWorkshop.id)
            })
            test("add check in checkpoint permission", async () => {
                let toAdd = addCheckpointPermission(checkIn.id)
                const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${fullTicket.id}/permissions`, toAdd)
                expect(validate(TicketPermissionSchema, resData)).toBe(true)
                expect(resData.ticket_id).toBe(fullTicket.id)
                expect(resData.checkpoint_id).toBe(checkIn.id)
            })
            test("add coffee break checkpoint permission", async () => {
                let toAdd = addCheckpointPermission(coffeeBreak.id)
                const resData = await post(owner, `/events/${event.id}/editions/${edition.id}/tickets/${fullTicket.id}/permissions`, toAdd)
                expect(validate(TicketPermissionSchema, resData)).toBe(true)
                expect(resData.ticket_id).toBe(fullTicket.id)
                expect(resData.checkpoint_id).toBe(coffeeBreak.id)
            })
        });
    });
});

let mug
let shirt
let vipTicketProduct
let fullTicketProduct
let standardTicketProduct
describe('products', () => {
    test("create mug", async () => {
        mug = await post(owner, `/events/${event.id}/editions/${edition.id}/products`, createMug)
        expect(validate(ProductSchema, mug)).toBe(true)
        expect(mug.edition_id).toBe(edition.id)
    })
    test("create shirt", async () => {
        shirt = await post(owner, `/events/${event.id}/editions/${edition.id}/products`, createShirt)
        expect(validate(ProductSchema, shirt)).toBe(true)
        expect(shirt.edition_id).toBe(edition.id)
    })
    describe('create ticket products', () => {
        test("create standard ticket", async () => {
            let toCreate = createTicketProduct(
                standardTicket.name,
                standardTicket.description,
                standardTicket.id,
                0,
                0,
                15,
                false,
                0,
            )
            const ticketProduct = await post(owner, `/events/${event.id}/editions/${edition.id}/products`, toCreate)
            standardTicketProduct = ticketProduct
            expect(validate(ProductSchema, ticketProduct)).toBe(true)
            expect(ticketProduct.edition_id).toBe(edition.id)
            expect(ticketProduct.ticket_id).toBe(standardTicket.id)
        })
        test("create vip ticket", async () => {
            let toCreate = createTicketProduct(
                vipTicket.name,
                vipTicket.description,
                vipTicket.id,
                1000,
                0,
                15,
                true,
                100,
            )
            const ticketProduct = await post(owner, `/events/${event.id}/editions/${edition.id}/products`, toCreate)
            vipTicketProduct = ticketProduct
            expect(validate(ProductSchema, ticketProduct)).toBe(true)
            expect(ticketProduct.edition_id).toBe(edition.id)
            expect(ticketProduct.ticket_id).toBe(vipTicket.id)
        })
        test("create full access ticket", async () => {
            let toCreate = createTicketProduct(
                fullTicket.name,
                fullTicket.description,
                fullTicket.id,
                2500,
                0,
                15,
                true,
                10,
            )
            const ticketProduct = await post(owner, `/events/${event.id}/editions/${edition.id}/products`, toCreate)
            fullTicketProduct = ticketProduct
            expect(validate(ProductSchema, ticketProduct)).toBe(true)
            expect(ticketProduct.edition_id).toBe(edition.id)
            expect(ticketProduct.ticket_id).toBe(fullTicket.id)
        })
    })
});

describe('purchase', () => {
    test("buy 2x VIP, 1x Full Access, 1x shirt, 2x standard", async () => {
        const cookies = await owner.defaults.jar.getCookies(process.env.BASE_URL)
        const cookieHeader = cookies.map(c => `${c.key}=${c.value}`).join('; ')

        const wsURL = process.env.BASE_URL
            .replace('http://', 'ws://')
            .replace('https://', 'wss://')

        const ws = new WebSocket(
            `${wsURL}/events/${event.id}/editions/${edition.id}/products/purchase`,
            { headers: { Cookie: cookieHeader } }
        )

        await new Promise((resolve, reject) => {
            let sessionID = null

            ws.on("open", () => {
                ws.send(JSON.stringify({
                    items: [
                        { product_id: vipTicketProduct.id, quantity: 2 },
                        { product_id: fullTicketProduct.id, quantity: 1 },
                        { product_id: shirt.id, quantity: 1 },
                        { product_id: standardTicketProduct.id, quantity: 2 },
                    ]
                }))
            })

            ws.on("message", async (data) => {
                const msg = JSON.parse(data)

                if (msg.type === "reservation_failed") return reject(new Error("reservation failed"))
                if (msg.type === "error") return reject(new Error(msg.payload))

                if (msg.type === "reservation_confirmed") {
                    sessionID = msg.payload.session_id
                    await post(owner, `/events/${event.id}/editions/${edition.id}/products/purchase/confirm`, {
                        session_id: sessionID,
                        payment_intent_id: msg.payload.payment_intent_id,
                    })
                }

                if (msg.type === "order_confirmed") {
                    resolve()
                }
            })

            ws.on("error", reject)
            ws.on("close", () => {
                if (!sessionID) reject(new Error("ws closed before reservation_confirmed"))
            })
        })

        ws.close()
    }, 10000)
})

describe('activity registration', () => {
    beforeAll(async () => {
        // wait for asynq grant permissions task to complete
        await new Promise(resolve => setTimeout(resolve, 2000))
    })

    test("register to rust activity", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${rustActivity.id}/register`)
    })

    test("register to kubernetes activity", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${kubernetesActivity.id}/register`)
    })

    test("register to premium workshop", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${premiumWorkshop.id}/register`)
    })

    test("unregister from kubernetes activity", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${kubernetesActivity.id}/unregister`)
    })

    test("re-register to kubernetes activity after unregister", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${kubernetesActivity.id}/register`)
    })
})

let attendanceRecord
describe('mark attendance', () => {
    test("get rust activity attendance record", async () => {
        const records = await get(owner, `/events/${event.id}/editions/${edition.id}/activities/${rustActivity.id}/records`)
        attendanceRecord = records[0]
        expect(attendanceRecord).toBeDefined()
        expect(attendanceRecord.status).toBe("registered")
    })

    test("mark rust activity attendance as completed", async () => {
        await post(owner, `/events/${event.id}/editions/${edition.id}/activities/${rustActivity.id}/records/${attendanceRecord.id}`)
        const records = await get(owner, `/events/${event.id}/editions/${edition.id}/activities/${rustActivity.id}/records`)
        const updated = records.find(r => r.id === attendanceRecord.id)
        expect(updated.status).toBe("completed")
    })
})