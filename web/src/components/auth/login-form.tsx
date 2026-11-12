import { useState } from "react"
import { Link, useNavigate, useSearchParams } from "react-router-dom"
import { zodResolver } from "@hookform/resolvers/zod"
import { LogIn01Icon as LogIn } from "@x1nx3r/altics-react-icons"
import { useForm, useWatch } from "react-hook-form"
import { z } from "zod"
import { api, ApiError } from "@/api/client"
import { useAuth } from "@/lib/auth"
import { translateFieldError, useTranslations } from "@/lib/i18n"
import { FormAlert } from "@/components/auth/form-alert"
import { PasswordInput } from "@/components/auth/password-input"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Spinner } from "@/components/ui/spinner"

const schema = z.object({
  email: z.email("validation.email"),
  password: z.string().min(1, "validation.required"),
  rememberMe: z.boolean(),
})

type LoginValues = z.infer<typeof schema>

function safeNextPath(next: string | null) {
  if (!next) return "/dashboard"
  if (!next.startsWith("/") || next.startsWith("//")) return "/dashboard"
  return next
}

export function LoginForm() {
  const t = useTranslations("auth.login")
  const tFields = useTranslations("auth.fields")
  const tValidation = useTranslations("validation")
  const tAll = useTranslations()
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const { setSession } = useAuth()
  const [formError, setFormError] = useState<string | null>(null)

  const errorParam = params.get("error")
  const initialError =
    errorParam === "account_inactive"
      ? tAll("errors.accountInactive")
      : errorParam === "account_pending"
        ? tAll("auth.messages.pendingApproval")
        : null
  const notice = params.get("deleted")
    ? tAll("auth.messages.accountDeleted")
    : params.get("reset")
      ? tAll("auth.messages.passwordReset")
      : null

  const form = useForm<LoginValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: "", password: "", rememberMe: true },
  })

  const rememberMe = useWatch({ control: form.control, name: "rememberMe" })
  const fieldError = (message?: string) => translateFieldError(message, tValidation)

  async function onSubmit(values: LoginValues) {
    setFormError(initialError)
    try {
      const result = await api.login(values)
      if ("twoFactorRequired" in result) {
        navigate("/two-factor")
        return
      }
      setSession(result)
      navigate(safeNextPath(params.get("next")), { replace: true })
    } catch (error) {
      setFormError(error instanceof ApiError ? tAll(error.code) : tAll("errors.unexpected"))
    }
  }

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} noValidate>
      <FieldGroup className="gap-4">
        <FormAlert message={notice} variant="default" />
        <FormAlert message={formError ?? initialError} />
        <Field data-invalid={!!form.formState.errors.email}>
          <FieldLabel htmlFor="email">{tFields("email")}</FieldLabel>
          <Input
            id="email"
            type="email"
            autoComplete="email"
            autoFocus
            aria-invalid={!!form.formState.errors.email}
            {...form.register("email")}
          />
          <FieldError>{fieldError(form.formState.errors.email?.message)}</FieldError>
        </Field>
        <Field data-invalid={!!form.formState.errors.password}>
          <div className="flex items-center justify-between gap-2">
            <FieldLabel htmlFor="password">{tFields("password")}</FieldLabel>
            <Link
              to="/forgot-password"
              className="text-xs text-muted-foreground underline-offset-4 hover:text-primary hover:underline"
            >
              {t("forgot")}
            </Link>
          </div>
          <PasswordInput
            id="password"
            autoComplete="current-password"
            aria-invalid={!!form.formState.errors.password}
            {...form.register("password")}
          />
          <FieldError>{fieldError(form.formState.errors.password?.message)}</FieldError>
        </Field>
        <Field orientation="horizontal">
          <Checkbox
            id="rememberMe"
            checked={rememberMe}
            onCheckedChange={(checked) => form.setValue("rememberMe", checked === true)}
          />
          <FieldLabel htmlFor="rememberMe" className="font-normal">
            {t("remember")}
          </FieldLabel>
        </Field>
        <Button type="submit" className="w-full" disabled={form.formState.isSubmitting}>
          {form.formState.isSubmitting ? <Spinner data-icon="inline-start" /> : <LogIn data-icon="inline-start" />}
          {t("submit")}
        </Button>
      </FieldGroup>
    </form>
  )
}
