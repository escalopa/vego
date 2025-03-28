"use client"

import { createContext, useContext, useState, useEffect, type ReactNode } from "react"
import { getUserInfo } from "@/lib/api"
import type { User } from "@/lib/types"

interface AuthContextType {
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
  setUser: (user: User | null) => void
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  isLoading: true,
  isAuthenticated: false,
  setUser: () => {},
})

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await getUserInfo()
        setUser(userData)
      } catch (error) {
        console.log("Not authenticated")
        setUser(null)
      } finally {
        setIsLoading(false)
      }
    }

    fetchUser()
  }, [])

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        setUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}

