import { loader } from 'fumadocs-core/source';
import { lucideIconsPlugin } from 'fumadocs-core/source/lucide-icons';
import { docs } from 'collections/server';
import { docsContentRoute, docsImageRoute, docsRoute } from './shared';
import { openapiPlugin, openapiSource, createOpenAPI } from "fumadocs-openapi/server";
import { inputs } from "@/lib/openapi";

const openapiFiles = {
  files: (
    await Promise.all(
      inputs.map((input) =>
        openapiSource(createOpenAPI({ input: [input.url] }), {
          baseDir: `${input.name}/api`,
        }),
      ),
    )
  ).flatMap((res) => res.files),
};


export const source = loader(
  {
    source: docs.toFumadocsSource(),
    openapi: openapiFiles,
  },
  {
    baseUrl: docsRoute,
    plugins: [lucideIconsPlugin(), openapiPlugin()],
  },
);

export function getPageImage(slugs: string[]) {
  const segments = [...slugs, 'image.webp'];

  return {
    segments,
    url: `${docsImageRoute}/${segments.join('/')}`,
  };
}

export function getPageMarkdownUrl(page: (typeof source)['$inferPage']) {
  const segments = [...page.slugs, 'content.md'];

  return {
    segments,
    url: `${docsContentRoute}/${segments.join('/')}`,
  };
}

export async function getLLMText(page: (typeof source)['$inferPage']) {
  if (page.type === 'openapi') {
    return JSON.stringify(page.data.getSchema().bundled, null, 2);
  }

  const processed = await page.data.getText('processed');
  return `# ${page.data.title} (${page.url})\n\n${processed}`;
}
