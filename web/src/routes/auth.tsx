import { Link } from "react-router-dom"
import { useTranslations } from "@/lib/i18n"
import { usePageMeta } from "@/lib/meta"
import { AuthCard } from "@/components/auth/auth-card"
import { ConfirmPasswordForm } from "@/components/auth/confirm-password-form"
import { ForgotPasswordForm } from "@/components/auth/forgot-password-form"
import { LoginForm } from "@/components/auth/login-form"
import { RegisterForm } from "@/components/auth/register-form"
import { ResetPasswordForm } from "@/components/auth/reset-password-form"
import { TwoFactorForm } from "@/components/auth/two-factor-form"
import { VerifyEmailCard } from "@/components/auth/verify-email-card"

function FooterLink({ to, children }: { to: string; children: React.ReactNode }) {
  return (
    <Link to={to} className="text-center font-medium text-primary underline-offset-4 hover:underline">
      {children}
    </Link>
  )
}

export function LoginPage() {
  const t = useTranslations("auth.login")
  usePageMeta(t("title"))
  return (
    <AuthCard
      title={t("title")}
      description={t("description")}
      footer={
        <p className="text-center">
          {t("noAccount")}{" "}
          <Link to="/register" className="font-medium text-primary underline-offset-4 hover:underline">
            {t("register")}
          </Link>
        </p>
      }
    >
      <LoginForm />
    </AuthCard>
  )
}

export function RegisterPage() {
  const t = useTranslations("auth.register")
  usePageMeta(t("title"))
  return (
    <AuthCard
      title={t("title")}
      description={t("description")}
      footer={
        <p className="text-center">
          {t("hasAccount")}{" "}
          <Link to="/login" className="font-medium text-primary underline-offset-4 hover:underline">
            {t("login")}
          </Link>
        </p>
      }
    >
      <RegisterForm />
    </AuthCard>
  )
}

export function ForgotPasswordPage() {
  const t = useTranslations("auth.forgot")
  usePageMeta(t("title"))
  return (
    <AuthCard title={t("title")} description={t("description")} footer={<FooterLink to="/login">{t("back")}</FooterLink>}>
      <ForgotPasswordForm />
    </AuthCard>
  )
}

export function ResetPasswordPage() {
  const t = useTranslations("auth.reset")
  usePageMeta(t("title"))
  return (
    <AuthCard title={t("title")} description={t("description")} footer={<FooterLink to="/login">{t("back")}</FooterLink>}>
      <ResetPasswordForm />
    </AuthCard>
  )
}

export function VerifyEmailPage() {
  return <VerifyEmailCard />
}

export function TwoFactorPage() {
  const t = useTranslations("auth.twoFactor")
  usePageMeta(t("title"))
  return (
    <AuthCard title={t("title")} description={t("description")} footer={<FooterLink to="/login">{t("backToLogin")}</FooterLink>}>
      <TwoFactorForm />
    </AuthCard>
  )
}

export function ConfirmPasswordPage() {
  const t = useTranslations("auth.confirmPassword")
  usePageMeta(t("title"))
  return (
    <AuthCard title={t("title")} description={t("description")}>
      <ConfirmPasswordForm />
    </AuthCard>
  )
}
