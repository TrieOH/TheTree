import type { ParentProps } from 'solid-js';
import { HydrationScript } from '@solidjs/web';

// Mirrors ThemeProvider: the "theme" key in localStorage, falling back to the OS.
const THEME_BOOTSTRAP = `(function(){try{var t=localStorage.getItem("theme");var d=t==="dark"||(t!=="light"&&matchMedia("(prefers-color-scheme: dark)").matches);var r=document.documentElement;r.classList.toggle("dark",d);r.classList.toggle("light",!d);r.style.colorScheme=d?"dark":"light"}catch(e){}})()`;

// The document shell — the new index.html: picked up by the src/Document.*
// convention, it wraps the app in the plugin's generated entries and must
// render the full <html>. Head tags go here. It is compiled only into the
// prerendered static shell and ships no bundle of its own (only the tiny
// theme bootstrap below, which has to run before first paint): in client mode
// <HydrationScript /> is stripped from the shell, and it activates when the
// app flips to SSR (`ssr: true` in vite.config.ts) — no document changes
// needed. Delete this file to fall back to the plugin's built-in shell.
export default function Document(props: ParentProps) {
  return (
    <html lang="pt-BR" class="overlay-scrollbar">
      <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="description" content="Descubra eventos, acompanhe programações e compre ingressos com a Univents." />
        <meta name="theme-color" content="#f5f5f5" media="(prefers-color-scheme: light)" />
        <meta name="theme-color" content="#0e1218" media="(prefers-color-scheme: dark)" />
        <meta name="color-scheme" content="light dark" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-title" content="Univents" />
        <meta name="mobile-web-app-capable" content="yes" />
        <link rel="manifest" href="/site.webmanifest" />
        <link rel="icon" href="/favicon.svg?v=2" type="image/svg+xml" />
        <link rel="icon" href="/favicon-96x96.png" type="image/png" sizes="96x96" />
        <link rel="apple-touch-icon" href="/apple-touch-icon.png" sizes="180x180" />
        <title>Univents</title>
        {/* Applied before first paint. ThemeProvider only runs once the bundle
            has loaded, so without this a dark-mode user sees a light page flash. */}
        <script innerHTML={THEME_BOOTSTRAP} />
        {/* The two fonts every page uses, fetched with the document instead of
            after the CSS, so text doesn't paint in a fallback and then reflow. */}
        <link rel="preload" href="/fonts/geist-variable.woff2" as="font" type="font/woff2" crossorigin="anonymous" />
        <link rel="preload" href="/fonts/space-grotesk-variable.woff2" as="font" type="font/woff2" crossorigin="anonymous" />
        <HydrationScript />
      </head>
      <body class="min-w-[320px] antialiased wrap:anywhere">
        <a class="sr-only focus:not-sr-only focus:fixed focus:left-4 focus:top-4 focus:z-100 focus:rounded-md focus:bg-background focus:px-4 focus:py-2" href="#main-content">
          Pular para o conteúdo principal
        </a>
        {props.children}
      </body>
    </html>
  );
}
