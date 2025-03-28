"use client"

import { useEffect, useRef, useState, useCallback } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Button, IconButton } from "@mui/material"
import { Mic, MicOff, Video, VideoOff, MessageSquare, X, Share2 } from "lucide-react"
import { useToast } from "@/components/ui/use-toast"
import { ChatBox } from "@/components/chat-box"
import { useWebRTC } from "@/hooks/use-webrtc"
import {
  CallPageContainer,
  VideoGrid,
  VideoContainer,
  VideoElement,
  ParticipantInfo,
  ParticipantAvatar,
  ParticipantName,
  VideoControls,
  ControlButton,
  ControlBar,
  ControlsGroup,
  ChatContainer,
  ChatHeader,
  LoadingContainer,
  LoadingSpinner,
} from "@/components/styled"

export default function CallPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { toast } = useToast()
  const roomId = searchParams.get("roomId")
  const token = searchParams.get("token")
  const [showChat, setShowChat] = useState(false)
  const [isMuted, setIsMuted] = useState(false)
  const [isVideoOff, setIsVideoOff] = useState(false)
  const localVideoRef = useRef<HTMLVideoElement>(null)
  const [connectionInitiated, setConnectionInitiated] = useState(false)

  const { localStream, remoteStreams, participants, messages, sendMessage, isConnected, connect, disconnect } =
    useWebRTC()

  // Only run once to initiate the connection
  useEffect(() => {
    if (!roomId || !token) {
      toast({
        title: "Error",
        description: "Missing room ID or token",
        variant: "destructive",
      })
      router.push("/")
      return
    }

    if (!connectionInitiated) {
      connect(roomId, token)
      setConnectionInitiated(true)
    }

    return () => {
      if (connectionInitiated) {
        disconnect()
      }
    }
  }, [roomId, token, router, toast, connect, disconnect, connectionInitiated])

  // Handle local video stream
  useEffect(() => {
    if (localStream && localVideoRef.current) {
      localVideoRef.current.srcObject = localStream
    }
  }, [localStream])

  const toggleMute = useCallback(() => {
    if (localStream) {
      localStream.getAudioTracks().forEach((track) => {
        track.enabled = !track.enabled
      })
      setIsMuted(!isMuted)
    }
  }, [localStream, isMuted])

  const toggleVideo = useCallback(() => {
    if (localStream) {
      localStream.getVideoTracks().forEach((track) => {
        track.enabled = !track.enabled
      })
      setIsVideoOff(!isVideoOff)
    }
  }, [localStream, isVideoOff])

  const handleLeaveRoom = useCallback(() => {
    disconnect()
    router.push("/")
  }, [disconnect, router])

  const copyRoomId = useCallback(() => {
    if (roomId) {
      navigator.clipboard.writeText(roomId)
      toast({
        title: "Success",
        description: "Room ID copied to clipboard",
      })
    }
  }, [roomId, toast])

  if (!isConnected) {
    return (
      <LoadingContainer>
        <LoadingSpinner />
      </LoadingContainer>
    )
  }

  return (
    <CallPageContainer>
      <div style={{ display: "flex", flex: 1 }}>
        <VideoGrid>
          {/* Local video */}
          <VideoContainer>
            <VideoElement ref={localVideoRef} autoPlay playsInline muted />
            <ParticipantInfo>
              <ParticipantAvatar
                src={participants.find((p) => p.isLocal)?.avatar || "/placeholder.svg?height=32&width=32"}
                alt="You"
              />
              <ParticipantName>You</ParticipantName>
            </ParticipantInfo>
            <VideoControls>
              <ControlButton size="small" onClick={toggleMute}>
                {isMuted ? <MicOff size={16} /> : <Mic size={16} />}
              </ControlButton>
              <ControlButton size="small" onClick={toggleVideo}>
                {isVideoOff ? <VideoOff size={16} /> : <Video size={16} />}
              </ControlButton>
            </VideoControls>
          </VideoContainer>

          {/* Remote videos */}
          {remoteStreams.map((stream, index) => {
            const participant = participants.find((p) => p.id === stream.id)
            return (
              <VideoContainer key={stream.id}>
                <VideoElement
                  autoPlay
                  playsInline
                  ref={(ref) => {
                    if (ref) ref.srcObject = stream.stream
                  }}
                />
                <ParticipantInfo>
                  <ParticipantAvatar
                    src={participant?.avatar || "/placeholder.svg?height=32&width=32"}
                    alt={participant?.name || `User ${index + 1}`}
                  />
                  <ParticipantName>{participant?.name || `User ${index + 1}`}</ParticipantName>
                </ParticipantInfo>
              </VideoContainer>
            )
          })}
        </VideoGrid>

        {/* Chat sidebar */}
        {showChat && (
          <ChatContainer>
            <ChatHeader>
              <h3>Chat</h3>
              <IconButton size="small" onClick={() => setShowChat(false)}>
                <X size={18} />
              </IconButton>
            </ChatHeader>
            <ChatBox messages={messages} participants={participants} onSendMessage={sendMessage} />
          </ChatContainer>
        )}
      </div>

      {/* Control bar */}
      <ControlBar>
        <Button variant="outlined" startIcon={<Share2 size={16} />} onClick={copyRoomId} style={{ color: "white" }}>
          Share
        </Button>

        <ControlsGroup>
          <IconButton
            color={isMuted ? "error" : "default"}
            onClick={toggleMute}
            style={{ backgroundColor: isMuted ? "rgba(244, 67, 54, 0.1)" : "rgba(255, 255, 255, 0.1)", color: "white" }}
          >
            {isMuted ? <MicOff /> : <Mic />}
          </IconButton>
          <IconButton
            color={isVideoOff ? "error" : "default"}
            onClick={toggleVideo}
            style={{
              backgroundColor: isVideoOff ? "rgba(244, 67, 54, 0.1)" : "rgba(255, 255, 255, 0.1)",
              color: "white",
            }}
          >
            {isVideoOff ? <VideoOff /> : <Video />}
          </IconButton>
          <IconButton
            color={showChat ? "primary" : "default"}
            onClick={() => setShowChat(!showChat)}
            style={{ backgroundColor: "rgba(255, 255, 255, 0.1)", color: "white" }}
          >
            <MessageSquare />
          </IconButton>
        </ControlsGroup>

        <Button variant="contained" color="error" onClick={handleLeaveRoom}>
          Leave Room
        </Button>
      </ControlBar>
    </CallPageContainer>
  )
}

