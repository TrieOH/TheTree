import handler from "@tanstack/react-start/server-entry";
import { TRACES_INGEST_PATH } from "@trieoh/front-core/tracing/constants";
import { handleTracesIngest } from "@trieoh/front-core/tracing/ingest";
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

    if (url.pathname === TRACES_INGEST_PATH && request.method === "POST") {
      return handleTracesIngest(request, env);
    }

    return handler.fetch(request);
  },
};
