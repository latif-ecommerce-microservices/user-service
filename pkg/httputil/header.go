package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	UserIdHeader         = "User-ID"
	OrganizationIdHeader = "Organization-ID"
	UserRoleHeader       = "User-Role"
)

func GetHeaderValue[T any](r *http.Request, key string) (T, error) {
	var zero T
	header := r.Header
	if header == nil {
		return zero, errors.New("header is nil")
	}

	data := header.Get(key)
	if data == "" {
		return zero, errors.New("header value not found")
	}

	var value any
	var err error

	switch any(zero).(type) {
	case int:
		value, err = strconv.Atoi(data)
		if err != nil {
			return zero, errors.New(fmt.Sprintf("%s should be int", key))
		}
	case string:
		value = data
	case float64:
		value, err = strconv.ParseFloat(data, 64)
		if err != nil {
			return zero, errors.New(fmt.Sprintf("%s should be float64", key))
		}
	default:
		return zero, errors.New("unsupported type")
	}

	return value.(T), nil
}

func GetUserIDFromHeader(r *http.Request) (string, error) {
	userID, err := GetHeaderValue[string](r, UserIdHeader)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetOrgIDFromHeader(r *http.Request) (string, error) {
	orgID, err := GetHeaderValue[string](r, OrganizationIdHeader)
	if err != nil {
		return "", err
	}

	return orgID, nil
}

func GetUserRoleFromHeader(r *http.Request) (string, error) {
	userRole, err := GetHeaderValue[string](r, UserRoleHeader)
	if err != nil {
		return "", err
	}

	return userRole, nil
}

func GetUserFullNameFromHeader(r *http.Request) (string, error) {
	fullName, err := GetHeaderValue[string](r, "User-Full-Name")
	if err != nil {
		return "", err
	}

	return fullName, nil
}
