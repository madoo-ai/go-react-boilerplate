import { useState } from "react"
import { CheckCircleIcon as CheckCircle2, Shield01Icon as ShieldAlert } from "@x1nx3r/altics-react-icons"
import { toast } from "sonner"
import { api, ApiError } from "@/api/client"
import { useTranslations } from "@/lib/i18n"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Field, FieldContent, FieldLabel } from "@/components/ui/field"
import { Switch } from "@/components/ui/switch"

export function RegistrationApprovalCard({ enabled, onChanged }: { enabled: boolean; onChanged: () => void }) {
  const t = useTranslations("users.approval")
  const tMessages = useTranslations("users.messages")
  const tAll = useTranslations()
  const [value, setValue] = useState(enabled)
  const [pending, setPending] = useState(false)

  async function update(next: boolean) {
    setValue(next)
    setPending(true)
    try {
      await api.setRegistrationApproval(next)
      toast.success(tMessages("approvalUpdated"))
      onChanged()
    } catch (error) {
      setValue(enabled)
      toast.error(error instanceof ApiError ? tAll(error.code) : tAll("errors.unexpected"))
    } finally {
      setPending(false)
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          {value ? <ShieldAlert /> : <CheckCircle2 />}
          {t("title")}
        </CardTitle>
        <CardDescription>{t("description")}</CardDescription>
      </CardHeader>
      <CardContent>
        <Field orientation="horizontal">
          <Switch
            id="registration-approval"
            checked={value}
            disabled={pending}
            onCheckedChange={(checked) => void update(checked)}
          />
          <FieldContent>
            <FieldLabel htmlFor="registration-approval">{t("toggle")}</FieldLabel>
            <p className="text-sm text-muted-foreground">{value ? t("enabled") : t("disabled")}</p>
          </FieldContent>
        </Field>
      </CardContent>
    </Card>
  )
}
