"use client"

import type React from "react"

import { useState, useRef, useEffect } from "react"
import { IconButton } from "@mui/material"
import { Send } from "lucide-react"
import type { Participant, ChatMessage } from "@/lib/types"
import {
  ChatMessages,
  MessageContainer,
  MessageAvatar,
  MessageContent,
  MessageHeader,
  MessageSender,
  MessageTime,
  MessageText,
  ChatInputContainer,
  ChatInput,
} from "@/components/styled"

interface ChatBoxProps {
  messages: ChatMessage[]
  participants: Participant[]
  onSendMessage: (content: string) => void
}

export function ChatBox({ messages, participants, onSendMessage }: ChatBoxProps) {
  const [message, setMessage] = useState("")
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  const handleSendMessage = () => {
    if (message.trim()) {
      onSendMessage(message)
      setMessage("")
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  const getParticipantById = (id: string) => {
    return participants.find((p) => p.id === id)
  }

  return (
    <>
      <ChatMessages>
        {messages.map((msg, index) => {
          const participant = getParticipantById(msg.senderId)
          return (
            <MessageContainer key={index}>
              <MessageAvatar
                src={participant?.avatar || "/placeholder.svg?height=32&width=32"}
                alt={participant?.name || "Unknown"}
              />
              <MessageContent>
                <MessageHeader>
                  <MessageSender>{participant?.name || "Unknown"}</MessageSender>
                  <MessageTime>
                    {new Date(msg.timestamp).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                  </MessageTime>
                </MessageHeader>
                <MessageText>{msg.content}</MessageText>
              </MessageContent>
            </MessageContainer>
          )
        })}
        <div ref={messagesEndRef} />
      </ChatMessages>

      <ChatInputContainer>
        <ChatInput
          placeholder="Type a message..."
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyDown={handleKeyDown}
          size="small"
          fullWidth
        />
        <IconButton color="primary" onClick={handleSendMessage}>
          <Send size={18} />
        </IconButton>
      </ChatInputContainer>
    </>
  )
}

