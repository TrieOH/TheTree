import type { BeforeMount } from "@monaco-editor/react";
import type { editor, languages, Position } from "monaco-editor";

export const ZED_LANGUAGE_ID = "zed";

export const beforeMount: BeforeMount = (monaco) => {
  monaco.editor.defineTheme("zed-dark", {
    base: "vs-dark",
    inherit: true,
    rules: [
      { token: "keyword", foreground: "6366f1" }, // indigo-500
      { token: "comment", foreground: "64748b", fontStyle: "italic" }, // slate-500
      { token: "identifier", foreground: "e2e8f0" }, // slate-200
      { token: "operator", foreground: "94a3b8" }, // slate-400
      { token: "delimiter", foreground: "94a3b8" }, // slate-400
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
      { token: "keyword", foreground: "4f46e5" }, // indigo-600
      { token: "comment", foreground: "94a3b8", fontStyle: "italic" }, // slate-400
      { token: "identifier", foreground: "1e293b" }, // slate-900
      { token: "operator", foreground: "64748b" }, // slate-500
      { token: "delimiter", foreground: "64748b" }, // slate-500
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
        [/\/\*/, "comment", "@comment"],
        [/\b(definition|relation|permission|with|and|or|not|nil)\b/, "keyword"],
        [/[a-zA-Z_][\w/]*/, "identifier"],
        [/[{}]/, "delimiter.bracket"],
        [/[+|&]/, "operator"],
        [/=/, "operator"],
        [/:/, "delimiter"],
        [/#/, "delimiter"],
      ],
      comment: [
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

      const suggestions = [
        {
          label: "definition",
          kind: monaco.languages.CompletionItemKind.Keyword,
          insertText: "definition ${1:name} {\n\t$0\n}",
          insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          range,
          detail: "Define um novo tipo",
          documentation: "Cria uma definição de tipo no schema Zed",
        },
        {
          label: "relation",
          kind: monaco.languages.CompletionItemKind.Keyword,
          insertText: "relation ${1:name}: ${2:type}",
          insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          range,
          detail: "Define uma relação",
          documentation: "Declara uma relação entre tipos",
        },
        {
          label: "permission",
          kind: monaco.languages.CompletionItemKind.Keyword,
          insertText: "permission ${1:name} = ${2:relation}",
          insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          range,
          detail: "Define uma permissão",
          documentation: "Cria uma permissão baseada em relações",
        },
      ];

      return { suggestions };
    },
  });
};