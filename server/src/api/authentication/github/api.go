package github

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/openmultiplayer/web/server/src/auth"
	"github.com/openmultiplayer/web/server/src/web"
)

type service struct {
	auth *auth.Authentication
	oa2  auth.OAuthProvider
}

func New(a *auth.Authentication, oa2 auth.OAuthProvider) *chi.Mux {
	rtr := chi.NewRouter()
	svc := service{
		auth: a,
		oa2:  oa2,
	}

	rtr.Get("/link", http.HandlerFunc(svc.link))
	rtr.Get("/callback", http.HandlerFunc(svc.callback))

	return rtr
}

type linkPayload struct {
	URL string `json:"url"`
}

func (s *service) link(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(linkPayload{URL: s.oa2.Link()})
}

func (s *service) callback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	user, err := s.oa2.Login(r.Context(), r.Form["code"][0])
	if err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	s.auth.EncodeAuthCookie(w, *user)
}