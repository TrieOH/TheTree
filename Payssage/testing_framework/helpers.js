import { config } from "dotenv";
import { resolve, dirname } from "path";
import { fileURLToPath } from "url";
import axios from "axios";
import { wrapper } from "axios-cookiejar-support";
import { CookieJar } from "tough-cookie";

const __dirname = dirname(fileURLToPath(import.meta.url));
config({ path: resolve(__dirname, ".env"), override: true });

export function createClient() {
    return wrapper(axios.create({
        baseURL: process.env.BASE_URL,
        jar: new CookieJar(),
        withCredentials: true,
    }));
}

export function createAPIKeyClient(apiKey) {
    return wrapper(axios.create({
        baseURL: process.env.BASE_URL,
        jar: new CookieJar(),
        withCredentials: true,
        headers: {
            "X-API-Key": apiKey,
        },
    }))
}

export async function loginAs(email, password) {
    const client = createClient();
    try {
        await client.post(
            `${process.env.GOAUTH_URL}/projects/${process.env.GOAUTH_PROJECT_ID}/login`,
            { email, password }
        );
    } catch (e) {
        console.error(`LOGIN failed:`, e.response?.status, JSON.stringify(e.response?.data, null, 2));
        throw e;
    }
    return client;
}

export async function post(client, url, body) {
    try {
        if (client === null) {
            const res = await axios.post(`${process.env.BASE_URL}${url}`, body)
            return res.data.data
        }
        const res = await client.post(url, body ?? {})
        return res.data.data
    } catch (e) {
        console.error(`POST ${url} failed:`, e.response?.status, JSON.stringify(e.response?.data, null, 2))
        throw e
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

export async function deleteReq(client, path) {
    try {
        const res = await client.delete(path)
        return res.data.data
    } catch (e) {
        console.error(`DELETE ${path} failed:`, e.response?.status, JSON.stringify(e.response?.data, null, 2))
        throw e
    }
}

export function validate(schema, data) {
    const result = schema.safeParse(data);
    if (!result.success) {
        console.error("schema errors:", JSON.stringify(result.error.format(), null, 2));
        return false;
    }
    return true;
}