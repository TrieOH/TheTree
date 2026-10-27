import { render } from "@solidjs/web";
import "./styles/app.css";
import App from "./App";

const root = document.getElementById("root");
if (root) {
  render(() => <App />, root);
}
