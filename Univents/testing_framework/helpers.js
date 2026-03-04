import "dotenv/config";
import axios from "axios";
import { wrapper } from "axios-cookiejar-support";
import { CookieJar } from "tough-cookie";
import Ajv from "ajv";
import addFormats from "ajv-formats"; // for uuid format
import { readFileSync } from "fs";
import { resolve, dirname } from "path";
import { fileURLToPath } from "url";
import {config} from "dotenv";

const __dirname = dirname(fileURLToPath(import.meta.url));
config({ path: resolve(__dirname, ".env"), override: true });

console.log("loaded from", resolve(__dirname, ".env"));
console.log("BASE_URL", process.env.BASE_URL);

const ajv = new Ajv();
addFormats(ajv);

export function fixture(path) {
    return JSON.parse(readFileSync(resolve(__dirname, `fixtures/${path}.json`), "utf-8"));
}

export function schema(path) {
    const s = JSON.parse(readFileSync(resolve(__dirname, `schemas/${path}.json`), "utf-8"));
    return ajv.compile(s);
}

export function createClient() {
    return wrapper(axios.create({
        baseURL: process.env.BASE_URL,
        jar: new CookieJar(),
        withCredentials: true,
    }));
}

export async function loginAsGoAuthAccount(email, password) {
    const client = createClient();
    try {
        await client.post(`${process.env.GOAUTH_URL}/auth/login`, { email, password });
    } catch (e) {
        console.log(e.response?.status, JSON.stringify(e.response?.data, null, 2));
        throw e;
    }
    return client;
}

export async function loginAs(email, password) {
    const client = createClient();
    try {
        await client.post(
            `${process.env.GOAUTH_URL}/projects/${process.env.GOAUTH_PROJECT_ID}/login`,
            { email, password }
        );
    } catch (e) {
        console.log(e.response?.status, JSON.stringify(e.response?.data, null, 2));
        throw e;
    }
    return client;
}

export async function post(client, path, body) {
    try {
        const res = await client.post(path, body);
        return res.data.data;
    } catch (e) {
        console.error(`POST ${path} failed:`, e.response?.status, JSON.stringify(e.response?.data, null, 2));
        throw e;
    }
}

export async function get(client, path) {
    try {
        const res = await client.get(path);
        return res.data.data;
    } catch (e) {
        console.error(`GET ${path} failed:`, e.response?.status, JSON.stringify(e.response?.data, null, 2));
        throw e;
    }
}

export function validate(schema, data) {
    const valid = schema(data);
    if (!valid) {
        console.error("schema errors:", JSON.stringify(schema.errors, null, 2));
    }
    return valid;
}