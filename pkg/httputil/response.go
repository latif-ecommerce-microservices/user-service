package httputil

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/latif-ecommerce-microservices/user-service/pkg/customerror"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
)

type Response struct {
	Data    any    `json:"data,omitempty"`
	Detail  string `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type PaginatedResponse[T any] struct {
	Message string         `json:"message,omitempty"`
	Meta    PaginationMeta `json:"meta,omitempty"`
	Data    T              `json:"data,omitempty"`
}

type PaginationMeta struct {
	Count       int64 `json:"count,omitempty"`
	TotalPage   int64 `json:"total_page,omitempty"`
	CurrentPage int64 `json:"current_page,omitempty"`
	PerPage     int64 `json:"per_page,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Detail  any    `json:"detail,omitempty"`
}

type ResponseWriter struct {
	W http.ResponseWriter
}

type BasePaginationFilter struct {
	Page           int64  `query:"page"`
	Size           int64  `query:"size"`
	Sort           string `query:"sort"`
	OrganizationID string `query:"organization_id"`
	Search         string `query:"search"`
}

func Writer(w http.ResponseWriter) ResponseWriter {
	return ResponseWriter{W: w}
}

func (res ResponseWriter) JSON(status int, data any) error {
	res.W.Header().Set("Content-Type", "application/json")
	res.W.WriteHeader(status)
	err := json.NewEncoder(res.W).Encode(data)

	return err
}

func BuildPaginatedResponse[T any](data T, totalData, page, size int64) PaginatedResponse[T] {
	totalPage := int64(0)

	if size > 0 && totalData > 0 {
		totalPage = int64(math.Ceil(float64(totalData) / float64(size)))
	} else {
		totalPage = 1
	}

	return PaginatedResponse[T]{
		Data: data,
		Meta: PaginationMeta{
			Count:       totalData,
			TotalPage:   totalPage,
			CurrentPage: page,
			PerPage:     size,
		},
	}
}

func WriteSuccessPaginatedResponse[T any](rs http.ResponseWriter, resp PaginatedResponse[T], message string) {
	writer := Writer(rs)
	writer.JSON(http.StatusOK, PaginatedResponse[T]{
		Data:    resp.Data,
		Meta:    resp.Meta,
		Message: message,
	})
}

func WriteSuccessResponse(rs http.ResponseWriter, data any, message string) {
	writer := Writer(rs)
	writer.JSON(http.StatusOK, Response{
		Detail:  "success",
		Message: message,
		Data:    data,
	})
}

func WriteBadRequestResponse(w http.ResponseWriter, message string, detail any) {
	writer := Writer(w)
	writer.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    string(customerror.ErrClient),
		Message: message,
		Detail:  detail,
	})
}

func WriteRawResponse(w http.ResponseWriter, message string, data any) {
	writer := Writer(w)
	writer.JSON(http.StatusOK, data)
}

func HandleError(writer http.ResponseWriter, log *logging.Logger, err error) {
	w := Writer(writer)
	switch e := err.(type) {
	case customerror.Error:
		buildCustomErrorResponse(w, e, log)

	case validator.ValidationErrors:
		buildValidationErrorResponse(w, e)

	default:
		buildGenericErrorResponse(w, log, err)
	}
}

func buildCustomErrorResponse(w ResponseWriter, err error, log *logging.Logger) {
	parsedError, ok := err.(customerror.Error)
	if !ok {
		return
	}

	apiError := Response{
		Code:    parsedError.Code(),
		Message: parsedError.Error(),
	}

	respCode := parsedError.HTTPStatus()

	if respCode >= 500 {
		sendLog(log, parsedError)
	}

	w.JSON(respCode, apiError)
}

func sendLog(log *logging.Logger, err customerror.Error) {
	var attr []slog.Attr
	if err.Locator() != nil {
		attr = append(attr, slog.Any("location", err.Locator().String()))
	}

	if err.StackTrace() != nil {
		attr = append(attr, slog.Any("trace", err.StackTrace()))
	}
	log.Error(err.BaseError().Error(), attr...)
}

func buildGenericErrorResponse(w ResponseWriter, log *logging.Logger, err error) {
	detailError := slog.String("error", err.Error())
	log.Error("an unexpected error is happening", detailError)

	errorResponse := Response{
		Code:    string(customerror.ErrUnexpected),
		Message: "internal server error",
		Data:    nil,
	}

	w.JSON(500, errorResponse)
}

func buildValidationErrorResponse(w ResponseWriter, validationErrors validator.ValidationErrors) {
	errors := make([]map[string]string, len(validationErrors))
	for i, fieldErr := range validationErrors {
		errors[i] = map[string]string{
			"field":   fieldErr.Field(),
			"tag":     fieldErr.Tag(),
			"value":   fmt.Sprintf("%v", fieldErr.Value()),
			"message": fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldErr.Field(), fieldErr.Tag()),
		}
	}

	response := Response{
		Code:    string(customerror.ErrClient),
		Message: "invalid request body",
		Data:    errors,
	}

	w.JSON(http.StatusBadRequest, response)
}

func CredentialHeaderMissing(w http.ResponseWriter, logger *logging.Logger, whichHeader string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	err := json.NewEncoder(w).Encode(Response{
		Message: "Some credentials is missing",
	})

	logger.Error("missing credential header", slog.String("missing_header", whichHeader))

	return err
}
