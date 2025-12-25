package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/latif-ecommerce-microservices/user-service/internal/service/user"
	"github.com/latif-ecommerce-microservices/user-service/pkg/httputil"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
)

type UserHandler struct {
	userService user.ServiceProvider
	validator   *validator.Validate
	logger      *logging.Logger
}

func NewUserHandler(
	userService user.ServiceProvider,
	validator *validator.Validate,
	logger *logging.Logger,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
		logger:      logger,
	}
}

func (h *UserHandler) RegisterRoutes(router chi.Router) {
	router.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetByID)
	})
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteBadRequestResponse(w, "invalid JSON body", err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		httputil.HandleError(w, h.logger, err)
		return
	}

	res, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		httputil.HandleError(w, h.logger, err)
		return
	}

	httputil.WriteSuccessResponse(w, res, "user created successfully")
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.WriteBadRequestResponse(w, "invalid user id", err.Error())
		return
	}

	res, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		httputil.HandleError(w, h.logger, err)
		return
	}

	httputil.WriteSuccessResponse(w, res, "user retrieved successfully")
}
