package middleware

import (
	"errors"
	"golang/internal/auth"
)

type Mid struct {
	auth *auth.Auth
}

/*
In this code, the Mid type is a middleware that depends on an auth.Auth instance, 
and the code is used to create and manage this middleware with proper error handling.
*/
func NewMid(a *auth.Auth) (Mid, error) {
	if a == nil {
		return Mid{}, errors.New("auth struct not provided")
	}
	return Mid{auth: a}, nil
}
