import { Icon } from "@speakeasy-api/moonshine";
import { cn } from "@/lib/utils";
import { Loader2Icon } from "lucide-react";

export default function ShareButton(props: {
  onClick: () => void;
  loading: boolean;
}) {
  return (
    <div className="mt-4 flex flex-shrink flex-row flex-wrap gap-3 items-stretch">
      {/* eslint-disable-next-line jsx-a11y/click-events-have-key-events */}
      <div
        key={"Share"}
        className={cn(
          "bg-foreground/5 hover:bg-foreground/10 text-foreground/80 relative flex cursor-pointer select-none flex-row items-center gap-1.5 whitespace-nowrap rounded-md border px-2.5 py-2 text-sm tracking-tight",
          props.loading && "cursor-not-allowed",
        )}
        onClick={props.loading ? undefined : props.onClick}
      >
        {props.loading ? (
          <Loader2Icon className="animate-spin" style={{ height: "75%" }} />
        ) : (
          <Icon
            name={"share"}
            className={cn(
              "stroke-primary relative size-4",
              "stroke-emerald-400",
            )}
            strokeWidth={1}
          />
        )}
        Share
      </div>
    </div>
  );
}
