package auth

import (
	"fmt"

	"github.com/escalopa/peer-cast/internal/domain"
	"golang.org/x/oauth2"
)

const (
	googleProvider = "google"
	githubProvider = "github"
	yandexProvider = "yandex"
)

type payload interface {
	ToUser() *domain.User
}

type provider struct {
	config   *oauth2.Config
	endpoint string
	payload  func() payload
}

type (
	googlePayload struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}

	githubPayload struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	yandexPayload struct {
		RealName        string `json:"real_name"`
		DefaultEmail    string `json:"default_email"`
		DefaultAvatarID string `json:"default_avatar_id"`
	}
)

func (p *googlePayload) ToUser() *domain.User {
	return &domain.User{
		Name:   p.Name,
		Email:  p.Email,
		Avatar: p.Picture,
	}
}

func (p *githubPayload) ToUser() *domain.User {
	return &domain.User{
		Name:   p.Name,
		Email:  p.Email,
		Avatar: p.AvatarURL,
	}
}

func (p *yandexPayload) ToUser() *domain.User {
	return &domain.User{
		Name:   p.RealName,
		Email:  p.DefaultEmail,
		Avatar: fmt.Sprintf("https://avatars.yandex.net/get-yapic/%s/islands-200", p.DefaultAvatarID),
	}
}
