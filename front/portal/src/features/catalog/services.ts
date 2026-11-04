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
}

export const SERVICES: readonly Service[] = [
  {
    id: "identityx",
    name: "IdentityX",
    description: "Multi-tenant authentication and authorization.",
    apiUrl: "https://api.identityx.trieoh.com",
    front: "https://dev.identityx.com.br",
    dir: "api/identityx",
  },
  {
    id: "payssage",
    name: "Payssage",
    description: "Payment routing on top of MercadoPago.",
    apiUrl: "https://api.payssage.trieoh.com",
    front: "https://dev.payssage.trieoh.com",
    dir: "api/payssage",
  },
  {
    id: "univents",
    name: "Univents",
    description: "Events and ticketing for universities.",
    apiUrl: "https://api.univents.com.br",
    front: "https://dev.univents.com.br",
    dir: "api/univents",
  },
  {
    id: "informd",
    name: "Informd",
    description: "Embeddable form builder.",
    apiUrl: "https://api.informd.trieoh.com",
    front: "https://dev.informd.trieoh.com",
    dir: "api/informd",
  },
];

export function findService(id: string): Service | undefined {
  return SERVICES.find((service) => service.id === id);
}

export const REPO_URL = "https://github.com/TrieOH/TheTree";
