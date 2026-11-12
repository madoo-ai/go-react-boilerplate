import {
  createContext,
  useCallback,
  useContext,
  useMemo,
} from "react"
import { useQuery, useQueryClient } from "@tanstack/react-query"
import { api, ApiError } from "@/api/client"
import type { SessionResponse, User } from "@/api/types"

type AuthValue = {
  session: SessionResponse | null
  user: User | null
  isLoading: boolean
  refresh: () => Promise<void>
  setSession: (session: SessionResponse | null) => void
  signOut: () => Promise<void>
}

const AuthContext = createContext<AuthValue | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const queryClient = useQueryClient()
  const query = useQuery({
    queryKey: ["session"],
    queryFn: api.session,
    retry: false,
    staleTime: 30_000,
  })

  const refresh = useCallback(async () => {
    await queryClient.invalidateQueries({ queryKey: ["session"] })
  }, [queryClient])

  const setSession = useCallback(
    (session: SessionResponse | null) => {
      queryClient.setQueryData(["session"], session)
    },
    [queryClient],
  )

  const signOut = useCallback(async () => {
    await api.logout()
    queryClient.setQueryData(["session"], null)
    queryClient.removeQueries({ queryKey: ["dashboard"] })
    queryClient.removeQueries({ queryKey: ["users"] })
    queryClient.removeQueries({ queryKey: ["roles"] })
  }, [queryClient])

  const value = useMemo<AuthValue>(
    () => ({
      session: query.data ?? null,
      user: query.data?.user ?? null,
      isLoading: query.isLoading,
      refresh,
      setSession,
      signOut,
    }),
    [query.data, query.isLoading, refresh, setSession, signOut],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error("useAuth must be used inside AuthProvider")
  return context
}

export function errorKey(error: unknown) {
  if (error instanceof ApiError) return error.code
  return "errors.unexpected"
}

export function fieldErrors(error: unknown): Record<string, string> {
  if (error instanceof ApiError) return error.fields
  return {}
}
