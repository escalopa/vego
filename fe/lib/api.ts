import axios from "axios"
import type { User } from "./types"

const api = axios.create({
  baseURL: process.env.BACKEND_URL ?? "http://localhost:8080/api",
  withCredentials: true,
})

// Add response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      const errorMessage = error.response.data?.error || "An error occurred"
      return Promise.reject(new Error(errorMessage))
    } else if (error.request) {
      // The request was made but no response was received
      return Promise.reject(new Error("No response from server"))
    } else {
      // Something happened in setting up the request that triggered an Error
      return Promise.reject(error)
    }
  },
)

export async function getUserInfo(): Promise<User> {
  try {
    const response = await api.get("/user/info")
    return response.data.user
  } catch (error) {

    throw error
  }
}

export async function logout(): Promise<void> {
  try {
    await api.post("/user/logout")
  } catch (error) {
    throw error
  }
}

export async function getOAuthUrl(provider: string): Promise<string> {
  try {
    const response = await api.get(`/oauth/${provider}`)
    return response.data.url || response.request.responseURL
  } catch (error) {
    throw error
  }
}

export async function handleOAuthCallback(provider: string, code: string): Promise<void> {
  try {
    await api.post(`/oauth/${provider}/callback`, { code })
    await new Promise((resolve) => setTimeout(resolve, 500))
  } catch (error) {
    throw error
  }
}

export async function joinRoom(roomId: string): Promise<string> {
  try {
    const response = await api.post(`/room/join/${roomId}`)
    return response.data.token
  } catch (error) {
    throw error
  }
}

