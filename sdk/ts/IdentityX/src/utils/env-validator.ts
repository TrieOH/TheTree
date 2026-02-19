import { env } from "../core/env";

export const validateProjectKey = () => {
  if (!env.PROJECT_KEY || env.PROJECT_KEY.trim() === "") {
    throw new Error(
      "[TRIEOH SDK] Project Key is missing. Please set PUBLIC_TRIEOH_AUTH_PROJECT_KEY, NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_KEY or VITE_TRIEOH_AUTH_PROJECT_KEY."
    );
  }
};

export const validateApiKey = () => {
  if (!env.API_KEY || env.API_KEY.trim() === "") {
    throw new Error(
      "[TRIEOH SDK] Private API Key is missing. This operation requires TRIEOH_AUTH_API_KEY to be set in a server-side environment."
    );
  }
};
