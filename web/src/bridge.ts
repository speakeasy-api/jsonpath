// bridge.ts
import OpenAPIWorker from "./openapi.web.worker.ts?worker";

const wasmWorker = new OpenAPIWorker();
let messageQueue: {
  resolve: Function;
  reject: Function;
  message: any;
  supercede: boolean;
}[] = [];
let isProcessing = false;

function processQueue() {
  if (isProcessing || messageQueue.length === 0) return;

  isProcessing = true;
  const { resolve, reject, message } = messageQueue.shift()!;

  wasmWorker.postMessage(message);
  wasmWorker.onmessage = (event: MessageEvent<any>) => {
    if (event.data.type.endsWith("Result")) {
      // Reject all superceded messages before resolving the current message
      const supercedeMessages = messageQueue.filter((_, index) => {
        return messageQueue.slice(index + 1).some((later) => later.supercede);
      });
      supercedeMessages.forEach((m) => m.reject(new Error("supercedeerror")));
      messageQueue = messageQueue.filter((m) => !supercedeMessages.includes(m));
      resolve(event.data.payload);
    } else if (event.data.type.endsWith("Error")) {
      reject(new Error(event.data.error));
    }
    isProcessing = false;
    processQueue();
  };
}

function sendMessage(message: any, supercede = false): Promise<any> {
  return new Promise((resolve, reject) => {
    messageQueue.push({ resolve, reject, message, supercede });
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

export type GetInfoMessage = {
  Request: {
    type: "GetInfo";
    payload: {
      openapi: string;
    };
  };
  Response:
    | {
        type: "GetInfoResult";
        payload: string;
      }
    | {
        type: "GetInfoError";
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

export function CalculateOverlay(
  from: string,
  to: string,
  supercede = false,
): Promise<any> {
  return sendMessage(
    {
      type: "CalculateOverlay",
      payload: { from, to },
    } satisfies CalculateOverlayMessage["Request"],
    supercede,
  );
}

export function ApplyOverlay(
  source: string,
  overlay: string,
  supercede = false,
): Promise<string> {
  return sendMessage(
    {
      type: "ApplyOverlay",
      payload: { source, overlay },
    } satisfies ApplyOverlayMessage["Request"],
    supercede,
  );
}

export function GetInfo(openapi: string, supercede = false): Promise<string> {
  return sendMessage(
    {
      type: "GetInfo",
      payload: { openapi },
    } satisfies GetInfoMessage["Request"],
    supercede,
  );
}
