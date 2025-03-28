import { styled } from "@mui/material/styles"
import { Box, Paper, Typography, Button, TextField, IconButton } from "@mui/material"

// Common styled components
export const PageContainer = styled(Box)(({ theme }) => ({
  display: "flex",
  minHeight: "100vh",
  flexDirection: "column",
  alignItems: "center",
  justifyContent: "center",
  padding: theme.spacing(2),
  backgroundColor: "#f5f5f5",
}))

export const ContentCard = styled(Paper)(({ theme }) => ({
  width: "100%",
  maxWidth: "500px",
  padding: theme.spacing(3),
  borderRadius: theme.shape.borderRadius,
  boxShadow: theme.shadows[2],
}))

export const CardHeader = styled(Box)(({ theme }) => ({
  marginBottom: theme.spacing(3),
  textAlign: "center",
}))

export const CardTitle = styled(Typography)(({ theme }) => ({
  fontSize: "1.5rem",
  fontWeight: 600,
}))

export const CardContent = styled(Box)(({ theme }) => ({
  display: "flex",
  flexDirection: "column",
  gap: theme.spacing(3),
}))

// Login page styled components
export const LoginIcon = styled("img")(({ theme }) => ({
  width: "64px",
  height: "64px",
  marginBottom: theme.spacing(2),
}))

export const OAuthButton = styled(Button)(({ theme }) => ({
  width: "100%",
  display: "flex",
  justifyContent: "flex-start",
  padding: theme.spacing(1.5),
  marginBottom: theme.spacing(1.5),
  textTransform: "none",
}))

export const OAuthIcon = styled("img")({
  width: "24px",
  height: "24px",
  marginRight: "12px",
})

// Home page styled components
export const UserInfoContainer = styled(Box)(({ theme }) => ({
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  marginBottom: theme.spacing(4),
}))

export const UserAvatar = styled("img")(({ theme }) => ({
  width: "80px",
  height: "80px",
  borderRadius: "50%",
  marginBottom: theme.spacing(1),
}))

export const UserName = styled(Typography)(({ theme }) => ({
  fontSize: "1.2rem",
  fontWeight: 500,
}))

export const RoomInputContainer = styled(Box)(({ theme }) => ({
  display: "flex",
  gap: theme.spacing(1),
  marginBottom: theme.spacing(2),
}))

// Call page styled components
export const CallPageContainer = styled(Box)({
  display: "flex",
  flexDirection: "column",
  height: "100vh",
  backgroundColor: "#0f172a",
})

export const VideoGrid = styled(Box)(({ theme }) => ({
  flex: 1,
  display: "grid",
  gridTemplateColumns: "repeat(auto-fit, minmax(300px, 1fr))",
  gap: theme.spacing(2),
  padding: theme.spacing(2),
  overflow: "auto",
}))

export const VideoContainer = styled(Box)(({ theme }) => ({
  position: "relative",
  backgroundColor: "#1e293b",
  borderRadius: theme.shape.borderRadius,
  overflow: "hidden",
  aspectRatio: "16/9",
}))

export const VideoElement = styled("video")({
  width: "100%",
  height: "100%",
  objectFit: "cover",
})

export const ParticipantInfo = styled(Box)(({ theme }) => ({
  position: "absolute",
  bottom: 0,
  left: 0,
  right: 0,
  display: "flex",
  alignItems: "center",
  padding: theme.spacing(1),
  backgroundColor: "rgba(0, 0, 0, 0.5)",
}))

export const ParticipantAvatar = styled("img")({
  width: "32px",
  height: "32px",
  borderRadius: "50%",
  marginRight: "8px",
})

export const ParticipantName = styled(Typography)({
  color: "white",
  fontSize: "0.875rem",
})

export const VideoControls = styled(Box)(({ theme }) => ({
  position: "absolute",
  top: theme.spacing(1),
  right: theme.spacing(1),
  display: "flex",
  gap: theme.spacing(0.5),
  opacity: 0,
  transition: "opacity 0.2s",
  ".MuiVideoContainer-root:hover &": {
    opacity: 1,
  },
}))

export const ControlButton = styled(IconButton)(({ theme }) => ({
  backgroundColor: "rgba(0, 0, 0, 0.5)",
  color: "white",
  "&:hover": {
    backgroundColor: "rgba(0, 0, 0, 0.7)",
  },
}))

export const ControlBar = styled(Box)(({ theme }) => ({
  display: "flex",
  justifyContent: "space-between",
  alignItems: "center",
  padding: theme.spacing(2),
  backgroundColor: "#1e293b",
}))

export const ControlsGroup = styled(Box)(({ theme }) => ({
  display: "flex",
  gap: theme.spacing(1),
}))

// Chat components
export const ChatContainer = styled(Box)(({ theme }) => ({
  width: "320px",
  backgroundColor: "white",
  borderLeft: `1px solid ${theme.palette.divider}`,
  display: "flex",
  flexDirection: "column",
  height: "100%",
}))

export const ChatHeader = styled(Box)(({ theme }) => ({
  display: "flex",
  justifyContent: "space-between",
  alignItems: "center",
  padding: theme.spacing(1.5),
  borderBottom: `1px solid ${theme.palette.divider}`,
}))

export const ChatMessages = styled(Box)(({ theme }) => ({
  flex: 1,
  padding: theme.spacing(2),
  overflowY: "auto",
  display: "flex",
  flexDirection: "column",
  gap: theme.spacing(2),
}))

export const MessageContainer = styled(Box)(({ theme }) => ({
  display: "flex",
  alignItems: "flex-start",
}))

export const MessageAvatar = styled("img")({
  width: "32px",
  height: "32px",
  borderRadius: "50%",
  marginRight: "8px",
  marginTop: "4px",
})

export const MessageContent = styled(Box)({
  flex: 1,
})

export const MessageHeader = styled(Box)(({ theme }) => ({
  display: "flex",
  alignItems: "baseline",
  marginBottom: theme.spacing(0.5),
}))

export const MessageSender = styled(Typography)(({ theme }) => ({
  fontWeight: 500,
  fontSize: "0.875rem",
  marginRight: theme.spacing(1),
}))

export const MessageTime = styled(Typography)(({ theme }) => ({
  fontSize: "0.75rem",
  color: theme.palette.text.secondary,
}))

export const MessageText = styled(Typography)(({ theme }) => ({
  fontSize: "0.875rem",
}))

export const ChatInputContainer = styled(Box)(({ theme }) => ({
  display: "flex",
  gap: theme.spacing(1),
  padding: theme.spacing(1.5),
  borderTop: `1px solid ${theme.palette.divider}`,
}))

export const ChatInput = styled(TextField)(({ theme }) => ({
  flex: 1,
}))

// Loading components
export const LoadingContainer = styled(Box)(({ theme }) => ({
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  minHeight: "100vh",
}))

export const LoadingSpinner = styled(Box)(({ theme }) => ({
  width: "48px",
  height: "48px",
  border: `4px solid ${theme.palette.divider}`,
  borderTop: `4px solid ${theme.palette.primary.main}`,
  borderRadius: "50%",
  animation: "spin 1s linear infinite",
  "@keyframes spin": {
    "0%": { transform: "rotate(0deg)" },
    "100%": { transform: "rotate(360deg)" },
  },
}))

