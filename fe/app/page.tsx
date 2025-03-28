"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { Button, IconButton } from "@mui/material"
import { LogOut, Video } from "lucide-react"
import { logout, joinRoom } from "@/lib/api"
import { useAuth } from "@/hooks/use-auth"
import { useToast } from "@/components/ui/use-toast"
import { v4 as uuidv4 } from "uuid"
import {
  PageContainer,
  ContentCard,
  CardHeader,
  CardTitle,
  CardContent,
  UserInfoContainer,
  UserAvatar,
  UserName,
  RoomInputContainer,
  LoadingContainer,
  LoadingSpinner,
} from "@/components/styled"
import { Input } from "@/components/ui/input"

export default function Home() {
  const router = useRouter()
  const { toast } = useToast()
  const { user, isLoading, isAuthenticated, setUser } = useAuth()
  const [roomId, setRoomId] = useState("")

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login")
    }
  }, [isLoading, isAuthenticated, router])

  const handleJoinRoom = async () => {
    if (!roomId) {
      toast({
        title: "Error",
        description: "Please enter a room ID",
        variant: "destructive",
      })
      return
    }

    try {
      const token = await joinRoom(roomId)
      router.push(`/call?roomId=${roomId}&token=${token}`)
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to join room",
        variant: "destructive",
      })
    }
  }

  const handleCreateRoom = async () => {
    const newRoomId = uuidv4()
    try {
      const token = await joinRoom(newRoomId)
      router.push(`/call?roomId=${newRoomId}&token=${token}`)
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to create room",
        variant: "destructive",
      })
    }
  }

  const handleLogout = async () => {
    try {
      await logout()
      setUser(null)
      router.push("/login")
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to logout",
        variant: "destructive",
      })
    }
  }

  if (isLoading || !isAuthenticated) {
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
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
            <CardTitle variant="h1">Vego</CardTitle>
            <IconButton onClick={handleLogout}>
              <LogOut />
            </IconButton>
          </div>
        </CardHeader>
        <CardContent>
          <UserInfoContainer>
            <UserAvatar src={user?.avatar || "https://www.svgrepo.com/show/532362/user.svg"} alt={user?.name} />
            <UserName variant="h2">Hi, {user?.name}</UserName>
          </UserInfoContainer>

          <div>
            <RoomInputContainer>
              <Input
                placeholder="Enter Room ID"
                value={roomId}
                onChange={(e) => setRoomId(e.target.value)}
                style={{ flex: 1 }}
              />
              <Button variant="contained" onClick={handleJoinRoom}>
                Join Room
              </Button>
            </RoomInputContainer>
            <Button variant="outlined" fullWidth startIcon={<Video size={16} />} onClick={handleCreateRoom}>
              Create a Room
            </Button>
          </div>
        </CardContent>
      </ContentCard>
    </PageContainer>
  )
}

