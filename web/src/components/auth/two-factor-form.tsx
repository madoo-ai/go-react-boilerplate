import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { Key01Icon as KeyRound, ShieldTickIcon as ShieldCheck } from "@x1nx3r/altics-react-icons"
import { REGEXP_ONLY_DIGITS } from "input-otp"
import { api, ApiError } from "@/api/client"
import { useAuth } from "@/lib/auth"
import { useTranslations } from "@/lib/i18n"
import { FormAlert } from "@/components/auth/form-alert"
import { Button } from "@/components/ui/button"
import { Field, FieldDescription, FieldGroup, FieldLabel } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { InputOTP, InputOTPGroup, InputOTPSlot } from "@/components/ui/input-otp"
import { Spinner } from "@/components/ui/spinner"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

export function TwoFactorForm() {
  const t = useTranslations("auth.twoFactor")
  const tFields = useTranslations("auth.fields")
  const tAll = useTranslations()
  const navigate = useNavigate()
  const { setSession } = useAuth()
  const [code, setCode] = useState("")
  const [backupCode, setBackupCode] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [pending, setPending] = useState<"totp" | "backup" | null>(null)

  async function verify(value: string, backup: boolean, kind: "totp" | "backup") {
    setError(null)
    setPending(kind)
    try {
      const session = await api.verifyTwoFactor(value, backup)
      setSession(session)
      navigate("/dashboard", { replace: true })
    } catch (caught) {
      setError(caught instanceof ApiError ? tAll(caught.code) : tAll("errors.invalidCode"))
      setPending(null)
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <FormAlert message={error} />
      <Tabs defaultValue="totp">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="totp">{t("totpTab")}</TabsTrigger>
          <TabsTrigger value="backup">{t("backupTab")}</TabsTrigger>
        </TabsList>
        <TabsContent value="totp">
          <form
            className="pt-4"
            onSubmit={(event) => {
              event.preventDefault()
              void verify(code, false, "totp")
            }}
          >
            <FieldGroup className="gap-4">
              <Field>
                <FieldLabel htmlFor="totp-code">{tFields("code")}</FieldLabel>
                <InputOTP
                  id="totp-code"
                  maxLength={6}
                  pattern={REGEXP_ONLY_DIGITS}
                  value={code}
                  onChange={setCode}
                  autoFocus
                >
                  <InputOTPGroup>
                    {Array.from({ length: 6 }).map((_, index) => (
                      <InputOTPSlot key={index} index={index} />
                    ))}
                  </InputOTPGroup>
                </InputOTP>
              </Field>
              <Button type="submit" className="w-full" disabled={code.length !== 6 || pending !== null}>
                {pending === "totp" ? <Spinner data-icon="inline-start" /> : <ShieldCheck data-icon="inline-start" />}
                {t("submit")}
              </Button>
            </FieldGroup>
          </form>
        </TabsContent>
        <TabsContent value="backup">
          <form
            className="pt-4"
            onSubmit={(event) => {
              event.preventDefault()
              void verify(backupCode, true, "backup")
            }}
          >
            <FieldGroup className="gap-4">
              <Field>
                <FieldLabel htmlFor="backup-code">{t("backupTitle")}</FieldLabel>
                <Input
                  id="backup-code"
                  autoComplete="one-time-code"
                  value={backupCode}
                  onChange={(event) => setBackupCode(event.target.value)}
                  autoFocus
                />
                <FieldDescription>{t("backupDescription")}</FieldDescription>
              </Field>
              <Button type="submit" className="w-full" disabled={backupCode.length === 0 || pending !== null}>
                {pending === "backup" ? <Spinner data-icon="inline-start" /> : <KeyRound data-icon="inline-start" />}
                {t("backupSubmit")}
              </Button>
            </FieldGroup>
          </form>
        </TabsContent>
      </Tabs>
    </div>
  )
}
