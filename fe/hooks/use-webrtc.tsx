"use client"

import { useState, useCallback, useRef, useEffect } from "react"
import type { Participant, ChatMessage } from "@/lib/types"

interface RemoteStream {
  id: string
  stream: MediaStream
}

interface WebRTCMessage {
  type: string
  from: string
  data: any
}

export function useWebRTC() {
  const [localStream, setLocalStream] = useState<MediaStream | null>(null)
  const [remoteStreams, setRemoteStreams] = useState<RemoteStream[]>([])
  const [participants, setParticipants] = useState<Participant[]>([])
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [isConnected, setIsConnected] = useState(false)

  const socketRef = useRef<WebSocket | null>(null)
  const peerConnectionsRef = useRef<Record<string, RTCPeerConnection>>({})
  const localUserIdRef = useRef<string>("")
  const isConnectingRef = useRef<boolean>(false)

  const connect = useCallback(
    async (roomId: string, token: string) => {
      // Prevent multiple connection attempts
      if (isConnectingRef.current || isConnected) {
        console.log("Connection already in progress or established")
        return
      }

      isConnectingRef.current = true

      try {
        // Get local media stream
        const stream = await navigator.mediaDevices.getUserMedia({
          video: true,
          audio: true,
        })
        setLocalStream(stream)

        // Connect to WebSocket server
        const wsEndpoint = process.env.WS_ENDPOINT ?? `${window.location.protocol === "https:" ? "wss:" : "ws:"}//localhost:8080/api/room/ws/`
        const wsUrl = wsEndpoint + `${roomId}?token=${token}`
        let socket = new WebSocket(wsUrl)

        socket.onopen = () => {
          console.log("WebSocket connection established")
          setIsConnected(true)
          isConnectingRef.current = false
        }

        socket.onmessage = (event) => {
          const message: WebRTCMessage = JSON.parse(event.data)
          handleWebSocketMessage(message)
        }

        socket.onerror = (error) => {
          socket = new WebSocket(wsUrl)
          console.error("WebSocket error:", error)
          isConnectingRef.current = false
        }

        socket.onclose = () => {
          console.log("WebSocket connection closed")
          setIsConnected(false)
          isConnectingRef.current = false
        }

        socketRef.current = socket
      } catch (error) {
        console.error("Error connecting to room:", error)
        isConnectingRef.current = false
      }
    },
    [isConnected],
  )

  const disconnect = useCallback(() => {
    // Close all peer connections
    Object.values(peerConnectionsRef.current).forEach((pc) => {
      pc.close()
    })
    peerConnectionsRef.current = {}

    // Close WebSocket connection
    if (socketRef.current) {
      socketRef.current.close()
      socketRef.current = null
    }

    // Stop local media tracks
    if (localStream) {
      localStream.getTracks().forEach((track) => {
        track.stop()
      })
      setLocalStream(null)
    }

    // Reset state
    setRemoteStreams([])
    setParticipants([])
    setMessages([])
    setIsConnected(false)
    isConnectingRef.current = false
  }, [localStream])

  const handleWebSocketMessage = useCallback((message: WebRTCMessage) => {
    const { type, from, data } = message

    switch (type) {
      case "info":
        // Handle initial room info with participants
        const infoData = data as { users: Array<{ inner_id: string; name: string; avatar: string }> }

        // Set local user ID and add participants
        localUserIdRef.current = from

        const participantsList: Participant[] = [
          // Add local user
          {
            id: from,
            name: "You",
            avatar: "",
            isLocal: true,
          },
          // Add remote users
          ...infoData.users.map((user) => ({
            id: user.inner_id,
            name: user.name,
            avatar: user.avatar,
            isLocal: false,
          })),
        ]

        setParticipants(participantsList)

        // Create peer connections for each participant
        infoData.users.forEach((user) => {
          createPeerConnection(user.inner_id, true)
        })
        break

      case "join":
        // Handle new participant joining
        const joinData = data as { name: string; avatar: string }

        // Add new participant to the list
        setParticipants((prev) => [
          ...prev,
          {
            id: from,
            name: joinData.name,
            avatar: joinData.avatar,
            isLocal: false,
          },
        ])

        // Create peer connection for the new participant
        createPeerConnection(from, true)
        break

      case "leave":
        // Handle participant leaving
        // Close and remove peer connection
        if (peerConnectionsRef.current[from]) {
          peerConnectionsRef.current[from].close()
          delete peerConnectionsRef.current[from]
        }

        // Remove participant from the list
        setParticipants((prev) => prev.filter((p) => p.id !== from))

        // Remove remote stream
        setRemoteStreams((prev) => prev.filter((s) => s.id !== from))
        break

      case "offer":
        // Handle incoming offer
        handleOffer(from, data.content)
        break

      case "answer":
        // Handle incoming answer
        handleAnswer(from, data.content)
        break

      case "ice-candidate":
        // Handle incoming ICE candidate
        handleIceCandidate(from, data.content)
        break

      case "chat-message":
        // Handle chat message
        const chatData = data as { content: string; ts: string }
        setMessages((prev) => [
          ...prev,
          {
            senderId: from,
            content: chatData.content,
            timestamp: new Date(chatData.ts).getTime(),
          },
        ])
        break
    }
  }, [])

  const createPeerConnection = useCallback(
    (peerId: string, isInitiator: boolean) => {
      if (!localStream) return

      // Check if we already have a connection for this peer
      if (peerConnectionsRef.current[peerId]) {
        console.log(`Peer connection for ${peerId} already exists`)
        return peerConnectionsRef.current[peerId]
      }

      // Create new RTCPeerConnection
      const pc = new RTCPeerConnection({
        iceServers: [{ urls: "stun:stun.l.google.com:19302" }, { urls: "stun:stun1.l.google.com:19302" }],
      })

      // Add local tracks to the connection
      localStream.getTracks().forEach((track) => {
        pc.addTrack(track, localStream)
      })

      // Handle ICE candidates
      pc.onicecandidate = (event) => {
        if (event.candidate) {
          sendWebRTCMessage(peerId, "ice-candidate", JSON.stringify(event.candidate))
        }
      }

      // Handle remote tracks
      pc.ontrack = (event) => {
        const stream = event.streams[0]

        setRemoteStreams((prev) => {
          // Check if we already have this stream
          if (prev.some((s) => s.id === peerId)) {
            return prev
          }

          return [...prev, { id: peerId, stream }]
        })
      }

      // Store the peer connection
      peerConnectionsRef.current[peerId] = pc

      // If we're the initiator, create and send an offer
      if (isInitiator) {
        pc.createOffer()
          .then((offer) => pc.setLocalDescription(offer))
          .then(() => {
            if (pc.localDescription) {
              sendWebRTCMessage(peerId, "offer", JSON.stringify(pc.localDescription))
            }
          })
          .catch((error) => console.error("Error creating offer:", error))
      }

      return pc
    },
    [localStream],
  )

  const handleOffer = useCallback(
    (peerId: string, offerStr: string) => {
      const offer = JSON.parse(offerStr)

      // Get or create peer connection
      let pc = peerConnectionsRef.current[peerId]
      if (!pc) {
        pc = createPeerConnection(peerId, false)
        if (!pc) return
      }

      // Set remote description and create answer
      pc.setRemoteDescription(new RTCSessionDescription(offer))
        .then(() => pc.createAnswer())
        .then((answer) => pc.setLocalDescription(answer))
        .then(() => {
          if (pc.localDescription) {
            sendWebRTCMessage(peerId, "answer", JSON.stringify(pc.localDescription))
          }
        })
        .catch((error) => console.error("Error handling offer:", error))
    },
    [createPeerConnection],
  )

  const handleAnswer = useCallback((peerId: string, answerStr: string) => {
    const answer = JSON.parse(answerStr)
    const pc = peerConnectionsRef.current[peerId]

    if (pc) {
      pc.setRemoteDescription(new RTCSessionDescription(answer)).catch((error) =>
        console.error("Error handling answer:", error),
      )
    }
  }, [])

  const handleIceCandidate = useCallback((peerId: string, candidateStr: string) => {
    const candidate = JSON.parse(candidateStr)
    const pc = peerConnectionsRef.current[peerId]

    if (pc) {
      pc.addIceCandidate(new RTCIceCandidate(candidate)).catch((error) =>
        console.error("Error adding ICE candidate:", error),
      )
    }
  }, [])

  const sendWebRTCMessage = useCallback((to: string, type: string, content: string) => {
    if (!socketRef.current || socketRef.current.readyState !== WebSocket.OPEN) return

    const message = {
      type,
      data: JSON.stringify({ to, content }),
    }

    socketRef.current.send(JSON.stringify(message))
  }, [])

  const sendMessage = useCallback((content: string) => {
    if (!socketRef.current || socketRef.current.readyState !== WebSocket.OPEN) return

    const message = {
      type: "chat-message",
      data: JSON.stringify({ content, ts: new Date().toISOString() }),
    }

    socketRef.current.send(JSON.stringify(message))

    // Add message to local state
    setMessages((prev) => [
      ...prev,
      {
        senderId: localUserIdRef.current,
        content,
        timestamp: Date.now(),
      },
    ])
  }, [])

  // Clean up on unmount
  useEffect(() => {
    return () => {
      disconnect()
    }
  }, [disconnect])

  return {
    localStream,
    remoteStreams,
    participants,
    messages,
    isConnected,
    connect,
    disconnect,
    sendMessage,
  }
}

