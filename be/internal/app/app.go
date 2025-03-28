package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/escalopa/vego/internal/domain"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type service interface {
	GetOAuthRedirectURL(provider string) (string, error)
	RegisterUser(ctx context.Context, provider string, code string) (*domain.Token, error)
	AuthenticateUser(ctx context.Context, token *domain.Token) (*domain.User, *domain.Token, error)
	CreateRoomToken(userID int64, roomID string) (string, error)
	AuthenticateWS(ctx context.Context, token string, roomID string) (*domain.User, error)
	HandleWS(user *domain.User, roomID string, conn *websocket.Conn)
}

type Config struct {
	Domain       string
	AllowOrigins []string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type App struct {
	cfg Config

	r   *gin.Engine
	srv service
	upg *websocket.Upgrader
}

func New(cfg Config, srv service) *App {
	kors := cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	a := &App{
		r:   gin.Default(),
		cfg: cfg,
		srv: srv,
		upg: &websocket.Upgrader{
			ReadBufferSize:  1024 * 1024,
			WriteBufferSize: 1024 * 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				for _, o := range cfg.AllowOrigins {
					if o == origin {
						return true
					}
				}
				return false
			},
		},
	}

	a.r.Use(kors)
	a.setup()

	return a
}

func (a *App) Run(address string) error {
	return a.r.Run(address)
}

func (a *App) setup() {
	a.r.GET("/api/health", a.health)

	userRoutes := a.r.Group("/api/user")
	userRoutes.Use(a.authMiddleware)
	{
		userRoutes.GET("/info", a.getUserInfo)
		userRoutes.POST("/logout", a.logout)
	}

	roomRoutes := a.r.Group("/api/room")
	roomRoutes.Use(a.authMiddleware)
	{
		roomRoutes.POST("/join/:room_id", a.joinRoom)
		roomRoutes.GET("/ws/:room_id", a.ws)
	}

	oauthRoutes := a.r.Group("/api/oauth")
	{
		oauthRoutes.GET("/:provider", a.oauthRedirect)
		oauthRoutes.POST("/:provider/callback", a.oauthCallback)
	}
}

func (a *App) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *App) getUserInfo(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (a *App) logout(c *gin.Context) {
	a.setTokenCookie(c, nil)
	c.JSON(http.StatusOK, gin.H{"message": "user logged out"})
}

func (a *App) joinRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	if _, err := uuid.Parse(roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corrupted room id"})
		return
	}

	user := a.user(c)
	token, err := a.srv.CreateRoomToken(user.UserID, roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "temporary cannot join room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *App) ws(c *gin.Context) {
	roomID := c.Param("room_id")
	if _, err := uuid.Parse(roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corrupted room id (uuid expected)"})
		return
	}

	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty room token"})
		return
	}

	user, err := a.srv.AuthenticateWS(c.Request.Context(), token, roomID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "corrupted room token"})
		return
	}

	conn, err := a.upg.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot upgrade connection"})
		return
	}

	a.srv.HandleWS(user, roomID, conn)
}

func (a *App) oauthRedirect(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty oauth provider"})
		return
	}

	url, err := a.srv.GetOAuthRedirectURL(provider)
	if err != nil {
		if errors.Is(err, domain.ErrOAuthUnsupportedProvider) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported oauth provider"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "temporary cannot redirect to oauth provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

type oauthCallbackBody struct {
	Code string `json:"code"`
}

func (a *App) oauthCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty oauth provider"})
		return
	}

	var body oauthCallbackBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corrupted request body"})
		return
	}

	if body.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty code"})
		return
	}

	token, err := a.srv.RegisterUser(c.Request.Context(), provider, body.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "temporary cannot register user"})
		return
	}

	a.setTokenCookie(c, token)
}

func (a *App) user(c *gin.Context) *domain.User {
	data, _ := c.Get("user")
	return data.(*domain.User)
}
