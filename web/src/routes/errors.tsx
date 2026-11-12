import { Link } from "react-router-dom"
import { Home01Icon as Home } from "@x1nx3r/altics-react-icons"
import { useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { Button } from "@/components/ui/button"

export function ErrorPage({ variant }: { variant: "notFound" | "forbidden" | "server" }) {
  const t = useTranslations("errorPages")
  const tCommon = useTranslations("common")
  usePageMeta(t(`${variant}.title`))

  return (
    <div className="grid min-h-svh place-items-center bg-background p-6">
      <div className="flex max-w-md flex-col items-center gap-5 text-center">
        <span className="grid size-14 place-items-center rounded-2xl bg-primary/10 text-primary ring-1 ring-primary/15">
          <Home />
        </span>
        <div className="flex flex-col gap-2">
          <h1 className="font-display text-4xl tracking-tight">{t(`${variant}.title`)}</h1>
          <p className="text-sm text-muted-foreground">{t(`${variant}.description`)}</p>
        </div>
        <div className="flex flex-wrap justify-center gap-2">
          <Button asChild>
            <Link to="/dashboard">{tCommon("backToDashboard")}</Link>
          </Button>
          <Button asChild variant="outline">
            <Link to="/">{tCommon("backToHome")}</Link>
          </Button>
        </div>
        {variant === "server" ? (
          <Button variant="ghost" onClick={() => window.location.reload()}>{t("server.retry")}</Button>
        ) : null}
      </div>
    </div>
  )
}
