package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type GoogleOAuthInterface interface {
	GetConfig() *oauth2.Config
	GenerateState() string
	GetAuthURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error)
}

type googleOAuthImpl struct {
	config *oauth2.Config
}

var GoogleOAuth GoogleOAuthInterface = NewGoogleOAuth()

func NewGoogleOAuth() GoogleOAuthInterface {
	redirectURI := env.AppEnv.GoogleRedirectURI
	if env.AppEnv.AppEnv == "production" {
		redirectURI = env.AppEnv.FEAppURL
	}

	cfg := &oauth2.Config{
		ClientID:     env.AppEnv.GoogleClientID,
		ClientSecret: env.AppEnv.GoogleClientSecret,
		RedirectURL:  redirectURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &googleOAuthImpl{config: cfg}
}

func (g *googleOAuthImpl) GetConfig() *oauth2.Config {
	return g.config
}

func (g *googleOAuthImpl) GenerateState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("oauth: crypto/rand.Read failed: %s", err.Error()))
	}
	return hex.EncodeToString(b)
}

func (g *googleOAuthImpl) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (g *googleOAuthImpl) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

func (g *googleOAuthImpl) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status: %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}
