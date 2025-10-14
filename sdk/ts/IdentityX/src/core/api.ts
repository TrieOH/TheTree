import { env } from "./env";

export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  module: string;
  timestamp: string;
  trace?: string[];
  data?: T;
}

export class Api {
  private baseURL: string;
  private apiKey: string;

  constructor(baseURL?: string) {
    this.baseURL = baseURL || env.BASE_URL;
    this.apiKey = env.API_KEY;
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

  async request<T = unknown>(
    path: string,
    options?: RequestInit
  ): Promise<ApiResponse<T>> {
    try {
      const res = await fetch(this.buildUrl(path), {
        ...options,
        headers: { ...this.headers, ...(options?.headers ?? {}) },
      });

      const data = await res.json();
      return data as ApiResponse<T>;
    } catch (error) {
      return {
        code: 503,
        message: "Network request failed — API may be offline.",
        module: "network",
        timestamp: new Date().toISOString(),
        trace: [(error as Error).message || "Unknown network error"],
      };
    }
  }

  get<T = unknown>(path: string) {
    return this.request<T>(path, { method: "GET" });
  }

  post<T = unknown>(path: string, body?: unknown) {
    return this.request<T>(path, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  put<T = unknown>(path: string, body?: unknown) {
    return this.request<T>(path, {
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  delete<T = unknown>(path: string) {
    return this.request<T>(path, { method: "DELETE" });
  }
}
