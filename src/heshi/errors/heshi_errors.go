package errors

import (
	"fmt"
	"strings"
)

type (
	NetworksError struct {
		HeshiError
	}

	DBError struct {
		HeshiError
	}
)

func NewNetworksError(msg string) NetworksError {
	return NetworksError{New(msg)}
}

func WrapNetworksError(err error, msg string) NetworksError {
	return NetworksError{Wrap(err, msg)}
}

func WrapNetworksResponseError(err error, URL string, statusCode int, body []byte, msg ...string) error {
	if err != nil {
		return err
	}
	c := fmt.Sprintf("%s; invalid response", URL)
	if len(msg) > 0 {
		c = msg[0]
	}
	m := fmt.Sprintf("%s; status code[%d]", c, statusCode)
	if len(body) != 0 {
		m += fmt.Sprintf(" with body[%s]", strings.TrimSpace(string(body)))
	}
	return NewNetworksError(m)
}

func UnWrapError(err error) error {
	switch err.(type) {
	case HeshiError:
		e := err.(HeshiError)
		return e.GetInner()
	default:
		return err
	}
}

func WrapDBError(err error, q string, params ...interface{}) DBError {
	return DBError{Wrapf(err, "SQL: "+q+", %s", params...)}
}

func WrapDBErrorf(err error, format string, a ...interface{}) DBError {
	return DBError{Wrapf(err, format, a...)}
}

func NewNetworksErrorf(format string, a ...interface{}) NetworksError {
	return NetworksError{Newf(format, a...)}
}
