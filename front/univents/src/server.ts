import handler from "@tanstack/react-start/server-entry";
import { handleStorageModerate, handleStorageUpload } from "./features/storage/api/storage-handlers";

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);

    if (url.pathname === "/storage/upload" && request.method === "POST") {
      return handleStorageUpload(request, env);
    }

    if (url.pathname === "/storage/moderate" && request.method === "POST") {
      return handleStorageModerate(request, env);
    }

    return handler.fetch(request);
  },
};
