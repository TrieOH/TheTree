# @trieoh/ui-solid

Componentes compartilhados dos frontends Solid da TrieOH. O showcase é um consumidor
da biblioteca: controles, links, superfícies e exemplos de código devem vir deste
pacote. HTML semântico de estrutura (main, section, header, nav) fica no aplicativo.

## Tema

```css
@import "tailwindcss";
@import "@trieoh/ui-solid/theme.css";
/* Caminho relativo à folha do consumidor; ajuste para o seu workspace. */
@source "../../../../lib/ts/ui-solid/src";
```

`theme.css` reúne os tokens claro/escuro e as fontes locais Geist e Space Grotesk
usadas pelo Univents/TrieOH. Adicione a classe `dark` ao elemento raiz para alternar
o tema. Evite redefinir cores, fontes e estilos de controles no consumidor.
O Tailwind do consumidor precisa mapear as cores semânticas (`--color-background:
var(--background)` etc.) no seu bloco `@theme`.

## Controles e composição

- `Button`: ações; `Link`: navegação com `href` e foco visível.
- `Input`, `Textarea`, `Field`: entradas e seus rótulos/descrições.
- `Combobox`: escolha de opções, inclusive propriedades de exemplos. Não substituir
  por um select nativo no showcase. O menu em portal mede o campo e sua própria
  altura, respeita o viewport e acompanha rolagem e redimensionamento.
- `Checkbox`: valor booleano controlado por `checked` e `onChange`, com rótulo,
  descrição e estado desabilitado.
- `SearchInput`: busca controlada por `value` e `onInput`, com nome acessível obrigatório.
- `Card`: superfície semântica para conteúdo, com borda e cores do tema.
- `CodeBlock`: destaque lexical de TSX, nome de arquivo opcional (`filename`), cópia do código original e feedback por toast. Monte `Toaster` uma vez na aplicação. Os tokens `--code-*` em `theme.css` precisam ser mapeados para `--color-code-*` no `@theme` do consumidor.
- `OverlayScrollbar`: rolagem do documento, renderizada em portal; em dispositivos com mouse, oculta a barra nativa. Não controla contêineres internos com `overflow`.

Exemplos editáveis e copiáveis estão em `front/showcase`, nas páginas de cada
componente. Antes de implementar um controle local, procure a implementação aqui;
quando faltar, adicione o componente neste pacote e seu exemplo na vitrine.

## Origem no Univents

- Combobox, SectionTabs, Accordion, Drawer, Tooltip, modais, formulário em etapas
  e OverlayScrollbar já são compartilhados com o Univents por reexports.
- Checkbox extrai o Toggle de `features/profile/ui/ProfileEditor.tsx`; o editor
  agora importa a mesma implementação da UI, sem manter uma cópia local.
- Card adota a superfície (borda, raio e sombra) de `widgets/ui/DashboardStatCard.tsx`.
  Os cards de domínio do Univents continuam sendo composições específicas.
- Não havia componentes genéricos independentes de CodeBlock, SearchInput ou Link
  nas pastas shared/ui e widgets/ui do Univents. Esses são novos no pacote.
- O exemplo MultiStepDialog precisa passar `values` e `onChange` e configurar
  `summary.items` na etapa de revisão; nomear uma etapa “Revisão” não cria um resumo.

## Verificação

```sh
pnpm -F @trieoh/ui-solid typecheck
pnpm -F @trieoh/ui-solid check
pnpm -F @trieoh/showcase build
```

A regressão de navegador está em `front/showcase/scripts/verify-ui.mjs`. Ela espera
um servidor Vite em `SHOWCASE_URL` (padrão `http://127.0.0.1:3004`) e `playwright-core`
disponível. `PLAYWRIGHT_MODULE` pode apontar para um módulo Playwright instalado
externamente, e `PLAYWRIGHT_CHROMIUM_EXECUTABLE` para um Chromium local. O script
valida alinhamento/largura do menu, seleção, edição de propriedades, código e tema
em 375, 471 e 1280 px.
