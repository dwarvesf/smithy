package gplus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/dwarvesf/smithy/backend/domain"
)

const (
	authURL         string = "https://accounts.google.com/o/oauth2/auth"
	tokenURL        string = "https://accounts.google.com/o/oauth2/token"
	endpointProfile string = "https://www.googleapis.com/oauth2/v2/userinfo"
	state                  = "state-token"
)

// NewProvider creates a new Google+ provider, and sets up important connection details.
// You should always call `NewProvider` to get a new Provider. Never try to create
// one manually.
func NewProvider(clientKey, secret string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey: clientKey,
		Secret:    secret,
	}
	p.Config = newConfig(p)
	return p
}

// Provider is the implementation of `auth.Provider` for accessing Google+.
type Provider struct {
	ClientKey   string
	Secret      string
	CallbackURL string
	HTTPClient  *http.Client
	Config      *oauth2.Config
}

func (p *Provider) Client() *http.Client {
	return HTTPClientWithFallBack(p.HTTPClient)
}

// GetAuthURL get url to direct to permission form
func (p *Provider) GetAuthURL() string {
	return p.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

//FetchUser will go to Google+ and access basic information about the user.
func (p *Provider) FetchUser(tok *oauth2.Token) (*domain.User, error) {
	userInfo := UserInfo{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    tok.Expiry.String(),
	}

	if userInfo.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return &domain.User{}, errors.New("gplus cannot get user information without accessToken")
	}

	response, err := p.Client().Get(endpointProfile + "?access_token=" + url.QueryEscape(userInfo.AccessToken))
	if err != nil {
		return &domain.User{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return &domain.User{}, fmt.Errorf("gplus responded with a %d trying to fetch user information", response.StatusCode)
	}

	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &domain.User{}, err
	}

	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&userInfo.RawData)
	if err != nil {
		return &domain.User{}, err
	}

	err = userFromReader(bytes.NewReader(bits), &userInfo)
	if err != nil {
		return &domain.User{}, err
	}

	return &domain.User{
		Email:          userInfo.Email,
		IsEmailAccount: true,
		Role:           "user",
	}, nil
}

func userFromReader(reader io.Reader, user *UserInfo) error {
	u := struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		FirstName string `json:"given_name"`
		LastName  string `json:"family_name"`
		Link      string `json:"link"`
		Picture   string `json:"picture"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.Name = u.Name
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.NickName = u.Name
	user.Email = u.Email
	user.AvatarURL = u.Picture
	user.UserID = u.ID

	return err
}

func newConfig(provider *Provider) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{
			"profile",
			"email",
			"https://www.googleapis.com/auth/plus.login",
		},
	}
}

//CompleteUserAuth complete the auth process
func (p *Provider) CompleteUserAuth(code, redirectURL string) (*domain.User, error) {
	// get new token and retry fetch
	p.Config.RedirectURL = redirectURL
	tok, err := p.Exchange(code)
	if err != nil {
		return &domain.User{}, err
	}
	return p.FetchUser(tok)
}

//Exchange to get new token
func (p *Provider) Exchange(code string) (*oauth2.Token, error) {
	token, err := p.Config.Exchange(ContextForClient(p.Client()), code)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.New("Invalid token received from provider")
	}

	return token, err
}

// ContextForClient provides a context for use with oauth2.
func ContextForClient(h *http.Client) context.Context {
	if h == nil {
		return context.Background()
	}
	return context.WithValue(context.Background(), oauth2.HTTPClient, h)
}

// HTTPClientWithFallBack to be used in all fetch operations.
func HTTPClientWithFallBack(h *http.Client) *http.Client {
	if h != nil {
		return h
	}
	return http.DefaultClient
}

// UserInfo contains the information common amongst most OAuth and OAuth2 providers.
// All of the "raw" datafrom the provider can be found in the `RawData` field.
type UserInfo struct {
	RawData           map[string]interface{}
	Provider          string
	Email             string
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	Description       string
	UserID            string
	AvatarURL         string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         string
}
