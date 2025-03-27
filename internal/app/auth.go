package app

import (
	"errors"
	"log"
	"net/http"

	"github.com/escalopa/peer-cast/internal/domain"
	"github.com/gin-gonic/gin"
)

const (
	accessTokenKey  = "X-Access-Token"
	refreshTokenKey = "X-Refresh-Token"
)

func (a *App) authMiddleware(c *gin.Context) {
	accessToken, err := c.Cookie(accessTokenKey)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := c.Cookie(refreshTokenKey)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if accessToken == "" && refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated user"})
		return
	}

	token := &domain.Token{Access: accessToken, Refresh: refreshToken}
	user, token, err := a.srv.AuthenticateUser(c.Request.Context(), token)
	if err != nil {
		log.Printf("db.AuthenticateUser: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "internal server error"})
		return
	}

	// if tokens were refreshed then set it back in the cookies
	if token != nil {
		a.setTokenCookie(c, token)
	}

	c.Set("user", user)
	c.Next()
}

func (a *App) setTokenCookie(c *gin.Context, token *domain.Token) {
	var (
		accessToken, refreshToken             string
		accessTokenExpiry, refreshTokenExpiry int
	)

	if token == nil { // delete the cookies
		accessTokenExpiry = -1
		refreshTokenExpiry = -1
	} else {
		accessToken = token.Access
		refreshToken = token.Refresh
		accessTokenExpiry = int(a.cfg.AccessTokenTTL.Seconds())
		refreshTokenExpiry = int(a.cfg.RefreshTokenTTL.Seconds())
	}

	// set access token cookie
	c.SetCookie(
		accessTokenKey,
		accessToken,
		accessTokenExpiry,
		"/",
		a.cfg.Domain,
		true,
		true,
	)

	// set refresh token cookie
	c.SetCookie(
		refreshTokenKey,
		refreshToken,
		refreshTokenExpiry,
		"/",
		a.cfg.Domain,
		true,
		true,
	)
}
