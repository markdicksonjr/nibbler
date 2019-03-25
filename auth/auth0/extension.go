package auth0

// TODO: expose requiresRole?

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/session"
	"golang.org/x/oauth2"
	"net/http"
)

type Extension struct {
	nibbler.Extension

	SessionExtension *session.Extension

	// required redirect URL
	LoggedInRedirectUrl string

	// optional override for callback endpoint from auth0
	RoutePathAuth0Callback *string

	// optional override for login redirect to auth0
	RoutePathAuth0Login *string

	// optional override for logout
	RoutePathAuth0Logout *string

	config *nibbler.Configuration

	OnLoginComplete  func(s *Extension, w http.ResponseWriter, r *http.Request) (allowRedirect bool, err error)
	OnLogoutComplete func(s *Extension, w http.ResponseWriter, r *http.Request) error
}

func (s *Extension) Init(app *nibbler.Application) error {

	// assert that we have the session extension
	if s.SessionExtension == nil {
		return errors.New("session extension was not provided to Auth0 extension")
	}

	s.config = app.GetConfiguration()

	// default to "/" if not set
	if len(s.LoggedInRedirectUrl) == 0 {
		s.LoggedInRedirectUrl = "/"
	}

	return nil
}

func (s *Extension) AddRoutes(app *nibbler.Application) error {
	bestCallbackUrl := s.RoutePathAuth0Callback
	if bestCallbackUrl == nil {
		value := "/callback"
		bestCallbackUrl = &value
	}

	bestLoginUrl := s.RoutePathAuth0Login
	if bestLoginUrl == nil {
		value := "/login"
		bestLoginUrl = &value
	}

	bestLogoutUrl := s.RoutePathAuth0Logout
	if bestLogoutUrl == nil {
		value := "/logout"
		bestLogoutUrl = &value
	}

	app.GetRouter().HandleFunc(*bestCallbackUrl, s.CallbackHandler)
	app.GetRouter().HandleFunc(*bestLoginUrl, s.LoginHandler)
	app.GetRouter().HandleFunc(*bestLogoutUrl, s.LogoutHandler)
	return nil
}

func (s *Extension) Destroy(app *nibbler.Application) error {
	return nil
}

func (s *Extension) CallbackHandler(w http.ResponseWriter, r *http.Request) {

	rawConfig := *s.config.Raw
	domain := rawConfig.Get("auth0", "domain").String("")

	conf := &oauth2.Config{
		ClientID:     rawConfig.Get("auth0", "client", "id").String(""),
		ClientSecret: rawConfig.Get("auth0", "client", "secret").String(""),
		RedirectURL:  rawConfig.Get("auth0", "callback", "url").String(""),
		Scopes:       []string{"openid", "profile"}, // TODO: make configurable
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}
	state := r.URL.Query().Get("state")
	sessionState, stateErr := s.SessionExtension.GetAttribute(r, "state")
	if stateErr != nil {
		http.Error(w, stateErr.Error(), http.StatusInternalServerError)
		return
	}

	sessionStateString, ok := sessionState.(string)
	if !ok || state != sessionStateString {
		http.Error(w, "Invalid state parameter", http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get the userInfo
	client := conf.Client(context.TODO(), token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var profile map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//tokenErr := s.SessionExtension.SetAttribute(w, r, "id_token", token.Extra("id_token"))
	//
	//if tokenErr != nil {
	//	http.Error(w, tokenErr.Error(), http.StatusInternalServerError)
	//	return
	//}

	accessTokenErr := s.SessionExtension.SetAttribute(w, r, "access_token", token.AccessToken)

	if accessTokenErr != nil {
		http.Error(w, accessTokenErr.Error(), http.StatusInternalServerError)
		return
	}

	delete(profile, "picture")
	delete(profile, "updated_at")

	if profileErr := s.SessionExtension.SetAttribute(w, r, "profile", profile); profileErr != nil {
		http.Error(w, profileErr.Error(), http.StatusInternalServerError)
		return
	}

	if s.OnLoginComplete == nil {
		http.Redirect(w, r, s.LoggedInRedirectUrl, http.StatusSeeOther)
		return
	}

	allowRedirect, err := s.OnLoginComplete(s, w, r)

	// if allowRedirect is false, we assume the caller handled the error
	if allowRedirect {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, s.LoggedInRedirectUrl, http.StatusSeeOther)
	}
}

func (s *Extension) LoginHandler(w http.ResponseWriter, r *http.Request) {
	rawConfig := *s.config.Raw
	domain := rawConfig.Get("auth0", "domain").String("")
	aud := rawConfig.Get("auth0", "audience").String("")

	conf := &oauth2.Config{
		ClientID:     rawConfig.Get("auth0", "client", "id").String(""),
		ClientSecret: rawConfig.Get("auth0", "client", "secret").String(""),
		RedirectURL:  rawConfig.Get("auth0", "callback", "url").String(""),
		Scopes:       []string{"openid", "profile"}, // TODO: make configurable
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	if aud == "" {
		aud = "https://" + domain + "/userinfo"
	}

	// generate random state
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	err := s.SessionExtension.SetAttribute(w, r, "state", state)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	audience := oauth2.SetAuthURLParam("audience", aud)
	url := conf.AuthCodeURL(state, audience)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Extension) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	s.SessionExtension.SetAttribute(w, r, "profile", nil)

	if s.OnLogoutComplete != nil {
		err := s.OnLogoutComplete(s, w, r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"result": "` + err.Error() + `"}`))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "ok"}`))
}

func (s *Extension) EnforceLoggedIn(routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, err := s.SessionExtension.GetAttribute(r, "profile")

		if err != nil {
			// TODO: log
			w.WriteHeader(http.StatusNotFound)
			//w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`404 page not found`))
			return
		}

		if profile == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`404 page not found`))
			return
		}

		routerFunc(w, r)
	}
}
