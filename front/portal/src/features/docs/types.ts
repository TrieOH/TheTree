/** What mdx/remark-adr lifts out of a doc's header. Every key is optional. */
export interface DocMeta {
  title?: string;
  /** e.g. "ADR-0019" */
  number?: string;
  /** proposed, accepted, superseded or deprecated */
  status?: string;
  /** When the decision was made, YYYY-MM-DD. */
  date?: string;
  /** When the decision was last changed in place, YYYY-MM-DD. */
  updated?: string;
  /** One line on what the doc covers. */
  summary?: string;
  /** The paths the decision applies to. */
  scope?: string[];
  /** ADR numbers this one replaces, changes, or is replaced by. */
  supersedes?: string[];
  amends?: string[];
  amended_by?: string[];
  superseded_by?: string;
  /** Reading order within a folder, e.g. "1". Docs without one come after, by title. */
  order?: string;
}
