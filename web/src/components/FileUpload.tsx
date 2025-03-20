import { useCallback, useRef } from "react";
import { Icon } from "@speakeasy-api/moonshine";
import type { Attachment } from "@speakeasy-api/moonshine/dist/components/PromptInput";
import { cn } from "@/lib/utils";

export default function FileUpload(props: {
  onFileUpload: (content: string) => void;
}) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileUpload = useCallback(
    (files: Attachment[]) => {
      if (files.length > 0) {
        // Convert bytes to string
        const buffer = new Uint8Array(files[0].bytes);
        const content = new TextDecoder().decode(buffer);
        props.onFileUpload(content);
      }
    },
    [props.onFileUpload],
  );

  const handleFiles = useCallback(
    async (files: File[]) => {
      const attachments: Attachment[] = await Promise.all(
        files.map(async (file) => ({
          id: crypto.randomUUID(),
          name: file.name,
          type: file.type,
          size: file.size,
          bytes: await file.arrayBuffer(),
        })),
      );
      handleFileUpload(attachments);
    },
    [handleFileUpload],
  );

  return (
    <div className="mt-4 flex flex-shrink flex-row flex-wrap gap-3 items-end">
      {/* eslint-disable-next-line jsx-a11y/no-static-element-interactions,jsx-a11y/click-events-have-key-events */}
      <div
        key={"Upload file"}
        className={cn(
          "bg-foreground/5 hover:bg-foreground/10 text-foreground/80 relative flex cursor-pointer select-none flex-row items-center gap-1.5 whitespace-nowrap rounded-md border px-2.5 py-2 text-sm tracking-tight",
        )}
        onClick={() => fileInputRef.current?.click()}
      >
        <input
          type="file"
          ref={fileInputRef}
          className="absolute inset-0 hidden h-full w-full"
          onChange={(e) => handleFiles(Array.from(e.target.files ?? []))}
          accept={[
            "application/yaml",
            "application/x-yaml",
            "application/json",
          ].join(",")}
        />
        <Icon
          name={"paperclip"}
          className={cn("stroke-primary relative size-4", "stroke-emerald-400")}
          strokeWidth={1}
        />
        Upload file
      </div>
    </div>
  );
}
