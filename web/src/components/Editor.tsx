import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import MonacoEditor, { Monaco } from "@monaco-editor/react";
import { editor } from "monaco-editor";
import { Progress } from "../../@/components/ui/progress.tsx";

export interface EditorComponentProps {
  readonly: boolean;
  value: string;
  loading?: boolean;
  title?: string;
  onChange: (
    value: string | undefined,
    ev: editor.IModelContentChangedEvent,
  ) => void;
}

const minLoadingTime = 150;

export function Editor(props: EditorComponentProps) {
  const editorRef = useRef<any>(null);
  const monacoRef = useRef<Monaco | null>(null);
  const [lastLoadingTime, setLastLoadingTime] = useState(minLoadingTime);
  const [progress, setProgress] = useState(100);

  const handleEditorDidMount = useCallback((editor: any, monaco: Monaco) => {
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
  }, []);

  const options: any = useMemo(
    () => ({
      readOnly: props.readonly,
      minimap: { enabled: false },
    }),
    [props.readonly],
  );
  const isLoading = useMemo(
    () => props.loading || progress < 100,
    [props.loading, progress],
  );

  useEffect(() => {
    if (props.loading) {
      const startTime = Date.now();

      return () => {
        const endTime = Date.now();
        setLastLoadingTime(Math.max(minLoadingTime, endTime - startTime));
      };
    }
  }, [props.loading]);

  useEffect(() => {
    if (isLoading) {
      const timer = setInterval(() => {
        setProgress((prevProgress) => {
          if (prevProgress >= 100) {
            clearInterval(timer);
            return 100;
          }
          return prevProgress + 10;
        });
      }, lastLoadingTime / 10);
      setProgress(0);
    }
  }, [isLoading]);

  const wrapperStyles = useMemo(() => {
    if (isLoading) {
      return {
        width: "100%",
        height: "100%",
        filter: "blur(1px)",
        position: "relative",
      } satisfies React.CSSProperties;
    }
    return {
      width: "100%",
      height: "100%",
      position: "relative",
    } satisfies React.CSSProperties;
  }, [isLoading]);

  return (
    <div style={wrapperStyles}>
      {progress < 100 && minLoadingTime < lastLoadingTime && (
        <Progress
          max={100}
          value={progress}
          className="absolute top-0 left-0 w-full h-1 z-10"
        />
      )}
      {props.title && (
        <div style={{ background: "#212015", padding: "1rem" }}>
          <h1 className="text-xl font-semibold leading-none tracking-tight">
            {props.title}
          </h1>
        </div>
      )}
      <MonacoEditor
        onMount={handleEditorDidMount}
        value={props.value}
        onChange={props.onChange}
        theme={"vscode-dark"}
        language="yaml"
        options={options}
      />
    </div>
  );
}
