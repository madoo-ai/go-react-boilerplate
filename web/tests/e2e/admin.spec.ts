import { expect, test } from "@playwright/test"

const BASE_URL = process.env.E2E_BASE_URL ?? "http://127.0.0.1:3000"

test.beforeEach(async ({ context, page }) => {
  await context.addCookies([
    { name: "locale", value: "en", url: BASE_URL },
    { name: "theme", value: "light", url: BASE_URL },
    { name: "color-theme", value: "default", url: BASE_URL },
  ])
  await page.goto("/login")
  await page.getByLabel("Email").fill("superadmin@example.com")
  await page.getByLabel("Password", { exact: true }).fill("password")
  await page.getByRole("button", { name: "Sign in" }).click()
  await expect(page).toHaveURL(/\/dashboard/)
})

test("create, verify and delete a user", async ({ page }) => {
  const email = `e2e-${Date.now()}@example.com`

  await page.goto("/admin/users/create")
  await page.getByLabel("Name").fill("E2E Person")
  await page.getByLabel("Email").fill(email)
  await page.getByLabel("Phone").fill("+62 812 0000 0000")
  await page.getByLabel("Job title").fill("QA Engineer")
  await page.getByLabel("Department").fill("Engineering")
  await page.getByLabel("Status").click()
  await page.getByRole("option", { name: "Pending" }).click()
  await page.getByLabel("Password", { exact: true }).fill("password")
  await page.getByRole("button", { name: "Create user" }).click()

  await expect(page).toHaveURL(/\/admin\/users$/)
  await expect(page.getByText(email)).toBeVisible()

  const row = page.getByRole("row", { name: new RegExp(email) })
  await row.getByRole("button", { name: "Open menu" }).click()
  await page.getByRole("menuitem", { name: "Verify" }).click()
  await expect(page.getByText("User verified.")).toBeVisible()

  await row.getByRole("button", { name: "Open menu" }).click()
  await page.getByRole("menuitem", { name: "Delete" }).click()
  await page.getByRole("alertdialog").getByRole("button", { name: "Delete user" }).click()
  await expect(page.getByText("User deleted.")).toBeVisible()
  await expect(page.getByText(email)).not.toBeVisible()
})

test("create, edit and delete a role", async ({ page }) => {
  const name = `e2e_role_${Date.now()}`

  await page.goto("/admin/roles")
  await page.getByRole("button", { name: "Create role" }).first().click()
  const createDialog = page.getByRole("dialog")
  await createDialog.getByLabel("Name").fill(name)
  await createDialog.getByLabel("Label").fill("E2E Role")
  await createDialog.getByLabel("Description").fill("Created by the E2E suite.")
  await createDialog.getByRole("checkbox", { name: "Admin access" }).click()
  await createDialog.getByRole("button", { name: "Create role" }).click()
  await expect(page.getByText("Role created.")).toBeVisible()
  await expect(page.getByText(name)).toBeVisible()

  const row = page.getByRole("row", { name: new RegExp(name) })
  await row.getByRole("button", { name: "Edit" }).click()
  const editDialog = page.getByRole("dialog")
  await editDialog.getByLabel("Label").fill("E2E Role Updated")
  await editDialog.getByRole("button", { name: "Save changes" }).click()
  await expect(page.getByText("Role updated.")).toBeVisible()
  await expect(page.getByText("E2E Role Updated").first()).toBeVisible()

  const updatedRow = page.getByRole("row", { name: new RegExp(name) })
  await updatedRow.getByRole("button", { name: "Delete" }).click()
  await page.getByRole("alertdialog").getByRole("button", { name: "Delete role" }).click()
  await expect(page.getByText("Role deleted.")).toBeVisible()
  await expect(page.getByText(name)).not.toBeVisible()
})

test("registration approval toggle writes the setting", async ({ page }) => {
  await page.goto("/admin/users")
  const toggle = page.getByRole("switch", { name: "Require approval" })
  const initial = await toggle.isChecked()
  await toggle.click()
  await expect(page.getByText("Registration approval updated.").first()).toBeVisible()
  await expect(toggle).toBeChecked({ checked: !initial })
  await toggle.click()
  await expect(page.getByText("Registration approval updated.").first()).toBeVisible()
  await expect(toggle).toBeChecked({ checked: initial })
})
