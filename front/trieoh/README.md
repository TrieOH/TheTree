# TrieOH Platform

Página institucional em Solid e Tailwind, usando `@trieoh/ui-solid` para os componentes, tokens e fontes.

- Desenvolvimento: `pnpm --filter @trieoh/trieoh-platform dev` (porta 3005).
- Build: `pnpm --filter @trieoh/trieoh-platform build`.
- Link de Components: `https://components.trieoh.com`, definido junto aos demais recursos em `src/App.tsx`.
- Deploy: tag `pages/trieoh/v0.1.0`, usando o workflow genérico de Pages.
- Criação inicial: projeto Pages / Direct Upload chamado `trieoh`, branch de produção `main`.

O projeto e o diretório de saída estão em `wrangler.jsonc`.

O tema acompanha o sistema até uma escolha manual, salva neste navegador.

A plataforma apresenta IdentityX, Univents, Informd e Payssage.
O Style Guide está em `http://styles.trieoh.com/`.
Os endereços dos sistemas estão na lista `products` de `src/App.tsx`.
