export interface User {
  user_id: number
  name: string
  email: string
  avatar: string
}

export interface Participant {
  id: string
  name: string
  avatar: string
  isLocal: boolean
}

export interface ChatMessage {
  senderId: string
  content: string
  timestamp: number
}

