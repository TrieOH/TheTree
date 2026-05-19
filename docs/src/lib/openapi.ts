import { createOpenAPI } from 'fumadocs-openapi/server';

function parseServicesEnv() {
  const raw = process.env.OPENAPI_SERVICES;

  if (!raw) return [];

  return raw.split(',').map((entry) => {
    const [name, url] = entry.split('|');

    if (!name || !url) return null;

    return {
      name: name.trim(),
      url: url.trim(),
    };
  }).filter((service): service is {
    name: string;
    url: string;
  } => service !== null);
}

export const services = parseServicesEnv();

async function resolveAvailableInputs() {
  const results = await Promise.allSettled(
    services.map((service) =>
      fetch(service.url, { signal: AbortSignal.timeout(3000) })
        .then((r) => {
          if (!r.ok) throw new Error();
          return service;
        })
    )
  );
  return results
    .filter((r): r is PromiseFulfilledResult<{ name: string; url: string }> => r.status === 'fulfilled')
    .map((r) => r.value);
}

export const inputs = await resolveAvailableInputs();

export const openapi = createOpenAPI({
  input: inputs.map(i => i.url),
});