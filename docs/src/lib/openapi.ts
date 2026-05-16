import { createOpenAPI } from 'fumadocs-openapi/server';

export const services = [
  { name: 'identityx', url: 'http://localhost:8080/swagger/doc.json' },
  { name: 'univents', url: 'http://localhost:8081/swagger/doc.json' },
  { name: 'payssage', url: 'http://localhost:8082/swagger/doc.json' },
  { name: 'informd', url: 'http://localhost:8083/swagger/doc.json' },
];

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