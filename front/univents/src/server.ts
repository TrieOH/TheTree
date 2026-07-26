import handler from "@tanstack/react-start/server-entry";
import {
  handleStorageImagePreprocess,
  handleStorageUpload,
} from "./features/storage/api/storage-handlers";

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);

    if (
      url.pathname === "/storage/image/preprocess" &&
      request.method === "POST"
    ) {
      return handleStorageImagePreprocess(request, env);
    }

    if (url.pathname === "/storage/upload" && request.method === "POST") {
      return handleStorageUpload(request, env);
    }

    return handler.fetch(request);
  },
};
