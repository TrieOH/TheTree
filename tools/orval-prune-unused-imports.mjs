#!/usr/bin/env node
/**
 * Orval post-generation hook.
 *
 * The `fetch` client sometimes emits a type import the generated file never
 * uses (for example `import type { Uuid } from './uuid'` in a params schema that
 * already refers to `ProjectIDQueryParameter`). Those imports fail
 * `noUnusedLocals` in the apps that type-check the generated sources, and
 * hand-editing generated files is not an option because `pnpm orval` rewrites
 * them.
 *
 * So we prune provably-unused imports from the generated schema files after
 * every run. Only files under `client/schemas/` are touched, and an import is
 * only dropped when none of its bindings appear anywhere else in the file.
 *
 * Orval runs this through `output.hooks.afterAllFilesWrite` in
 * `orval.config.ts`, passing the generated directories as arguments.
 */
import { readFileSync, readdirSync, statSync, writeFileSync } from "node:fs"
import { join } from "node:path"

/** Matches a whole `import` / `import type` statement, capturing its bindings. */
const IMPORT_RE = /^import\s+(?:type\s+)?\{([^}]*)\}\s+from\s+["'][^"']+["'];?[ \t]*$/gm

/** Only generated schema files are candidates. */
const isSchemaFile = (filePath) =>
  /[/\\]client[/\\]schemas[/\\][^/\\]+\.ts$/.test(filePath) && !filePath.endsWith("index.ts")

/** Orval hands us directories; collect every `.ts` file below them. */
function collect(dir) {
  let entries
  try {
    entries = readdirSync(dir)
  } catch {
    return isSchemaFile(dir) ? [dir] : []
  }
  return entries.flatMap((entry) => {
    const path = join(dir, entry)
    return statSync(path).isDirectory() ? collect(path) : [path]
  })
}

function pruneFile(filePath) {
  const source = readFileSync(filePath, "utf8")
  let changed = false

  const next = source.replace(IMPORT_RE, (statement, bindings) => {
    const names = bindings
      .split(",")
      .map((entry) => entry.trim())
      .filter(Boolean)
      // `Foo as Bar` is used under `Bar`.
      .map((entry) => {
        const [, original, alias] = /^(\w+)\s+as\s+(\w+)$/.exec(entry) ?? []
        return { original: original ?? entry, alias: alias ?? entry }
      })

    const body = source.replace(statement, "")
    const used = names.filter(({ alias }) => new RegExp(`\\b${alias}\\b`).test(body))

    if (used.length === names.length) return statement
    changed = true

    if (used.length === 0) return ""
    return statement.replace(
      bindings,
      ` ${used.map(({ original, alias }) => (original === alias ? original : `${original} as ${alias}`)).join(", ")} `,
    )
  })

  if (!changed) return false
  // Collapse the blank line a fully-removed import can leave behind.
  writeFileSync(filePath, next.replace(/\n{3,}/g, "\n\n").replace(/\n+$/, "\n"))
  return true
}

const targets = [...new Set(process.argv.slice(2).flatMap(collect).filter(isSchemaFile))]
const pruned = targets.filter(pruneFile)

if (pruned.length > 0) {
  console.log(`orval: pruned unused imports in ${pruned.length} schema file(s)`)
}
