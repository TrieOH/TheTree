export interface StorageAdapter {
  getItem(key: string): string | null;
  setItem(key: string, value: string): void;
  removeItem(key: string): void;
}

export const browserStorage: StorageAdapter = {
  getItem: (key) => (typeof window !== "undefined" ? localStorage.getItem(key) : null),
  setItem: (key, value) => {
    if (typeof window !== "undefined") localStorage.setItem(key, value);
  },
  removeItem: (key) => {
    if (typeof window !== "undefined") localStorage.removeItem(key);
  },
};

export const sessionBrowserStorage: StorageAdapter = {
  getItem: (key) => (typeof window !== "undefined" ? sessionStorage.getItem(key) : null),
  setItem: (key, value) => {
    if (typeof window !== "undefined") sessionStorage.setItem(key, value);
  },
  removeItem: (key) => {
    if (typeof window !== "undefined") sessionStorage.removeItem(key);
  },
};

export interface CookieOptions {
  expires?: string;
  path?: string;
  domain?: string | null;
  secure?: boolean;
  sameSite?: "Lax" | "None" | "Strict";
}

export const cookieStorage = {
  get: (name: string): string | null => {
    if (typeof window === "undefined") return null;
    const nameEQ = name + "=";
    const ca = document.cookie.split(";");
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === " ") c = c.substring(1, c.length);
      if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
  },

  set: (name: string, value: string, options: CookieOptions = {}) => {
    if (typeof window === "undefined") return;

    const {
      expires,
      path = "/",
      domain,
      secure = window.location.protocol === "https:",
      sameSite = secure ? "None" : "Lax",
    } = options;

    const cookieParts = [
      `${name}=${value}`,
      domain ? `Domain=${domain}` : "",
      `Path=${path}`,
      `SameSite=${sameSite}`,
      secure ? "Secure" : "",
      expires ? `expires=${expires}` : "",
    ];

    document.cookie = cookieParts.filter(Boolean).join("; ");
  },

  remove: (name: string, domain?: string | null) => {
    cookieStorage.set(name, "", {
      expires: "Thu, 01 Jan 1970 00:00:00 GMT",
      domain,
    });
  },
};
