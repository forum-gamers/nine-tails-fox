package helpers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func NewAppError(code codes.Code, msg string) error {
	return status.Error(code, msg)
}
