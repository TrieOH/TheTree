// MDX compiles to `jsxDEV` in development; the runtime is the same.
import { jsx } from "./jsx-runtime";

export { Fragment } from "./jsx-runtime";
export const jsxDEV = (type: Parameters<typeof jsx>[0], props: Parameters<typeof jsx>[1]) =>
  jsx(type, props);
