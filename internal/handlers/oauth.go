package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ashil-poojary/gostruct/internal/config"
	"github.com/ashil-poojary/gostruct/internal/db"
	"github.com/ashil-poojary/gostruct/internal/models"
	"github.com/ashil-poojary/gostruct/internal/repository"
	"github.com/ashil-poojary/gostruct/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func NewMicrosoftOAuthHandler(cfg config.MicrosoftOAuthConfig) (gin.HandlerFunc, gin.HandlerFunc) {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURI,
		Scopes:       []string{"User.Read"},
		Endpoint:     microsoft.AzureADEndpoint(cfg.Tenant),
	}

	loginHandler := func(c *gin.Context) {
		url := oauthConfig.AuthCodeURL("random-state", oauth2.AccessTypeOffline)
		c.Redirect(http.StatusFound, url)
		// utils.Success(c, "redirecting to Microsoft", gin.H{"redirect_url": url})
	}

	callbackHandler := func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			utils.Error(c, http.StatusBadRequest, "missing 'code' in query params")
			return
		}

		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "token exchange failed: "+err.Error())
			return
		}

		client := oauthConfig.Client(context.Background(), token)
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

		userRepo := repository.NewUserRepository(db.DB1, db.DB2)
		user, err := userRepo.GetByEmail(email)
		if err != nil || user == nil {
			user = &models.User{
				Username:  msUser.DisplayName,
				Email:     email,
				Role:      utils.Constants.RoleUser,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := userRepo.CreateUser(user); err != nil {
				utils.Error(c, http.StatusInternalServerError, "failed to create user")
				return
			}
		}

		jwtToken, err := utils.GenerateAccessToken(int(user.ID), user.Email, user.Role)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to generate access token")
			return
		}

		utils.Success(c, "Microsoft login successful", gin.H{
			"access_token": jwtToken,
		})
	}

	return loginHandler, callbackHandler
}
