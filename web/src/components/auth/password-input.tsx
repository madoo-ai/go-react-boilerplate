import { useState } from "react"
import { EyeIcon as Eye, EyeOffIcon as EyeOff } from "@x1nx3r/altics-react-icons"
import { useTranslations } from "@/lib/i18n"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group"

export function PasswordInput({ ...props }: React.ComponentProps<typeof InputGroupInput>) {
  const t = useTranslations("common")
  const [visible, setVisible] = useState(false)

  return (
    <InputGroup>
      <InputGroupInput type={visible ? "text" : "password"} {...props} />
      <InputGroupAddon align="inline-end">
        <InputGroupButton
          type="button"
          size="icon-xs"
          onClick={() => setVisible((current) => !current)}
          aria-label={visible ? t("hidePassword") : t("showPassword")}
        >
          {visible ? <EyeOff /> : <Eye />}
        </InputGroupButton>
      </InputGroupAddon>
    </InputGroup>
  )
}
