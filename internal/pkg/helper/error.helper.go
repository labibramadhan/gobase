package helper

import (
	"database/sql"
	"errors"

	pkgErr "clodeo.tech/public/go-universe/pkg/err"
	"google.golang.org/grpc/status"
)

/*
If there's only one message (index 0), it is the general error message
*/
func getMessage(s []string) string {
	msg := ""
	for i, v := range s {
		if i == 0 {
			msg = v
			break
		}
	}

	return msg
}

func NewBadRequestErr(err error, s ...string) pkgErr.CustErr {
	msg := getMessage(s)
	return pkgErr.NewError(pkgErr.ErrorOpts{
		Cause:   err,
		Message: msg,
		Type:    pkgErr.ErrBadRequest,
	})
}

func NewUnauthorizedErr(err error, s ...string) pkgErr.CustErr {
	msg := getMessage(s)
	return pkgErr.NewError(pkgErr.ErrorOpts{
		Cause:   err,
		Message: msg,
		Type:    pkgErr.ErrUnauthorized,
	})
}

func NewInternalServerErr(err error, s ...string) pkgErr.CustErr {
	msg := getMessage(s)
	return pkgErr.NewError(pkgErr.ErrorOpts{
		Cause:   err,
		Message: msg,
		Type:    pkgErr.ErrInternal,
	})
}

func NewForbiddenErr(err error, s ...string) pkgErr.CustErr {
	msg := getMessage(s)
	return pkgErr.NewError(pkgErr.ErrorOpts{
		Cause:   err,
		Message: msg,
		Type:    pkgErr.ErrForbidden,
	})
}

func NewNotFoundErr(err error, s ...string) pkgErr.CustErr {
	msg := getMessage(s)
	return pkgErr.NewError(pkgErr.ErrorOpts{
		Cause:   err,
		Message: msg,
		Type:    pkgErr.ErrNotFound,
	})
}

func NewSqlErr(err error, s ...string) pkgErr.CustErr {
	msg := getMessage(s)
	errOpts := pkgErr.ErrorOpts{
		Cause:   err,
		Message: msg,
		Type:    pkgErr.ErrInternal,
	}

	if errors.Is(err, sql.ErrNoRows) {
		errOpts.Message = "data not found"
		errOpts.Type = pkgErr.ErrNotFound
	}

	return pkgErr.NewError(errOpts)
}

func NewGrpcError(err error) error {
	e := pkgErr.GetError(err)
	return status.Error(e.GRPCCode, e.Error())
}
