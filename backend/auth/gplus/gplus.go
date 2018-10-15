package gplus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	auth "github.com/dwarvesf/smithy/backend/auth"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
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
		ClientKey:   clientKey,
		Secret:      secret,
		CallbackURL: "postmessage",
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
	return auth.HTTPClientWithFallBack(p.HTTPClient)
}

// GetAuthURL get url to direct to permission form
func (p *Provider) GetAuthURL() string {
	return p.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

//FetchUser will go to Google+ and access basic information about the user.
func (p *Provider) FetchUser(tok *oauth2.Token) (backendConfig.Email, error) {
	user := auth.UserInfo{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    tok.Expiry.String(),
	}

	email := backendConfig.Email{}

	if user.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return email, errors.New("gplus cannot get user information without accessToken")
	}

	response, err := p.Client().Get(endpointProfile + "?access_token=" + url.QueryEscape(user.AccessToken))
	if err != nil {
		return email, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return email, fmt.Errorf("gplus responded with a %d trying to fetch user information", response.StatusCode)
	}

	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return email, err
	}

	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
	if err != nil {
		return email, err
	}

	err = userFromReader(bytes.NewReader(bits), &user)

	email = backendConfig.Email{
		ID:           user.UserID,
		Name:         user.Email,
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		ExpiresAt:    tok.Expiry.String(),
	}

	return email, err
}

func userFromReader(reader io.Reader, user *auth.UserInfo) error {
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
func (p *Provider) CompleteUserAuth(code string) (backendConfig.Email, error) {
	// get new token and retry fetch
	tok, err := p.Exchange(code)
	if err != nil {
		return backendConfig.Email{}, err
	}
	return p.FetchUser(tok)
}

//Exchange to get new token
func (p *Provider) Exchange(code string) (*oauth2.Token, error) {
	token, err := p.Config.Exchange(auth.ContextForClient(p.Client()), code)
	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.New("Invalid token received from provider")
	}

	return token, err
}
