import { useCallback, useState } from "react";
import "./App.css";
import { Editor } from "./components/Editor.tsx";
import { editor } from "monaco-editor";

function App() {
  const [original, setOriginal] = useState("");

  const onChangeOriginal = useCallback(
    (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      setOriginal(value || "");
    },
    [],
  );

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
