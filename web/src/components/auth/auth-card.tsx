import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
} from "@/components/ui/card"

export function AuthCard({
  title,
  description,
  children,
  footer,
}: {
  title: string
  description?: string
  children: React.ReactNode
  footer?: React.ReactNode
}) {
  return (
    <Card className="relative gap-4 overflow-hidden">
      <span
        aria-hidden
        className="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-brand-indigo via-brand-magenta to-brand-ruby"
      />
      <CardHeader className="pt-5">
        <h1 className="font-display text-3xl leading-tight tracking-tight">{title}</h1>
        {description ? <CardDescription>{description}</CardDescription> : null}
      </CardHeader>
      <CardContent>{children}</CardContent>
      {footer ? (
        <CardFooter className="flex flex-col gap-2 border-t border-border/60 text-sm text-muted-foreground">
          {footer}
        </CardFooter>
      ) : null}
    </Card>
  )
}
