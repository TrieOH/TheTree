import type { BeforeMount } from "@monaco-editor/react";
import type { editor, languages, Position } from "monaco-editor";

export const ZED_LANGUAGE_ID = "zed";

export const beforeMount: BeforeMount = (monaco) => {
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
        },
        {
          label: "relation",
          kind: monaco.languages.CompletionItemKind.Keyword,
          insertText: "relation ${1:name}: ${2:type}",
          insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          range,
        },
        {
          label: "permission",
          kind: monaco.languages.CompletionItemKind.Keyword,
          insertText: "permission ${1:name} = ${2:relation}",
          insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          range,
        },
      ];

      return { suggestions };
    },
  });
};