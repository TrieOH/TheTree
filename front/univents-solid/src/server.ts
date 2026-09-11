import { TRACES_INGEST_PATH } from "@trieoh/front-core/tracing/constants";
import { handleTracesIngest } from "@trieoh/front-core/tracing/ingest";
import { BFF_PATH, handleBffRequest } from "./features/auths/api/bff-handler";
import {
  handleStorageImagePreprocess,
  handleStorageUpload,
} from "./features/storage/api/storage-handlers";


export default {
  // `ctx` is optional so tests and callers that only care about routing can
  // skip it; the runtime always provides it.
  async fetch(request: Request, env: Env, ctx?: ExecutionContext): Promise<Response> {
    const url = new URL(request.url);
    if (request.method === "POST") {
      if (url.pathname === BFF_PATH) {
        return handleBffRequest(request, env, ctx);
      }

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
