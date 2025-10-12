import "dotenv/config";

export const env = {
  API_KEY: process.env.TRIEOH_AUTH_API_KEY || "",
  BASE_URL: "https://api.default.com", // I need to change
};
