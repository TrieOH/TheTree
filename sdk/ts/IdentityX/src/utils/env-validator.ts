import { env } from "../core/env";

export const validateProjectKey = () => {
  if (!env.PROJECT_ID || env.PROJECT_ID.trim() === "") {
    throw new Error(
      "[TRIEOH SDK] Project ID is missing. Please set PUBLIC_TRIEOH_AUTH_PROJECT_ID, NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_ID or VITE_TRIEOH_AUTH_PROJECT_ID."
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
