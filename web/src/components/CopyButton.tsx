"use client";

import * as React from "react";
import { DropdownMenuTriggerProps } from "@radix-ui/react-dropdown-menu";
import { CheckIcon, ClipboardIcon } from "lucide-react";

import { cn } from "@/lib/utils";
import { Button, ButtonProps } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";

interface CopyButtonProps extends ButtonProps {
  value: string;
  src?: string;
}

export async function copyToClipboardWithMeta(value: string) {
  navigator.clipboard.writeText(value);
}

export function CopyButton({
  value,
  className,
  variant = "ghost",
  ...props
}: CopyButtonProps) {
  const [hasCopied, setHasCopied] = React.useState(false);

  React.useEffect(() => {
    setTimeout(() => {
      setHasCopied(false);
    }, 5000);
  }, [hasCopied]);

  return (
    <Button
      size="icon"
      variant={variant}
      className={cn(
        "flex grow relative border-none bg-foreground/5 hover:bg-foreground/10 text-foreground/80   cursor-pointer select-none flex-row items-center gap-1.5 whitespace-nowrap rounded-md border px-2.5 py-2 text-sm tracking-tight",
        className,
      )}
      onClick={() => {
        copyToClipboardWithMeta(value);
        setHasCopied(true);
      }}
      style={{ background: "transparent" }}
      {...props}
    >
      <Input
        readOnly
        value={value}
        className="w-full font-mono font-semibold text-sm bg-zinc-900"
      />
      <span className=" h-full aspect-square rounded-md flex items-center justify-center  text-zinc-50 bg-zinc-900 hover:bg-zinc-700 hover:text-zinc-50 [&_svg]:h-3 [&_svg]:w-3">
        <span className="sr-only aspect-square">Copy</span>
        {hasCopied ? <CheckIcon /> : <ClipboardIcon />}
      </span>
    </Button>
  );
}

interface CopyWithClassNamesProps extends DropdownMenuTriggerProps {
  value: string;
  classNames: string;
  className?: string;
}

export function CopyWithClassNames({
  value,
  classNames,
  className,
}: CopyWithClassNamesProps) {
  const [hasCopied, setHasCopied] = React.useState(false);

  React.useEffect(() => {
    setTimeout(() => {
      setHasCopied(false);
    }, 2000);
  }, [hasCopied]);

  const copyToClipboard = React.useCallback((value: string) => {
    copyToClipboardWithMeta(value);
    setHasCopied(true);
  }, []);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          size="icon"
          variant="ghost"
          className={cn(
            "relative z-10 h-6 w-6 text-zinc-50 hover:bg-zinc-700 hover:text-zinc-50",
            className,
          )}
        >
          {hasCopied ? (
            <CheckIcon className="h-3 w-3" />
          ) : (
            <ClipboardIcon className="h-3 w-3" />
          )}
          <span className="sr-only">Copy</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => copyToClipboard(value)}>
          Component
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => copyToClipboard(classNames)}>
          Classname
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
