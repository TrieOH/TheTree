import { Editor } from "@monaco-editor/react";
import { beforeMount, ZED_LANGUAGE_ID } from "../lib/language-config";

interface PropsI {
  value: string
  onChange: (newValue: string | undefined) => void
}

export default function SchemaEditor({
  value,
  onChange
}: PropsI) {
  return (
    <Editor
      height="100%"
      language={ZED_LANGUAGE_ID}
      value={value}
      onChange={onChange}
      beforeMount={beforeMount}
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