package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"
	"github.com/latif-ecommerce-microservices/user-service/internal/service/auth"
	"github.com/latif-ecommerce-microservices/user-service/pkg/httputil"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
	"net/http"
)

type AuthHandler struct {
	authService auth.ServiceProvider
	validator   *validator.Validate
	logger      *logging.Logger
}

func NewAuthHandler(
	authService auth.ServiceProvider,
	validator *validator.Validate,
	logger *logging.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
		logger:      logger,
	}
}

func (h *AuthHandler) RegisterRoutes(router chi.Router) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		//r.Post("/logout", h.Logout)
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteBadRequestResponse(w, "invalid JSON body", err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		httputil.HandleError(w, h.logger, err)
		return
	}

	res, err := h.authService.Login(r.Context(), req)
	if err != nil {
		httputil.HandleError(w, h.logger, err)
		return
	}

	httputil.WriteSuccessResponse(w, res, "Login successful")
}

//func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
//	var req dto.LogoutRequest
//
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		httputil.WriteBadRequestResponse(w, "invalid JSON body", err.Error())
//		return
//	}
//
//	if err := h.validator.Struct(req); err != nil {
//		httputil.HandleError(w, h.logger, err)
//		return
//	}
//
//	res, err := h.authService.Logout(r.Context(), req)
//	if err != nil {
//		httputil.HandleError(w, h.logger, err)
//		return
//	}
//
//	httputil.WriteSuccessResponse(w, res, "Login successful")
//}
