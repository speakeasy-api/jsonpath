import { useCallback, useEffect, useState } from "react";
import "./App.css";
import { Editor } from "./components/Editor.tsx";
import { editor } from "monaco-editor";
import { ApplyOverlay, CalculateOverlay } from "./bridge.ts";

function Playground() {
  const [ready, setReady] = useState(false);
  const [original, setOriginal] = useState("");
  const [changed, setChanged] = useState("");
  const [result, setResult] = useState("");
  const [error, setError] = useState("");
  useEffect(() => {
    (async () => {
      setReady(true);
    })();
  }, []);
  const onChangeA = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setOriginal(value || "");
        const res = await CalculateOverlay(value || "", changed);
        setResult(res);
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      }
    },
    [changed, original],
  );
  const onChangeB = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setChanged(value || "");
        const res = await CalculateOverlay(original, value || "");
        setResult(res);
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      }
    },
    [changed, original],
  );

  const onChangeC = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setResult(value || "");
        const res = await ApplyOverlay(original, value || "");
        setChanged(res);
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      }
    },
    [changed, original],
  );

  if (!ready) {
    return <div>Loading...</div>;
  }

  return (
    <>
      {error && <div>{error}</div>}
      <div style={{ height: "100vh", width: "30vw" }}>
        <Editor readonly={false} value={original} onChange={onChangeA} />
      </div>
      <div style={{ height: "100vh", width: "30vw" }}>
        <Editor readonly={false} value={changed} onChange={onChangeB} />
      </div>
      <div style={{ height: "100vh", width: "30vw" }}>
        <Editor readonly={false} value={result} onChange={onChangeC} />
      </div>
    </>
  );
}

export default Playground;
