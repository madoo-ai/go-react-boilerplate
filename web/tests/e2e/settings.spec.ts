import { expect, test } from "@playwright/test"
import { totpCode } from "./totp"

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

test("profile update is reflected in the session", async ({ page }) => {
  await page.goto("/settings/profile")
  const jobTitle = page.getByLabel("Job title")
  const original = await jobTitle.inputValue()
  await jobTitle.fill("E2E Operator")
  await page.getByRole("button", { name: "Save profile" }).click()
  await expect(page.getByText("Profile updated.")).toBeVisible()
  await page.reload()
  await expect(page.getByLabel("Job title")).toHaveValue("E2E Operator")

  await page.getByLabel("Job title").fill(original)
  await page.getByRole("button", { name: "Save profile" }).click()
  await expect(page.getByText("Profile updated.")).toBeVisible()
})

test("enable and disable TOTP two-factor authentication", async ({ page }) => {
  await page.goto("/settings/security")
  await page.getByRole("button", { name: "Enable 2FA" }).click()
  const enableDialog = page.getByRole("dialog")
  await enableDialog.getByLabel("Confirm your password.").fill("password")
  await enableDialog.getByRole("button", { name: "Enable 2FA" }).click()

  const secret = await enableDialog.locator("code").textContent()
  expect(secret).toBeTruthy()

  const input = enableDialog.locator("#totp-code")
  await input.click()
  await page.keyboard.type(totpCode(secret!.trim()))
  await enableDialog.getByRole("button", { name: "Verify code" }).click()
  await expect(page.getByText("Two-factor authentication enabled.")).toBeVisible()
  await expect(page.getByText("Enabled", { exact: true })).toBeVisible()

  await page.getByRole("button", { name: "Disable 2FA" }).click()
  const disableDialog = page.getByRole("dialog")
  await disableDialog.getByLabel("Confirm your password.").fill("password")
  await disableDialog.getByRole("button", { name: "Disable 2FA" }).click()
  await expect(page.getByText("Two-factor authentication disabled.")).toBeVisible()
})
