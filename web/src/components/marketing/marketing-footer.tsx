import { Link } from "react-router-dom"
import { AppLogo } from "@/components/app-shell/app-logo"
import { useTranslations } from "@/lib/i18n"

export function MarketingFooter() {
  const t = useTranslations("landing.footer")

  return (
    <footer className="border-t border-border/60">
      <div className="mx-auto grid max-w-7xl gap-10 px-6 py-14 md:grid-cols-[1.2fr_1fr]">
        <div className="flex flex-col gap-4">
          <AppLogo />
          <p className="max-w-xs text-sm leading-relaxed text-muted-foreground">{t("tagline")}</p>
        </div>
        <div className="grid grid-cols-2 gap-8 text-sm">
          <div className="flex flex-col gap-3">
            <p className="font-medium">{t("product")}</p>
            <a href="#features" className="text-muted-foreground transition-colors hover:text-foreground">{t("features")}</a>
            <a href="#themes" className="text-muted-foreground transition-colors hover:text-foreground">{t("themes")}</a>
            <a href="#stack" className="text-muted-foreground transition-colors hover:text-foreground">{t("stack")}</a>
          </div>
          <div className="flex flex-col gap-3">
            <p className="font-medium">{t("access")}</p>
            <Link to="/login" className="text-muted-foreground transition-colors hover:text-foreground">{t("signIn")}</Link>
            <Link to="/register" className="text-muted-foreground transition-colors hover:text-foreground">{t("create")}</Link>
          </div>
        </div>
      </div>
      <div className="border-t border-border/60">
        <div className="mx-auto flex max-w-7xl flex-col gap-2 px-6 py-6 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
          <p>{t("built")}</p>
          <p>{t("license")}</p>
        </div>
      </div>
    </footer>
  )
}
