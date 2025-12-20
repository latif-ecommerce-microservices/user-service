package customerror

import "errors"

func GetBaseError(err error) error {
	var ce Error
	if errors.As(err, &ce) {
		return ce.BaseError()
	}

	return err
}
