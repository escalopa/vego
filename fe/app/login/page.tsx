"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/hooks/use-auth"
import { getOAuthUrl } from "@/lib/api"
import {
  PageContainer,
  ContentCard,
  CardHeader,
  CardTitle,
  CardContent,
  LoginIcon,
  OAuthButton,
  OAuthIcon,
  LoadingContainer,
  LoadingSpinner,
} from "@/components/styled"

export default function Login() {
  const router = useRouter()
  const { isAuthenticated, isLoading } = useAuth()

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.push("/")
    }
  }, [isLoading, isAuthenticated, router])

  const handleOAuthLogin = async (provider: string) => {
    try {
      const url = await getOAuthUrl(provider)
      console.log("OAuth URL:", url)
      window.location.href = url
    } catch (error) {
      console.error("OAuth error:", error)
    }
  }

  if (isLoading) {
    return (
      <LoadingContainer>
        <LoadingSpinner />
      </LoadingContainer>
    )
  }

  return (
    <PageContainer>
      <ContentCard>
        <CardHeader>
          <CardTitle variant="h1">Welcome to Vego</CardTitle>
        </CardHeader>
        <CardContent>
          <div style={{ display: "flex", justifyContent: "center", marginBottom: "24px" }}>
            <LoginIcon src="https://www.svgrepo.com/show/17636/video-call.svg" alt="Video call icon" />
          </div>

          <div>
            <OAuthButton variant="outlined" color="primary" onClick={() => handleOAuthLogin("google")}>
              <OAuthIcon src="https://svgrepo.com/show/475656/google-color.svg" alt="Google" />
              Continue with Google
            </OAuthButton>

            <OAuthButton variant="outlined" onClick={() => handleOAuthLogin("github")}>
              <OAuthIcon src="https://svgrepo.com/show/512317/github-142.svg" alt="GitHub" />
              Continue with GitHub
            </OAuthButton>

            <OAuthButton variant="outlined" onClick={() => handleOAuthLogin("yandex")}>
              <OAuthIcon src="https://svgrepo.com/show/197976/yandex.svg" alt="Yandex" />
              Continue with Yandex
            </OAuthButton>
          </div>
        </CardContent>
      </ContentCard>
    </PageContainer>
  )
}

