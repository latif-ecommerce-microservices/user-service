package grpc

import (
	"errors"
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/latif-ecommerce-microservices/user-service/pkg/customerror"
)

func toGRPCError(err error) error {
	var ce customerror.Error
	if errors.As(err, &ce) {
		st := status.New(mapGRPCCode(ce), ce.Error())

		st, _ = st.WithDetails(
			&errdetails.ErrorInfo{
				Reason: ce.Code(),
			},
		)

		return st.Err()
	}

	return status.Error(codes.Internal, "internal server error")
}

func mapGRPCCode(err customerror.Error) codes.Code {
	switch err.HTTPStatus() {
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusBadRequest:
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}
