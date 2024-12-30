import { useCallback, useEffect, useState } from "react";
import "./App.css";
import { Editor } from "./components/Editor.tsx";
import { editor } from "monaco-editor";
import speakeasyWASM from "./assets/wasm/lib.wasm?url";

declare var CalculateOverlay: any;
declare var ApplyOverlay: any;

function App() {
  const [ready, setReady] = useState(false);
  const [original, setOriginal] = useState("");
  const [changed, setChanged] = useState("");
  const [result, setResult] = useState("");
  useEffect(() => {
    (async () => {
      const go = new Go();
      const result = await WebAssembly.instantiateStreaming(
        fetch(speakeasyWASM),
        go.importObject,
      );
      go.run(result.instance);
      setReady(true);
    })();
  }, []);
  const onChangeA = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      setOriginal(value || "");
      const res = await CalculateOverlay(value, changed);
      setResult(res);
    },
    [changed, original],
  );
  const onChangeB = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      setChanged(value || "");
      const res = await CalculateOverlay(original, value);
      setResult(res);
    },
    [changed, original],
  );

  const onChangeC = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      setResult(value || "");
      const res = await ApplyOverlay(original, value);
      setChanged(res);
    },
    [changed, original],
  );

  if (!ready) {
    return <div>Loading...</div>;
  }

  return (
    <>
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

export default App;
