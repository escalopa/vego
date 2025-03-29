import OAuthCallbackClient from "./client"

export default function OAuthCallback({ params }: { params: { provider: string } }) {
  return <OAuthCallbackClient provider={params.provider} />
}

