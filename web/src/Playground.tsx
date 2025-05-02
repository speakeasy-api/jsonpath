import {
  ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import "./App.css";
import { Editor } from "./components/Editor";
import { editor, MarkerSeverity } from "monaco-editor";
import { ApplyOverlay, CalculateOverlay, GetInfo } from "./bridge";
import { Alert } from "@speakeasy-api/moonshine";
import { blankOverlay, emptyOverlay, petstore } from "./defaults";
import speakeasyWhiteLogo from "./assets/speakeasy-white.svg";
import openapiLogo from "./assets/openapi.svg";
import { compress, decompress } from "@/compress";
import {
  ImperativePanelGroupHandle,
  Panel,
  PanelGroup,
  PanelResizeHandle,
} from "react-resizable-panels";
import posthog from "posthog-js";
import { useDebounceCallback, useMediaQuery } from "usehooks-ts";
import FileUpload from "@/components/FileUpload";
import ShareButton from "@/components/ShareButton";
import {
  arraysEqual,
  formatDocument,
  guessDocumentLanguage,
} from "./lib/utils";
import ShareDialog, { ShareDialogHandle } from "./components/ShareDialog";
import { parse as yamlParse } from "yaml";

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

const tryHandlePageTitle = ({
  title,
  version,
}: {
  title: string;
  version: string;
}) => {
  try {
    const pageTitle = `${title} ${version} | Speakeasy OpenAPI Overlay Playground`;
    if (document.title !== pageTitle) {
      document.title = pageTitle;
    }
  } catch (e: unknown) {
    console.error(e);
  }
};

function Playground() {
  const [ready, setReady] = useState(false);

  const original = useRef(petstore);
  const implicitShare = useRef(false);
  const originalLang = useRef<DocumentLanguage>("yaml");
  const changed = useRef("");
  const [changedLoading, setChangedLoading] = useState(false);
  const [applyOverlayMode, setApplyOverlayMode] = useState<
    "original+overlay" | "jsonpathexplorer"
  >("original+overlay");
  let appliedPanelTitle = "Original + Edits";
  if (applyOverlayMode == "jsonpathexplorer") {
    appliedPanelTitle = "JSONPath Explorer";
  }
  const result = useRef(blankOverlay);
  const [resultLoading, setResultLoading] = useState(false);
  const [error, setError] = useState("");
  const [shareUrlLoading, setShareUrlLoading] = useState(false);
  const [overlayMarkers, setOverlayMarkers] = useState<editor.IMarkerData[]>(
    [],
  );
  const isSmallScreen = useMediaQuery("(max-width: 768px)");
  const clearError = useCallback(() => {
    setError("");
  }, []);
  const defaultLayout = useMemo(
    () => (isSmallScreen ? [20, 60, 20] : [30, 40, 30]),
    [],
  );

  const onChangeOverlay = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setChangedLoading(true);
        result.current = value || "";
        const response = await ApplyOverlay(
          original.current,
          result.current,
          true,
        );
        if (response.type == "success") {
          setApplyOverlayMode("original+overlay");
          changed.current = formatDocument(response.result);
          setError("");
          setOverlayMarkers([]);
          const info = await GetInfo(changed.current, false);
          tryHandlePageTitle(JSON.parse(info));
        } else if (response.type == "incomplete") {
          setApplyOverlayMode("jsonpathexplorer");

          if (originalLang.current == "json") {
            // !TODO: this is a hack to get around the fact
            //  that the json parser only returns yaml.
            const obj = yamlParse(response.result);
            changed.current = JSON.stringify(obj, null, 2);
          } else {
            changed.current = response.result;
          }

          setError("");
          setOverlayMarkers([]);
        } else if (response.type == "error") {
          setApplyOverlayMode("jsonpathexplorer");
          setOverlayMarkers([
            {
              startLineNumber: response.line,
              endLineNumber: response.line,
              startColumn: response.col,
              endColumn: response.col + 1000, // end of line
              message: response.error,
              severity: MarkerSeverity.Error, // Use MarkerSeverity from Monaco
            },
          ]);
        }
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      } finally {
        setChangedLoading(false);
      }
    },
    [],
  );

  const onChangeOverlayDebounced = useDebounceCallback(onChangeOverlay, 500);

  const shareDialogRef = useRef<ShareDialogHandle>(null);
  const lastSharedStart = useRef<string>("");

  const getShareUrl = useCallback(async () => {
    if (!shareDialogRef.current) return;

    try {
      setShareUrlLoading(true);
      const info = await GetInfo(original.current, false);
      const start = JSON.stringify({
        result: result.current,
        original: original.current,
        info: info,
      });

      const alreadySharedThis = lastSharedStart.current === start;
      if (alreadySharedThis) {
        shareDialogRef.current.setOpen(true);
        return;
      }

      const response = await fetch("/api/share", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: await compress(start),
      });

      if (response.ok) {
        const base64Data = await response.json();
        const currentUrl = new URL(window.location.href);

        currentUrl.hash = "";
        currentUrl.searchParams.set("s", base64Data);

        if (!implicitShare.current) {
          lastSharedStart.current = start;
          shareDialogRef.current.setUrl(currentUrl.toString());
          shareDialogRef.current.setOpen(true);
        }

        history.pushState(null, "", currentUrl.toString());
        posthog.capture("overlay.speakeasy.com:share", {
          openapi: JSON.parse(info),
        });
      } else {
        setError("Failed to create share URL");
      }
    } catch (e: any) {
      setError("Couldn't create share url: " + e.message);
    } finally {
      setShareUrlLoading(false);
    }
  }, [original, result]);

  useEffect(() => {
    (async () => {
      const urlParams = new URLSearchParams(window.location.search);
      const hash = urlParams.get("s");

      if (hash) {
        try {
          // base64 decode the download url
          const downloadUrl = atob(hash);
          const blob = await fetch(downloadUrl);
          if (!blob.body) {
            throw new Error("No body");
          }
          const decompressedData = await decompress(blob.body);
          const decompressed: { original: string; result: string } =
            JSON.parse(decompressedData);

          original.current = decompressed.original;
          result.current = decompressed.result;

          await onChangeOverlay(result.current, {} as any);
        } catch (error: any) {
          console.error("invalid share url:", error.message);
        }
      } else {
        try {
          const changedNew = await ApplyOverlay(
            original.current,
            result.current,
            false,
          );
          if (changedNew.type == "success") {
            changed.current = formatDocument(changedNew.result);
          }
        } catch (e: unknown) {
          if (e instanceof Error) {
            setError(e.message);
          } else {
            setError(JSON.stringify(e));
          }
          console.error(e);
        }
      }
      setReady(true);
    })();
  }, []);

  const onChangeA = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setResultLoading(true);
        original.current = value || "";
        originalLang.current = guessDocumentLanguage(original.current);
        const res = await CalculateOverlay(
          value || "",
          changed.current,
          result.current,
          true,
        );
        result.current = res;
        const info = await GetInfo(original.current, false);
        tryHandlePageTitle(JSON.parse(info));
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      } finally {
        setResultLoading(false);
      }
    },
    [],
  );

  const onChangeADebounced = useDebounceCallback(onChangeA, 500);

  const onChangeB = useCallback(
    async (value: string | undefined, _: editor.IModelContentChangedEvent) => {
      try {
        setResultLoading(true);
        changed.current = value || "";
        result.current = await CalculateOverlay(
          original.current,
          value || "",
          result.current,
          true,
        );
        const info = await GetInfo(changed.current, false);
        tryHandlePageTitle(JSON.parse(info));
        setError("");
      } catch (e: unknown) {
        if (e instanceof Error) {
          setError(e.message);
        }
      } finally {
        setResultLoading(false);
      }
    },
    [],
  );

  const onChangeBDebounced = useDebounceCallback(onChangeB, 500);

  const ref = useRef<ImperativePanelGroupHandle>(null);

  const maxLayout = useCallback((index: number) => {
    const panelGroup = ref.current;
    if (!panelGroup) return;

    const currentLayout = panelGroup?.getLayout();

    if (!arraysEqual(currentLayout, defaultLayout)) {
      panelGroup.setLayout(defaultLayout);
      return;
    }

    const baseWidth = 10;
    const maxedWidth = 80;
    const desiredWidths = Array(3).fill(baseWidth);

    if (index < desiredWidths.length && index >= 0) {
      desiredWidths[index] = maxedWidth;
    }

    // Reset each Panel to 50% of the group's width
    panelGroup.setLayout(desiredWidths);
  }, []);

  const onFileUpload = useCallback(
    async (content: string) => {
      setResultLoading(true);
      implicitShare.current = true;
      original.current = content;
      changed.current = content;
      result.current = emptyOverlay;
      await getShareUrl();
      setResultLoading(false);
      setOverlayMarkers([]);
    },
    [original, changed, result],
  );

  const clickShareButton = useCallback(async () => {
    implicitShare.current = false;
    await getShareUrl();
  }, [getShareUrl, implicitShare]);

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
      {isSmallScreen ? (
        <Alert variant="info">
          For proper user experience, please use a desktop device
        </Alert>
      ) : null}
      <div style={{ width: "100vw" }}>
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
                <p className="text-sm text-muted-foreground pt-2">
                  (Upload an OpenAPI spec and track edits as an overlay or write an overlay directly)
                </p>
              </div>
            </div>
            <div className="flex flex-1 flex-row-reverse">
              <div className="flex flex-col gap-4 justify-between">
                <div className="flex gap-x-2">
                  <span>
                    <Link href="https://www.speakeasy.com?utm_source=overlay.speakeasy.com">
                      Made by the team at
                      <div className="sr-only ml-2">Speakeasy</div>
                      <img
                        className="inline-block h-3 w-auto align-baseline ml-2"
                        src={speakeasyWhiteLogo}
                        alt=""
                      />
                    </Link>
                  </span>
                  <span className="before:pe-2 before:content-['•']">
                    <Link href="https://github.com/speakeasy-api/jsonpath">
                      GitHub
                    </Link>
                  </span>
                  <span className="before:pe-2 before:content-['•']">
                    <Link href="https://github.com/OAI/Overlay-Specification">
                      OpenAPI Overlay
                    </Link>
                  </span>
                </div>
                <div className="flex gap-x-2 justify-end">
                  <FileUpload onFileUpload={onFileUpload} />
                  <ShareButton
                    onClick={clickShareButton}
                    loading={shareUrlLoading}
                  />
                  <ShareDialog ref={shareDialogRef} />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      {error && (
        <Alert onDismiss={clearError} variant={"error"}>
          {error.split("\n").length > 1 ? (
            <>
              {error.split("\n")[0]}
              <br />
              <div className="text-left whitespace-pre">
                <pre>{error.split("\n").slice(1).join("\n")}</pre>
              </div>
            </>
          ) : (
            error
          )}
        </Alert>
      )}
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          width: "100%",
          maxWidth: "100vw",
          justifyContent: "space-between",
          gap: "1rem",
          overflow: "hidden",
        }}
      >
        <PanelGroup direction="horizontal" ref={ref}>
          <Panel defaultSize={defaultLayout[0]} minSize={10}>
            <div style={{ height: "calc(100vh - 50px)" }}>
              <Editor
                readonly={false}
                value={original.current}
                onChange={onChangeADebounced}
                title="Original"
                index={0}
                maxOnClick={maxLayout}
                language={originalLang.current}
              />
            </div>
          </Panel>
          <PanelResizeHandle />
          <Panel defaultSize={defaultLayout[1]} minSize={10}>
            <div style={{ height: "calc(100vh - 50px)" }}>
              <Editor
                readonly={false}
                original={
                  applyOverlayMode == "original+overlay"
                    ? original.current
                    : undefined
                }
                value={changed.current}
                onChange={onChangeBDebounced}
                loading={changedLoading}
                title={appliedPanelTitle}
                index={1}
                maxOnClick={maxLayout}
                language={originalLang.current}
              />
            </div>
          </Panel>
          <PanelResizeHandle />
          <Panel defaultSize={defaultLayout[2]} minSize={10}>
            <div style={{ height: "calc(100vh - 50px)" }}>
              <Editor
                readonly={false}
                value={result.current}
                onChange={onChangeOverlayDebounced}
                loading={resultLoading}
                markers={overlayMarkers}
                title={"Overlay"}
                index={2}
                maxOnClick={maxLayout}
                language="yaml"
              />
            </div>
          </Panel>
        </PanelGroup>
      </div>
    </div>
  );
}

export default Playground;
