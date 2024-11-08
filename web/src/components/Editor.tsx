import { useRef } from "react";
import MonacoEditor, { Monaco } from "@monaco-editor/react";
import { editor } from "monaco-editor";

export interface EditorComponentProps {
  readonly: boolean;
  value: string;
  onChange: (
    value: string | undefined,
    ev: editor.IModelContentChangedEvent,
  ) => void;
}

export function Editor(props: EditorComponentProps) {
  const editorRef = useRef<any>(null);
  const monacoRef = useRef<any>(null);

  function handleEditorDidMount(editor: any, monaco: Monaco) {
    editorRef.current = editor;
    monacoRef.current = monaco;

    const options = {
      base: "vs-dark",
      renderSideBySide: false,
      inherit: true,
      rules: [
        {
          foreground: "F3F0E3",
          token: "string",
        },
        {
          foreground: "679FE1",
          token: "type",
        },
      ],
      colors: {
        "editor.foreground": "#F3F0E3",
        "editor.background": "#212015",
        "editorCursor.foreground": "#679FE1",
        "editor.lineHighlightBackground": "#1D2A3A",
        "editorLineNumber.foreground": "#6368747F",
        "editorLineNumber.activeForeground": "#FBE331",
        "editor.inactiveSelectionBackground": "#FF3C742D",
        "diffEditor.removedTextBackground": "#FF3C741A",
        "diffEditor.insertedTextBackground": "#1D2A3A",
      },
    };
    // @ts-ignore
    monaco.editor.defineTheme("speakeasy", options);
    monaco.editor.setTheme("speakeasy");
  }

  const options: any = {
    readOnly: props.readonly,
    minimap: { enabled: false },
  };

  return (
    <MonacoEditor
      onMount={handleEditorDidMount}
      value={props.value}
      onChange={props.onChange}
      theme={"vscode-dark"}
      language="yaml"
      options={options}
    />
  );
}
