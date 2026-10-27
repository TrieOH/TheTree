import { render } from "@solidjs/web";
import App from "./App";
import { applyTheme, readThemePreference } from "./theme";
import "./styles.css";

applyTheme(
  readThemePreference() ??
    (window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light"),
);

const root = document.getElementById("root");
if (root) render(() => <App />, root);
