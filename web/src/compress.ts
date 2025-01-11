export const compress = async (decodedString: string): Promise<Blob> => {
  try {
    const stream = new Blob([decodedString]).stream();
    const compressedStream: ReadableStream<Uint8Array> = stream.pipeThrough(
      new CompressionStream("gzip"),
    );
    const reader = compressedStream.getReader();
    const chunks = [];
    for (
      let chunk = await reader.read();
      !chunk.done;
      chunk = await reader.read()
    ) {
      chunks.push(chunk.value);
    }
    const blob = new Blob(chunks);

    return blob;
  } catch (e: any) {
    throw new Error("failed to compress string: " + e.message);
  }
};

export const decompress = async (
  stream: ReadableStream<Uint8Array>,
): Promise<string> => {
  const decompressedStream = stream.pipeThrough(
    new DecompressionStream("gzip"),
  );
  const decompressedReader = decompressedStream.getReader();
  const decompressedChunks = [];
  for (
    let chunk = await decompressedReader.read();
    !chunk.done;
    chunk = await decompressedReader.read()
  ) {
    decompressedChunks.push(chunk.value);
  }

  const blob = new Blob(decompressedChunks);
  return blob.text();
};
