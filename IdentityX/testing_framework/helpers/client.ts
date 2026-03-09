import axios, { AxiosInstance } from "axios";
import { wrapper } from "axios-cookiejar-support";
import { CookieJar } from "tough-cookie";

// ============================================================================
// AUTH CONTEXT - Holds auth state for a logged-in session
// ============================================================================

export interface AuthContext {
    accessToken: string;
    refreshToken: string;
}

// ============================================================================
// CLIENT - Thin wrapper around axios with optional auth
// ============================================================================

export class Client {
    readonly http: AxiosInstance;
    readonly auth: AuthContext | null;
    private readonly _email: string;
    private readonly _password: string;

    constructor(
        http: AxiosInstance,
        auth: AuthContext | null = null,
        email = "",
        password = ""
    ) {
        this.http = http;
        this.auth = auth;
        this._email = email;
        this._password = password;
    }

    withCredentials(email: string, password: string): Client {
        return new Client(this.http, this.auth, email, password);
    }

    withAuth(auth: AuthContext): Client {
        return new Client(this.http, auth, this._email, this._password);
    }

    // ----------------
    // Auth operations
    // ----------------

    async register(): Promise<Client> {
        await post(this, "/auth/register", {
            email: this._email,
            password: this._password,
        });
        return this;
    }

    async projectRegister(projectId: string): Promise<Client> {
        await post(this, `/projects/${projectId}/register`, {
            email: this._email,
            password: this._password,
        });
        return this;
    }

    async login(): Promise<Client> {
        const auth = await loginWithCredentials(
            this.http,
            this._email,
            this._password
        );
        return this.withAuth(auth);
    }

    async projectLogin(projectId: string): Promise<Client> {
        const auth = await projectLoginWithCredentials(
            this.http,
            projectId,
            this._email,
            this._password
        );
        return this.withAuth(auth);
    }

    async logout(): Promise<void> {
        if (!this.auth) throw new Error("logout called on unauthenticated client");
        try {
            await this.http.post("/auth/logout", {}, {
                headers: authHeaders(this.auth),
            });
        } catch (e: any) {
            throw e;
        }
    }

    async get<T = any>(url: string): Promise<T> {
        return get(this, url);
    }

    async post<T = any>(url: string, body: unknown = {}): Promise<T> {
        return post(this, url, body);
    }

    async patch<T = any>(url: string, body: unknown = {}): Promise<T> {
        return patch(this, url, body);
    }

    async put<T = any>(url: string, body: unknown = {}): Promise<T> {
        return put(this, url, body);
    }

    async del<T = any>(url: string): Promise<T> {
        return del(this, url);
    }

    async rawGet(url: string): Promise<RawResponse> {
        return raw(this, "GET", url);
    }

    async rawPost(url: string, body?: unknown): Promise<RawResponse> {
        return raw(this, "POST", url, body);
    }

    async rawDelete(url: string): Promise<RawResponse> {
        return raw(this, "DELETE", url);
    }

    async refresh(): Promise<Client> {
        if (!this.auth) throw new Error("refresh called on unauthenticated client");
        const auth = await refreshTokens(this.http, this.auth.refreshToken);
        return this.withAuth(auth);
    }
}

// ============================================================================
// FACTORY - Creates a fresh client (no auth)
// ============================================================================

export function createClient(): Client {
    const jar = new CookieJar();
    const http = wrapper(
        axios.create({
            baseURL: process.env.BASE_URL,
            jar,
            withCredentials: true,
        })
    );
    return new Client(http);
}

// ============================================================================
// AUTH HELPERS - Extract cookies from responses
// ============================================================================

function extractAuthCookies(headers: Record<string, any>): AuthContext {
    const setCookie: string[] = [].concat(headers["set-cookie"] ?? []);

    const get = (name: string) => {
        const entry = setCookie.find((c) => c.startsWith(`${name}=`));
        if (!entry) throw new Error(`Missing cookie: ${name}`);
        return entry.split(";")[0].slice(name.length + 1);
    };

    return {
        accessToken: get("access_token"),
        refreshToken: get("refresh_token"),
    };
}

async function loginWithCredentials(
    http: AxiosInstance,
    email: string,
    password: string
): Promise<AuthContext> {
    try {
        const res = await http.post("/auth/login", { email, password });
        return extractAuthCookies(res.headers);
    } catch (e: any) {
        throw e;
    }
}

async function projectLoginWithCredentials(
    http: AxiosInstance,
    projectId: string,
    email: string,
    password: string
): Promise<AuthContext> {
    try {
        const res = await http.post(`/projects/${projectId}/login`, { email, password });
        return extractAuthCookies(res.headers);
    } catch (e: any) {
        throw e;
    }
}

async function refreshTokens(
    http: AxiosInstance,
    refreshToken: string
): Promise<AuthContext> {
    try {
        const res = await http.post(
            "/auth/refresh",
            {},
            { headers: { Cookie: `refresh_token=${refreshToken}` } }
        );
        return extractAuthCookies(res.headers);
    } catch (e: any) {
        throw e;
    }
}

// ============================================================================
// REQUEST HELPERS - post / get / del / patch
// ============================================================================

function authHeaders(auth: AuthContext | null): Record<string, string> {
    if (!auth) return {};
    return {
        Cookie: `access_token=${auth.accessToken}; refresh_token=${auth.refreshToken}`,
    };
}

export async function post<T = any>(
    client: Client,
    url: string,
    body: unknown = {}
): Promise<T> {
    try {
        const res = await client.http.post(url, body, {
            headers: authHeaders(client.auth),
        });
        return res.data.data as T;
    } catch (e: any) {
        throw e;
    }
}

export async function get<T = any>(client: Client, url: string): Promise<T> {
    try {
        const res = await client.http.get(url, {
            headers: authHeaders(client.auth),
        });
        return res.data.data as T;
    } catch (e: any) {
        throw e;
    }
}

export async function patch<T = any>(
    client: Client,
    url: string,
    body: unknown = {}
): Promise<T> {
    try {
        const res = await client.http.patch(url, body, {
            headers: authHeaders(client.auth),
        });
        return res.data.data as T;
    } catch (e: any) {
        throw e;
    }
}

export async function put<T = any>(
    client: Client,
    url: string,
    body: unknown = {}
): Promise<T> {
    try {
        const res = await client.http.put(url, body, {
            headers: authHeaders(client.auth),
        });
        return res.data.data as T;
    } catch (e: any) {
        throw e;
    }
}

export async function del<T = any>(client: Client, url: string): Promise<T> {
    try {
        const res = await client.http.delete(url, {
            headers: authHeaders(client.auth),
        });
        return res.data.data as T;
    } catch (e: any) {
        throw e;
    }
}

// ============================================================================
// RAW REQUEST - For tests that need to assert on status codes / error bodies
// ============================================================================

export interface RawResponse {
    status: number;
    data: any;
}

export async function raw(
    client: Client,
    method: "GET" | "POST" | "PATCH" | "PUT" | "DELETE",
    url: string,
    body?: unknown
): Promise<RawResponse> {
    try {
        const res = await client.http.request({
            method,
            url,
            data: body,
            headers: authHeaders(client.auth),
            validateStatus: () => true, // never throw on HTTP errors
        });
        return { status: res.status, data: res.data };
    } catch (e: any) {
        console.error(`${method} ${url} failed:`, e.message);
        throw e;
    }
}