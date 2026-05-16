import { createOpenAPI } from 'fumadocs-openapi/server';

const services = [
  'http://localhost:8080/swagger/doc.json', // identity-x
  'http://localhost:8081/swagger/doc.json', // univents
  'http://localhost:8082/swagger/doc.json', // payssage
  'http://localhost:8083/swagger/doc.json', // informd
];

async function resolveAvailableInputs() {
  const results = await Promise.allSettled(
    services.map((url) =>
      fetch(url, { signal: AbortSignal.timeout(3000) })
        .then((r) => { if (!r.ok) throw new Error(); return url; })
    )
  );
  return results
    .filter((r): r is PromiseFulfilledResult<string> => r.status === 'fulfilled')
    .map((r) => r.value);
}

export const inputs = await resolveAvailableInputs();

export const openapi = inputs.length > 0
  ? createOpenAPI({ input: inputs })
  : createOpenAPI({ input: [] });