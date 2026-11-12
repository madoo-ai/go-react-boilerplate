import type {
  AccountOverview,
  DashboardResponse,
  LoginResponse,
  ProfileInput,
  RoleInput,
  RolesResponse,
  SessionResponse,
  TwoFactorSetup,
  UserEditResponse,
  UserInput,
  UsersResponse,
} from "./types"

export class ApiError extends Error {
  code: string
  fields: Record<string, string>
  status: number

  constructor(code: string, status: number, fields: Record<string, string> = {}) {
    super(code)
    this.name = "ApiError"
    this.code = code
    this.status = status
    this.fields = fields
  }
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const method = init.method ?? "GET"
  const headers = new Headers(init.headers)
  if (method !== "GET" && method !== "HEAD") {
    headers.set("X-UI-Request", "1")
  }
  if (init.body) {
    headers.set("Content-Type", "application/json")
  }
  const response = await fetch(path, {
    ...init,
    headers,
    credentials: "same-origin",
  })
  if (response.status === 204) return undefined as T
  const body = await response.json().catch(() => ({}))
  if (!response.ok) {
    throw new ApiError(
      typeof body.error === "string" ? body.error : "errors.unexpected",
      response.status,
      body.fields ?? {},
    )
  }
  return body as T
}

function query(params: Record<string, string | number | undefined>) {
  const search = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== "") search.set(key, String(value))
  }
  const value = search.toString()
  return value ? `?${value}` : ""
}

export const api = {
  session: async (): Promise<SessionResponse | null> => {
    try {
      return await request<SessionResponse>("/api/auth/session")
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) return null
      throw error
    }
  },
  register: (body: { name: string; email: string; password: string }) =>
    request<{ status: string; emailSent: boolean }>("/api/auth/register", {
      method: "POST",
      body: JSON.stringify(body),
    }),
  login: (body: { email: string; password: string; rememberMe: boolean }) =>
    request<LoginResponse>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify(body),
    }),
  logout: () => request<void>("/api/auth/logout", { method: "POST" }),
  verifyEmail: (token: string) =>
    request<{ ok: boolean }>("/api/auth/verify-email", {
      method: "POST",
      body: JSON.stringify({ token }),
    }),
  requestVerification: () =>
    request<{ ok: boolean }>("/api/auth/verify-email/request", { method: "POST" }),
  forgotPassword: (email: string) =>
    request<{ ok: boolean }>("/api/auth/forgot-password", {
      method: "POST",
      body: JSON.stringify({ email }),
    }),
  resetPassword: (token: string, password: string) =>
    request<{ ok: boolean }>("/api/auth/reset-password", {
      method: "POST",
      body: JSON.stringify({ token, password }),
    }),
  verifyTwoFactor: (code: string, backup: boolean) =>
    request<SessionResponse>("/api/auth/two-factor/verify", {
      method: "POST",
      body: JSON.stringify({ code, backup }),
    }),
  confirmPassword: (password: string) =>
    request<{ ok: boolean }>("/api/auth/confirm-password", {
      method: "POST",
      body: JSON.stringify({ password }),
    }),
  dashboard: () => request<DashboardResponse>("/api/dashboard"),
  users: (params: { page?: number; search?: string; role?: string; status?: string }) =>
    request<UsersResponse>(`/api/users${query(params)}`),
  createUser: (body: UserInput) =>
    request<{ ok: boolean }>("/api/users", { method: "POST", body: JSON.stringify(body) }),
  user: (id: string) => request<UserEditResponse>(`/api/users/${id}`),
  updateUser: (id: string, body: UserInput) =>
    request<{ ok: boolean }>(`/api/users/${id}`, { method: "PATCH", body: JSON.stringify(body) }),
  deleteUser: (id: string) => request<void>(`/api/users/${id}`, { method: "DELETE" }),
  verifyUser: (id: string) =>
    request<{ ok: boolean }>(`/api/users/${id}/verify`, { method: "POST" }),
  setRegistrationApproval: (enabled: boolean) =>
    request<{ ok: boolean }>("/api/admin/registration-approval", {
      method: "PATCH",
      body: JSON.stringify({ enabled }),
    }),
  roles: () => request<RolesResponse>("/api/roles"),
  createRole: (body: RoleInput) =>
    request<{ ok: boolean }>("/api/roles", { method: "POST", body: JSON.stringify(body) }),
  updateRole: (id: string, body: RoleInput) =>
    request<{ ok: boolean }>(`/api/roles/${id}`, { method: "PATCH", body: JSON.stringify(body) }),
  deleteRole: (id: string) => request<void>(`/api/roles/${id}`, { method: "DELETE" }),
  accountOverview: () => request<AccountOverview>("/api/settings/account"),
  updateProfile: (body: ProfileInput) =>
    request<SessionResponse>("/api/settings/profile", {
      method: "PATCH",
      body: JSON.stringify(body),
    }),
  changePassword: (body: { currentPassword: string; newPassword: string }) =>
    request<{ ok: boolean }>("/api/settings/password", {
      method: "POST",
      body: JSON.stringify(body),
    }),
  enableTwoFactor: (password: string) =>
    request<TwoFactorSetup>("/api/settings/two-factor/enable", {
      method: "POST",
      body: JSON.stringify({ password }),
    }),
  confirmTwoFactor: (code: string) =>
    request<{ ok: boolean }>("/api/settings/two-factor/confirm", {
      method: "POST",
      body: JSON.stringify({ code }),
    }),
  disableTwoFactor: (password: string) =>
    request<{ ok: boolean }>("/api/settings/two-factor/disable", {
      method: "POST",
      body: JSON.stringify({ password }),
    }),
  deleteAccount: () => request<void>("/api/settings/account", { method: "DELETE" }),
}
