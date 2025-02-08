import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import MonacoEditor, { Monaco, DiffEditor } from "@monaco-editor/react";
import { editor, Uri } from "monaco-editor";
import { Progress } from "../../@/components/ui/progress";
import { Icon } from "@speakeasy-api/moonshine";
import { Button } from "@/components/ui/button";
import IModelContentChangedEvent = editor.IModelContentChangedEvent;
import ITextModel = editor.ITextModel;
import ICodeEditor = editor.ICodeEditor;
import IStandaloneDiffEditor = editor.IStandaloneDiffEditor;

export interface EditorComponentProps {
  readonly: boolean;
  value: string;
  original?: string;
  loading?: boolean;
  title: string;
  markers?: editor.IMarkerData[];
  index: number;
  maxOnClick?: (index: number) => void;
  onChange: (
    value: string | undefined,
    ev: editor.IModelContentChangedEvent,
  ) => void;
  language: DocumentLanguage;
}

const minLoadingTime = 150;

export function Editor(props: EditorComponentProps) {
  const editorRef = useRef<ICodeEditor | IStandaloneDiffEditor | null>(null);
  const monacoRef = useRef<Monaco | null>(null);
  const modelRef = useRef<ITextModel | null>(null);
  const [lastLoadingTime, setLastLoadingTime] = useState(minLoadingTime);
  const [progress, setProgress] = useState(100);

  const encodedTitle = useMemo(() => {
    return btoa(props.title);
  }, [props.title]);

  const EditorComponent =
    props.original === undefined ? MonacoEditor : DiffEditor;

  const onChange = useCallback(
    (value: string, event: IModelContentChangedEvent) => {
      props.onChange(value, event);
    },
    [props.onChange],
  );

  const handleEditorWillMount = useCallback((monaco: Monaco) => {
    monacoRef.current = monaco;
    const matchesURI = (uri: Uri | undefined) => {
      return uri?.path.includes(encodedTitle);
    };
    monaco.editor.onDidCreateModel((model) => {
      if (!matchesURI(model.uri)) {
        return;
      }
      if (props.original && !model.uri.path.includes("modified")) {
        return;
      }

      modelRef.current = model;
      modelRef.current.onDidChangeContent((event) => {
        if (editorRef.current?.hasTextFocus()) {
          onChange(model.getValue(), event);
        }
      });
    });
  }, []);

  const handleEditorDidMount = useCallback(
    (editor: ICodeEditor | IStandaloneDiffEditor, monaco: Monaco) => {
      editorRef.current = editor;

      const options = {
        base: "vs-dark",
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
      } satisfies editor.IStandaloneThemeData;
      monaco.editor.defineTheme("speakeasy", options);
      monaco.editor.setTheme("speakeasy");
    },
    [onChange],
  );

  const options: any = useMemo(
    () => ({
      readOnly: props.readonly,
      minimap: { enabled: false },
      automaticLayout: true,
      renderSideBySide: false,
    }),
    [props.readonly],
  );
  const isLoading = useMemo(
    () => props.loading || progress < 100,
    [props.loading, progress],
  );

  const onMaxClick = useCallback(() => {
    props.maxOnClick?.(props.index);
  }, [props.maxOnClick, props.index]);

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

  useEffect(() => {
    if (modelRef?.current) {
      monacoRef.current?.editor?.setModelMarkers(
        modelRef?.current,
        "diagnostics",
        props.markers || [],
      );
    }
  }, [props.markers]);

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
            {props.maxOnClick ? (
              <Button
                onClick={onMaxClick}
                className="bg-transparent"
                size="icon"
                variant="ghost"
              >
                <Icon name="maximize" />
              </Button>
            ) : null}
          </h1>
        </div>
      )}
      <EditorComponent
        onMount={handleEditorDidMount}
        beforeMount={handleEditorWillMount}
        original={props.original}
        modified={props.value}
        value={props.value}
        path={encodedTitle}
        originalModelPath={encodedTitle + "/original"}
        modifiedModelPath={encodedTitle + "/modified"}
        theme={"vscode-dark"}
        language={props.language}
        options={options}
      />
    </div>
  );
}
