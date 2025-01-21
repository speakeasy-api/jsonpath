// openapi.web.worker.ts
import speakeasyWASM from "./assets/wasm/lib.wasm?url";
import "./assets/wasm/wasm_exec.js";
import type {
  CalculateOverlayMessage,
  ApplyOverlayMessage,
  GetInfoMessage,
} from "./bridge";

const _wasmExecutors = {
  CalculateOverlay: (..._: any): any => false,
  ApplyOverlay: (..._: any): any => false,
  GetInfo: (..._: any): any => false,
} as const;

type MessageHandlers = {
  [K in keyof typeof _wasmExecutors]: (payload: any) => Promise<any>;
};

const messageHandlers: MessageHandlers = {
  CalculateOverlay: async (
    payload: CalculateOverlayMessage["Request"]["payload"],
  ) => {
    return exec("CalculateOverlay", payload.from, payload.to, payload.existing);
  },
  ApplyOverlay: async (payload: ApplyOverlayMessage["Request"]["payload"]) => {
    return exec("ApplyOverlay", payload.source, payload.overlay);
  },
  GetInfo: async (payload: GetInfoMessage["Request"]["payload"]) => {
    return exec("GetInfo", payload.openapi);
  },
};

let instantiated = false;

async function Instantiate() {
  if (instantiated) {
    return;
  }
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch(speakeasyWASM),
    go.importObject,
  );
  go.run(result.instance);
  for (const funcName of Object.keys(_wasmExecutors)) {
    // @ts-ignore
    if (!globalThis[funcName]) {
      throw new Error("missing expected function " + funcName);
    }
    // @ts-ignore
    _wasmExecutors[funcName] = globalThis[funcName];
  }
  instantiated = true;
}

async function exec(funcName: keyof typeof _wasmExecutors, ...args: any) {
  if (!instantiated) {
    await Instantiate();
  }
  if (!_wasmExecutors[funcName]) {
    throw new Error("not defined");
  }
  return _wasmExecutors[funcName](...args);
}

self.onmessage = async (
  event: MessageEvent<
    CalculateOverlayMessage["Request"] | ApplyOverlayMessage["Request"]
  >,
) => {
  const { type, payload } = event.data;
  try {
    const handler = messageHandlers[type];
    if (handler) {
      const result = await handler(payload);
      self.postMessage({ type: `${type}Result`, payload: result });
    } else {
      throw new Error(`Unknown message type: ${type}`);
    }
  } catch (err: any) {
    if (err && err.message) {
      self.postMessage({ type: `${type}Error`, error: err.message });
    } else {
      self.postMessage({ type: `${type}Error`, error: "unknown error" });
    }
  }
};
