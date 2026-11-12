import { cn } from "@/lib/utils"
import { Loading03Icon as Loader2Icon } from "@x1nx3r/altics-react-icons";

function Spinner({ className, ...props }: React.ComponentProps<"svg">) {
  return (
    <Loader2Icon data-slot="spinner" role="status" aria-label="Loading" className={cn("size-4 animate-spin", className)} {...props} />
  )
}

export { Spinner }
