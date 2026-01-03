package customerror

type Code string

var (
	ErrUnexpected   Code = "unexpected_error"
	ErrClient       Code = "client_error"
	ErrNotFound     Code = "data_not_found"
	ErrUnauthorized Code = "unauthorized"
)

var (
	InternalServerError    = NewInternalServerError(ErrUnexpected, "Something Went Wrong")
	EmailAlreadyExistError = NewClientError(ErrClient, "Email Already Exist")
	UserNotFoundError      = NewNotFoundError(ErrNotFound, "User Not Found")
	InvalidCredentialError = NewUnauthorizedError(ErrUnauthorized, "Invalid email or password")
)
