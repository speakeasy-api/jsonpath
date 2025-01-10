import { ReactNode, useCallback, useEffect, useState } from "react";
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

const originalOpenAPI = atomWithHash("originalOpenAPI", petstore, {
  setHash: "replaceState",
});
const changedOpenAPI = atomWithHash("changedOpenAPI", petstore, {
  setHash: "replaceState",
});
const overlay = atomWithHash("overlay", blankOverlay, {
  setHash: "replaceState",
});

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

  const Link = ({ children, href }: { children: ReactNode; href: string }) => (
    <a
      className="border-b border-transparent pb-[2px] transition-all duration-200 hover:border-current "
      href={href}
      target="_blank"
      rel="noreferrer"
      style={{ color: "#FBE331" }}
    >
      {children}
    </a>
  );

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
              <div className="flex items-center pr-4">
                <img
                  src={openapiLogo}
                  alt="OpenAPI Logo"
                  className="h-12 w-12 shrink-0 grow-0 origin-center rotate-180 rounded-full"
                />
              </div>
              <div className="grow-1">
                <h1 className="text-xl font-semibold leading-none tracking-tight">
                  OpenAPI Overlay Playground
                </h1>
                <p className="max-w-prose text-sm text-muted-foreground pt-2">
                  The{" "}
                  <Link href="https://github.com/OAI/Overlay-Specification">
                    OpenAPI Overlay Specification
                  </Link>{" "}
                  lets you update arbitrary values in a YAML document using{" "}
                  <Link href="https://datatracker.ietf.org/doc/rfc9535/">
                    jsonpath
                  </Link>
                  .
                </p>
              </div>
            </div>
            <div className="flex flex-1 flex-row-reverse">
              <ul className="flex gap-x-2">
                <li>
                  <Link href="https://www.speakeasy.com?utm_source=overlay.speakeasy.com">
                    Made by the team at
                    <div className="sr-only ml-2">Speakeasy</div>
                    <picture>
                      <source
                        srcSet={speakeasyWhiteLogo}
                        media="(prefers-color-scheme: dark)"
                      />
                      <img
                        className="inline-block h-3 w-auto align-baseline ml-2"
                        src={speakeasyWhiteLogo}
                        alt=""
                      />
                    </picture>
                  </Link>
                </li>
                <li className="before:pe-2 before:content-['•']">
                  <Link href="https://github.com/speakeasy-api/jsonpath">
                    GitHub
                  </Link>
                </li>
                <li className="before:pe-2 before:content-['•']">
                  <Link href="https://github.com/OAI/Overlay-Specification">
                    OpenAPI Overlay
                  </Link>
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
          gap: "1rem",
          overflow: "hidden",
        }}
      >
        <div style={{ height: "calc(100vh - 50px)", width: "33vw" }}>
          <Editor
            readonly={false}
            value={original}
            onChange={onChangeA}
            title="Original"
          />
        </div>
        <div style={{ height: "calc(100vh - 50px)", width: "33vw" }}>
          <Editor
            readonly={false}
            value={changed}
            onChange={onChangeB}
            loading={changedLoading}
            title={"Original + Overlay"}
          />
        </div>
        <div style={{ height: "calc(100vh - 50px)", width: "33vw" }}>
          <Editor
            readonly={false}
            value={result}
            onChange={onChangeC}
            loading={resultLoading}
            title={"Overlay"}
          />
        </div>
      </div>
    </div>
  );
}

export default Playground;
