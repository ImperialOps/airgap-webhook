package main

type ApiError struct {
    code int
    error string
}

func NewApiError(code int, error string) *ApiError {
    return &ApiError{
        code: code,
        error: error,
    }
}

func (e *ApiError) Error() string {
    return e.error
}

func (e *ApiError) Code() int {
    return e.code
}

