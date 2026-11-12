import { useState } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { zodResolver } from "@hookform/resolvers/zod"
import { ShieldTickIcon as ShieldCheck } from "@x1nx3r/altics-react-icons"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { api, ApiError } from "@/api/client"
import { translateFieldError, useTranslations } from "@/lib/i18n"
import { FormAlert } from "@/components/auth/form-alert"
import { PasswordInput } from "@/components/auth/password-input"
import { Button } from "@/components/ui/button"
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Spinner } from "@/components/ui/spinner"

const schema = z.object({ password: z.string().min(1, "validation.required") })
type ConfirmValues = z.infer<typeof schema>

function safeNextPath(next: string | null) {
  if (!next) return "/settings/profile"
  if (!next.startsWith("/") || next.startsWith("//")) return "/settings/profile"
  return next
}

export function ConfirmPasswordForm() {
  const t = useTranslations("auth.confirmPassword")
  const tFields = useTranslations("auth.fields")
  const tValidation = useTranslations("validation")
  const tAll = useTranslations()
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const [formError, setFormError] = useState<string | null>(null)

  const form = useForm<ConfirmValues>({
    resolver: zodResolver(schema),
    defaultValues: { password: "" },
  })

  const fieldError = (message?: string) => translateFieldError(message, tValidation)

  async function onSubmit(values: ConfirmValues) {
    setFormError(null)
    try {
      await api.confirmPassword(values.password)
      navigate(safeNextPath(params.get("next")), { replace: true })
    } catch (error) {
      setFormError(error instanceof ApiError ? tAll(error.code) : tAll("errors.unexpected"))
    }
  }

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} noValidate>
      <FieldGroup className="gap-4">
        <FormAlert message={formError} />
        <Field data-invalid={!!form.formState.errors.password}>
          <FieldLabel htmlFor="password">{tFields("password")}</FieldLabel>
          <PasswordInput
            id="password"
            autoComplete="current-password"
            autoFocus
            aria-invalid={!!form.formState.errors.password}
            {...form.register("password")}
          />
          <FieldError>{fieldError(form.formState.errors.password?.message)}</FieldError>
        </Field>
        <Button type="submit" className="w-full" disabled={form.formState.isSubmitting}>
          {form.formState.isSubmitting ? <Spinner data-icon="inline-start" /> : <ShieldCheck data-icon="inline-start" />}
          {t("submit")}
        </Button>
      </FieldGroup>
    </form>
  )
}
