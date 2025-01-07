import React, { useCallback, useEffect, useMemo, useState } from "react";
import "./App.css";
import { Editor } from "./components/Editor.tsx";
import { editor } from "monaco-editor";
import { ApplyOverlay, CalculateOverlay } from "./bridge.ts";
import { Alert, Badge, PageHeader } from "@speakeasy-api/moonshine";
import { blankOverlay, petstore } from "./defaults.ts";

function Playground() {
  const [ready, setReady] = useState(false);
  const [original, setOriginal] = useState(petstore);
  const [executing, setExecuting] = useState(false);
  const [changed, setChanged] = useState(petstore);
  const [result, setResult] = useState(blankOverlay);
  const [error, setError] = useState("");
  const isLoading = useMemo(() => !ready && !executing, [ready, executing]);
  useEffect(() => {
    (async () => {
      setReady(true);
    })();
  }, []);
  const onChangeA = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setExecuting(true);
        setOriginal(value || "");
        const res = await CalculateOverlay(value || "", changed, true);
        setResult(res);
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      } finally {
        setExecuting(false);
      }
    },
    [changed, original],
  );
  const onChangeB = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setExecuting(true);
        setChanged(value || "");
        const res = await CalculateOverlay(original, value || "", true);
        setResult(res);
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      } finally {
        setExecuting(false);
      }
    },
    [changed, original],
  );

  const onChangeC = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setExecuting(true);
        setResult(value || "");
        const res = await ApplyOverlay(original, value || "", true);
        setChanged(res);
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      } finally {
        setExecuting(false);
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
          <div style={{ width: "100%" }}>
            <div className="mt-2">
              {isLoading && <Badge variant="default">Loading</Badge>}
            </div>
          </div>
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
        <div style={{ height: "100vh", width: "33vw" }}>
          <Editor readonly={false} value={original} onChange={onChangeA} />
        </div>
        <div style={{ height: "100vh", width: "33vw" }}>
          <Editor readonly={false} value={changed} onChange={onChangeB} />
        </div>
        <div style={{ height: "100vh", width: "33vw" }}>
          <Editor readonly={false} value={result} onChange={onChangeC} />
        </div>
      </div>
    </div>
  );
}

export default Playground;
