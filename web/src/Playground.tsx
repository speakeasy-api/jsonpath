import React, { useCallback, useEffect, useState } from "react";
import "./App.css";
import { Editor } from "./components/Editor.tsx";
import { editor } from "monaco-editor";
import { ApplyOverlay, CalculateOverlay } from "./bridge.ts";
import { Alert, PageHeader } from "@speakeasy-api/moonshine";
import { blankOverlay, petstore } from "./defaults.ts";
import { useAtom } from "jotai";
import { atomWithHash } from "jotai-location";

const originalOpenAPI = atomWithHash("originalOpenAPI", petstore);
const changedOpenAPI = atomWithHash("changedOpenAPI", petstore);
const overlay = atomWithHash("overlay", blankOverlay);

function Playground() {
  const [ready, setReady] = useState(false);

  const [original, setOriginal] = useAtom(originalOpenAPI);
  const [changed, setChanged] = useAtom(changedOpenAPI);
  const [changedLoading, setChangedLoading] = useState(false);
  const [result, setResult] = useAtom(overlay);
  const [resultLoading, setResultLoading] = useState(false);
  const [error, setError] = useState("");
  useEffect(() => {
    (async () => {
      setReady(true);
    })();
  }, []);
  const onChangeA = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setResultLoading(true);
        setOriginal(value || "");
        const res = await CalculateOverlay(value || "", changed, true);
        setResult(res);
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      } finally {
        setResultLoading(false);
      }
    },
    [changed, original],
  );
  const onChangeB = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setResultLoading(true);
        setChanged(value || "");
        const res = await CalculateOverlay(original, value || "", true);
        setResult(res);
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      } finally {
        setResultLoading(false);
      }
    },
    [changed, original],
  );

  const onChangeC = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setChangedLoading(true);
        setResult(value || "");
        const res = await ApplyOverlay(original, value || "", true);
        setChanged(res);
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      } finally {
        setChangedLoading(false);
      }
    },
    [changed, original],
  );

  if (!ready) {
    return "";
  }

  return (
    <div
      style={{
        display: "flex",
        color: "white",
        flexDirection: "column",
        alignItems: "center",
        minHeight: "100vh",
        width: "100%",
      }}
    >
      <div style={{ paddingBottom: "1rem", width: "100vw" }}>
        <PageHeader
          image="https://avatars.githubusercontent.com/u/91446104?s=200&v=4"
          subtitle="Best in class API tooling for robust SDKs, Terraform Providers and End to End Testing. OpenAPI Native."
          title="Speakeasy OpenAPI Overlay Playground"
        >
          <div style={{ width: "100%" }} />
        </PageHeader>
      </div>
      {error && <Alert variant={"error"}>{error}</Alert>}
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          width: "100%",
          justifyContent: "space-between",
        }}
      >
        <div style={{ height: "calc(100vh - 50px)", width: "33vw" }}>
          <Editor readonly={false} value={original} onChange={onChangeA} />
        </div>
        <div style={{ height: "calc(100vh - 50px)", width: "33vw" }}>
          <Editor
            readonly={false}
            value={changed}
            onChange={onChangeB}
            loading={changedLoading}
          />
        </div>
        <div style={{ height: "calc(100vh - 50px)", width: "33vw" }}>
          <Editor
            readonly={false}
            value={result}
            onChange={onChangeC}
            loading={resultLoading}
          />
        </div>
      </div>
    </div>
  );
}

export default Playground;
