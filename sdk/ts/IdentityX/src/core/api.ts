import { env } from "./env";

class Api {
  private baseURL: string;
  private apiKey: string;

  constructor(baseURL?: string) {
    this.baseURL = baseURL || env.BASE_URL;
    this.apiKey = env.API_KEY;
    console.log(this.apiKey);
    if (!this.apiKey) {
      console.warn("[TRIEOH SDK] API_KEY not found, verify your .env file");
      throw new Error("[TRIEOH SDK] API_KEY not found, verify your .env file");
    }
  }

  private get headers() {
    return {
      "Content-Type": "application/json",
      ...(this.apiKey ? { Authorization: `Bearer ${this.apiKey}` } : {}),
    };
  }

  private buildUrl(path: string) {
    return `${this.baseURL.replace(/\/$/, "")}/${path.replace(/^\//, "")}`;
  }

  async request(path: string, options?: RequestInit) {
    const res = await fetch(this.buildUrl(path), {
      ...options,
      headers: { ...this.headers, ...(options?.headers ?? {}) },
    });

    if (!res.ok)
      throw new Error(`Request failed (${res.status}): ${res.statusText}`);
    return res.json();
  }

  get(path: string) {
    return this.request(path, { method: "GET" });
  }

  post(path: string, body?: unknown) {
    return this.request(path, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  put(path: string, body?: unknown) {
    return this.request(path, {
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  delete(path: string) {
    return this.request(path, { method: "DELETE" });
  }
}

export const api = new Api(); // Default instance
export { Api };
