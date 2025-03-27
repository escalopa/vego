package domain

type (
	User struct {
		UserID int64  `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}

	UserTokenPayload struct {
		UserID int64  `json:"user_id"`
		Email  string `json:"email"`
	}

	RoomTokenPayload struct {
		UserID int64  `json:"user_id"`
		RoomID string `json:"room_id"`
	}

	Token struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}
)
