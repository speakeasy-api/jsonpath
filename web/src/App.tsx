import { useCallback, useEffect, useState } from "react";
import "./App.css";
import { Editor } from "./components/Editor.tsx";
import { editor } from "monaco-editor";
import speakeasyWASM from "./assets/wasm/lib.wasm?url";

declare var CalculateOverlay: any;

function App() {
  const [ready, setReady] = useState(false);
  const [original, setOriginal] = useState("");
  const [changed, setChanged] = useState("");
  const [result, setResult] = useState("");
  useEffect(() => {
    (async () => {
      // const body = await wasmResponse.body;
      // // tostring
      // if (!body) {
      //   return;
      // }
      // const reader = body.getReader();
      // const txt = await reader.read();
      // const decoder = new TextDecoder();
      // const decoded = decoder.decode(txt.value);
      // console.log(decoded);
      // const arr = new Uint8Array(decoded.length);
      // for (let i = 0; i < decoded.length; i++) {
      //   arr[i] = decoded.charCodeAt(i);
      // }
      const go = new Go();
      const result = await WebAssembly.instantiateStreaming(
        fetch(speakeasyWASM),
        go.importObject,
      );
      go.run(result.instance);
      console.log(result);
      // await go.run(result.instance);
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
        <Editor readonly={true} value={result} onChange={console.log} />
      </div>
    </>
  );
}

export default App;
