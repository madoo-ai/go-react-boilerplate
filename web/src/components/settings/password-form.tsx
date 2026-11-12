import { useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { Key01Icon as KeyRound } from "@x1nx3r/altics-react-icons"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"
import { api, ApiError } from "@/api/client"
import { translateFieldError, useTranslations } from "@/lib/i18n"
import { FormAlert } from "@/components/auth/form-alert"
import { PasswordInput } from "@/components/auth/password-input"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Spinner } from "@/components/ui/spinner"

const schema = z
  .object({
    currentPassword: z.string().min(1, "validation.required"),
    newPassword: z.string().min(8, "validation.passwordLength").max(128, "validation.passwordLength"),
    confirmPassword: z.string().min(1, "validation.required"),
  })
  .refine((values) => values.newPassword === values.confirmPassword, {
    path: ["confirmPassword"],
    message: "validation.passwordMismatch",
  })

type PasswordValues = z.infer<typeof schema>

export function PasswordForm() {
  const t = useTranslations("settings.password")
  const tValidation = useTranslations("validation")
  const tAll = useTranslations()
  const [formError, setFormError] = useState<string | null>(null)

  const form = useForm<PasswordValues>({
    resolver: zodResolver(schema),
    defaultValues: { currentPassword: "", newPassword: "", confirmPassword: "" },
  })

  const fieldError = (message?: string) => translateFieldError(message, tValidation)

  async function onSubmit(values: PasswordValues) {
    setFormError(null)
    try {
      await api.changePassword({
        currentPassword: values.currentPassword,
        newPassword: values.newPassword,
      })
      form.reset()
      toast.success(t("saved"))
    } catch (error) {
      setFormError(error instanceof ApiError ? tAll(error.code) : tAll("errors.unexpected"))
    }
  }

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} noValidate>
      <Card>
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
          <CardDescription>{t("description")}</CardDescription>
        </CardHeader>
        <CardContent>
          <FieldGroup className="max-w-xl">
            <FormAlert message={formError} />
            <Field data-invalid={!!form.formState.errors.currentPassword}>
              <FieldLabel htmlFor="current-password">{t("current")}</FieldLabel>
              <PasswordInput id="current-password" autoComplete="current-password" aria-invalid={!!form.formState.errors.currentPassword} {...form.register("currentPassword")} />
              <FieldError>{fieldError(form.formState.errors.currentPassword?.message)}</FieldError>
            </Field>
            <Field data-invalid={!!form.formState.errors.newPassword}>
              <FieldLabel htmlFor="new-password">{t("new")}</FieldLabel>
              <PasswordInput id="new-password" autoComplete="new-password" aria-invalid={!!form.formState.errors.newPassword} {...form.register("newPassword")} />
              <FieldError>{fieldError(form.formState.errors.newPassword?.message)}</FieldError>
            </Field>
            <Field data-invalid={!!form.formState.errors.confirmPassword}>
              <FieldLabel htmlFor="confirm-password">{t("confirm")}</FieldLabel>
              <PasswordInput id="confirm-password" autoComplete="new-password" aria-invalid={!!form.formState.errors.confirmPassword} {...form.register("confirmPassword")} />
              <FieldError>{fieldError(form.formState.errors.confirmPassword?.message)}</FieldError>
            </Field>
          </FieldGroup>
        </CardContent>
        <CardFooter className="justify-end border-t">
          <Button type="submit" disabled={form.formState.isSubmitting}>
            {form.formState.isSubmitting ? <Spinner data-icon="inline-start" /> : <KeyRound data-icon="inline-start" />}
            {t("save")}
          </Button>
        </CardFooter>
      </Card>
    </form>
  )
}
