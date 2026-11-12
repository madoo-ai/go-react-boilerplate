export type UserStatus = "pending" | "active" | "suspended"

export type Role = {
  id: string
  name: string
  label: string
  description: string | null
  isSystem: boolean
  permissions: string[]
  userCount: number
}

export type User = {
  id: string
  name: string
  email: string
  emailVerified: boolean
  image: string | null
  status: UserStatus
  phoneNumber: string | null
  jobTitle: string | null
  department: string | null
  twoFactorEnabled: boolean
  role: Role | null
  permissions: string[]
  createdAt: string
  updatedAt: string
}

export type Session = {
  id: string
  createdAt: string
  updatedAt: string
  expiresAt: string
  ipAddress: string | null
  userAgent: string | null
}

export type SessionResponse = {
  session: Session
  user: User
}

export type ListUser = User & {
  canManage: boolean
  canDelete: boolean
}

export type UsersResponse = {
  users: ListUser[]
  total: number
  page: number
  pageCount: number
  pageSize: number
  roles: Role[]
  assignableRoles: Role[]
  registrationApproval: boolean
  superadminCount: number
}

export type UserEditResponse = {
  user: User
  assignableRoles: Role[]
  canManage: boolean
}

export type PermissionInfo = {
  key: string
  label: string
}

export type RolesResponse = {
  roles: Role[]
  permissions: PermissionInfo[]
}

export type RoleDatum = {
  name: string
  label: string
  users: number
}

export type SignupDatum = {
  date: string
  users: number
}

export type DashboardResponse = {
  totalUsers: number
  activeUsers: number
  pendingUsers: number
  rolesCount: number
  roleData: RoleDatum[]
  signupData: SignupDatum[]
  recentUsers: User[]
}

export type AccountOverview = {
  canDelete: boolean
  confirmed: boolean
}

export type TwoFactorSetup = {
  secret: string
  uri: string
  qrDataUrl: string
  backupCodes: string[]
}

export type LoginResponse =
  | SessionResponse
  | {
      twoFactorRequired: true
    }

export type UserInput = {
  name: string
  email: string
  phone: string
  jobTitle: string
  department: string
  role: string
  status: UserStatus
  password: string
}

export type RoleInput = {
  name: string
  label: string
  description: string
  permissions: string[]
}

export type ProfileInput = {
  name: string
  email: string
  phone: string
  jobTitle: string
  department: string
}
