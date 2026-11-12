import { expect, test } from "@playwright/test"

const headers = { "X-UI-Request": "1" }

test("registration approval blocks and then admits a new account", async ({ request }) => {
  const email = `approval-${Date.now()}@example.com`

  const login = await request.post("/api/auth/login", {
    headers,
    data: { email: "superadmin@example.com", password: "password", rememberMe: true },
  })
  expect(login.ok()).toBeTruthy()

  try {
    const enable = await request.patch("/api/admin/registration-approval", {
      headers,
      data: { enabled: true },
    })
    expect(enable.ok()).toBeTruthy()

    const register = await request.post("/api/auth/register", {
      headers,
      data: { name: "Approval Test", email, password: "password" },
    })
    expect(register.status()).toBe(201)
    expect((await register.json()).status).toBe("pending")

    const blocked = await request.post("/api/auth/login", {
      headers,
      data: { email, password: "password", rememberMe: true },
    })
    expect(blocked.status()).toBe(403)
    expect((await blocked.json()).error).toBe("auth.messages.pendingApproval")

    const list = await request.get(`/api/users?search=${encodeURIComponent(email)}`)
    expect(list.ok()).toBeTruthy()
    const found = (await list.json()).users.find((user: { email: string }) => user.email === email)
    expect(found).toBeTruthy()

    const verify = await request.post(`/api/users/${found.id}/verify`, { headers })
    expect(verify.ok()).toBeTruthy()

    const admitted = await request.post("/api/auth/login", {
      headers,
      data: { email, password: "password", rememberMe: true },
    })
    expect(admitted.ok()).toBeTruthy()
    expect((await admitted.json()).user.email).toBe(email)
  } finally {
    // The admitted user's login replaced the API cookie jar. Sign back in as
    // the superadmin before the cleanup writes.
    await request.post("/api/auth/login", {
      headers,
      data: { email: "superadmin@example.com", password: "password", rememberMe: true },
    })
    await request.patch("/api/admin/registration-approval", {
      headers,
      data: { enabled: false },
    })
    const list = await request.get(`/api/users?search=${encodeURIComponent(email)}`)
    if (list.ok()) {
      const found = (await list.json()).users.find((user: { email: string }) => user.email === email)
      if (found) await request.delete(`/api/users/${found.id}`, { headers })
    }
  }
})
