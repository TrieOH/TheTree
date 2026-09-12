#!/usr/bin/env node
/**
 * Keeps framework boundaries honest across the workspace.
 *
 * The convention, inherited from the identityx SDKs:
 *
 *   @trieoh/<domain>          framework-free — runs in Node, a Worker, a browser
 *   @trieoh/<domain>-react    React only
 *   @trieoh/<domain>-solid    Solid only
 *
 * Why this exists: a neutral package that depends on React leaks React into
 * every consumer, so a Solid app ends up installing it "without needing to".
 * The name already says which framework a package is for; this makes the
 * manifest agree with the name, on every commit.
 *
 * Hard errors are reserved for `dependencies`/`optionalDependencies`, because
 * those are what a consumer installs transitively. Framework packages in
 * `devDependencies` (tests, storybook) or `peerDependencies` (the correct place
 * for a library) are reported as warnings.
 *
 * Usage: node tools/check-package-boundaries.mjs [--quiet]
 */

import { existsSync, readFileSync, readdirSync, statSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const quiet = process.argv.includes("--quiet");

/** Names that only exist in one framework's world. */
const REACT_ONLY = /^(@types\/)?react(-dom)?$|^react-|^@tanstack\/react-|^@base-ui\/react$|^lucide-react$|^sonner$|^recharts$/;
const SOLID_ONLY = /^solid-js$|^@solidjs\/|-solid$|^@tanstack\/solid-|^eslint-plugin-solid$/;

/**
 * Framework of a package, from its name. Anything without a suffix is neutral:
 * it must work in every runtime, so it may not depend on either framework.
 */
function frameworkOf(name) {
  if (/-react$/.test(name)) return "react";
  if (/-solid$/.test(name)) return "solid";
  return "neutral";
}

/**
 * Expanding one `*` per pattern is enough for the shapes this workspace uses
 * (`front/*`, `lib/ts/*`, `sdk/ts/*`), and avoids a glob dependency. Only the
 * `packages:` block is read — later sections use list syntax too.
 */
function workspaceDirs() {
  const workspace = readFileSync(join(root, "pnpm-workspace.yaml"), "utf8");
  const lines = workspace.split("\n");
  const start = lines.findIndex((line) => /^packages:\s*$/.test(line));
  const patterns = [];

  for (let index = start + 1; start !== -1 && index < lines.length; index += 1) {
    const line = lines[index];
    if (/^\S/.test(line)) break; // next top-level key
    const match = line.match(/^\s*-\s*'([^']+)'/);
    if (match) patterns.push(match[1]);
  }

  const dirs = [];

  for (const pattern of patterns) {
    const star = pattern.indexOf("*");
    if (star === -1) {
      dirs.push(pattern);
      continue;
    }
    const parent = join(root, pattern.slice(0, star));
    if (!existsSync(parent)) continue;
    for (const entry of readdirSync(parent)) {
      const candidate = join(parent, entry);
      if (statSync(candidate).isDirectory() && !entry.startsWith(".")) {
        dirs.push(join(pattern.slice(0, star), entry));
      }
    }
  }

  return dirs;
}

function readManifest(dir) {
  try {
    return JSON.parse(readFileSync(join(root, dir, "package.json"), "utf8"));
  } catch {
    return null;
  }
}

function matches(dep, regex) {
  if (regex.test(dep)) return true;
  // Scoped names: a React package may be pulled in as `@scope/react-foo`.
  const unscoped = dep.startsWith("@") ? dep.slice(dep.indexOf("/") + 1) : dep;
  return unscoped !== dep && regex.test(unscoped);
}

// --- collect ------------------------------------------------------------

const packages = new Map();
for (const dir of workspaceDirs()) {
  const manifest = readManifest(dir);
  if (manifest?.name) packages.set(manifest.name, { dir, manifest });
}

const MANIFEST_SECTIONS = [
  "dependencies",
  "optionalDependencies",
  "peerDependencies",
  "devDependencies",
];

/** Framework a dependency belongs to: a workspace package answers for itself. */
const depFramework = (dep) => {
  const local = packages.get(dep);
  if (local) return local.resolved ?? "neutral";
  if (matches(dep, REACT_ONLY)) return "react";
  if (matches(dep, SOLID_ONLY)) return "solid";
  return "neutral";
};

const resolveFramework = (entry) => {
  const react = entry.used.filter((item) => item.framework === "react");
  const solid = entry.used.filter((item) => item.framework === "solid");
  const mixed = entry.used.filter((item) => item.framework === "mixed");
  entry.resolved = (react.length > 0 && solid.length > 0) || mixed.length > 0
    ? "mixed"
    : react.length > 0
      ? "react"
      : solid.length > 0
        ? "solid"
        : "neutral";
  entry.why = {
    react: react.map((item) => `\`${item.dep}\``).join(", "),
    solid: solid.map((item) => `\`${item.dep}\``).join(", "),
    mixed: mixed.map((item) => `\`${item.dep}\``).join(", "),
  };
};

// A framework reaches a package either directly or through another workspace
// package, so resolve until nothing changes (two passes here, capped for safety).
for (let pass = 0; pass < 5; pass += 1) {
  let changed = false;

  for (const entry of packages.values()) {
    const used = [];
    for (const section of MANIFEST_SECTIONS) {
      for (const dep of Object.keys(entry.manifest[section] ?? {})) {
        const framework = depFramework(dep);
        if (framework !== "neutral") used.push({ dep, section, framework });
      }
    }

    const before = JSON.stringify(entry.used ?? null);
    entry.used = used;
    resolveFramework(entry);
    if (before !== JSON.stringify(used)) changed = true;
  }

  if (!changed) break;
}

/**
 * A library is something another workspace package depends on: its dependencies
 * get installed into its consumers, so its framework leaks. An app is a leaf —
 * it may use a framework, it just may not use two.
 */
const dependedOn = new Set();
for (const { manifest } of packages.values()) {
  for (const section of MANIFEST_SECTIONS) {
    for (const dep of Object.keys(manifest[section] ?? {})) {
      if (packages.has(dep)) dependedOn.add(dep);
    }
  }
}

const isLibrary = (name, manifest) => dependedOn.has(name) || !!manifest.exports;

// --- check --------------------------------------------------------------

const errors = [];
const warnings = [];

for (const [name, { dir, manifest, resolved, used, why }] of packages) {
  const named = frameworkOf(name);
  const library = isLibrary(name, manifest);
  const add = (list, message) => list.push({ name, dir, message });

  if (resolved === "mixed") {
    add(
      errors,
      "pulls in both frameworks — React via " +
      `${why.react || "(nothing directly)"}, Solid via ${why.solid || "(nothing directly)"}` +
      `${why.mixed ? `, already-mixed ${why.mixed}` : ""}. ` +
      "One of them has to go, or be pushed into a `-react`/`-solid` package.",
    );
    continue;
  }

  // The name is the contract. A neutral-named library that uses a framework
  // installs that framework into every consumer, whatever its name says.
  if (library && resolved !== "neutral" && named !== resolved) {
    add(
      errors,
      `named \`${name}\` but uses ${resolved} (${resolved === "react" ? why.react : why.solid}): ` +
      `rename it \`${name}-${resolved}\`, or move those imports into a \`-${resolved}\` package.`,
    );
  }

  for (const { dep, section, framework } of used) {
    if (framework === "mixed") {
      add(
        errors,
        `depends on \`${dep}\` (${section}), which itself mixes React and Solid`,
      );
      continue;
    }

    if (framework !== resolved) {
      add(errors, `is ${resolved} but depends on ${framework} \`${dep}\` in ${section}`);
      continue;
    }

    // Same framework, wrong section: an app picks the version, not the library.
    if (section === "dependencies" && library) {
      add(
        warnings,
        `\`${dep}\` should be a peerDependency, not a dependency, in a ${named} library`,
      );
    }
  }
}

// --- report -------------------------------------------------------------

const grouped = new Map();
for (const item of errors) {
  if (!grouped.has(item.name)) grouped.set(item.name, { dir: item.dir, items: [] });
  grouped.get(item.name).items.push(item.message);
}

if (grouped.size > 0) {
  console.error("\nPackage boundary violations:\n");
  for (const [name, { dir, items }] of grouped) {
    console.error(`  ${name}  (${dir})`);
    for (const message of items) console.error(`    ✗ ${message}`);
    console.error("");
  }
}

if (!quiet && warnings.length > 0) {
  console.log(`Warnings (${warnings.length}):\n`);
  for (const { name, message } of warnings) console.log(`  ${name}: ${message}`);
  console.log("");
}

console.log(
  `checked ${packages.size} packages · ` +
  `${errors.length} error(s) · ${warnings.length} warning(s)`,
);

if (errors.length > 0) {
  console.error(
    "\nA package's framework is its name suffix: unsuffixed = framework-free, " +
    "`-react` = React only, `-solid` = Solid only.\n",
  );
  process.exit(1);
}
