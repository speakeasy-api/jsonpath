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
      existing: string;
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

export type QueryJSONPathMessage = {
  Request: {
    type: "QueryJSONPath";
    payload: {
      source: string;
      jsonpath: string;
    };
  };
  Response:
    | {
        type: "QueryJSONPathResult";
        payload: string;
      }
    | {
        type: "QueryJSONPathError";
        error: string;
      };
};

export function CalculateOverlay(
  from: string,
  to: string,
  existing: string,
  supercede = false,
): Promise<any> {
  return sendMessage(
    {
      type: "CalculateOverlay",
      payload: { from, to, existing },
    } satisfies CalculateOverlayMessage["Request"],
    supercede,
  );
}

type IncompleteOverlayErrorMessage = {
  type: "incomplete";
  line: number;
  col: number;
  result: string;
};

type JSONPathErrorMessage = {
  type: "error";
  line: number;
  col: number;
  error: string;
};

type ApplyOverlayResultMessage = {
  type: "success";
  result: string;
};

type ApplyOverlaySuccess =
  | ApplyOverlayResultMessage
  | IncompleteOverlayErrorMessage
  | JSONPathErrorMessage;

export async function ApplyOverlay(
  source: string,
  overlay: string,
  supercede = false,
): Promise<ApplyOverlaySuccess> {
  const result = await sendMessage(
    {
      type: "ApplyOverlay",
      payload: { source, overlay },
    } satisfies ApplyOverlayMessage["Request"],
    supercede,
  );
  return JSON.parse(result);
}

export function QueryJSONPath(
  source: string,
  jsonpath: string,
  supercede = false,
): Promise<string> {
  return sendMessage(
    {
      type: "QueryJSONPath",
      payload: { source, jsonpath },
    } satisfies QueryJSONPathMessage["Request"],
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
