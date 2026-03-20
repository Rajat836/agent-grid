import { Loader2Icon } from "lucide-react";

import { cn } from "@/lib/utils";

function Spinner({ className, ...props }: React.ComponentProps<"svg">) {
  return (
    <output
      aria-live="polite"
      aria-busy="true"
      className="inline-flex items-center justify-center"
    >
      <Loader2Icon
        aria-label="Loading"
        className={cn("size-4 animate-spin", className)}
        {...props}
      />
    </output>
  );
}

export { Spinner };
