import { Editor } from "@monaco-editor/react";
import { beforeMount, ZED_LANGUAGE_ID } from "../lib/language-config";
import type { editor } from "monaco-editor";

interface PropsI {
  value: string
  onChange: (newValue: string | undefined) => void
  theme?: "zed-dark" | "zed-light"
  onMount?: (editor: editor.IStandaloneCodeEditor) => void;
}

export default function SchemaEditor({
  value,
  onChange,
  onMount,
  theme = "zed-dark"
}: PropsI) {
  return (
    <Editor
      height="100%"
      language={ZED_LANGUAGE_ID}
      value={value}
      onChange={onChange}
      onMount={onMount}
      beforeMount={beforeMount}
      theme={theme}
      options={{
        minimap: { enabled: false },
        fontSize: 14,
        tabSize: 2,
        scrollBeyondLastLine: false,
        formatOnPaste: true,
      }}
    />
  )
}