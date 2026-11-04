import { Callout } from "./Callout";
import { Card, Cards } from "./Cards";
import { DocTabs, Tab } from "./DocTabs";
import { elements } from "./elements";
import { Mermaid } from "./Mermaid";
import { Problems } from "./Problems";
import { Steps } from "./Steps";

/** Everything a doc can use without importing it. */
export const docComponents = {
  ...elements,
  Callout,
  Card,
  Cards,
  Mermaid,
  Problems,
  Steps,
  Tab,
  Tabs: DocTabs,
};
