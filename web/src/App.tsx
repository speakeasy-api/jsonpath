import { useCallback, useEffect, useState } from "react";
import "./App.css";
import { Editor } from "./components/Editor.tsx";
import { editor } from "monaco-editor";
import speakeasyWASM from "./assets/wasm/lib.wasm?url";
function App() {
  const [ready, setReady] = useState(false);
  const [original, setOriginal] = useState("");
  useEffect(() => {
    (async () => {
      const wasmResponse = await fetch(speakeasyWASM);
      const body = await wasmResponse.body;
      // tostring
      if (!body) {
        return;
      }
      const reader = body.getReader();
      const txt = await reader.read();
      const decoder = new TextDecoder();
      const decoded = decoder.decode(txt.value);
      console.log(decoded);
      const go = new Go();
      const wasmArray = await wasmResponse.arrayBuffer();
      const result = (await WebAssembly.instantiate(
        wasmArray,
        go.importObject,
      )) as any;
      await go.run(result.instance);
      setReady(true);
    })();
  }, []);
  const onChangeOriginal = useCallback(
    (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      setOriginal(value || "");
    },
    [],
  );

  if (!ready) {
    return <div>Loading...</div>;
  }

  return (
    <>
      <div style={{ height: "100vh", width: "30vw" }}>
        <Editor readonly={false} value={original} onChange={onChangeOriginal} />
      </div>
      <div style={{ height: "100vh", width: "30vw" }}>
        <Editor readonly={false} value={original} onChange={onChangeOriginal} />
      </div>
      <div style={{ height: "100vh", width: "30vw" }}>
        <Editor readonly={false} value={original} onChange={onChangeOriginal} />
      </div>
    </>
  );
}

export default App;
