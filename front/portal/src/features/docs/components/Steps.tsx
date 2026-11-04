import type { JSX } from "@solidjs/web";

/**
 * Numbered steps: every `###` heading inside starts the next one.
 *
 *   <Steps>
 *   ### Install the tools
 *   …
 *   ### Start the APIs
 *   …
 *   </Steps>
 */
export function Steps(props: { children?: JSX.Element }) {
  return <div class="doc-steps my-6">{props.children}</div>;
}
