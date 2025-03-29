"use client"

import { useEffect, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { PageContainer, ContentCard, CardContent, LoadingSpinner } from "@/components/styled"
import { Typography, Button } from "@mui/material"
import { handleOAuthCallback } from "@/lib/api"
import { useAuth } from "@/hooks/use-auth"

export default function OAuthCallbackClient({ provider }: { provider: string }) {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { checkAuth } = useAuth()
  const code = searchParams.get("code")
  const [error, setError] = useState<string | null>(null)
  const [isProcessing, setIsProcessing] = useState(true)
  const [hasAttemptedAuth, setHasAttemptedAuth] = useState(false)

  useEffect(() => {
    if (!code || hasAttemptedAuth) {
      if (!code) {
        setError("Authorization code is missing")
      }
      setIsProcessing(false)
      return
    }

    const processOAuthCallback = async () => {
      try {
        setHasAttemptedAuth(true)
        await handleOAuthCallback(provider, code)
        // Check if we're actually authenticated after setting tokens
        const isAuthenticated = await checkAuth()
        if (isAuthenticated) {
          router.push("/")
        } else {
          setError("Failed to authenticate")
        }
      } catch (err: any) {
        setError(err.message || "Failed to authenticate")
      } finally {
        setIsProcessing(false)
      }
    }

    processOAuthCallback()
  }, [router, provider, code, checkAuth, hasAttemptedAuth])

  const handleRetry = () => {
    router.push(`/login`)
  }

  if (error) {
    return (
      <PageContainer>
        <ContentCard>
          <CardContent>
            <div style={{ display: "flex", flexDirection: "column", alignItems: "center", padding: "24px" }}>
              <Typography variant="h6" style={{ marginBottom: "16px", color: "#d32f2f" }}>
                Authentication Error
              </Typography>
              <Typography style={{ marginBottom: "16px", textAlign: "center" }}>{error}</Typography>
              <Button variant="contained" color="primary" onClick={handleRetry}>
                Return to Login
              </Button>
            </div>
          </CardContent>
        </ContentCard>
      </PageContainer>
    )
  }

  return (
    <PageContainer>
      <ContentCard>
        <CardContent>
          <div style={{ display: "flex", flexDirection: "column", alignItems: "center", padding: "24px" }}>
            <Typography variant="h6" style={{ marginBottom: "16px" }}>
              Authentication in progress
            </Typography>
            <Typography variant="body2" style={{ marginBottom: "16px", textAlign: "center" }}>
              Authenticating with {provider.charAt(0).toUpperCase() + provider.slice(1)}...
            </Typography>
            <LoadingSpinner />
          </div>
        </CardContent>
      </ContentCard>
    </PageContainer>
  )
} 