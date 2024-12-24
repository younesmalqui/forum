package config

type CustomError struct {
	message    string
	isInternal bool
}

func (e *CustomError) Error() string {
	return e.message
}

func (e *CustomError) IsInternal() bool {
	return e.isInternal
}

func NewError(err error) *CustomError {
	if err == nil {
		return nil
	}
	return &CustomError{
		message:    err.Error(),
		isInternal: false,
	}
}

func NewInternalError(err error) *CustomError {
	if err == nil {
		return nil
	}
	return &CustomError{
		message:    err.Error(),
		isInternal: true,
	}
}
