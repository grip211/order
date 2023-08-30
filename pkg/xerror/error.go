package xerror

type PublicError struct {
	Err error
}

func ThrowPublicError(context string, err error) error {
	return Wrap(context, &PublicError{
		Err: err,
	})
}

func (e *PublicError) Error() string {
	return e.Err.Error()
}

type PrivateError struct {
	Err error
}

func ThrowPrivateError(context string, err error) error {
	return Wrap(context, &PrivateError{
		Err: err,
	})
}

func (e *PrivateError) Error() string {
	return e.Err.Error()
}
