package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/escalopa/peer-cast/internal/config"
	"github.com/escalopa/peer-cast/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
)

type OAuthProvider struct {
	providers map[string]*provider
}

func NewOAuthProvider(cfg config.OAuthConfig) *OAuthProvider {
	op := &OAuthProvider{providers: make(map[string]*provider)}

	op.providers[googleProvider] = &provider{
		config: &oauth2.Config{
			ClientID:     cfg.Google.ClientID,
			ClientSecret: cfg.Google.ClientSecret,
			RedirectURL:  cfg.Google.RedirectURL,
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		endpoint: cfg.Google.UserEndpoint,
		payload:  func() payload { return &googlePayload{} },
	}

	op.providers[githubProvider] = &provider{
		config: &oauth2.Config{
			ClientID:     cfg.GitHub.ClientID,
			ClientSecret: cfg.GitHub.ClientSecret,
			RedirectURL:  cfg.GitHub.RedirectURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
		endpoint: cfg.GitHub.UserEndpoint,
		payload:  func() payload { return &githubPayload{} },
	}

	op.providers[yandexProvider] = &provider{
		config: &oauth2.Config{
			ClientID:     cfg.Yandex.ClientID,
			ClientSecret: cfg.Yandex.ClientSecret,
			RedirectURL:  cfg.Yandex.RedirectURL,
			Scopes:       []string{"login:info", "login:email", "login:avatar"},
			Endpoint:     yandex.Endpoint,
		},
		endpoint: cfg.Yandex.UserEndpoint,
		payload:  func() payload { return &yandexPayload{} },
	}

	return op
}

func (op *OAuthProvider) GetRedirectURL(provider string) (string, error) {
	p, exists := op.providers[provider]
	if !exists {
		return "", domain.ErrOAuthUnsupportedProvider
	}

	return p.config.AuthCodeURL(provider), nil
}

func (op *OAuthProvider) HandleCallback(ctx context.Context, provider string, code string) (*domain.User, error) {
	p, exists := op.providers[provider]
	if !exists {
		return nil, domain.ErrOAuthUnsupportedProvider
	}

	// get unique token for user's data retrieval
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		log.Printf("oauthConfig.Exchange err: %v", err)
		return nil, domain.ErrOAuthExchange
	}

	// fetch user's data
	client := p.config.Client(ctx, token)
	resp, err := client.Get(p.endpoint)
	if err != nil {
		return nil, err
	}
	defer func(b io.ReadCloser) { _ = b.Close() }(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch user info: %s", resp.Status)
	}

	// parse user's data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dst := p.payload()
	err = json.Unmarshal(body, dst)
	if err != nil {
		return nil, err
	}

	return dst.ToUser(), nil
}
