import { useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { Pencil01Icon as Pencil, PlusIcon as Plus, Save01Icon as Save } from "@x1nx3r/altics-react-icons"
import { useForm, useWatch } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"
import { api, ApiError } from "@/api/client"
import type { PermissionInfo, Role } from "@/api/types"
import { translateFieldError, useTranslations } from "@/lib/i18n"
import { FormAlert } from "@/components/auth/form-alert"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Field, FieldDescription, FieldError, FieldGroup, FieldLabel, FieldLegend, FieldSet } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Spinner } from "@/components/ui/spinner"
import { Textarea } from "@/components/ui/textarea"

const roleSchema = z.object({
  name: z.string().trim().min(2, "validation.required").max(60, "validation.max").regex(/^[a-z0-9_]+$/, "validation.roleName"),
  label: z.string().trim().min(1, "validation.required").max(80, "validation.max"),
  description: z.string().trim().max(240).optional().or(z.literal("")),
  permissions: z.array(z.string()),
})

type RoleFormValues = z.infer<typeof roleSchema>

function RoleDialog({
  mode,
  role,
  permissions,
  trigger,
  onChanged,
}: {
  mode: "create" | "edit"
  role?: Role
  permissions: PermissionInfo[]
  trigger: React.ReactNode
  onChanged: () => void
}) {
  const t = useTranslations("roles")
  const tForm = useTranslations("roles.form")
  const tActions = useTranslations("roles.actions")
  const tCommon = useTranslations("common")
  const tValidation = useTranslations("validation")
  const tPermissionNames = useTranslations("permissions")
  const tAll = useTranslations()
  const [open, setOpen] = useState(false)
  const [formError, setFormError] = useState<string | null>(null)

  const form = useForm<RoleFormValues>({
    resolver: zodResolver(roleSchema),
    defaultValues: {
      name: role?.name ?? "",
      label: role?.label ?? "",
      description: role?.description ?? "",
      permissions: role?.permissions ?? [],
    },
  })

  const fieldError = (message?: string) => translateFieldError(message, tValidation)
  const selectedPermissions = useWatch({ control: form.control, name: "permissions" }) ?? []
  const lockName = role?.name === "user"

  async function onSubmit(values: RoleFormValues) {
    setFormError(null)
    const input = {
      name: values.name,
      label: values.label,
      description: values.description ?? "",
      permissions: values.permissions,
    }
    try {
      if (mode === "create") {
        await api.createRole(input)
        toast.success(t("messages.created"))
      } else if (role) {
        await api.updateRole(role.id, input)
        toast.success(t("messages.updated"))
      }
      setOpen(false)
      onChanged()
    } catch (error) {
      setFormError(error instanceof ApiError ? tAll(error.code) : tAll("errors.unexpected"))
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle>{mode === "create" ? t("create") : t("edit")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>
        <form className="flex flex-col gap-4" onSubmit={form.handleSubmit(onSubmit)} noValidate>
          <FormAlert message={formError} />
          <FieldGroup>
            <div className="grid gap-4 sm:grid-cols-2">
              <Field data-invalid={!!form.formState.errors.name}>
                <FieldLabel htmlFor="role-name">{tForm("name")}</FieldLabel>
                <Input id="role-name" disabled={lockName} aria-invalid={!!form.formState.errors.name} {...form.register("name")} />
                {lockName ? <FieldDescription>{tForm("nameHint")}</FieldDescription> : null}
                <FieldError>{fieldError(form.formState.errors.name?.message)}</FieldError>
              </Field>
              <Field data-invalid={!!form.formState.errors.label}>
                <FieldLabel htmlFor="role-label">{tForm("label")}</FieldLabel>
                <Input id="role-label" aria-invalid={!!form.formState.errors.label} {...form.register("label")} />
                <FieldError>{fieldError(form.formState.errors.label?.message)}</FieldError>
              </Field>
            </div>
            <Field>
              <FieldLabel htmlFor="role-description">{tForm("description")}</FieldLabel>
              <Textarea id="role-description" rows={3} {...form.register("description")} />
            </Field>
            <FieldSet>
              <FieldLegend>{tForm("permissions")}</FieldLegend>
              <FieldGroup>
                {permissions.map((permission) => {
                  const checked = selectedPermissions.includes(permission.key)
                  return (
                    <Field key={permission.key} orientation="horizontal">
                      <Checkbox
                        id={`permission-${permission.key}`}
                        checked={checked}
                        onCheckedChange={(value) => {
                          const next = value
                            ? [...selectedPermissions, permission.key]
                            : selectedPermissions.filter((key) => key !== permission.key)
                          form.setValue("permissions", next, { shouldValidate: true })
                        }}
                      />
                      <FieldLabel htmlFor={`permission-${permission.key}`} className="font-normal">
                        {tPermissionNames(permission.key)}
                      </FieldLabel>
                    </Field>
                  )
                })}
              </FieldGroup>
            </FieldSet>
          </FieldGroup>
          <DialogFooter>
            <DialogClose asChild>
              <Button type="button" variant="outline">{tCommon("cancel")}</Button>
            </DialogClose>
            <Button type="submit" disabled={form.formState.isSubmitting}>
              {form.formState.isSubmitting ? <Spinner data-icon="inline-start" /> : <Save data-icon="inline-start" />}
              {mode === "create" ? tActions("create") : tActions("save")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export function CreateRoleDialog({ permissions, onChanged }: { permissions: PermissionInfo[]; onChanged: () => void }) {
  const t = useTranslations("roles")
  return (
    <RoleDialog
      mode="create"
      permissions={permissions}
      onChanged={onChanged}
      trigger={
        <Button>
          <Plus data-icon="inline-start" />
          {t("create")}
        </Button>
      }
    />
  )
}

export function EditRoleDialog({ role, permissions, onChanged }: { role: Role; permissions: PermissionInfo[]; onChanged: () => void }) {
  const tCommon = useTranslations("common")
  return (
    <RoleDialog
      mode="edit"
      role={role}
      permissions={permissions}
      onChanged={onChanged}
      trigger={
        <Button variant="ghost" size="icon-sm" aria-label={tCommon("edit")}>
          <Pencil />
        </Button>
      }
    />
  )
}
