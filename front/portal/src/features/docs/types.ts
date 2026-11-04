/** What mdx/remark-adr lifts out of a doc's header. Every key is optional. */
export interface DocMeta {
  title?: string;
  /** e.g. "ADR-0019" */
  number?: string;
  status?: string;
  date?: string;
  scope?: string;
  supersedes?: string;
}
