import { useCallback, useEffect, useState } from "react";
import "./App.css";
import { Editor } from "./components/Editor.tsx";
import { editor } from "monaco-editor";
import { ApplyOverlay, CalculateOverlay, GetInfo } from "./bridge.ts";
import { Alert } from "@speakeasy-api/moonshine";
import { blankOverlay, petstore } from "./defaults.ts";
import { useAtom } from "jotai";
import { atomWithHash } from "jotai-location";
import speakeasyWhiteLogo from "./assets/speakeasy-white.svg";
import openapiLogo from "./assets/openapi.svg";

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

  useEffect(() => {
    const tryHandlePageTitle = async () => {
      try {
        const info = await GetInfo(original);
        const { title, version } = JSON.parse(info);
        const pageTitle = `${title} ${version} | Speakeasy OpenAPI Overlay Playground`;
        if (document.title !== pageTitle) {
          document.title = pageTitle;
        }
      } catch (e: unknown) {
        console.error(e);
      }
    };
    tryHandlePageTitle();
  }, [original]);

  if (!ready) {
    return "";
  }

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        minHeight: "100vh",
        width: "100%",
      }}
    >
      <div style={{ paddingBottom: "1rem", width: "100vw" }}>
        <div className="border-b border-muted p-4 md:p-6 text-left">
          <div className="flex gap-2">
            <div className="flex flex-1">
              <div className="flex items-center pr-2">
                <img
                  src={openapiLogo}
                  alt="OpenAPI Logo"
                  className="h-12 w-12 shrink-0 grow-0 origin-center rotate-180 rounded-full"
                />
              </div>
              <div className="grow-1">
                <h1 className="text-xl font-semibold leading-none tracking-tight">
                  <a
                    className="underline hover:no-underline pr-1"
                    href="https://github.com/OAI/Overlay-Specification"
                  >
                    OpenAPI Overlay
                  </a>
                  Playground
                </h1>
                <p className="max-w-prose text-sm text-muted-foreground pt-2">
                  The OpenAPI Overlay Specification lets you update arbitrary
                  values in an YAML document using{" "}
                  <a
                    className="border-b border-transparent pb-[2px] transition-all duration-200 hover:border-current"
                    href="https://datatracker.ietf.org/doc/rfc9535/"
                  >
                    jsonpath
                  </a>
                  .
                </p>
              </div>
            </div>
            <div className="flex flex-1 flex-row-reverse">
              <ul className="flex gap-x-2">
                <li>
                  <a
                    className="border-b border-transparent pb-[2px] transition-all duration-200 hover:border-current"
                    href="https://www.speakeasy.com"
                  >
                    Made by the team at{" "}
                    <span className="sr-only">Speakeasy</span>
                    <picture>
                      <source
                        srcSet={speakeasyWhiteLogo}
                        media="(prefers-color-scheme: dark)"
                      />
                      <img
                        className="inline-block h-3 w-auto align-baseline"
                        src={speakeasyWhiteLogo}
                        alt=""
                      />
                    </picture>
                  </a>
                </li>
                <li className="before:pe-2 before:content-['•']">
                  <a
                    className="border-b border-transparent pb-[2px] transition-all duration-200 hover:border-current"
                    href="https://github.com/speakeasy-api/jsonpath"
                  >
                    GitHub
                  </a>
                </li>
                <li className="before:pe-2 before:content-['•']">
                  <a
                    className="border-b border-transparent pb-[2px] transition-all duration-200 hover:border-current"
                    href="https://github.com/OAI/Overlay-Specification"
                  >
                    OpenAPI Overlay
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </div>
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
