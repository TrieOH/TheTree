// The spec types informd's json.RawMessage fields (answers, default values,
// select options, …) as objects, but they carry any JSON value, so the app
// keeps them loose and narrows where it reads them.
// biome-ignore lint/suspicious/noExplicitAny: arbitrary JSON
export type RawJSON = any;
