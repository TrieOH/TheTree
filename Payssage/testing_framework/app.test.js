import { describe, test, beforeAll, expect } from "vitest";
import {get, loginAs, post, validate} from "./helpers.js";
import WebSocket from "ws"

let user;
beforeAll(async () => {
    user = await loginAs(process.env.OWNER_EMAIL, process.env.OWNER_PASSWORD);
});
