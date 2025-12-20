package customerror

type Code string

var (
	ErrUnexpected Code = "unexpected_error"
	ErrClient     Code = "client_error"
	ErrNotFound   Code = "data_not_found"
)

var (
	InternalServerError    = NewInternalServerError(ErrUnexpected, "Something Went Wrong")
	EmailAlreadyExistError = NewClientError(ErrClient, "Email Already Exist")
)
