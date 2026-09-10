import { TRACES_INGEST_PATH } from "@trieoh/front-core/tracing/constants";
import { handleTracesIngest } from "@trieoh/front-core/tracing/ingest";
import {
  handleStorageImagePreprocess,
  handleStorageUpload,
} from "./features/storage/api/storage-handlers";


export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);
    if (request.method === "POST") {
      if (url.pathname === "/storage/image/preprocess") {
        return handleStorageImagePreprocess(request, env);
      }

      if (url.pathname === "/storage/upload") {
        return handleStorageUpload(request, env);
      }

      if (url.pathname === TRACES_INGEST_PATH) {
        return handleTracesIngest(request, env);
      }
    }
    return env.ASSETS.fetch(request);
  },
};
