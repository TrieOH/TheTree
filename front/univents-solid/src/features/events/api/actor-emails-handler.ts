import { privateJsonResponse } from "@/shared/lib/http-cache";

export interface ActorEmailsRequest {
  actorIds: string[];
}

export interface ActorRecord {
  email: string | null;
  pfp_url: string | null;
}

export type ActorEmailsMap = Record<string, ActorRecord>;

// In-memory cache per worker isolate with 5-minute TTL to protect IdentityX rate limits
const ACTOR_CACHE = new Map<string, { record: ActorRecord; expiresAt: number }>();
const CACHE_TTL_MS = 5 * 60 * 1000; // 5 minutes

export function clearActorCache(): void {
  ACTOR_CACHE.clear();
}

/**
 * Cloudflare Worker handler to securely fetch actor emails and profiles (pfp_url) from IdentityX
 * without exposing IDENTITYX_ACCESS_API_KEY to the client browser.
 * Includes in-memory caching and single-flight per actor to prevent rate-limiting.
 */
export async function handleActorEmailsRequest(
  request: Request,
  env: Env,
): Promise<Response> {
  if (request.method !== "POST") {
    return new Response("Method not allowed", { status: 405 });
  }

  try {
    const body = (await request.json().catch(() => null)) as ActorEmailsRequest | null;
    if (!body || !Array.isArray(body.actorIds)) {
      return privateJsonResponse(
        { error: "Invalid request: actorIds array required" },
        { status: 400 },
      );
    }

    const actorIds = [...new Set(body.actorIds.filter(Boolean))];
    if (actorIds.length === 0) {
      return privateJsonResponse({});
    }

    const now = Date.now();
    const resultEntries: [string, ActorRecord][] = [];
    const missingIds: string[] = [];

    // Check in-memory cache first
    for (const actorId of actorIds) {
      const cached = ACTOR_CACHE.get(actorId);
      if (cached && cached.expiresAt > now) {
        resultEntries.push([actorId, cached.record]);
      } else {
        missingIds.push(actorId);
      }
    }

    // If everything was cached, return immediately
    if (missingIds.length === 0) {
      return privateJsonResponse(Object.fromEntries(resultEntries));
    }

    const baseURL = (env.VITE_AUTH_API_URL || "https://identityx.test").replace(/\/+$/, "");
    const projectId = env.VITE_TRIEOH_AUTH_PROJECT_ID;
    const apiKey = env.IDENTITYX_ACCESS_API_KEY;

    if (!projectId || !apiKey) {
      console.warn("[actor-emails] missing projectId or apiKey in server env");
      return privateJsonResponse(Object.fromEntries(resultEntries));
    }

    const fetchedResults = await Promise.all(
      missingIds.map(async (actorId): Promise<[string, ActorRecord] | null> => {
        try {
          const [actorRes, profileRes] = await Promise.all([
            fetch(`${baseURL}/projects/${projectId}/actors/${actorId}`, {
              method: "GET",
              headers: {
                "X-API-Key": apiKey,
                Accept: "application/json",
              },
            }).catch(() => null),
            fetch(`${baseURL}/projects/${projectId}/actors/${actorId}/profile`, {
              method: "GET",
              headers: {
                "X-API-Key": apiKey,
                Accept: "application/json",
              },
            }).catch(() => null),
          ]);

          let email: string | null = null;
          let pfp_url: string | null = null;

          if (actorRes?.ok) {
            const json = (await actorRes.json().catch(() => null)) as {
              data?: { email?: string | null };
            } | null;
            email = json?.data?.email ?? null;
          }

          if (profileRes?.ok) {
            const profileJson = (await profileRes.json().catch(() => null)) as {
              data?: { pfp_url?: string | null };
            } | null;
            pfp_url = profileJson?.data?.pfp_url ?? null;
          }

          if (!email && !pfp_url) {
            console.warn("[actor-emails] lookup found no email or profile", { actorId });
            return null;
          }

          const record: ActorRecord = { email, pfp_url };
          ACTOR_CACHE.set(actorId, { record, expiresAt: Date.now() + CACHE_TTL_MS });
          return [actorId, record];
        } catch (err) {
          console.warn("[actor-emails] error fetching actor", { actorId, err });
          return null;
        }
      }),
    );

    for (const item of fetchedResults) {
      if (item !== null) {
        resultEntries.push(item);
      }
    }

    const emailsMap: ActorEmailsMap = Object.fromEntries(resultEntries);
    return privateJsonResponse(emailsMap);
  } catch (err) {
    console.error("[actor-emails] unhandled exception", err);
    return privateJsonResponse({ error: "Internal server error" }, { status: 500 });
  }
}
