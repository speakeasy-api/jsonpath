// bridge.ts
import OpenAPIWorker from "./openapi.web.worker.ts?worker";

const wasmWorker = new OpenAPIWorker();
let messageQueue: { resolve: Function; reject: Function; message: any }[] = [];
let isProcessing = false;

function processQueue() {
  if (isProcessing || messageQueue.length === 0) return;

  isProcessing = true;
  const { resolve, reject, message } = messageQueue.shift()!;

  wasmWorker.postMessage(message);
  wasmWorker.onmessage = (event: MessageEvent<any>) => {
    if (event.data.type.endsWith("Result")) {
      resolve(event.data.payload);
    } else if (event.data.type.endsWith("Error")) {
      reject(new Error(event.data.error));
    }
    isProcessing = false;
    processQueue();
  };
}

function sendMessage(message: any): Promise<any> {
  return new Promise((resolve, reject) => {
    messageQueue.push({ resolve, reject, message });
    processQueue();
  });
}

export type CalculateOverlayMessage = {
  Request: {
    type: "CalculateOverlay";
    payload: {
      from: string;
      to: string;
    };
  };
  Response:
    | {
        type: "CalculateOverlayResult";
        payload: string;
      }
    | {
        type: "CalculateOverlayError";
        error: string;
      };
};

export type ApplyOverlayMessage = {
  Request: {
    type: "ApplyOverlay";
    payload: {
      source: string;
      overlay: string;
    };
  };
  Response:
    | {
        type: "ApplyOverlayResult";
        payload: string;
      }
    | {
        type: "ApplyOverlayError";
        error: string;
      };
};

export function CalculateOverlay(from: string, to: string): Promise<any> {
  return sendMessage({
    type: "CalculateOverlay",
    payload: { from, to },
  } satisfies CalculateOverlayMessage["Request"]);
}

export function ApplyOverlay(source: string, overlay: string): Promise<string> {
  return sendMessage({
    type: "ApplyOverlay",
    payload: { source, overlay },
  } satisfies ApplyOverlayMessage["Request"]);
}
