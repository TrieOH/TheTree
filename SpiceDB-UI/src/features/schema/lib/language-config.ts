import type { BeforeMount } from "@monaco-editor/react";
import type { editor, languages, Position } from "monaco-editor";

export const ZED_LANGUAGE_ID = "zed";

const KEYWORD_DOCS: Record<string, { detail: string; documentation: string }> = {
  definition: {
    detail: "definition <n> { ... }",
    documentation:
      "Declares a new resource type in the SpiceDB schema.\n\n" +
      "**Example:**\n```\ndefinition document {\n  relation owner: user\n  permission view = owner\n}\n```",
  },
  relation: {
    detail: "relation <n>: <type>[#<relation>] [| ...]",
    documentation:
      "Declares a relation between two types. Supports multiple types with `|` and " +
      "sub-relation references with `#`.\n\n" +
      "**Examples:**\n```\nrelation owner: user\nrelation viewer: user | group#member\nrelation parent: folder\n```",
  },
  permission: {
    detail: "permission <n> = <expr>",
    documentation:
      "Defines a computed permission from relations and logical operators.\n\n" +
      "**Available operators:** `+` (union), `&` (intersection), `-` (exclusion)\n\n" +
      "**Example:**\n```\npermission edit = owner + editor\npermission view = viewer & member\npermission delete = owner - banned\n```",
  },
};

export const beforeMount: BeforeMount = (monaco) => {
  monaco.editor.defineTheme("zed-dark", {
    base: "vs-dark",
    inherit: true,
    rules: [
      { token: "keyword", foreground: "6366f1" },
      { token: "comment", foreground: "64748b", fontStyle: "italic" },
      { token: "identifier", foreground: "e2e8f0" },
      { token: "type", foreground: "38bdf8" },
      { token: "operator", foreground: "94a3b8" },
      { token: "delimiter", foreground: "94a3b8" },
      { token: "delimiter.arrow", foreground: "fb923c" },
      { token: "string", foreground: "86efac" },
      { token: "number", foreground: "fbbf24" },
    ],
    colors: {
      "editor.background": "#00000000",
      "editor.foreground": "#e2e8f0",
      "editorLineNumber.foreground": "#475569",
      "editor.lineHighlightBackground": "#1e293b",
      "editorCursor.foreground": "#6366f1",
    },
  });

  monaco.editor.defineTheme("zed-light", {
    base: "vs",
    inherit: true,
    rules: [
      { token: "keyword", foreground: "4f46e5" },
      { token: "comment", foreground: "94a3b8", fontStyle: "italic" },
      { token: "identifier", foreground: "1e293b" },
      { token: "type", foreground: "0284c7" },
      { token: "operator", foreground: "64748b" },
      { token: "delimiter", foreground: "64748b" },
      { token: "delimiter.arrow", foreground: "ea580c" },
      { token: "string", foreground: "16a34a" },
      { token: "number", foreground: "d97706" },
    ],
    colors: {
      "editor.background": "#00000000",
      "editor.foreground": "#1e293b",
      "editorLineNumber.foreground": "#94a3b8",
      "editor.lineHighlightBackground": "#f1f5f9",
      "editorCursor.foreground": "#4f46e5",
    },
  });

  if (monaco.languages.getLanguages().some((l: languages.ILanguageExtensionPoint) => l.id === ZED_LANGUAGE_ID)) return;

  monaco.languages.register({ id: ZED_LANGUAGE_ID });

  monaco.languages.setMonarchTokensProvider(ZED_LANGUAGE_ID, {
    keywords: ["definition", "relation", "permission", "with", "and", "or", "not", "nil"],

    tokenizer: {
      root: [
        [/\/\/.*$/, "comment"],
        [/\/\*/, "comment", "@blockComment"],

        [/"([^"\\]|\\.)*$/, "string.invalid"],
        [/'([^'\\]|\\.)*$/, "string.invalid"],

        [/\b\d+(\.\d+)?\b/, "number"],
        [/->/, "delimiter.arrow"],

        [/\b(definition|relation|permission|with|and|or|not|nil)\b/, "keyword"],

        [/[a-zA-Z_][\w/]*/, "identifier"],
        [/[{}()]/, "delimiter.bracket"],
        [/[+\-&|]/, "operator"],
        [/=/, "operator"],
        [/:/, "delimiter"],
        [/#/, "delimiter"],
        [/,/, "delimiter"],
      ],

      blockComment: [
        [/[^/*]+/, "comment"],
        [/\*\//, "comment", "@pop"],
        [/[/*]/, "comment"],
      ],
    },
  });

  monaco.languages.registerCompletionItemProvider(ZED_LANGUAGE_ID, {
    provideCompletionItems: (model: editor.ITextModel, position: Position) => {
      const word = model.getWordUntilPosition(position);
      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endColumn: word.endColumn,
      };

      const CK = monaco.languages.CompletionItemKind;
      const SNIPPET = monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet;

      const suggestions: languages.CompletionItem[] = [
        {
          label: "definition",
          kind: CK.Keyword,
          insertText: "definition ${1:name} {\n\t$0\n}",
          insertTextRules: SNIPPET,
          range,
          detail: KEYWORD_DOCS.definition.detail,
          documentation: { value: KEYWORD_DOCS.definition.documentation },
        },
        {
          label: "relation",
          kind: CK.Keyword,
          insertText: "relation ${1:name}: ${2:type}",
          insertTextRules: SNIPPET,
          range,
          detail: KEYWORD_DOCS.relation.detail,
          documentation: { value: KEYWORD_DOCS.relation.documentation },
        },
        {
          label: "permission",
          kind: CK.Keyword,
          insertText: "permission ${1:name} = ${2:relation}",
          insertTextRules: SNIPPET,
          range,
          detail: KEYWORD_DOCS.permission.detail,
          documentation: { value: KEYWORD_DOCS.permission.documentation },
        },
      ];

      return { suggestions };
    },
  });

  monaco.languages.registerHoverProvider(ZED_LANGUAGE_ID, {
    provideHover: (model: editor.ITextModel, position: Position) => {
      const word = model.getWordAtPosition(position);
      if (!word) return null;

      const token = word.word;

      const doc = KEYWORD_DOCS[token];
      return {
        range: new monaco.Range(
          position.lineNumber,
          word.startColumn,
          position.lineNumber,
          word.endColumn
        ),
        contents: [
          { value: `**\`${doc.detail}\`**` },
          { value: doc.documentation },
        ],
      };
    },
  });
};