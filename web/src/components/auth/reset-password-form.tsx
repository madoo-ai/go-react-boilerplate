import { useState } from "react"
import { Link, useNavigate, useSearchParams } from "react-router-dom"
import { zodResolver } from "@hookform/resolvers/zod"
import { Key01Icon as KeyRound } from "@x1nx3r/altics-react-icons"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { api, ApiError } from "@/api/client"
import { translateFieldError, useTranslations } from "@/lib/i18n"
import { FormAlert } from "@/components/auth/form-alert"
import { PasswordInput } from "@/components/auth/password-input"
import { Button } from "@/components/ui/button"
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Spinner } from "@/components/ui/spinner"

const schema = z
  .object({
    password: z.string().min(8, "validation.passwordLength").max(128, "validation.passwordLength"),
    confirmPassword: z.string().min(1, "validation.required"),
  })
  .refine((values) => values.password === values.confirmPassword, {
    path: ["confirmPassword"],
    message: "validation.passwordMismatch",
  })

type ResetValues = z.infer<typeof schema>

export function ResetPasswordForm() {
  const t = useTranslations("auth.reset")
  const tFields = useTranslations("auth.fields")
  const tValidation = useTranslations("validation")
  const tAll = useTranslations()
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const token = params.get("token") ?? undefined
  const [formError, setFormError] = useState<string | null>(null)

  const form = useForm<ResetValues>({
    resolver: zodResolver(schema),
    defaultValues: { password: "", confirmPassword: "" },
  })

  const fieldError = (message?: string) => translateFieldError(message, tValidation)

  async function onSubmit(values: ResetValues) {
    if (!token) return
    setFormError(null)
    try {
      await api.resetPassword(token, values.password)
      navigate("/login?reset=1", { replace: true })
    } catch (error) {
      setFormError(error instanceof ApiError ? tAll(error.code) : t("invalid"))
    }
  }

  if (!token) {
    return (
      <div className="flex flex-col gap-4">
        <FormAlert message={t("invalid")} />
        <Button asChild variant="outline" className="w-full">
          <Link to="/forgot-password">{t("requestNew")}</Link>
        </Button>
      </div>
    )
  }

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} noValidate>
      <FieldGroup className="gap-4">
        <FormAlert message={formError} />
        <Field data-invalid={!!form.formState.errors.password}>
          <FieldLabel htmlFor="password">{tFields("password")}</FieldLabel>
          <PasswordInput id="password" autoComplete="new-password" autoFocus aria-invalid={!!form.formState.errors.password} {...form.register("password")} />
          <FieldError>{fieldError(form.formState.errors.password?.message)}</FieldError>
        </Field>
        <Field data-invalid={!!form.formState.errors.confirmPassword}>
          <FieldLabel htmlFor="confirmPassword">{tFields("confirmPassword")}</FieldLabel>
          <PasswordInput id="confirmPassword" autoComplete="new-password" aria-invalid={!!form.formState.errors.confirmPassword} {...form.register("confirmPassword")} />
          <FieldError>{fieldError(form.formState.errors.confirmPassword?.message)}</FieldError>
        </Field>
        <Button type="submit" className="w-full" disabled={form.formState.isSubmitting}>
          {form.formState.isSubmitting ? <Spinner data-icon="inline-start" /> : <KeyRound data-icon="inline-start" />}
          {t("submit")}
        </Button>
      </FieldGroup>
    </form>
  )
}
