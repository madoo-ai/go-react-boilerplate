import { expect, test } from "@playwright/test"

const BASE_URL = process.env.E2E_BASE_URL ?? "http://127.0.0.1:3000"

test.beforeEach(async ({ context }) => {
  await context.addCookies([
    { name: "locale", value: "en", url: BASE_URL },
    { name: "theme", value: "light", url: BASE_URL },
    { name: "color-theme", value: "default", url: BASE_URL },
  ])
})

async function login(page: import("@playwright/test").Page, email: string) {
  await page.goto("/login")
  await page.getByLabel("Email").fill(email)
  await page.getByLabel("Password", { exact: true }).fill("password")
  await page.getByRole("button", { name: "Sign in" }).click()
}

test.describe("public surface", () => {
  test("landing page renders and switches palette", async ({ page }) => {
    await page.goto("/")
    await expect(page.getByRole("heading", { level: 1 })).toContainText("Start with the parts")
    await expect(page.getByRole("link", { name: "Open the demo" })).toBeVisible()
    await page.getByRole("button", { name: "Dracula", exact: true }).click()
    await expect(page.locator("html")).toHaveAttribute("data-color-theme", "dracula")
  })

  test("prerendered login page has its title", async ({ page }) => {
    await page.goto("/login")
    await expect(page).toHaveTitle(/Sign in/)
    await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible()
  })
})

test.describe("authenticated surface", () => {
  test("superadmin sees the dashboard and admin pages", async ({ page }) => {
    await login(page, "superadmin@example.com")
    await expect(page).toHaveURL(/\/dashboard/)
    await expect(page.getByRole("heading", { name: "Dashboard" })).toBeVisible()
    await expect(page.getByText("Total users")).toBeVisible()

    await page.goto("/admin/users")
    await expect(page.getByRole("heading", { name: "Users" })).toBeVisible()
    await expect(page.getByText("superadmin@example.com")).toBeVisible()

    await page.goto("/admin/roles")
    await expect(page.getByRole("heading", { name: "Roles" })).toBeVisible()
    await expect(page.getByText("Super Admin", { exact: true }).first()).toBeVisible()

    await page.goto("/settings/profile")
    await expect(page.getByRole("heading", { name: "Settings" })).toBeVisible()
    await expect(page.getByLabel("Name")).toHaveValue("Super Admin")
  })

  test("HR cannot open role management", async ({ page }) => {
    await login(page, "hr@example.com")
    await expect(page).toHaveURL(/\/dashboard/)
    await page.goto("/admin/roles")
    await expect(page.getByText("You do not have permission to view this page.")).toBeVisible()
  })

  test("regular user gets the self-service dashboard", async ({ page }) => {
    await login(page, "user@example.com")
    await expect(page).toHaveURL(/\/dashboard/)
    await expect(page.getByText("Regular User").first()).toBeVisible()
    await expect(page.getByText("Total users")).not.toBeVisible()
  })
})

test.describe("auth gates", () => {
  test("protected routes redirect to login", async ({ page }) => {
    await page.goto("/admin/users")
    await expect(page).toHaveURL(/\/login\?next=/)
  })

  test("wrong password shows a translated error", async ({ page }) => {
    await page.goto("/login")
    await page.getByLabel("Email").fill("superadmin@example.com")
    await page.getByLabel("Password", { exact: true }).fill("wrong-password")
    await page.getByRole("button", { name: "Sign in" }).click()
    await expect(page.getByText("Invalid email or password.")).toBeVisible()
  })

  test("sign out returns to login", async ({ page }) => {
    await login(page, "superadmin@example.com")
    await expect(page).toHaveURL(/\/dashboard/)
    await page.getByRole("button", { name: /Super Admin/ }).click()
    await page.getByRole("menuitem", { name: "Sign out" }).click()
    await expect(page).toHaveURL(/\/login/)
  })
})
