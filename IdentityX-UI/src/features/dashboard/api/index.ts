// import { authFetcher } from "@/shared/lib/api/fetch";
// import { createClientOnlyFn } from "@tanstack/react-start";
// import {  } from "@trieoh/node-auth-sdk/react";

// // /projects/${env.API_KEY}/api-keys/rotate
// export const rotateApiKey = createClientOnlyFn((project_id: string) => {

//   return authFetcher<{api_key: string}>(`/projects/${project_id}/api-keys/rotate`, {
//     method: "POST",
//     headers: { "Content-Type": "application/json" }, // it's already used in the lib per default
//   });
// });