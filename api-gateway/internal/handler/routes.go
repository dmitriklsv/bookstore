package handler

import (
	"net/http"

	"github.com/Levap123/api_gateway/internal/middlwares"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) InitRoutes() http.Handler {
	r := httprouter.New()
	r.Handler(http.MethodPost, "/auth/sign-up", middlwares.CheckErrorMiddlware(h.signUp))
	r.Handler(http.MethodPost, "/auth/sign-in", middlwares.CheckErrorMiddlware(h.signIn))
	r.Handler(http.MethodGet, "/api/user", middlwares.CheckErrorMiddlware(h.getMe))
	r.Handler(http.MethodGet, "/api/user/:user_id", middlwares.CheckErrorMiddlware(h.getByID))

	return r
}
