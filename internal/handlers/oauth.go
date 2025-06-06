package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ashil-poojary/gostruct/internal/config"
	"github.com/ashil-poojary/gostruct/internal/models"
	"github.com/ashil-poojary/gostruct/internal/repository"
	"github.com/ashil-poojary/gostruct/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type OAuthHandler struct {
	UserRepo         *repository.UserRepository
	RefreshTokenRepo *repository.AuthRepo
	Config           *oauth2.Config
}

func NewOAuthHandler(userRepo *repository.UserRepository, rtRepo *repository.AuthRepo, cfg config.MicrosoftOAuthConfig) *OAuthHandler {
	return &OAuthHandler{
		UserRepo:         userRepo,
		RefreshTokenRepo: rtRepo,
		Config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURI,
			Scopes:       []string{"User.Read"},
			Endpoint:     microsoft.AzureADEndpoint(cfg.Tenant),
		},
	}
}

func (h *OAuthHandler) MicrosoftLogin(c *gin.Context) {
	log.Println("Redirecting to Microsoft OAuth login")

	if h.Config.ClientID == "" || h.Config.ClientSecret == "" || h.Config.RedirectURL == "" {
		utils.Error(c, http.StatusInternalServerError, "Microsoft OAuth configuration is incomplete")
		log.Printf("OAuth Config: %+v", h.Config)
		return
	}

	url := h.Config.AuthCodeURL("random-state", oauth2.AccessTypeOffline)

	c.Redirect(http.StatusFound, url)
}

func (h *OAuthHandler) MicrosoftCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.Error(c, http.StatusBadRequest, "missing 'code' in query params")
		return
	}

	token, err := h.Config.Exchange(context.Background(), code)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "token exchange failed: "+err.Error())
		return
	}

	client := h.Config.Client(context.Background(), token)
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to fetch user info")
		return
	}
	defer resp.Body.Close()

	var msUser struct {
		DisplayName       string `json:"displayName"`
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
		ID                string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&msUser); err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to parse user info")
		return
	}

	email := msUser.Mail
	if email == "" {
		email = msUser.UserPrincipalName
	}

	user, err := h.UserRepo.GetByEmail(email)
	if err != nil || user == nil {
		user = &models.User{
			Username:  msUser.DisplayName,
			Email:     email,
			Role:      utils.Constants.RoleUser,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := h.UserRepo.CreateUser(user); err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to create user")
			return
		}
	}

	accessToken, err := utils.GenerateAccessToken(int(user.ID), user.Email, user.Role)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(int(user.ID), user.Email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()
	rtExpiry := time.Now().Add(7 * 24 * time.Hour)

	if err := h.RefreshTokenRepo.SaveRefreshToken(refreshToken, user.ID, userAgent, ip, rtExpiry); err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to save refresh token")
		return
	}

	// Redirect to frontend with token and email as query params
	redirectURL := fmt.Sprintf("http://localhost:3000?token=%s&email=%s", accessToken, user.Email)
	c.Redirect(http.StatusFound, redirectURL)
}
