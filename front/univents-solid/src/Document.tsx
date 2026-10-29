import type { ParentProps } from 'solid-js';
import { HydrationScript } from '@solidjs/web';
import { WaveSpinnerLoading } from '@/shared/ui/loader/WaveSpinnerLoading';

// Mirrors ThemeProvider: the "theme" key in localStorage, falling back to the OS.
const THEME_BOOTSTRAP = `(function(){try{var t=localStorage.getItem("theme");var d=t==="dark"||(t!=="light"&&matchMedia("(prefers-color-scheme: dark)").matches);var r=document.documentElement;r.classList.toggle("dark",d);r.classList.toggle("light",!d);r.style.colorScheme=d?"dark":"light"}catch(e){}})()`;

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
        <HydrationScript />
      </head>
      <body class="min-w-[320px] antialiased wrap:anywhere">
        <a class="sr-only focus:not-sr-only focus:fixed focus:left-4 focus:top-4 focus:z-100 focus:rounded-md focus:bg-background focus:px-4 focus:py-2" href="#main-content">
          Pular para o conteúdo principal
        </a>
        {/* Client-only shell: until the bundle has booted the page would be
            blank. This paints with the HTML + CSS and is removed by
            RouterHost (App.tsx) once the first route is on screen. It mirrors
            the AuthProvider fallback so the two hand over without a flash. */}
        <div id="boot" class="fixed inset-0 z-200 flex items-center justify-center bg-background">
          <WaveSpinnerLoading text="Carregando..." />
        </div>
        {props.children}
      </body>
    </html>
  );
}
