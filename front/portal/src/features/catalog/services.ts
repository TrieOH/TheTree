// The four APIs as the portal describes them. Status, success rate and p95
// come from stored Cloudflare and Grafana rollups in later phases; nothing
// here is ever fetched from the services themselves.

export type ServiceId = "identityx" | "payssage" | "univents" | "informd";

export interface Service {
  id: ServiceId;
  name: string;
  description: string;
  apiUrl: string;
  front: string;
  dir: string;
  /** Where the analytics come from (never called directly). */
  analytics: {
    /** Cloudflare zone holding the API host. */
    zoneId: string;
    apiHost: string;
    /** Worker in front of the API container, and the front's Worker. */
    apiWorker: string;
    frontWorker: string;
    neonProject: string;
  };
}

export const SERVICES: readonly Service[] = [
  {
    id: "identityx",
    name: "IdentityX",
    description: "Multi-tenant authentication and authorization.",
    apiUrl: "https://api.identityx.trieoh.com",
    front: "https://dev.identityx.com.br",
    dir: "api/identityx",
    analytics: {
      zoneId: "6fa6890c8818fbdec24d20546ab06f08",
      apiHost: "api.identityx.trieoh.com",
      apiWorker: "identityx-api",
      frontWorker: "identityx",
      neonProject: "small-mud-79083906",
    },
  },
  {
    id: "payssage",
    name: "Payssage",
    description: "Payment routing on top of MercadoPago.",
    apiUrl: "https://api.payssage.trieoh.com",
    front: "https://dev.payssage.trieoh.com",
    dir: "api/payssage",
    analytics: {
      zoneId: "6fa6890c8818fbdec24d20546ab06f08",
      apiHost: "api.payssage.trieoh.com",
      apiWorker: "payssage-api",
      frontWorker: "payssage",
      neonProject: "curly-heart-99301388",
    },
  },
  {
    id: "univents",
    name: "Univents",
    description: "Events and ticketing for universities.",
    apiUrl: "https://api.univents.com.br",
    front: "https://dev.univents.com.br",
    dir: "api/univents",
    analytics: {
      zoneId: "d9faa8d4b3c1255a18993cd7aea4732f",
      apiHost: "api.univents.com.br",
      apiWorker: "univents-api",
      frontWorker: "univents",
      neonProject: "bold-mud-67143707",
    },
  },
  {
    id: "informd",
    name: "Informd",
    description: "Embeddable form builder.",
    apiUrl: "https://api.informd.trieoh.com",
    front: "https://dev.informd.trieoh.com",
    dir: "api/informd",
    analytics: {
      zoneId: "6fa6890c8818fbdec24d20546ab06f08",
      apiHost: "api.informd.trieoh.com",
      apiWorker: "informd-api",
      frontWorker: "informd",
      neonProject: "little-fog-25621543",
    },
  },
];

export function findService(id: string): Service | undefined {
  return SERVICES.find((service) => service.id === id);
}

export const REPO_URL = "https://github.com/TrieOH/TheTree";
