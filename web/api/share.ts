import { put } from "@vercel/blob";
import { createHash } from "crypto";

const MAX_DATA_SIZE = 5 * 1024 * 1024; // 5MB
const AllowedOrigin = process.env.VERCEL_URL
  ? `https://${process.env.VERCEL_URL}`
  : "http://localhost";
const productionOrigin = "https://overlay.speakeasy.com";

export function POST(request: Request) {
  const origin = request.headers.get("Origin");

  if (
    !origin ||
    (new URL(origin).hostname != new URL(AllowedOrigin).hostname &&
      !origin.startsWith(productionOrigin))
  ) {
    return new Response("Unauthorized", { status: 403 });
  }

  return new Promise<Response>((resolve, reject) => {
    const body: ReadableStream | null = request.body;
    if (!body) {
      reject(new Error("No body"));
      return;
    }

    const reader = body.getReader();
    const chunks: Uint8Array[] = [];

    const readData = (): void => {
      reader.read().then(({ done, value }) => {
        if (done) {
          processData(chunks).then(resolve).catch(reject);
        } else {
          chunks.push(value);
          if (
            chunks.reduce((acc, chunk) => acc + chunk.length, 0) > MAX_DATA_SIZE
          ) {
            reject(new Error("Data exceeds the maximum allowed size of 5MB"));
          } else {
            readData();
          }
        }
      });
    };

    readData();
  });
}

async function processData(chunks: Uint8Array[]): Promise<Response> {
  const data = new Uint8Array(
    chunks.reduce((acc, chunk) => acc + chunk.length, 0),
  );
  let offset = 0;
  for (const chunk of chunks) {
    data.set(chunk, offset);
    offset += chunk.length;
  }

  const hash = createHash("sha256").update(data).digest("base64").slice(0, 7);
  const key = `share-urls/${hash}`;

  const result = await put(key, data, {
    addRandomSuffix: false,
    access: "public",
  });
  const downloadURL = result.downloadUrl;
  const encodedDownloadURL = Buffer.from(downloadURL.trim()).toString("base64");

  return new Response(JSON.stringify(encodedDownloadURL), {
    status: 200,
    headers: { "Content-Type": "application/json" },
  });
}
