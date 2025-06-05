package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ashil-poojary/gostruct/internal/config"
	"github.com/ashil-poojary/gostruct/internal/db"
	"github.com/ashil-poojary/gostruct/internal/models"
	"github.com/ashil-poojary/gostruct/internal/repository"
	"github.com/ashil-poojary/gostruct/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func NewMicrosoftOAuthHandler(cfg config.MicrosoftOAuthConfig) (http.HandlerFunc, http.HandlerFunc) {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURI,
		Scopes:       []string{"User.Read"},
		Endpoint:     microsoft.AzureADEndpoint(cfg.Tenant),
	}

	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		log.Println("[OAuth] Starting Microsoft login flow")
		url := oauthConfig.AuthCodeURL("random-state", oauth2.AccessTypeOffline)
		log.Printf("[OAuth] Redirecting to Microsoft: %s\n", url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}

	callbackHandler := func(w http.ResponseWriter, r *http.Request) {
		log.Println("[OAuth] Handling Microsoft callback")

		code := r.URL.Query().Get("code")
		if code == "" {
			log.Println("[OAuth] Missing 'code' in query parameters")
			http.Error(w, "Missing 'code' in query params", http.StatusBadRequest)
			return
		}
		log.Printf("[OAuth] Received auth code: %s\n", code)

		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			log.Printf("[OAuth] Token exchange failed: %v\n", err)
			http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("[OAuth] Token exchange successful")

		client := oauthConfig.Client(context.Background(), token)
		resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
		if err != nil {
			log.Printf("[OAuth] Failed to get user info from Microsoft Graph: %v\n", err)
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		log.Println("[OAuth] Successfully fetched user info from Microsoft Graph")

		var msUser struct {
			DisplayName       string `json:"displayName"`
			Mail              string `json:"mail"`
			UserPrincipalName string `json:"userPrincipalName"`
			ID                string `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&msUser); err != nil {
			log.Printf("[OAuth] Failed to parse user info: %v\n", err)
			http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
			return
		}
		log.Printf("[OAuth] Parsed user info: %+v\n", msUser)

		email := msUser.Mail
		if email == "" {
			email = msUser.UserPrincipalName
		}
		log.Printf("[OAuth] Resolved user email: %s\n", email)

		userRepo := repository.NewUserRepository(db.DB1, db.DB2)

		user, err := userRepo.GetByEmail(email)
		if err != nil {
			log.Printf("[OAuth] User not found, creating new user for %s\n", email)
			user = &models.User{
				Username:  msUser.DisplayName,
				Email:     email,
				Role:      utils.Constants.RoleUser,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := userRepo.CreateUser(user); err != nil {
				log.Printf("[OAuth] Failed to create user: %v\n", err)
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
			log.Printf("[OAuth] Created new user: %s\n", email)
		} else {
			log.Printf("[OAuth] Existing user found: %s\n", email)
		}

		jwtToken, err := utils.GenerateAccessToken(int(user.ID), user.Email, user.Role)
		if err != nil {
			log.Printf("[OAuth] Failed to generate JWT token: %v\n", err)
			http.Error(w, "JWT generation failed", http.StatusInternalServerError)
			return
		}
		log.Printf("[OAuth] JWT token generated for user %s\n", user.Email)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"token": jwtToken}); err != nil {
			log.Printf("[OAuth] Failed to write response: %v\n", err)
		} else {
			log.Println("[OAuth] Login flow completed successfully")
		}
	}

	return loginHandler, callbackHandler
}
