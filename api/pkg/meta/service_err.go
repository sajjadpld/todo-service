package meta

import (
	"errors"
	st "microservice/internal/server/http/status"
)

type (
	IServiceErr interface {
		Data(items map[string]any) *Error
		SetErr(err string) *Error
		Error() string
	}

	Error struct {
		Msg    st.HttpMappedStatus
		Err    error
		Detail map[string]any
	}
)

func ServiceErr(msg st.HttpMappedStatus, err ...error) *Error {
	se := &Error{Msg: msg, Err: errors.New("")}
	if len(err) > 0 {
		se.Err = err[0]
	}

	return se
}

func (svc *Error) Data(items map[string]any) *Error {
	svc.Detail = items
	return svc
}

func (svc *Error) SetErr(err string) *Error {
	svc.Err = errors.New(err)
	return svc
}

func (svc *Error) Error() string {
	return svc.Err.Error()
}
