import type { ParentProps } from "solid-js";
import { HydrationScript } from "@solidjs/web";

// Mirrors ThemeProvider: dark unless localStorage says "light".
const THEME_BOOTSTRAP = `(function(){var d=true;try{d=localStorage.getItem("theme")!=="light"}catch(e){}var r=document.documentElement;r.classList.toggle("dark",d);r.classList.toggle("light",!d);r.style.colorScheme=d?"dark":"light"})()`;

export default function Document(props: ParentProps) {
  return (
    <html lang="en" class="dark">
      <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="description" content="Docs, services and the state of TheTree in one place." />
        <meta name="theme-color" content="#0e1218" />
        <meta name="color-scheme" content="dark light" />
        <link rel="icon" href="/favicon.svg" type="image/svg+xml" />
        <title>TrieOH dev</title>
        <script innerHTML={THEME_BOOTSTRAP} />
        <HydrationScript />
      </head>
      <body class="min-w-[320px] antialiased wrap:anywhere">
        <a
          class="sr-only focus:not-sr-only focus:fixed focus:left-4 focus:top-4 focus:z-100 focus:rounded-md focus:bg-background focus:px-4 focus:py-2"
          href="#main-content"
        >
          Skip to content
        </a>
        {props.children}
      </body>
    </html>
  );
}
